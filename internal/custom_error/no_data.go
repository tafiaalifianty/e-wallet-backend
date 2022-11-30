package custom_error

import "fmt"

type NoDataFound struct {
	DataType string
}

func (e NoDataFound) Error() string {
	return fmt.Sprintf("No %s data found", e.DataType)
}
