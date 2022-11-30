package custom_error

type AuthHeaderNotAvailable struct {
}

func (e AuthHeaderNotAvailable) Error() string {
	return "Authorization Header is Not Available"
}
