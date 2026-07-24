package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	dir := t.TempDir()
	source := `package sample
import (
  "time"
  "github.com/MaratBR/openlibrary/internal/app"
)
// go2tsdef:generate
type Example struct {
  CreatedAt time.Time ` + "`json:\"createdAt\"`" + `
  Count app.Int64String ` + "`json:\"count,omitempty\"`" + `
  ID int64 ` + "`json:\"id,string\"`" + `
  Unix time.Time ` + "`go2tsdef:\"number\"`" + `
  Ignored string ` + "`json:\"-\"`" + `
  Labels []string
}
// go2tsdef:generate
type ID int64
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		Src:          []string{"./**/*.go"},
		Dest:         map[string]string{"default": "generated/types.ts"},
		InsertBefore: "import type { Brand } from './brand';\n\ntype ExternalID = Brand<string, 'ExternalID'>;",
		Types: map[string]string{
			"github.com/MaratBR/openlibrary/internal/app:Int64String": "string",
		},
	}
	if err := Generate(cfg, dir); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "generated/types.ts"))
	if err != nil {
		t.Fatal(err)
	}
	wantParts := []string{"import type { Brand } from './brand';", "type ExternalID = Brand<string, 'ExternalID'>;", "/** Generated from `model.go`. */\nexport interface Example", "createdAt: string;", "count?: string;", "id: string;", "Unix: number;", "Labels: Array<string>;", "/** Generated from `model.go`. */\nexport type ID = number;"}
	for _, want := range wantParts {
		if !strings.Contains(string(got), want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
	if strings.Index(string(got), "type ExternalID") > strings.Index(string(got), "export interface Example") {
		t.Errorf("insert_before content was not emitted before definitions:\n%s", got)
	}
	if strings.Contains(string(got), "Ignored") {
		t.Errorf("ignored JSON field was emitted:\n%s", got)
	}
}

func TestRunReadsInsertBeforeFromYAML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte("package sample\n// go2tsdef:generate\ntype Value string\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	config := "src:\n  - ./*.go\ndest:\n  default: out.ts\ninsert_before: |\n  const prefix = 'arbitrary TypeScript';\n"
	configPath := filepath.Join(dir, "go2tsdef.yaml")
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run(configPath); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "out.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "const prefix = 'arbitrary TypeScript';\n\n/** Generated from `model.go`. */\nexport type Value") {
		t.Fatalf("unexpected output:\n%s", got)
	}
}

func TestMarkedTypeWithinGroupedDeclaration(t *testing.T) {
	dir := t.TempDir()
	source := `package sample
type (
  // go2tsdef:generate
  Selected string
  NotSelected string
)
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Generate(Config{Src: []string{"."}, Dest: map[string]string{"default": "out.ts"}}, dir); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "out.ts"))
	if !strings.Contains(string(got), "Selected") || strings.Contains(string(got), "NotSelected") {
		t.Fatalf("unexpected output:\n%s", got)
	}
}

func TestGenerateRoutesTypesToDestinations(t *testing.T) {
	dir := t.TempDir()
	source := `package sample
// go2tsdef:generate
type DefaultType string

// go2tsdef:generate admin
type AdminType struct { ID int }
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		Src:  []string{"."},
		Dest: map[string]string{"default": "default.ts", "admin": "admin.ts"},
	}
	if err := Generate(cfg, dir); err != nil {
		t.Fatal(err)
	}
	defaultOutput, _ := os.ReadFile(filepath.Join(dir, "default.ts"))
	adminOutput, _ := os.ReadFile(filepath.Join(dir, "admin.ts"))
	if !strings.Contains(string(defaultOutput), "DefaultType") || strings.Contains(string(defaultOutput), "AdminType") {
		t.Fatalf("unexpected default destination:\n%s", defaultOutput)
	}
	if !strings.Contains(string(adminOutput), "AdminType") || strings.Contains(string(adminOutput), "DefaultType") {
		t.Fatalf("unexpected admin destination:\n%s", adminOutput)
	}
}

func TestGenerateOverrideType(t *testing.T) {
	dir := t.TempDir()
	source := `package sample
// go2tsdef:generate
// go2tsdef:override_type "info" | "warning" | "error"
type NotificationType uint8
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Generate(Config{Src: []string{"."}, Dest: map[string]string{"default": "out.ts"}}, dir); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "out.ts"))
	want := `export type NotificationType = "info" | "warning" | "error";`
	if !strings.Contains(string(got), want) {
		t.Fatalf("output missing %q:\n%s", want, got)
	}
}

func TestGenerateRejectsUnknownDestination(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte("package sample\n// go2tsdef:generate missing\ntype Value string\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Generate(Config{Src: []string{"."}, Dest: map[string]string{"default": "out.ts"}}, dir)
	if err == nil || !strings.Contains(err.Error(), `unknown destination "missing"`) {
		t.Fatalf("expected unknown destination error, got %v", err)
	}
}
