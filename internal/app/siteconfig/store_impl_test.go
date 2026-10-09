package siteconfig

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type configRow struct {
	data []byte
	err  error
}

func (r configRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*dest[0].(*[]byte) = r.data
	return nil
}

type configDB struct {
	dal.DB
	raw      []byte
	readErr  error
	writeErr error
	writes   int
	saved    []byte
}

func (d *configDB) QueryRow(context.Context, string, ...any) pgx.Row {
	return configRow{d.raw, d.readErr}
}
func (d *configDB) Exec(_ context.Context, _ string, args ...any) (pgconn.CommandTag, error) {
	d.writes++
	if d.writeErr == nil {
		d.saved = append([]byte(nil), args[0].([]byte)...)
	}
	return pgconn.CommandTag{}, d.writeErr
}

func TestIndependentMigrations(t *testing.T) {
	type value struct{ Number int }
	var active value
	definition := Definition[value]{Name: "test", Field: func(*ConfigData) *value { return &active }, Migrations: []func(*value) error{
		func(v *value) error { v.Number += 2; return nil },
		func(v *value) error { v.Number *= 3; return nil },
	}}
	changed, err := definition.Load(&ConfigData{}, json.RawMessage(`{"Version":1,"Data":{"Number":4}}`))
	if err != nil || !changed || active.Number != 12 {
		t.Fatalf("migration result: %+v, %v, %v", active, changed, err)
	}
	for _, raw := range []string{`{"Version":3,"Data":{}}`, `{"Version":-1,"Data":{}}`, `null`, `{"Version":0}`, `{"Data":null}`, `{"Data":42}`} {
		if _, err := definition.Load(&ConfigData{}, json.RawMessage(raw)); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
}

func TestLoadAndSave(t *testing.T) {
	db := &configDB{raw: []byte(`{"FontConfiguration":{"Version":0,"Data":{"MaxPerChapter":4}},"CaptchaSettings":{"Version":0,"Data":{"Type":"custom"}},"Unknown":{"Version":99,"Data":{}}}`)}
	s := NewStore(db, nil)
	if err := s.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s.Get().FontConfiguration.MaxPerChapter != 4 || s.Get().CaptchaSettings.Type != "custom" || db.writes != 1 {
		t.Fatalf("unexpected snapshot or writes: %+v, %d", s.Get(), db.writes)
	}
	var rows map[string]settingsRow
	if err := json.Unmarshal(db.saved, &rows); err != nil {
		t.Fatal(err)
	}
	if rows["FontConfiguration"].Version != 1 || rows["CaptchaSettings"].Version != 0 {
		t.Fatal("versions are not independent")
	}
	if _, exists := rows["Unknown"]; exists {
		t.Fatal("unknown settings included in upsert")
	}
	if err := s.Save(context.Background(), false); err != nil {
		t.Fatal(err)
	}
	if db.writes != 1 {
		t.Fatal("unchanged configuration saved")
	}
	s.Get().PasswordRequirements.MinLength = 12
	db.writeErr = errors.New("write failed")
	if err := s.Save(context.Background(), false); err == nil {
		t.Fatal("write failure ignored")
	}
	db.writeErr = nil
	if err := s.Save(context.Background(), false); err != nil {
		t.Fatal(err)
	}
	if db.writes != 3 {
		t.Fatal("failed write was not retried")
	}
	if err := s.Save(context.Background(), true); err != nil {
		t.Fatal(err)
	}
	if db.writes != 4 {
		t.Fatal("forced save skipped")
	}
}

func TestLoadFailurePreservesActive(t *testing.T) {
	for _, tc := range []struct {
		raw      string
		writeErr error
	}{
		{raw: `{"FontConfiguration":{"Version":2,"Data":{}}}`},
		{raw: `{"FontConfiguration":{"Version":0,"Data":{}}}`, writeErr: errors.New("write failed")},
		{raw: `{"PasswordRequirements":{"Version":0,"Data":null}}`},
	} {
		db := &configDB{raw: []byte(tc.raw), writeErr: tc.writeErr}
		s := NewStore(db, nil)
		s.Get().PasswordRequirements.MinLength = 17
		if err := s.Load(context.Background()); err == nil {
			t.Fatal("expected load failure")
		}
		if s.Get().PasswordRequirements.MinLength != 17 {
			t.Fatal("failed load replaced active snapshot")
		}
	}
}

func TestEmptyDatabaseUsesDefaults(t *testing.T) {
	db := &configDB{readErr: pgx.ErrNoRows}
	s := NewStore(db, nil)
	if err := s.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s.Get().FontConfiguration.MaxPerChapter != 10 || db.writes != 0 {
		t.Fatal("empty database did not use defaults")
	}
}
