package custom_error

import "fmt"

type FailedToUpdateData struct {
	DataType string
}

func (e FailedToUpdateData) Error() string {
	return fmt.Sprintf("Failed to update %s data", e.DataType)
}
