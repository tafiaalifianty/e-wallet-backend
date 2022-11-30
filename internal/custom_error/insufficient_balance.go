package custom_error

type InsufficientBalance struct {
}

func (e InsufficientBalance) Error() string {
	return "Wallet's balance is insufficient"
}
