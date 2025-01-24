package app

import (
	"errors"
	"net/url"
)

type UploadURIType int

const (
	UploadURITypePublicMinio UploadURIType = iota
	UploadURITypePrivateMinio
	UploadURITypeExternal
	UploadURITypeUnknown = 0xffff
)

func (t UploadURIType) String() string {
	switch t {
	case UploadURITypePublicMinio:
		return "public-minio"
	case UploadURITypePrivateMinio:
		return "private-minio"
	case UploadURITypeExternal:
		return "external"
	default:
		return "unknown"
	}
}

func parseUploadURIType(s string) UploadURIType {
	switch s {
	case "public-minio":
		return UploadURITypePublicMinio
	case "private-minio":
		return UploadURITypePrivateMinio
	case "external":
		return UploadURITypeExternal
	default:
		return UploadURITypeUnknown
	}
}

type UploadURI struct {
	Type UploadURIType
	Key  string
}

const (
	uploadURISchema = "ol-file"
)

var (
	ErrInvalidUploadURI = errors.New("invalid upload uri schema")
)

func (uu UploadURI) URL() *url.URL {
	u := new(url.URL)
	u.Scheme = uploadURISchema
	u.Host = uu.Type.String()

	switch uu.Type {
	case UploadURITypePublicMinio, UploadURITypePrivateMinio:
		u.Path = uu.Key
		break
	case UploadURITypeExternal:
		u.Path = "/"
		u.RawQuery = "key=" + url.QueryEscape(uu.Key)
		break
	case UploadURITypeUnknown:
		u.Path = "/"
		u.RawQuery = "key=" + url.QueryEscape(uu.Key)
		break
	}

	return u
}

func (uu UploadURI) String() string {
	return uu.URL().String()
}

func ParseUploadURI(s string) (*UploadURI, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != uploadURISchema {
		return nil, ErrInvalidUploadURI
	}

	uriType := parseUploadURIType(u.Host)

	var key string

	switch uriType {
	case UploadURITypePublicMinio, UploadURITypePrivateMinio:
		key = u.Path
		if key == "" {
			key = "/"
		}
		break
	case UploadURITypeExternal:
		key = u.Query().Get("key")
		break
	case UploadURITypeUnknown:
		key = u.Query().Get("key")
		break
	}

	return &UploadURI{
		Type: uriType,
		Key:  key,
	}, nil
}
