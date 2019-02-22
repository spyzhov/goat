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

// Colors
const (
	//Header    = "\033[95m"
	//OkBlue    = "\033[94m"
	OkGreen = "\033[92m"
	//Warning   = "\033[93m"
	//Fail      = "\033[91m"
	EndColor  = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
)

const (
	defaultN = "[y/" + Bold + "N" + EndColor + "]"
	defaultY = "[" + Bold + "Y" + EndColor + "/n]"
	default0 = "[" + Bold + "0" + EndColor + "-%d]"
)

func New() *Console {
	return &Console{
		input: bufio.NewReader(os.Stdin),
	}
}

func (c *Console) Print(ask string, args ...interface{}) (int, error) {
	return fmt.Printf(ask+"\n", args...)
}

func (c *Console) Scanln(ask string, args ...interface{}) (string, error) {
	fmt.Printf(ask, args...)

	input, err := c.input.ReadString('\n')
	return strings.TrimSpace(input), err
}

func (c *Console) Prompt(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" "+defaultN+": ", args...)
	if err != nil {
		panic(err)
	}
	return len(str) > 0 && strings.ToLower(str)[0] == 'y'
}

func (c *Console) PromptY(ask string, args ...interface{}) bool {
	str, err := c.Scanln(ask+" "+defaultY+": ", args...)
	if err != nil {
		panic(err)
	}
	return !(len(str) > 0 && strings.ToLower(str)[0] == 'n')
}

func (c *Console) PromptInt(ask string, max int, args ...interface{}) int {
	for {
		str, err := c.Scanln(ask+" "+default0+": ", append(args, max)...)
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

func Wrap(msg string, color string) string {
	return color + msg + EndColor
}
