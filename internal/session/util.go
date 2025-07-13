package session

type redactedString string

func (c redactedString) String() string {
	return "[REDACTED]"
}

func (c redactedString) MarshalJSON() ([]byte, error) {
	return []byte("\"[REDACTED]\""), nil
}
