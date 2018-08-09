package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Console struct {
	Input *bufio.Reader
}

func New() *Console {
	return &Console{
		Input: bufio.NewReader(os.Stdin),
	}
}

func (c *Console) Scanln(ask string, args ...interface{}) (string, error) {
	fmt.Printf(ask, args...)

	input, err := c.Input.ReadString('\n')
	return strings.TrimSpace(input), err
}

func (c *Console) Prompt(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" [y/N]: ", args...)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(str) == "y"
}

func (c *Console) PromptY(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" [Y/n]: ", args...)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(str) != "n"
}

func (c *Console) PromptInt(ask string, args ...interface{}) int {
	for {
		str, err := c.Scanln(ask+" [0-9]: ", args...)
		if err != nil {
			panic(err)
		}
		if str == "" {
			return 0
		}
		result, err := strconv.Atoi(str)
		if err == nil {
			return result

		}
		fmt.Printf("Enter only numbers...\n")
	}
}

func (c *Console) PromptIntP(ask string, args ...interface{}) (result int) {
	for {
		result = c.PromptInt(ask, args...)
		if result >= 0 {
			return result
		}
		fmt.Printf("Enter only positive numbers...\n")
	}
}

func (c *Console) PromptIntN(ask string, args ...interface{}) (result int) {
	for {
		result = c.PromptInt(ask, args...)
		if result <= 0 {
			return result
		}
		fmt.Printf("Enter only negative numbers...\n")
	}
}
