package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleInput() string {
	fmt.Println("Type prompt (empty Enter to finish):")

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	return scanner.Text()
}

func Confirm() (bool, error) {

	answers := make(map[string]bool)
	answers["y"] = true
	answers["n"] = false
	fmt.Print("Подтвердите y/n: ")
	reader := bufio.NewReader(os.Stdin)
	body, _ := reader.ReadString('\n')

	result := strings.ToLower(strings.TrimSpace(body))

	command, ok := answers[result]
	if !ok {
		return false, nil
	}
	return command, nil
}
