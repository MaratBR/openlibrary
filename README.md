# Prerequisite

## libvips installation

1. Download libvips: https://github.com/libvips/libvips/releases
2. Install dependencies:

For debian-based:
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

### libvips - gotchas (at least for me, idk)

1. Make sure that `/usr/local/lib` is in `LD_LIBRARY_PATH`

# Install go tools

```
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/mitranim/gow@latest

# Requires protoc installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

```

