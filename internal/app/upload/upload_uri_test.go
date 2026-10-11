package upload_test

import (
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/upload"
)

func TestUploadURIParsing_Public(t *testing.T) {
	testUploadURIParsingCase(t, "ol-file://public-minio/book/123123423452345.jpeg", upload.UploadURITypePublicMinio, "/book/123123423452345.jpeg")

}

func TestUploadURIParsing_Private(t *testing.T) {
	testUploadURIParsingCase(t, "ol-file://private-minio/book/123123423452345.jpeg", upload.UploadURITypePrivateMinio, "/book/123123423452345.jpeg")
}

func TestUploadURIParsing_External(t *testing.T) {
	testUploadURIParsingCase(t, "ol-file://external/?key=%2Fbook%2F123123423452345.jpeg", upload.UploadURITypeExternal, "/book/123123423452345.jpeg")
}

func testUploadURIParsingCase(t *testing.T, uri string, expectedType upload.UploadURIType, expectedKey string) {
	uploadURL, err := upload.ParseUploadURI(uri)
	if err != nil {
		t.Fatal(err)
	}
	if uploadURL.String() != uri {
		t.Errorf("Expected string to be %s, got %s", uri, uploadURL.String())
	}
	if uploadURL.Type != expectedType {
		t.Errorf("Expected bucket to be %s, got %s", expectedType.String(), uploadURL.Type.String())
	}
	if uploadURL.Key != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, uploadURL.Key)
	}
}
