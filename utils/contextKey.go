package utils

import (
	"fmt"
	"net/http"
)

type ContextKey string

const UserContextKey ContextKey = "user"
const JWTContextKey ContextKey = "jwt"
const SchemaValidatorContextKey ContextKey = "parsedBody"

func ReadContextValue[T any](r *http.Request, key ContextKey) (T, error) {
	value, ok := r.Context().Value(key).(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("context value for key %s not found or type mismatch", key)
	}
	return value, nil
}
