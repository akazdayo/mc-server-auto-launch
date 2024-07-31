package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type server struct {
	stop       chan struct{}
	wg         sync.WaitGroup
	isRunning  chan bool
	controlURL chan string
	serverIP   chan string
}

func NewServer(isRunning chan bool, controlURL chan string, serverIP chan string) *server {
	return &server{
		make(chan struct{}),
		sync.WaitGroup{},
		isRunning,
		controlURL,
		serverIP,
	}
}

func (s *server) LaunchMinecraft(path string) {
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
	go s.getOutput(stdout, checkOutput)

	//終了時処理
	<-s.stop
	fmt.Println("Stopping Minecraft")
	io.WriteString(stdin, "stop\n") // stopコマンドを送信
	stdin.Close()
	cmd.Wait()

	fmt.Println("Minecraft has stopped")

}

func (s *server) LaunchSSNet(path string) {
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
	go s.getOutput(stdout, checkOutput)

	// 終了時処理
	<-s.stop
	stdout.Close()
	fmt.Println("Secure Share Net has stopped")
}

func (s *server) getOutput(stdout io.Reader, callback func(string) string) {
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

func (s *server) QuitServer() {
	fmt.Println("Stopping server")
	close(s.stop)
	s.wg.Wait()
}

// 出力を加工する関数
func checkOutput(output string) string {
	//if strings.Contains(output, "コントロールURL") {
	//} else if strings.Contains(output, "公開開始") {
	//}
	return fmt.Sprintf("[%s] %s\n", time.Now(), output)
}
