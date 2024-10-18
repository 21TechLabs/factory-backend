package utils

type RaiseError struct {
	Message string
}

func (err RaiseError) Error() string {
	return err.Message
}
