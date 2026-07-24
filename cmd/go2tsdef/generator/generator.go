// Package generator converts selected Go type declarations to TypeScript.
package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	generateMarker = "go2tsdef:generate"
	overrideMarker = "go2tsdef:override_type"
	defaultDest    = "default"
)

// Config describes the source files, destination, and Go-to-TypeScript overrides.
type Config struct {
	Src          []string          `yaml:"src"`
	Dest         map[string]string `yaml:"dest"`
	Types        map[string]string `yaml:"types"`
	InsertBefore string            `yaml:"insert_before"`
}

type declaration struct {
	name     string
	typeExpr ast.Expr
	imports  map[string]string
	source   string
	dest     string
	override string
}

// Run reads a YAML configuration file and writes its generated TypeScript file.
// Relative paths in the configuration are relative to the configuration file.
func Run(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	base, err := filepath.Abs(filepath.Dir(configPath))
	if err != nil {
		return err
	}
	return Generate(cfg, base)
}

// Generate scans cfg.Src and writes each configured destination. Relative paths use baseDir.
func Generate(cfg Config, baseDir string) error {
	if len(cfg.Src) == 0 {
		return errors.New("config src must contain at least one path or glob")
	}
	if len(cfg.Dest) == 0 {
		return errors.New("config dest must contain at least one destination")
	}
	for name, path := range cfg.Dest {
		if strings.TrimSpace(name) == "" || strings.TrimSpace(path) == "" {
			return errors.New("config dest names and paths must not be empty")
		}
	}
	files, err := sourceFiles(cfg.Src, baseDir)
	if err != nil {
		return err
	}
	decls, err := parseDeclarations(files)
	if err != nil {
		return err
	}
	byDest := make(map[string][]declaration, len(cfg.Dest))
	seen := make(map[string]map[string]string, len(cfg.Dest))
	for _, d := range decls {
		if _, ok := cfg.Dest[d.dest]; !ok {
			return fmt.Errorf("type %s in %s references unknown destination %q", d.name, d.source, d.dest)
		}
		if seen[d.dest] == nil {
			seen[d.dest] = make(map[string]string)
		}
		if previous, ok := seen[d.dest][d.name]; ok {
			return fmt.Errorf("duplicate generated type %s in destination %q: %s and %s", d.name, d.dest, previous, d.source)
		}
		seen[d.dest][d.name] = d.source
		byDest[d.dest] = append(byDest[d.dest], d)
	}
	for name, path := range cfg.Dest {
		out, err := render(byDest[name], cfg.Types, cfg.InsertBefore, baseDir)
		if err != nil {
			return err
		}
		dest := resolve(baseDir, path)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("create destination %q directory: %w", name, err)
		}
		if err := os.WriteFile(dest, out, 0o644); err != nil {
			return fmt.Errorf("write destination %q: %w", name, err)
		}
	}
	return nil
}

func sourceFiles(patterns []string, base string) ([]string, error) {
	set := make(map[string]struct{})
	for _, raw := range patterns {
		pattern := resolve(base, raw)
		info, err := os.Stat(pattern)
		if err == nil && info.IsDir() {
			err = filepath.WalkDir(pattern, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if !d.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
					set[path] = struct{}{}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("scan %s: %w", raw, err)
			}
			continue
		}
		re, err := globRegexp(filepath.Clean(pattern))
		if err != nil {
			return nil, fmt.Errorf("invalid source pattern %q: %w", raw, err)
		}
		root := globRoot(pattern)
		if _, err := os.Stat(root); err != nil {
			return nil, fmt.Errorf("source pattern %q: %w", raw, err)
		}
		err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !d.IsDir() && re.MatchString(filepath.Clean(path)) && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				set[path] = struct{}{}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", raw, err)
		}
	}
	files := make([]string, 0, len(set))
	for file := range set {
		files = append(files, file)
	}
	sort.Strings(files)
	return files, nil
}

func resolve(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(base, path))
}

func globRoot(pattern string) string {
	i := strings.IndexAny(pattern, "*?[")
	if i < 0 {
		return filepath.Dir(pattern)
	}
	prefix := pattern[:i]
	root := filepath.Dir(prefix)
	if root == "." && filepath.IsAbs(pattern) {
		return string(filepath.Separator)
	}
	return root
}

