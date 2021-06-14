package console

import (
	"bufio"
	"github.com/name5566/leaf/log"
	"os"
	"strings"
)
var reader = bufio.NewReader(os.Stdin)
var commands = map[string]func(args...interface{}) {}

func Register(command string,commandFunc func(args...interface{})) {
	commands[command] = commandFunc
}

func Start() {
	go run()
}

func run() {
	for {
		read_line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		read_line = strings.TrimSuffix(read_line[:len(read_line)-1], "\r")
		contents := strings.Fields(read_line)

		if len(contents) < 2 {
			continue
		}

		command := contents[0]
		args_strings := contents[1:]
		args := make([]interface{}, len(args_strings))
		for k, v := range args_strings {
			args[k] = v
		}

		if commandFunc,ok := commands[command];ok {
			commandFunc(args...)
		}else{
			log.Error("command <%v> not find",command)
		}
	}
}
