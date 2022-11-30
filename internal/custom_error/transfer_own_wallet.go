package custom_error

type CannotTransferToOwnWallet struct {
}

func (e CannotTransferToOwnWallet) Error() string {
	return "Cannot Transfer to Own Wallet"
}
