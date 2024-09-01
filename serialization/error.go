package serialization

type DeserializationError struct {
	Err          error
	ErrorMessage string
	ErrorStatus  int
}

// Error is a helper to get error message
func (e *DeserializationError) Error() string {
	if e.Err == nil {
		return "error in deserialization"
	}
	return e.Err.Error()
}

// Unwrap is a helper to get error
func (e *DeserializationError) Unwrap() error {
	return e.Err
}
