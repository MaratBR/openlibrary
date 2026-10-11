package upload

import "go.uber.org/fx"

var FXModule = fx.Module("upload", fx.Provide(NewUploadServiceFromApplicationConfig))
