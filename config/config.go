package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ProjectID string

	Database struct {
		DBType string
		DBName string
	}

	Server struct {
		Address int
	}
}

var C Config

func ReadConf() {
	Config := &C

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	// viper.AddConfigPath(filepath.Join("$GOPATH", "src", "clicker-pro", "config"))
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalln(err)
	}

	if projectId := os.Getenv("PROJECT_ID"); projectId != "" {
		C.ProjectID = projectId
	} else {
		os.Setenv("PROJECT_ID", C.ProjectID)
	}
}
