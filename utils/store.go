package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
)

func StoreCreds(response *serializers.Store) error {
	responseBytes, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(constants.TempStoreFile, responseBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func LoadCreds() (*serializers.Store, error) {
	tempStoreFile, err := os.Open(constants.TempStoreFile)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return &serializers.Store{}, nil
		}

		return nil, err
	}

	defer tempStoreFile.Close()
	byteValue, err := ioutil.ReadAll(tempStoreFile)
	if err != nil {
		return nil, err
	}

	if len(byteValue) == 0 {
		return &serializers.Store{}, nil
	}

	var store *serializers.Store
	if err := json.Unmarshal(byteValue, &store); err != nil {
		return nil, err
	}

	return store, nil
}
