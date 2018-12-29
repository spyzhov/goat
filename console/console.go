package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Console struct {
	input *bufio.Reader
}

func New() *Console {
	return &Console{
		input: bufio.NewReader(os.Stdin),
	}
}

func (c *Console) Scanln(ask string, args ...interface{}) (string, error) {
	fmt.Printf(ask, args...)

	input, err := c.input.ReadString('\n')
	return strings.TrimSpace(input), err
}

func (c *Console) Prompt(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" [y/N]: ", args...)
	if err != nil {
		panic(err)
	}
	return len(str) > 0 && strings.ToLower(str)[0] == 'y'
}

func (c *Console) PromptY(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" [Y/n]: ", args...)
	if err != nil {
		panic(err)
	}
	return !(len(str) > 0 && strings.ToLower(str)[0] == 'n')
}

func (c *Console) PromptInt(ask string, max int, args ...interface{}) int {
	for {
		args = append(args, max)
		str, err := c.Scanln(ask+" [0-%d]: ", args...)
		if err != nil {
			panic(err)
		}
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		result, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Enter only numbers...\n")
			continue
		}
		if result < 0 {
			fmt.Printf("Enter only positive numbers...\n")
			continue
		}
		if result > max {
			fmt.Printf("Enter only numbers lower or equal than %d...\n", max)
			continue
		}
		return result
	}
}
