package olhttp

import (
	"encoding/json"
	"net/http"
)

var (
	errTypeJSONDecode = httpErrors.NewType("json_decode")
	ErrJSONDecodeEOF  = errTypeJSONDecode.New("EOF - http body appears to be empty")
)

func ReadJSONBody(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(v)
	if err != nil {
		if err.Error() == "EOF" {
			return ErrJSONDecodeEOF
		}

		return err
	}
	return nil
}