func globRegexp(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(pattern); {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				i += 2
				if i < len(pattern) && os.IsPathSeparator(pattern[i]) {
					i++
					b.WriteString("(?:.*[/\\\\])?")
				} else {
					b.WriteString(".*")
				}
			} else {
				i++
				b.WriteString("[^/\\\\]*")
			}
		case '?':
			i++
			b.WriteString("[^/\\\\]")
		case '[':
			j := strings.IndexByte(pattern[i+1:], ']')
			if j < 0 {
				return nil, errors.New("unterminated character class")
			}
			j += i + 1
			b.WriteString(pattern[i : j+1])
			i = j + 1
		default:
			b.WriteString(regexp.QuoteMeta(string(pattern[i])))
			i++
		}
	}
	b.WriteString("$")
	return regexp.Compile(b.String())
}

func parseDeclarations(files []string) ([]declaration, error) {
	fset := token.NewFileSet()
	var result []declaration
	for _, path := range files {
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		imports := make(map[string]string)
		for _, imp := range file.Imports {
			p, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			name := filepath.Base(p)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			imports[name] = p
		}
		for _, d := range file.Decls {
			gen, ok := d.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts := spec.(*ast.TypeSpec)
				dest, marked, err := generateDestination(gen.Doc, ts.Doc)
				if err != nil {
					return nil, fmt.Errorf("parse directives for %s in %s: %w", ts.Name.Name, path, err)
				}
				if !marked {
					continue
				}
				override, err := typeOverride(gen.Doc, ts.Doc)
				if err != nil {
					return nil, fmt.Errorf("parse directives for %s in %s: %w", ts.Name.Name, path, err)
				}
				result = append(result, declaration{name: ts.Name.Name, typeExpr: ts.Type, imports: imports, source: path, dest: dest, override: override})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].name < result[j].name })
	return result, nil
}

func generateDestination(groups ...*ast.CommentGroup) (string, bool, error) {
	dest := defaultDest
	found := false
	for _, directive := range commentDirectives(generateMarker, groups...) {
		if found {
			return "", false, errors.New("multiple go2tsdef:generate directives")
		}
		found = true
		if directive != "" {
			dest = directive
		}
	}
	return dest, found, nil
}

func typeOverride(groups ...*ast.CommentGroup) (string, error) {
	directives := commentDirectives(overrideMarker, groups...)
	if len(directives) > 1 {
		return "", errors.New("multiple go2tsdef:override_type directives")
	}
	if len(directives) == 0 {
		return "", nil
	}
	if directives[0] == "" {
		return "", errors.New("go2tsdef:override_type requires a TypeScript definition")
	}
	return directives[0], nil
}

func commentDirectives(marker string, groups ...*ast.CommentGroup) []string {
	var result []string
	for _, group := range groups {
		if group == nil {
			continue
		}
		for _, comment := range group.List {
			text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
			if text == marker {
				result = append(result, "")
			} else if strings.HasPrefix(text, marker+" ") || strings.HasPrefix(text, marker+"\t") {
				result = append(result, strings.TrimSpace(text[len(marker):]))
			}
		}
	}
	return result
}

func render(decls []declaration, overrides map[string]string, insertBefore, baseDir string) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("// Code generated by go2tsdef. DO NOT EDIT.\n\n")
	if insertBefore != "" {
		b.WriteString(insertBefore)
		if !strings.HasSuffix(insertBefore, "\n") {
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}
	for _, d := range decls {
		source := d.source
		if relative, err := filepath.Rel(baseDir, source); err == nil {
			source = relative
		}
		fmt.Fprintf(&b, "/** Generated from `%s`. */\n", filepath.ToSlash(source))
		conv := converter{imports: d.imports, overrides: overrides}
		if d.override != "" {
			fmt.Fprintf(&b, "export type %s = %s;\n\n", d.name, d.override)
		} else if st, ok := d.typeExpr.(*ast.StructType); ok {
			fmt.Fprintf(&b, "export interface %s {\n", d.name)
			if err := conv.fields(&b, st.Fields.List); err != nil {
				return nil, fmt.Errorf("type %s: %w", d.name, err)
			}
			b.WriteString("}\n\n")
		} else {
			t, err := conv.tsType(d.typeExpr)
			if err != nil {
				return nil, fmt.Errorf("type %s: %w", d.name, err)
			}
			fmt.Fprintf(&b, "export type %s = %s;\n\n", d.name, t)
		}
	}
	return b.Bytes(), nil
}

