package main

import (
	"fmt"
	"net"
	"os"

	database "github.com/startdusk/tiny-redis/api/db"
	"github.com/startdusk/tiny-redis/cluster"
	"github.com/startdusk/tiny-redis/config"
	"github.com/startdusk/tiny-redis/db"
	"github.com/startdusk/tiny-redis/lib/logger"
	"github.com/startdusk/tiny-redis/resp/handler"
	"github.com/startdusk/tiny-redis/tcp"
)

const configFile string = "../../redis.conf"
const timeFormat = "2006-01-02"

var defaultProp = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// using telnet for test tcp server
// ```$ telnet localhost 6379```
// and send message
// and quit
// ctrl + [ and enter input `quit`
func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "tinyredis",
		Ext:        "log",
		TimeFormat: timeFormat,
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProp
	}

	var store database.Database
	if config.Properties.Self != "" && len(config.Properties.Peers) > 0 {
		store = cluster.NewDatabase()
	} else {
		store = db.NewStandaloneDatabase(
			config.Properties.Databases,
			config.Properties.AppendFilename,
			config.Properties.AppendOnly)
	}

	if err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: net.JoinHostPort(config.Properties.Bind, fmt.Sprintf("%d", config.Properties.Port)),
		}, handler.NewHandler(store)); err != nil {
		logger.Error(err)
	}
}
