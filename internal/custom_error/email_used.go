package custom_error

type EmailAlreadyUsed struct {
}

func (e EmailAlreadyUsed) Error() string {
	return "Email is already used"
}
