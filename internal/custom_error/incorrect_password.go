package custom_error

type WrongPassword struct {
}

func (e WrongPassword) Error() string {
	return "Password is incorrect"
}
