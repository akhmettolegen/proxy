package configs

import (
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type Config struct {
	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":8080"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:"/proxy"`
	FilesDir   string `long:"files-directory" env:"FILES_DIR" description:"Directory where all static files are located" required:"false" default:"/usr/share/proxy"`

	DSName string `short:"n" long:"ds" env:"DATASTORE" description:"DataStore name (format: dgraph/null)" required:"false" default:"mongo"`
	DSDB   string `short:"d" long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name (format: inventory)" required:"false" default:"proxy"`
	DSURL  string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: mongodb://localhost:27017)" required:"false" default:"mongodb://localhost:27017"`

	InDebugMode bool `long:"in-debug-mode" env:"DEBUG" description:"debug mode"`
	IsTesting   bool `long:"testing" env:"TESTING" description:"testing mode"`
}

func ConfigWithParsedFlags() *Config {
	config := new(Config)
	p := flags.NewParser(config, flags.Default)

	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] Error while parsing config options:", err)
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return config
}
