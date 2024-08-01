package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/akazdayo/mc-server-auto-launch/discord"
)

type Server struct {
	disbot     *discord.Discord
	stop       chan struct{}
	wg         sync.WaitGroup
	isRunning  chan bool
	controlURL chan string
	serverIP   chan string
}

func NewServer(isRunning chan bool, controlURL chan string, serverIP chan string) *Server {
	return &Server{
		discord.NewDiscord(),
		make(chan struct{}),
		sync.WaitGroup{},
		isRunning,
		controlURL,
		serverIP,
	}
}

func (s *Server) LaunchMinecraft(path string) {
	s.wg.Add(1)
	defer s.wg.Done()
	fmt.Println("Starting Minecraft")
	cmd := exec.Command("sh", path)

	// 入出力を取得
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	// コマンドを開始
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		return
	}

	// 出力を表示
	go s.getOutput(stdout, s.checkSSNOutput)
	go s.sendCommand(stdin)
	go s.saveCommand(stdin)

	//終了時処理
	<-s.stop
	fmt.Println("Stopping Minecraft")
	io.WriteString(stdin, "stop\n") // stopコマンドを送信
	stdin.Close()
	cmd.Wait()

	fmt.Println("Minecraft has stopped")

}

func (s *Server) LaunchSSNet(path string) {
	s.wg.Add(1)
	defer s.wg.Done()
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
	go s.getOutput(stdout, s.checkSSNOutput)

	// 終了時処理
	<-s.stop
	stdout.Close()
	fmt.Println("Secure Share Net has stopped")
}

func (s *Server) getOutput(stdout io.Reader, callback func(string) string) {
	reader := bufio.NewReader(stdout)
	for { // 書いたのCopilotだから理解できてない。後で調べよう
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading stdout: %s\n", err)
			return
		}
		// 処理
		processedLine := callback(line)
		fmt.Print("test ", processedLine)
	}
}

func (s *Server) sendCommand(stdin io.WriteCloser) {
	var command string
	for {
		fmt.Scan(&command)
		fmt.Println("Command: ", command)
		io.WriteString(stdin, command+"\n")
	}
}

func (s *Server) saveCommand(stdin io.WriteCloser) {
	for {
		time.Sleep(120 * time.Second)
		io.WriteString(stdin, "save-all\n")
	}
}

func (s *Server) QuitServer() {
	fmt.Println("Stopping server")
	close(s.stop)
	s.wg.Wait()
}

// 出力を加工する関数
func (s *Server) checkSSNOutput(output string) string {
	if strings.Contains(output, "コントロールURL") {
		s.disbot.Send(output)
	} else if strings.Contains(output, "公開開始") {
		s.disbot.Send(output)
	}
	return fmt.Sprintf("[%s] %s\n", time.Now(), output)
}
