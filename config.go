package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Configuration holds application configuration. Values will be pulled from
// environment variables, prefixed by keyPrefix. Default values can be added
// via tags.
type Configuration struct {
	Listen    string `config:"tcp://:8080"`
	Host      string `config:"localhost:8080"`
	AdminUser string `config:""`
	AdminPass string `config:""`
	Cert      string `config:""`
	Key       string `config:""`
	Scheme    string `config:"http"`
	Public    string `config:"public"`
	UseTus    string `config:"false"`
	TusHost   string `config:"localhost:1080"`

	S3Endpoint string `config:"http://127.0.0.1:9000"`
	S3Bucket   string `config:"lfs-data-store"`
	S3Region   string `config:"eu-west-1"`
	TraceURL   string `config:"http://localhost:14268"`
	LogLevel   string `config:"info"`
}

func (c *Configuration) IsHTTPS() bool {
	return strings.Contains(Config.Scheme, "https")
}

func (c *Configuration) IsPublic() bool {
	switch Config.Public {
	case "1", "true", "TRUE", "public":
		return true
	}
	return false
}

func (c *Configuration) IsUsingTus() bool {
	switch Config.UseTus {
	case "1", "true", "TRUE":
		return true
	}
	return false
}

// Config is the global app configuration
var Config = &Configuration{}

const keyPrefix = "LFS"

func init() {
	te := reflect.TypeOf(Config).Elem()
	ve := reflect.ValueOf(Config).Elem()

	for i := 0; i < te.NumField(); i++ {
		sf := te.Field(i)
		name := sf.Name
		field := ve.FieldByName(name)

		envVar := strings.ToUpper(fmt.Sprintf("%s_%s", keyPrefix, name))
		env := os.Getenv(envVar)
		tag := sf.Tag.Get("config")

		if env == "" && tag != "" {
			env = tag
		}

		field.SetString(env)
	}

	if port := os.Getenv("PORT"); port != "" {
		// If $PORT is set, override LFS_LISTEN. This is useful for deploying to Heroku.
		Config.Listen = "tcp://:" + port
	}
}
