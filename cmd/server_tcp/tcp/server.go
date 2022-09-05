package tcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/streadway/amqp"

	"github.com/art-injener/iot-platform/pkg/models/rmq"
)

// ServerTCP - реализация
type ServerTCP struct {
	Addr   string
	server net.Listener
	rmq    *amqp.Channel
	queue  amqp.Queue
}

// Run Запуск сервера.
// Операция блокирующая. Сервер в бесконечном цикле ждёт новых подключений
func (t *ServerTCP) Run() (err error) {
	t.server, err = net.Listen("tcp", t.Addr)

	if err != nil {
		return err
	}

	defer t.Close()

	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	if err != nil {
		return err
	}
	defer conn.Close()

	t.rmq, err = conn.Channel()
	if err != nil {
		return err
	}

	t.queue, err = t.rmq.QueueDeclare("device_info", true, false, false, false, nil)
	if err != nil {
		return err
	}

	return t.waitNewConnections()
}

// Закрытие сервера
func (t *ServerTCP) Close() (err error) {
	defer t.rmq.Close()
	return t.server.Close()
}

// функция ожидания новых подключений.
// Операция блокирующая. Сервер в бесконечном цикле ждёт новых подключений
func (t *ServerTCP) waitNewConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			return err
		}
		// обработчик каждого нового подключения запускается в новой горутине
		go t.handleConnection(conn)
	}
}

func (t *ServerTCP) handleConnection(conn net.Conn) {
	defer conn.Close()

	name := conn.RemoteAddr().String()
	log.Printf("%+v connected to server\n", name)

	// Будем прослушивать все сообщения разделенные \n
	message, err := bufio.NewReader(conn).ReadString('*')
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Распечатываем полученое сообщение
	fmt.Print("Message Received:", string(message), "\n")

	msg, err := rmq.NewDeviceMessage(message)
	if err != nil {
		log.Println(err.Error())
		return
	}

	rawMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = t.rmq.Publish("", t.queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         rawMsg,
		})

	if err != nil {
		log.Fatalf("basic.publish: %v", err)
	}
}
