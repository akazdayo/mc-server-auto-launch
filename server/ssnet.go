package server

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/akazdayo/mc-server-auto-launch/discord"
)

type SSNet struct {
	disbot *discord.Discord
	stop   chan struct{}
}

func NewSSNet() *SSNet {
	return &SSNet{
		discord.NewDiscord(),
		make(chan struct{}),
	}
}

func (s *SSNet) LaunchSSNet(path string) {
	fmt.Println("Starting Secure Share Net")
	cmd := exec.Command("sh", path)

	// 標準出力を取得
	stdout, _ := cmd.StdoutPipe()

	// コマンドを開始
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		return
	}

	// 標準出力をリアルタイムで読み取る
	go getOutput(stdout, s.checkSSNOutput)

	// 終了時処理
	<-s.stop
	stdout.Close()
	fmt.Println("Secure Share Net has stopped")
}

// 出力を加工する関数
func (s *SSNet) checkSSNOutput(output string) string {
	if strings.Contains(output, "コントロールURL") {
		//s.disbot.Send(output)
	} else if strings.Contains(output, "公開開始") {
		//s.disbot.Send(output)
	}
	return fmt.Sprintf("[%s] %s\n", time.Now(), output)
}

func (s *SSNet) Quit() {
	close(s.stop)
}
