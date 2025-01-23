package olresponse

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joomcode/errorx"
)

func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

type jsonErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
	Code    string `json:"code"`
}

func writeErrorJSON(err error, w http.ResponseWriter) {
	var resp jsonErrorResponse

	if errx, ok := err.(*errorx.Error); ok {
		resp.Message = errx.Message()
		resp.Cause = errx.Cause().Error()
		resp.Code = errx.Type().FullName()
	} else {
		resp.Code = "UNKNOWN"
		resp.Message = err.Error()
	}

	WriteJSON(w, resp)
}
