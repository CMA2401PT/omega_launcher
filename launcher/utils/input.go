package utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func GetInput() string {
	buf := bufio.NewReader(os.Stdin)
	l, _, _ := buf.ReadLine()
	return string(strings.TrimSpace(string(l)))
}

func GetValidInput() string {
	for {
		s := GetInput()
		if s == "" {
			pterm.Error.Println("无效输入，输入不能为空")
			continue
		}
		return s
	}
}

func GetInputYN() bool {
	for {
		s := GetInput()
		if strings.HasPrefix(s, "Y") || strings.HasPrefix(s, "y") {
			return true
		} else if strings.HasPrefix(s, "N") || strings.HasPrefix(s, "n") {
			return false
		}
		pterm.Error.Println("无效输入，输入应该为 y 或者 n")
	}
}
