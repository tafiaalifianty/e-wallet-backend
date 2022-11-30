package custom_error

import "fmt"

type FailedToCreateData struct {
	DataType string
}

func (e FailedToCreateData) Error() string {
	return fmt.Sprintf("Failed to create %s data", e.DataType)
}
