package app_test

import (
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
)

func TestMoveArrEl(t *testing.T) {
	arr := []int{
		42,
		12,
		873,
		123,
		1000000012,
	}

	newPosiiton, ok := app.MoveItem(arr, 873, 0)
	if ok == false {
		t.Error("MoveItem failed")
	}

	if newPosiiton != 0 {
		t.Error("newPosiiton is not 0")
	}

	if arr[0] != 873 {
		t.Error("873 is not in position 0")
	}
	if arr[1] != 42 {
		t.Error("42 is not in position 1")
	}
	if arr[2] != 12 {
		t.Error("12 is not in position 2")
	}
	if arr[3] != 123 {
		t.Error("123 is not in position 3")
	}
	if arr[4] != 1000000012 {
		t.Error("1000000012 is not in position 4")
	}
}
