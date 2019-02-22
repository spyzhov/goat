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
		str, err := c.Scanln(ask+" [0-%d]: ", append(args, max)...)
		if err != nil {
			panic(err)
		}
		str = strings.TrimSpace(str)
		if str == "" {
			return 0
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

func (c *Console) Select(ask string, args ...interface{}) int {
	for {
		str := ask + "\n"
		for i := len(args); i > 0; i-- {
			str += fmt.Sprintf("%2d) %s\n", i, args[i-1])
		}
		str += " 0) No one...\nPlease, select"
		return c.PromptInt(str, len(args))
	}
}
