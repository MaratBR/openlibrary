package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractAndGenerateAllowedFamilies(t *testing.T) {
	archive := makeArchive(t, map[string]string{
		"fonts-main/ofl/tangerine/METADATA.pb":           metadata("Tangerine", "Tangerine-Regular.ttf"),
		"fonts-main/ofl/tangerine/Tangerine-Regular.ttf": "ofl-font",
		"fonts-main/apache/roboto/METADATA.pb":           metadata("Roboto", "Roboto-Regular.ttf"),
		"fonts-main/apache/roboto/Roboto-Regular.ttf":    "apache-font",
		"fonts-main/ufl/ubuntu/METADATA.pb":              metadata("Ubuntu", "Ubuntu-Regular.ttf"),
		"fonts-main/ufl/ubuntu/Ubuntu-Regular.ttf":       "disallowed-font",
		"fonts-main/ofl/tangerine/OFL.txt":               "license text",
	})
	stage := t.TempDir()
	require.NoError(t, extractAllowed(bytes.NewReader(archive), stage))
	require.NoError(t, generateFamilies(stage))

	font, err := os.ReadFile(filepath.Join(stage, "Tangerine", "files", "Tangerine-Regular.ttf"))
	require.NoError(t, err)
	require.Equal(t, "ofl-font", string(font))
	css, err := os.ReadFile(filepath.Join(stage, "Tangerine", "include.css"))
	require.NoError(t, err)
	require.Contains(t, string(css), `url("/_api/fonts/google?name=Tangerine&file=Tangerine-Regular.ttf")`)
	metadataJSON, err := os.ReadFile(filepath.Join(stage, "Tangerine", "metadata.json"))
	require.NoError(t, err)
	require.JSONEq(t, `{
		"name":"Tangerine",
		"source":"google",
		"license":"SIL Open Font License 1.1",
		"includeCss":"/_api/fonts/google/include?name=Tangerine",
		"externalLink":"https://fonts.google.com/specimen/Tangerine"
	}`, string(metadataJSON))
	apacheMetadataJSON, err := os.ReadFile(filepath.Join(stage, "Roboto", "metadata.json"))
	require.NoError(t, err)
	require.Contains(t, string(apacheMetadataJSON), `"license": "Apache License 2.0"`)
	_, err = os.Stat(filepath.Join(stage, "ubuntu"))
	require.ErrorIs(t, err, os.ErrNotExist)
	_, err = os.Stat(filepath.Join(stage, "Tangerine", "OFL.txt"))
	require.ErrorIs(t, err, os.ErrNotExist)
}

func metadata(family, filename string) string {
	return "name: " + strconvQuote(family) + "\nfonts {\n  name: " + strconvQuote(family) + "\n  style: \"normal\"\n  weight: 400\n  filename: " + strconvQuote(filename) + "\n}\n"
}

func strconvQuote(s string) string { return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"` }

func makeArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	tw := tar.NewWriter(gz)
	for name, contents := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(contents)), Typeflag: tar.TypeReg}))
		_, err := io.WriteString(tw, contents)
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	return out.Bytes()
}
