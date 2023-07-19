package main

import (
	"errors"
	"os"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/scripts"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		logger.Error("failed to load the config", zap.Error(err))
		return
	}

	// TODO: Add config validation
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case constants.ClearStore:
			err = scripts.ClearStore()
		case constants.InitUsers:
			err = scripts.InitUsers(config, logger)
		case constants.CreateChannels:
			err = scripts.CreateChannels(config, logger)
		case constants.CreateChats:
			err = scripts.CreateChats(config, logger)
		default:
			err = errors.New("invalid arguments")
		}
	}
	if err != nil {
		logger.Error("failed to run the script", zap.String("arg", args[1]), zap.Error(err))
	}

	_ = logger.Sync()
}
