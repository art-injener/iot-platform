package main

import "github.com/art-injener/iot-platform/cmd/server_tcp/tcp"

func main() {
	// читает данные из tcp и складывает в rabbitMQ

	var srv = tcp.ServerTCP{
		Addr: ":9000",
	}

	err := srv.Run()
	if err != nil {
		return
	}
}
