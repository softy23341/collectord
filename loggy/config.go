package loggy

import (
	"os"

	"github.com/BurntSushi/toml"
)

// DecodeTomlFile -- allow using environment variables in *.toml
// usage: “$VARIABLE_NAME“ or “${VARIABLE_NAME}“
// example: url = "amqp://$USER:$PASSWORD@rabbitmq/"
func DecodeTomlFile(configPath string, config interface{}) error {
	content, err := replaceEnvsFile(configPath)
	if err != nil {
		return err
	}
	_, err = toml.Decode(content, config)
	return err
}

func replaceEnvsFile(configPath string) (string, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return replaceEnvs(string(content)), nil
}

func replaceEnvs(content string) string {
	return os.ExpandEnv(content)
}
