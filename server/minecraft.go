package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type Minecraft struct {
	Stop  chan int
	stdin io.WriteCloser
	wg    sync.WaitGroup
}

type ServerInterface interface {
	Quit() int
}

func NewMinecraft() *Minecraft {
	return &Minecraft{
		make(chan int),
		nil,
		sync.WaitGroup{},
	}
}

func (s *Minecraft) LaunchMinecraft(path string) {
	s.wg.Add(1)
	defer s.wg.Done()
	fmt.Println("Starting Minecraft")
	cmd := exec.Command("sh", path)

	// 入出力を取得
	s.stdin, _ = cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	// コマンドを開始
	cmd.Start()

	// 出力を表示
	go getOutput(stdout, s.checkSSNOutput)

	// コマンドを送信
	go s.saveCommand()

	//終了時処理
	cmd.Wait()
	s.stdin.Close()

	fmt.Println("Minecraft has stopped")
}

func getOutput(stdout io.Reader, callback func(string) string) {
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

func (s *Minecraft) SendCommand(command string) {
	io.WriteString(s.stdin, command+"\n")
}

func (s *Minecraft) saveCommand() {
	for {
		time.Sleep(120 * time.Second)
		io.WriteString(s.stdin, "save-all\n")
	}
}

func (s *Minecraft) Quit(status int) {
	fmt.Println("Stopping")
	if status == 0 {
		io.WriteString(s.stdin, "stop\n") // stopコマンドを送信
	}
	s.wg.Wait()
	s.Stop <- status
	fmt.Println("Stopped")

}

func (s *Minecraft) checkSSNOutput(output string) string {
	return fmt.Sprintf("[%s] %s\n", time.Now(), output)
}
