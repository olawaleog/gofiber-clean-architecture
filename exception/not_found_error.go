package exception

type NotFoundError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.Message
}

func (b BadRequestError) Error() string {
	return b.Message
}
