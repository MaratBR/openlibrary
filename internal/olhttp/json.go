package olhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errTypeJSONDecode = httpErrors.NewType("json_decode")
	ErrJSONDecodeEOF  = errTypeJSONDecode.New("EOF - http body appears to be empty")
)

func ReadJSONBody(r *http.Request, v any) error {
	return readJSONBody(json.NewDecoder(r.Body), v)
}

// ReadJSONBodyMax reads one JSON value from a size-limited request body.
// Unknown fields and additional JSON values are rejected.
func ReadJSONBodyMax(w http.ResponseWriter, r *http.Request, v any, maxBytes int64) error {
	return readJSONBody(json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBytes)), v)
}

func readJSONBody(dec *json.Decoder, v any) error {
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return ErrJSONDecodeEOF
		}
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("request body must contain only one JSON value")
		}
		return err
	}
	return nil
}
