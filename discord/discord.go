package discord

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Discord struct {
	token     string
	channelID string
}

func NewDiscord() *Discord {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	token := os.Getenv("DISCORD_TOKEN")
	channelID := os.Getenv("DISCORD_CHANNEL_ID")

	if token == "" || channelID == "" {
		fmt.Println("Environment variables DISCORD_TOKEN or DISCORD_CHANNEL_ID are not set")
		os.Exit(1)
	}

	return &Discord{
		token:     token,
		channelID: channelID,
	}
}

func (b *Discord) Send(message string) {
	dg, _ := discordgo.New("Bot " + b.token)
	_, err := dg.ChannelMessageSend(b.channelID, message)
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}
