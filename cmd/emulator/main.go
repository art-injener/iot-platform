package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/art-injener/iot-platform/internal/config"
	"github.com/art-injener/iot-platform/internal/imitation/httpserver/server"
	lg "github.com/art-injener/iot-platform/pkg/logger"
)

// имитирует данные и записывает в tcp порт

func main() { // читаем конфигурационные настройки
	cfg, err := config.GetConfig("./configs")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cfg.Log = lg.NewConsole(cfg.LogLevel == config.DebugLevel)

	// Graceful Shutdown механиз завершения работы приложения
	// канал возврата кода завершения приложения
	var returnCode = make(chan int)

	// канал для передеачи в горутины сигнала о необходимости завершить свою работу
	var finish = make(chan struct{})

	// канал подтверждения остановки горутин
	var done = make(chan struct{})

	// обработчик сигналов ОС
	go signalHandler(returnCode, finish, done)

	var web *server.WebApp

	web, err = server.NewApp(cfg)
	if err != nil {
		cfg.Log.Error().Msgf("Error create imitation. Error  in : %s", err)
		os.Exit(1)
	}
	web.Run(finish)

	// Подтверждение об остановке и корректном завершении программы
	done <- struct{}{}

	web.Stop()
	os.Exit(<-returnCode)
}

// signalHandler - обработчик сигналов операционной системы
func signalHandler(returnCode chan int, finish chan struct{}, done chan struct{}) {

	// signals  - канал для перехвата системных сигналов завершения
	var signals = make(chan os.Signal, 1)

	// делаем подписку
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)

	// блокируемся до прихода системного сигнала завершения приложения
	sig := <-signals

	fmt.Printf("caught sig: %+v", sig)

	// посылаем остальным горутинам сообщение о неоходимости завершиться
	finish <- struct{}{}

	// ждем когда остальные горутины дадут подтвержение о своем завершении,
	// если в течении 30 секунд они не ответили - возвращем значение кода != 0
	select {
	case <-time.After(30 * time.Second):
		returnCode <- 1
	case <-done:
		returnCode <- 0
	}
}
