package public

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestModerationUserIDsFilter(t *testing.T) {
	first := uuid.Must(uuid.NewV4())
	second := uuid.Must(uuid.NewV4())
	ids, err := moderationUserIDsFilter(first.String() + ", " + second.String() + "," + first.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != first || ids[1] != second {
		t.Fatalf("expected ordered distinct IDs, got %#v", ids)
	}
	if ids, err = moderationUserIDsFilter("  "); err != nil || ids != nil {
		t.Fatalf("expected empty filter, got %#v, %v", ids, err)
	}
}

func TestModerationUserIDsFilterRejectsInvalidAndExcessiveValues(t *testing.T) {
	if _, err := moderationUserIDsFilter("not-a-uuid"); err == nil {
		t.Fatal("expected invalid UUID error")
	}
	value := ""
	for range 21 {
		if value != "" {
			value += ","
		}
		value += uuid.Must(uuid.NewV4()).String()
	}
	if _, err := moderationUserIDsFilter(value); err == nil {
		t.Fatal("expected excessive user count error")
	}
}
