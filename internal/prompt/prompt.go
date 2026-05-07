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

func AskChoice(question string, choices []string) (string, error) {
	fmt.Println(question)
	for i, c := range choices {
		fmt.Printf("  %d) %s\n", i+1, c)
	}

	for {
		answer, err := Ask("Enter number")
		if err != nil {
			return "", err
		}

		var idx int
		_, err = fmt.Sscanf(answer, "%d", &idx)
		if err == nil && idx > 0 && idx <= len(choices) {
			return choices[idx-1], nil
		}
		fmt.Println("  Invalid choice, try again")
	}
}
