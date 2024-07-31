package main

import (
	"fmt"
	"time"

	"github.com/akazdayo/mc-server-auto-launch/server"
)

const (
	DiscordToken     = "Discord Token Here"
	DiscordChannelId = "Discord Channel ID Here"
	ServerPath       = "Server Path Here"
	SSNetPath        = "SSNet Path Here"
)

func main() {
	isRunning := make(chan bool)
	controlURL := make(chan string)
	serverIP := make(chan string)

	s := server.NewServer(isRunning, controlURL, serverIP)
	go s.LaunchMinecraft("./test.sh")
	go s.LaunchSSNet("./ssnet.sh")
	time.Sleep(10 * time.Second)
	fmt.Println("Stopping server")
	s.QuitServer()
	time.Sleep(3 * time.Second)
	fmt.Println("Server has stopped")
}
