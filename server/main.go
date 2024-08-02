package server

import (
	"fmt"
	"strings"
)

func SendCommand(input []string, m *Minecraft, s *SSNet) {
	switch input[0] {
	case "mc":
		if len(input) < 2 {
			fmt.Println("Please enter a command")
			return
		}
		switch input[1] {
		case "stop":
			m.Quit(0)
		case "restart":
			m.Quit(1)
		case "cmd":
			if len(input) >= 3 {
				m.SendCommand(strings.Join(input[2:], " "))
			}
		}
	case "ssn":
	}
}
