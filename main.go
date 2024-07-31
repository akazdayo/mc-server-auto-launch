package main

import "github.com/akazdayo/mc-server-auto-launch/server"

const (
	DiscordToken     = "Discord Token Here"
	DiscordChannelId = "Discord Channel ID Here"
	ServerPath       = "Server Path Here"
	SSNetPath        = "SSNet Path Here"
)

func main() {
	quit := make(chan bool)
	isRunning := make(chan bool)
	controlURL := make(chan string)
	serverIP := make(chan string)

	s := server.NewServer(quit, isRunning, controlURL, serverIP)
	go s.LaunchMinecraft("")
	go s.LaunchSSNet("")
}
