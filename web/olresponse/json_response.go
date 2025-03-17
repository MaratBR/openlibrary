package olresponse

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	data          any
	notifications []Notification
}

func NewJSONResponse(data any) *JSON {
	return &JSON{
		data: data,
	}
}

func (r *JSON) AddNotification(notif Notification) {
	r.notifications = append(r.notifications, notif)
}

func (r *JSON) Write(w http.ResponseWriter) {
	if len(r.notifications) > 0 {
		flashesJson, err := json.Marshal(r.notifications)
		if err == nil {
			w.Header().Add("x-flash", string(flashesJson))
		}
	}
	WriteJSON(w, r.data)
}

func WriteOKResponse(w http.ResponseWriter) {
	WriteJSONResponse(w, "ok")
}

func WriteJSONResponse(w http.ResponseWriter, data any) {
	w.WriteHeader(200)
	NewJSONResponse(data).Write(w)
}
