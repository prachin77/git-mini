package utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform.")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ShowClientIpAdd() {
	cmd := exec.Command("ipconfig")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error fetching client IP address:", err)
		return
	}

	re := regexp.MustCompile(`IPv4 Address.*?: (\d+\.\d+\.\d+\.\d+)`)
	matches := re.FindSubmatch(output)

	if len(matches) > 1 {
		fmt.Println("Your IP Address is:", string(matches[1]))
	} else {
		fmt.Println("IPv4 Address not found.")
	}
}
