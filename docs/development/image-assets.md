# Raster image asset workflow

Use this workflow whenever a raster asset needs a downsized PNG fallback and
modern AVIF and WebP variants. It applies to PNG and JPEG source images in any
frontend asset directory, including
[`web/frontend/embed-assets/img`](../../web/frontend/embed-assets/img).

## Source assets

Keep source files in an `original_assets/` directory alongside the generated
assets. The browser-facing variants remain in the owning asset directory. For
example:

```text
web/frontend/embed-assets/img/
├── original_assets/books.png
├── original_assets/books_dark.png
├── books.avif
├── books.webp
├── books.png
├── books_dark.avif
├── books_dark.webp
└── books_dark.png
```

Do not overwrite the source image. Preserve alpha channels for PNG inputs;
transparent artwork requires them.

## Generate the variants

The required local tools are `vips` (resizing and WebP) and `avifenc` (from
`libavif-tools`, alpha-safe AVIF encoding). Run from the repository root,
setting the variables for the asset being generated:

```sh
asset_dir='web/frontend/embed-assets/img'
source="$asset_dir/original_assets/example.png" # PNG or JPEG
name='example'
scale='0.6906077348' # target width / source width; e.g. 1000 / 1448

vips resize "$source" "$asset_dir/$name.png" "$scale" --kernel lanczos3
vips resize "$source" "$asset_dir/$name.webp[Q=90,alpha_q=100,effort=6]" "$scale" --kernel lanczos3
avifenc --qcolor 85 --qalpha 100 --speed 6 "$asset_dir/$name.png" "$asset_dir/$name.avif"
```

`alpha_q=100` and `--qalpha 100` keep alpha lossless. Use `avifenc` for PNGs
with transparency: the FFmpeg `libaom-av1` encoder available in this project
accepts only non-alpha pixel formats, which would flatten transparent artwork.
FFmpeg remains suitable for opaque image conversion.

Calculate `scale` as `target_width / source_width`; the target height follows
the original aspect ratio. For JPEG source files, use the `.jpg` or `.jpeg`
source extension but still generate the three variants above. Confirm the
result before committing:

```sh
magick identify "$asset_dir/$name.png" "$asset_dir/$name.webp"
avifdec --info "$asset_dir/$name.avif"
```

## Reference the files

For markup images, use a `<picture>` element with AVIF, WebP, then PNG/JPEG.
For CSS background images, which cannot use `<picture>`, use `image-set()` and
list candidates in this order:

1. AVIF (`type("image/avif")`)
2. WebP (`type("image/webp")`)
3. PNG (`type("image/png")`)

The browser selects the first format it supports, while the generated PNG
remains the universal fallback. If an asset has theme variants, apply the same
ordering to each variant. See
[`home.scss`](../../web/frontend/src/common/style/home.scss) for a CSS
background example.

## Verify

Run the normal frontend verification after stylesheet changes:

```sh
pnpm run build
git diff --check
```

If TypeScript has an unrelated baseline error but the stylesheet needs
validation, `pnpm exec vite build --mode production` checks Vite and Sass
without running `tsc`. Build output under `dist/` is ignored and must not be
committed.
