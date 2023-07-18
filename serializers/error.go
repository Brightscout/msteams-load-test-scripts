package serializers

import "fmt"

type GraphAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *GraphAPIError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}
