package custom_error

type InvalidToken struct {
}

func (e InvalidToken) Error() string {
	return "Request token is invalid"
}
