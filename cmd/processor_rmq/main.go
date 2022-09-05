package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/art-injener/iot-platform/pkg/models/rmq"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("device_info", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {

			dm := rmq.DeviceMessageModel{}

			err := json.Unmarshal(d.Body, &dm)
			if err != nil {
				return
			}
			log.Printf("Received a message: %+v", dm)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			}
		}
	}()

	// Stop for program termination
	<-stopChan
}
