#!/bin/bash

# Directory containing the JPG images
INPUT_DIR="./cats"

# Output directories
WEBP_DIR="$INPUT_DIR/webp"
AVIF_DIR="$INPUT_DIR/avif"

mkdir -p "$WEBP_DIR" "$AVIF_DIR"

# Loop over each .jpg file
for jpg_file in "$INPUT_DIR"/*.jpg; do
    filename=$(basename "$jpg_file" .jpg)

    # Convert to WebP
    cwebp -quiet "$jpg_file" -o "$WEBP_DIR/${filename}.webp" && \
        echo "Converted $jpg_file -> $WEBP_DIR/${filename}.webp"

    # Convert to AVIF
    avifenc "$jpg_file" "$AVIF_DIR/${filename}.avif" && \
        echo "Converted $jpg_file -> $AVIF_DIR/${filename}.avif"
done

echo "All JPGs converted to WebP and AVIF."
