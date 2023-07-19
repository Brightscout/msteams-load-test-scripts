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

	if err := config.IsConnectionConfigurationValid(); err != nil {
		logger.Error("Error in validating the connection configuration", zap.Error(err))
		return
	}

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case constants.ClearStore:
			err = scripts.ClearStore()
		case constants.InitUsers:
			if err := config.IsUsersConfigurationValid(); err != nil {
				logger.Error("Error in validating the user configuration", zap.Error(err))
				break
			}

			err = scripts.InitUsers(config, logger)
		case constants.CreateChannels:
			if err := config.IsChannelsConfigurationValid(); err != nil {
				logger.Error("Error in validating the channel configuration", zap.Error(err))
				break
			}

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
