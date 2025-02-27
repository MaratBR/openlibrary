# AO3 fanific ids

30383217
22073215
18111758
26026525
40718202
7991047

# libvips installation

1. Download libvips: https://github.com/libvips/libvips/releases
2. Install dependencies:
```bash
apt install \
    libjpeg8 libjpeg8-dev \
    libpng16-16 libpng-dev \
    libwebp7 libwebp-dev
```
3. Install libvips
```
./scripts/install-vips.sh
```

# Install go tools

```
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/mitranim/gow@latest

# Requires protoc installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latestoaded
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

```

# protobuf

Install protoc.


```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```