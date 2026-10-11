package app

// Compatibility with existing callers; new feature code belongs in upload.
import (
	"github.com/MaratBR/openlibrary/internal/app/upload"
	"github.com/knadh/koanf/v2"
)

type UploadConfig = upload.UploadConfig

func NewUploadServiceFromApplicationConfig(cfg *koanf.Koanf) *UploadService {
	return upload.NewUploadServiceFromApplicationConfig(cfg)
}

type UploadService = upload.UploadService

func NewUploadService(cfg UploadConfig) *UploadService { return upload.NewUploadService(cfg) }

type UploadURIType = upload.UploadURIType

const UploadURITypePublicMinio = upload.UploadURITypePublicMinio
const UploadURITypePrivateMinio = upload.UploadURITypePrivateMinio
const UploadURITypeExternal = upload.UploadURITypeExternal
const UploadURITypeUnknown = upload.UploadURITypeUnknown

type UploadURI = upload.UploadURI

var ErrInvalidUploadURI = upload.ErrInvalidUploadURI

func ParseUploadURI(s string) (*UploadURI, error) { return upload.ParseUploadURI(s) }
