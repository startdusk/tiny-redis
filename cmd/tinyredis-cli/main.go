package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

var (
	historyFn = filepath.Join(os.TempDir(), ".liner_example_history")
	names     = []string{"john", "james", "mary", "nancy"}
)

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range names {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	for {
		command, err := line.Prompt("127.0.0.1:5200>")
		command = strings.TrimSpace(command)
		if err != nil {
			if errors.Is(err, liner.ErrPromptAborted) {
				log.Print("Aborted")
			} else {
				log.Print("Error reading line: ", err)
			}
		}

		if strings.ToLower(command) == "exit" {
			break
		} else if strings.ToLower(command) == "help" {
			printHelper()
		} else {
			log.Print("Got command: ", command)
			line.AppendHistory(command)
		}
	}

	if f, err := os.Create(historyFn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func printHelper() {
	helpText := `
Thanks for using TinyRedis
And the command is 
TinyRedis-cli
To get help about command:
	Type: "help <command>" for help on command
To quit:
	<ctrl+c> or <exit>
	`

	fmt.Println(helpText)
}
