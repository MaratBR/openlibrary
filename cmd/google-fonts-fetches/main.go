// Command google-fonts-fetches downloads the SIL OFL and Apache 2.0 families
// from the official Google Fonts repository.
package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultSource = "https://github.com/google/fonts/archive/refs/heads/main.tar.gz"

type fontFace struct {
	Filename string
	Style    string
	Weight   int
}

type fontMetadata struct {
	Name         string `json:"name"`
	Source       string `json:"source"`
	License      string `json:"license"`
	IncludeCSS   string `json:"includeCss"`
	ExternalLink string `json:"externalLink"`
}

func main() {
	output := flag.String("output", "oldata/fonts/google", "destination directory")
	source := flag.String("source", defaultSource, "Google Fonts tar.gz URL")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if err := run(ctx, http.DefaultClient, *source, *output); err != nil {
		fmt.Fprintln(os.Stderr, "google-fonts-fetches:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, client *http.Client, source, output string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: %s", source, resp.Status)
	}

	parent := filepath.Dir(output)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll("temp", 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp("temp", "google-fonts-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)

	if err := extractAllowed(resp.Body, stage); err != nil {
		return err
	}
	if err := generateFamilies(stage); err != nil {
		return err
	}
	if err := os.RemoveAll(output); err != nil {
		return err
	}
	if err := os.Rename(stage, output); err != nil {
		return err
	}
	return nil
}

func extractAllowed(src io.Reader, stage string) error {
	gz, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read archive: %w", err)
		}
		parts := strings.Split(path.Clean(h.Name), "/")
		if len(parts) != 4 || (parts[1] != "ofl" && parts[1] != "apache") || h.Typeflag != tar.TypeReg {
			continue
		}
		name := parts[3]
		ext := strings.ToLower(filepath.Ext(name))
		if name != "METADATA.pb" && ext != ".ttf" && ext != ".otf" && ext != ".woff2" {
			continue
		}
		if path.Base(name) != name {
			continue
		}
		dir := filepath.Join(stage, parts[2])
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(f, io.LimitReader(tr, h.Size))
		closeErr := f.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		license := "SIL Open Font License 1.1"
		if parts[1] == "apache" {
			license = "Apache License 2.0"
		}
		if err := os.WriteFile(filepath.Join(dir, ".license"), []byte(license), 0o644); err != nil {
			return err
		}
	}
}

func generateFamilies(stage string) error {
	entries, err := os.ReadDir(stage)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		sourceDir := filepath.Join(stage, entry.Name())
		family, faces, err := parseMetadata(filepath.Join(sourceDir, "METADATA.pb"))
		if err != nil {
			slog.Error("failed to parse metadata", "sourceDir", sourceDir, "err", err)
			continue
		}
		if !safeComponent(family) {
			return fmt.Errorf("unsafe family name %q", family)
		}
		license, err := os.ReadFile(filepath.Join(sourceDir, ".license"))
		if err != nil {
			return err
		}
		familyDir := filepath.Join(stage, family)
		if familyDir != sourceDir {
			if err := os.Rename(sourceDir, familyDir); err != nil {
				return err
			}
		}
		if err := os.Remove(filepath.Join(familyDir, "METADATA.pb")); err != nil {
			return err
		}
		if err := os.Remove(filepath.Join(familyDir, ".license")); err != nil {
			return err
		}
		filesDir := filepath.Join(familyDir, "files")
		if err := os.Mkdir(filesDir, 0o755); err != nil {
			return err
		}
		sort.Slice(faces, func(i, j int) bool { return faces[i].Filename < faces[j].Filename })
		var css strings.Builder
		for _, face := range faces {
			if !safeComponent(face.Filename) {
				return fmt.Errorf("unsafe filename %q", face.Filename)
			}
			if err := os.Rename(filepath.Join(familyDir, face.Filename), filepath.Join(filesDir, face.Filename)); err != nil {
				return err
			}
			fmt.Fprintf(&css, "@font-face {\n  font-family: %s;\n  font-style: %s;\n  font-weight: %d;\n  font-display: swap;\n  src: url(%q);\n}\n\n",
				strconv.Quote(family), face.Style, face.Weight,
				"/_api/fonts/google?name="+url.QueryEscape(family)+"&file="+url.QueryEscape(face.Filename))
		}
		if err := os.WriteFile(filepath.Join(familyDir, "include.css"), []byte(css.String()), 0o644); err != nil {
			return err
		}
		metadata := fontMetadata{
			Name:         family,
			Source:       "google",
			License:      string(license),
			IncludeCSS:   "/_api/fonts/google/include?name=" + url.QueryEscape(family),
			ExternalLink: "https://fonts.google.com/specimen/" + url.PathEscape(family),
		}
		metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return err
		}
		metadataJSON = append(metadataJSON, '\n')
		if err := os.WriteFile(filepath.Join(familyDir, "metadata.json"), metadataJSON, 0o644); err != nil {
			return err
		}

		slog.Info("processed", "sourceDir", sourceDir)
	}
	return nil
}

func parseMetadata(filename string) (string, []fontFace, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	var family string
	var faces []fontFace
	var current *fontFace
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if family == "" && strings.HasPrefix(line, "name: ") {
			family, _ = strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(line, "name:")))
		}
		if line == "fonts {" {
			current = &fontFace{Style: "normal", Weight: 400}
			continue
		}
		if current == nil {
			continue
		}
		if line == "}" {
			if current.Filename != "" {
				faces = append(faces, *current)
			}
			current = nil
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		switch key {
		case "filename":
			current.Filename, _ = strconv.Unquote(value)
		case "style":
			current.Style, _ = strconv.Unquote(value)
		case "weight":
			current.Weight, _ = strconv.Atoi(value)
		}
	}
	if err := s.Err(); err != nil {
		return "", nil, err
	}
	if family == "" || len(faces) == 0 {
		return "", nil, errors.New("incomplete METADATA.pb")
	}
	return family, faces, nil
}

func safeComponent(value string) bool {
	return value != "" && value != "." && value != ".." && filepath.Base(value) == value && !strings.ContainsAny(value, "/\\\x00\r\n")
}
