package app

import (
	"context"
	"errors"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/siteconfig"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type setupTestDB struct {
	DB
	admin bool
	err   error
	reads int
}
type setupTestRow struct{ db *setupTestDB }

func (r setupTestRow) Scan(dest ...any) error {
	if r.db.err != nil {
		return r.db.err
	}
	*dest[0].(*bool) = r.db.admin
	return nil
}
func (d *setupTestDB) QueryRow(context.Context, string, ...any) pgx.Row {
	d.reads++
	return setupTestRow{d}
}

type setupTestConfig struct {
	siteconfig.Store
	data  siteconfig.ConfigData
	err   error
	saves int
}

func (c *setupTestConfig) Get() *siteconfig.ConfigData      { return &c.data }
func (c *setupTestConfig) Save(context.Context, bool) error { c.saves++; return c.err }

func TestSetupStateRules(t *testing.T) {
	for _, admin := range []bool{false, true} {
		for _, flag := range []bool{false, true} {
			db := &setupTestDB{admin: admin}
			cfg := &setupTestConfig{}
			cfg.data.SiteSetup.IsSetup = flag
			s := NewSetupService(db, cfg)
			state, err := s.GetSetupState()
			if err != nil || state.AdminUserExist != admin || state.IsSetup != (admin && flag) {
				t.Fatalf("admin=%v flag=%v: %+v %v", admin, flag, state, err)
			}
		}
	}
}

func TestSetupCachingAndFreshState(t *testing.T) {
	db := &setupTestDB{}
	cfg := &setupTestConfig{}
	s := NewSetupService(db, cfg)
	if done, err := s.IsSetup(); done || err != nil {
		t.Fatal(done, err)
	}
	db.admin = true
	cfg.data.SiteSetup.IsSetup = true
	if done, err := s.IsSetup(); done || err != nil || db.reads != 1 {
		t.Fatal(done, err, db.reads)
	}
	if state, err := s.GetSetupState(); !state.IsSetup || err != nil || db.reads != 2 {
		t.Fatal(state, err, db.reads)
	}
}

func TestSetupCheckRetriesErrors(t *testing.T) {
	db := &setupTestDB{admin: true, err: errors.New("offline")}
	cfg := &setupTestConfig{}
	cfg.data.SiteSetup.IsSetup = true
	s := NewSetupService(db, cfg)
	if _, err := s.IsSetup(); err == nil {
		t.Fatal("expected error")
	}
	db.err = nil
	if done, err := s.IsSetup(); !done || err != nil || db.reads != 2 {
		t.Fatal(done, err, db.reads)
	}
}

func TestSetupSkipsExistingAdminAndInvalidatesCache(t *testing.T) {
	db := &setupTestDB{admin: true}
	cfg := &setupTestConfig{}
	s := NewSetupService(db, cfg)
	if _, err := s.IsSetup(); err != nil {
		t.Fatal(err)
	}
	cfg.err = errors.New("save failed")
	if err := s.Setup(context.Background(), SetupCommand{SiteName: "Library"}); err == nil {
		t.Fatal("expected save error")
	}
	if cfg.data.SiteSetup.IsSetup {
		t.Fatal("failed save changed setup flag")
	}
	cfg.err = nil
	if err := s.Setup(context.Background(), SetupCommand{SiteName: " Library "}); err != nil {
		t.Fatal(err)
	}
	if done, err := s.IsSetup(); !done || err != nil {
		t.Fatal(done, err)
	}
	if cfg.data.SiteSetup.SiteName != "Library" {
		t.Fatal(cfg.data.SiteSetup)
	}
	if err := s.Setup(context.Background(), SetupCommand{SiteName: "Again"}); err == nil {
		t.Fatal("setup ran twice")
	}
}

type setupTestTx struct {
	pgx.Tx
	db        *setupTestDB
	writes    int
	committed bool
}

func (tx *setupTestTx) QueryRow(context.Context, string, ...any) pgx.Row {
	return setupTestRow{&setupTestDB{}}
}
func (tx *setupTestTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	tx.writes++
	return pgconn.CommandTag{}, nil
}
func (tx *setupTestTx) Commit(context.Context) error {
	tx.committed = true
	tx.db.admin = true
	return nil
}
func (tx *setupTestTx) Rollback(context.Context) error { return nil }

type setupCreationDB struct {
	*setupTestDB
	tx *setupTestTx
}

func (d *setupCreationDB) Begin(context.Context) (pgx.Tx, error) { return d.tx, nil }

func TestSetupCreatesMissingAdmin(t *testing.T) {
	db := &setupCreationDB{setupTestDB: &setupTestDB{}}
	db.tx = &setupTestTx{db: db.setupTestDB}
	cfg := &setupTestConfig{}
	s := NewSetupService(db, cfg)
	if err := s.Setup(context.Background(), SetupCommand{SiteName: "Library", AdminUsername: "owner", AdminPassword: "secret"}); err != nil {
		t.Fatal(err)
	}
	if !db.tx.committed || db.tx.writes != 2 || !cfg.data.SiteSetup.IsSetup {
		t.Fatalf("admin creation failed: tx=%+v config=%+v", db.tx, cfg.data)
	}
	if done, err := s.IsSetup(); !done || err != nil {
		t.Fatal(done, err)
	}
}

func TestSetupRejectsInvalidInput(t *testing.T) {
	for _, command := range []SetupCommand{
		{SiteName: " ", AdminUsername: "owner", AdminPassword: "password123"},
		{SiteName: "Library", AdminUsername: "!", AdminPassword: "password123"},
		{SiteName: "Library", AdminUsername: "owner"},
		{SiteName: "Library", AdminUsername: "owner", AdminPassword: "short"},
	} {
		cfg := &setupTestConfig{}
		cfg.data.PasswordRequirements.MinLength = 8
		s := NewSetupService(&setupTestDB{}, cfg)
		if err := s.Setup(context.Background(), command); !errors.Is(err, ErrSetupInvalidInput) {
			t.Fatalf("expected validation error for %+v: %v", command, err)
		}
		if cfg.saves != 0 {
			t.Fatal("invalid input saved")
		}
	}
}
