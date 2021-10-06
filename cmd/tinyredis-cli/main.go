package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterh/liner"

	redis "github.com/startdusk/tiny-redis"
)

var (
	historyFn = filepath.Join(os.TempDir(), ".liner_example_history")
	names     = []string{"john", "james", "mary", "nancy"}
)

func main() {
	defaultExpiration, _ := time.ParseDuration("0.5h")
	gcInterval, _ := time.ParseDuration("3s")
	cache := redis.NewCache(defaultExpiration, gcInterval)

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
			cmd := strings.Fields(command)
			switch cmd[0] {
			case "set":
				if len(cmd) < 4 {
					cache.Set(cmd[1], cmd[2], redis.NoExpiration)
				} else {
					expiration, err := time.ParseDuration(cmd[3])
					if err != nil {
						fmt.Println("Error time")
						break
					}
					cache.Set(cmd[1], cmd[2], expiration)
				}
				fmt.Println("OK")
			case "get":
				key := cmd[1]
				if v, ok := cache.Get(key); ok {
					fmt.Println(v)
				} else {
					fmt.Println("Not found key", key)
				}
			case "delete":
				key := cmd[1]
				if _, ok := cache.Get(key); ok {
					cache.Delete(key)
					fmt.Printf("Delete key: %v successfully\n", key)
				} else {
					fmt.Println("Not found key", key)
				}
			case "all":
				all := cache.All()
				for k, v := range all {
					fmt.Printf("The key %v is: %v \n", k, v.Value())
				}
			}
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
