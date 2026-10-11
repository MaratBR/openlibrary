package value_test

import (
	"encoding/json"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/value"
)

func TestNullableJSON(t *testing.T) {
	for _, test := range []struct {
		name string
		data value.Nullable[int]
		json string
	}{
		{"absent", value.Null[int](), "null"},
		{"zero", value.Value(0), "0"},
		{"present", value.Value(42), "42"},
	} {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := json.Marshal(test.data)
			if err != nil || string(encoded) != test.json {
				t.Fatalf("Marshal = %s, %v; want %s", encoded, err, test.json)
			}
			decoded := value.Value(99)
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded != test.data {
				t.Fatalf("Unmarshal = %+v; want %+v", decoded, test.data)
			}
		})
	}
}

func TestNullableRejectsInvalidJSON(t *testing.T) {
	var decoded value.Nullable[int]
	if err := json.Unmarshal([]byte(`"invalid"`), &decoded); err == nil {
		t.Fatal("accepted string as nullable integer")
	}
	if decoded.Valid {
		t.Fatal("invalid input marked value as present")
	}
}
