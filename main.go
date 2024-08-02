package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akazdayo/mc-server-auto-launch/server"
)

func killTimer(m *server.Minecraft) {
	for {
		time.Sleep(10 * time.Second)
		now := time.Now()
		if now.Hour() >= 1 && now.Hour() < 5 {
			m.Quit(0)
			time.Sleep(20 * time.Second)
			fmt.Println("Server has stopped")
		}
	}
}

func main() {
	m := server.NewMinecraft()
	s := server.NewSSNet()

	go m.LaunchMinecraft(os.Args[1])
	go s.LaunchSSNet(os.Args[2])
	go killTimer(m)

	for {
		status := <-m.Stop
		if status >= 1 {
			time.Sleep(20 * time.Second)
			m = server.NewMinecraft()
			go m.LaunchMinecraft(os.Args[1])
		} else {
			break
		}
	}
	s.Quit()
	time.Sleep(20 * time.Second)
	fmt.Println("Server has stopped")
}
