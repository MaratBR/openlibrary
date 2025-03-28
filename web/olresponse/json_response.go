package olresponse

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joomcode/errorx"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

type APIBody interface {
	WriteHttpBody(w http.ResponseWriter)
}

type APIResponse struct {
	body          APIBody
	notifications []Notification
}

func (r *APIResponse) AddNotification(notif Notification) {
	r.notifications = append(r.notifications, notif)
}

func (r *APIResponse) Write(w http.ResponseWriter) {
	addNotifications(w, r.notifications)
	r.body.WriteHttpBody(w)
}

func NewAPIResponse(data any) *APIResponse {
	return createAPIResponse(&apiData{data: data})
}

func NewAPIResponseOK() *APIResponse {
	return NewAPIResponse("ok")
}

func NewAPIError(err error) *APIResponse {
	return createAPIResponse(&apiError{err: err})
}

func createAPIResponse(body APIBody) *APIResponse {
	return &APIResponse{
		body: body,
	}
}

type apiData struct {
	data any
}

func (r *apiData) WriteHttpBody(w http.ResponseWriter) {
	writeJSON(w, r.data)
}

type apiError struct {
	err error
}

type jsonErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
	Code    string `json:"code"`
}

func (r *apiError) WriteHttpBody(w http.ResponseWriter) {
	var resp jsonErrorResponse

	if errx, ok := r.err.(*errorx.Error); ok {
		resp.Message = errx.Message()
		resp.Cause = errx.Cause().Error()
		resp.Code = errx.Type().FullName()
	} else {
		resp.Code = "UNKNOWN"
		resp.Message = r.err.Error()
	}

	writeJSON(w, resp)
}

func addNotifications(w http.ResponseWriter, notifications []Notification) {
	if len(notifications) > 0 {
		flashesJson, err := json.Marshal(notifications)
		if err == nil {
			w.Header().Add("x-flash", string(flashesJson))
		}
	}
}
