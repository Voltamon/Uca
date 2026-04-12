package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Ask(question string) (string, error) {
	fmt.Print(question + ": ")
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

func AskRequired(question string) (string, error) {
	for {
		answer, err := Ask(question)
		if err != nil {
			return "", err
		}
		if answer != "" {
			return answer, nil
		}
		fmt.Println("  This field is required")
	}
}

func AskDefault(question string, defaultVal string) (string, error) {
	answer, err := Ask(fmt.Sprintf("%s [%s]", question, defaultVal))
	if err != nil {
		return "", err
	}
	if answer == "" {
		return defaultVal, nil
	}
	return answer, nil
}
