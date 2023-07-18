package scripts

import (
	"os"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
)

func ClearStore() error {
	if err := os.Truncate(constants.TempStoreFile, 0); err != nil {
		return err
	}

	return nil
}
