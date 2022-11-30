package custom_error

type FailedToGetInfoFromToken struct {
}

func (e FailedToGetInfoFromToken) Error() string {
	return "Failed to get info from token"
}
