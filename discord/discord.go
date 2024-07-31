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

func NewDiscord() *Discord { // 関数名をエクスポートするように変更
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	return &Discord{
		os.Getenv("DISCORD_TOKEN"),
		os.Getenv("DISCORD_CHANNEL_ID"),
	}
}

func (b *Discord) Send(message string) {
	dg, _ := discordgo.New("Bot " + b.token)
	_, err := dg.ChannelMessageSend(b.channelID, message)
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}
