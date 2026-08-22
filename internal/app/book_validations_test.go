package app

import (
	"strings"
	"testing"
)

func TestValidateChapterName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "Chapter 1"},
		{name: "empty", value: "   ", wantErr: true},
		{name: "seventy unicode characters", value: strings.Repeat("界", 70)},
		{name: "too long", value: strings.Repeat("a", 201), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChapterName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateChapterName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
