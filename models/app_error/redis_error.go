package app_error

type RedisError struct {
	Message string
}

func (e RedisError) Error() string {
	return e.Message
}
