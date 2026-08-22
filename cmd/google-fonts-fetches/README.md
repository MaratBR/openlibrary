# Google Fonts fetcher

Downloads the official Google Fonts repository and replaces
`oldata/fonts/google` with families from its `ofl` (SIL Open Font License) and
`apache` (Apache License 2.0) directories. Other license directories are not
extracted.

```sh
go run ./cmd/google-fonts-fetches
```

Use `-output` to select another destination or `-source` to use a compatible
Google Fonts repository tarball (useful for mirrors and tests). Each family is
written as `<family>/include.css`, `<family>/metadata.json`, and
`<family>/files/*`. Download staging directories are created under `./temp`.
