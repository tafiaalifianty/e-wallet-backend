package custom_error

type InvalidRequestBody struct {
}

func (e InvalidRequestBody) Error() string {
	return "Request body is invalid"
}
