package olhttp

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSONBodyIsStrict(t *testing.T) {
	type input struct {
		Name string `json:"name"`
	}

	for _, test := range []struct {
		name string
		body string
	}{
		{name: "unknown field", body: `{"name":"one","extra":true}`},
		{name: "multiple values", body: `{"name":"one"} {"name":"two"}`},
		{name: "trailing invalid data", body: `{"name":"one"} nope`},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(test.body))
			if err := ReadJSONBody(r, &input{}); err == nil {
				t.Fatal("expected decoding error")
			}
		})
	}
}

func TestReadJSONBodyRejectsEmptyAndOversizedBodies(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(""))
	if err := ReadJSONBody(r, &struct{}{}); !errors.Is(err, ErrJSONDecodeEOF) {
		t.Fatalf("expected empty body error, got %v", err)
	}

	r = httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"long"}`))
	w := httptest.NewRecorder()
	if err := ReadJSONBodyMax(w, r, &struct {
		Name string `json:"name"`
	}{}, 5); err == nil {
		t.Fatal("expected body size error")
	}
}
