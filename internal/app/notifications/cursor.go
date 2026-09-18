package notifications

import (
	"encoding/json"
	"strconv"
	"time"
)

type Cursor int64

func (v Cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *Cursor) UnmarshalJSON(data []byte) error {
	var decoded string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var (
		err error
		c   Cursor
	)
	c, err = ParseCursor(decoded)
	if err != nil {
		return err
	}
	*v = c
	return nil
}

func (c Cursor) String() string {
	return strconv.FormatInt(int64(c), 16)
}

func ParseCursor(s string) (Cursor, error) {
	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return Cursor(i), nil
}

func getCursorFromTime(t time.Time) Cursor {
	return Cursor(t.UnixMilli())
}

func getTimeFromCursor(c Cursor) time.Time {
	return time.UnixMilli(int64(c))
}
