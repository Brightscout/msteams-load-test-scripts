package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
)

func LoadConfig() (*serializers.Config, error) {
	configFile, err := os.Open(constants.ConfigFile)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config *serializers.Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return config, nil
}
