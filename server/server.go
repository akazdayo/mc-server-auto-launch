package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type server struct {
	stop        chan struct{}
	mc_stopped  chan bool
	ssn_stopped chan bool
	isRunning   chan bool
	controlURL  chan string
	serverIP    chan string
}

func NewServer(isRunning chan bool, controlURL chan string, serverIP chan string) *server {
	return &server{
		make(chan struct{}),
		make(chan bool),
		make(chan bool),
		isRunning,
		controlURL,
		serverIP,
	}
}

func (s *server) LaunchMinecraft(path string) {
	fmt.Println("Starting Minecraft")
	cmd := exec.Command("sh", path)
	stdin, _ := cmd.StdinPipe()

	//終了時処理
	<-s.stop
	io.WriteString(stdin, "stop\n") // stopコマンドを送信
	stdin.Close()
	cmd.Wait()

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("結果: %s\n", out)
	fmt.Println("Minecraft has stopped")
	s.mc_stopped <- true
}

func (s *server) LaunchSSNet(path string) {
	fmt.Println("Starting Secure Share Net")
	cmd := exec.Command("sh", path)

	// 標準出力を取得
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating StdoutPipe: %s\n", err)
		return
	}

	// コマンドを開始
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		return
	}

	// 標準出力をリアルタイムで読み取る
	reader := bufio.NewReader(stdout)
	select {
	case <-s.stop:
		break

	default:
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading stdout: %s\n", err)
			return
		}
		// 出力を加工
		processedLine := checkOutput(line)
		fmt.Println(processedLine)
	}

	// 終了時処理
	fmt.Println("Stopping Secure Share Net")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println(out)

	fmt.Println("Secure Share Net has stopped")
	s.ssn_stopped <- true
}

func (s *server) QuitServer() {
	close(s.stop)
	<-s.mc_stopped
	<-s.ssn_stopped
}

// 出力を加工する関数
func checkOutput(output string) string {
	//if strings.Contains(output, "コントロールURL") {
	//} else if strings.Contains(output, "公開開始") {
	//}
	return fmt.Sprintf("[%s] %s\n", time.Now(), output)
}

/*
func main() {
	stop := make(chan bool)
	go LaunchMinecraft("./test.sh", stop)
	go LaunchSSNet("./ssnet.sh", stop)
	time.Sleep(10 * time.Second)
	stop <- true
	time.Sleep(2 * time.Second)
}
*/