type converter struct{ imports, overrides map[string]string }

func (c converter) fields(b *bytes.Buffer, fields []*ast.Field) error {
	for _, f := range fields {
		if len(f.Names) == 0 {
			continue
		}
		// Unexported fields are never part of Go's JSON representation.
		if !ast.IsExported(f.Names[0].Name) {
			continue
		}
		name, optional, skip := fieldName(f)
		if skip {
			continue
		}
		t := reflectTag(f, "go2tsdef")
		if t == "-" {
			continue
		}
		if t == "" {
			var err error
			t, err = c.tsType(f.Type)
			if err != nil {
				return err
			}
		}
		if optional {
			name += "?"
		}
		fmt.Fprintf(b, "  %s: %s;\n", quoteProperty(name), t)
	}
	return nil
}

func fieldName(f *ast.Field) (string, bool, bool) {
	name := f.Names[0].Name
	jsonTag := reflectTag(f, "json")
	if jsonTag == "-" {
		return "", false, true
	}
	parts := strings.Split(jsonTag, ",")
	if len(parts) > 0 && parts[0] != "" {
		name = parts[0]
	}
	optional := false
	for _, p := range parts[1:] {
		if p == "omitempty" || p == "omitzero" {
			optional = true
		}
	}
	return name, optional, false
}

func reflectTag(f *ast.Field, key string) string {
	if f.Tag == nil {
		return ""
	}
	raw, err := strconv.Unquote(f.Tag.Value)
	if err != nil {
		return ""
	}
	return reflect.StructTag(raw).Get(key)
}

func (c converter) tsType(expr ast.Expr) (string, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		if v, ok := c.overrides[x.Name]; ok {
			return v, nil
		}
		switch x.Name {
		case "string":
			return "string", nil
		case "bool":
			return "boolean", nil
		case "byte", "rune", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128":
			return "number", nil
		case "any", "interface{}":
			return "unknown", nil
		default:
			return x.Name, nil
		}
	case *ast.SelectorExpr:
		alias, ok := x.X.(*ast.Ident)
		if !ok {
			return "unknown", nil
		}
		path := c.imports[alias.Name]
		keys := []string{path + ":" + x.Sel.Name, path + "." + x.Sel.Name, alias.Name + "." + x.Sel.Name}
		for _, key := range keys {
			if v, ok := c.overrides[key]; ok {
				return v, nil
			}
		}
		if path == "time" && x.Sel.Name == "Time" {
			return "string", nil
		}
		return x.Sel.Name, nil
	case *ast.StarExpr:
		return c.tsType(x.X)
	case *ast.ArrayType:
		t, err := c.tsType(x.Elt)
		if err != nil {
			return "", err
		}
		return "Array<" + t + ">", nil
	case *ast.MapType:
		k, err := c.tsType(x.Key)
		if err != nil {
			return "", err
		}
		v, err := c.tsType(x.Value)
		if err != nil {
			return "", err
		}
		if k != "string" && k != "number" {
			k = "string"
		}
		return "Record<" + k + ", " + v + ">", nil
	case *ast.InterfaceType:
		return "unknown", nil
	case *ast.StructType:
		var b bytes.Buffer
		b.WriteString("{\n")
		if err := c.fields(&b, x.Fields.List); err != nil {
			return "", err
		}
		b.WriteString("}")
		return b.String(), nil
	case *ast.IndexExpr:
		base, err := c.tsType(x.X)
		if err != nil {
			return "", err
		}
		arg, err := c.tsType(x.Index)
		if err != nil {
			return "", err
		}
		return base + "<" + arg + ">", nil
	case *ast.IndexListExpr:
		base, err := c.tsType(x.X)
		if err != nil {
			return "", err
		}
		args := make([]string, len(x.Indices))
		for i, arg := range x.Indices {
			args[i], err = c.tsType(arg)
			if err != nil {
				return "", err
			}
		}
		return base + "<" + strings.Join(args, ", ") + ">", nil
	default:
		return "unknown", nil
	}
}

var tsIdentifier = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*\??$`)

func quoteProperty(name string) string {
	if tsIdentifier.MatchString(name) {
		return name
	}
	optional := strings.HasSuffix(name, "?")
	name = strings.TrimSuffix(name, "?")
	q := strconv.Quote(name)
	if optional {
		q += "?"
	}
	return q
}
