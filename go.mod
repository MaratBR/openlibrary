module github.com/MaratBR/openlibrary

go 1.23.0

toolchain go1.23.3

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/PuerkitoBio/goquery v1.10.1
	github.com/a-h/templ v0.3.833
	github.com/alexedwards/argon2id v1.0.0
	github.com/doug-martin/goqu/v9 v9.19.0
	github.com/eduardolat/goeasyi18n v1.3.0
	github.com/elastic/go-elasticsearch/v9 v9.0.1
	github.com/ggicci/httpin v0.19.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-faker/faker/v4 v4.5.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/h2non/bimg v1.1.9
	github.com/hjson/hjson-go/v4 v4.4.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/joomcode/errorx v1.2.0
	github.com/k3a/html2text v1.2.1
	github.com/knadh/koanf/parsers/toml/v2 v2.1.0
	github.com/knadh/koanf/providers/file v1.1.2
	github.com/knadh/koanf/v2 v2.1.2
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mileusna/useragent v1.3.5
	github.com/minio/minio-go/v7 v7.0.80
	github.com/pelletier/go-toml/v2 v2.2.3
	github.com/redis/go-redis/v9 v9.7.0
	github.com/urfave/cli/v3 v3.0.0-beta1
	golang.org/x/net v0.33.0
	golang.org/x/text v0.22.0
	golang.org/x/time v0.8.0
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.7.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/ggicci/owl v0.8.2 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/elastic/go-elasticsearch/v9 => github.com/MaratBR/go-elasticsearch/v9 v9.0.1
