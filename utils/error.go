package utils

import (
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func NormalizeGraphAPIError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *odataerrors.ODataError:
		if terr := e.GetError(); terr != nil {
			return &serializers.GraphAPIError{
				Code:    *terr.GetCode(),
				Message: *terr.GetMessage(),
			}
		}
	default:
		return &serializers.GraphAPIError{
			Code:    "",
			Message: err.Error(),
		}
	}

	return nil
}
