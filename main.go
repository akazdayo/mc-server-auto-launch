package main

import (
	"fmt"
	"os"
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
	go s.LaunchMinecraft(os.Args[1])
	go s.LaunchSSNet(os.Args[2])
	time.Sleep(5 * time.Second)
	s.QuitServer()
	time.Sleep(2 * time.Second)
	fmt.Println("Server has stopped")
}
