package custom_error

import "fmt"

type AmountNotInRange struct {
	Minimum int
	Maximum int
}

func (e AmountNotInRange) Error() string {
	return fmt.Sprintf(
		"Valid amount is between %d and %d",
		e.Minimum,
		e.Maximum,
	)
}
