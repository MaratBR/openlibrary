#!/bin/bash

# Total number of images to download
TOTAL_IMAGES=200

# Batch size for parallel downloads
BATCH_SIZE=10

# Output directory
OUTPUT_DIR="./cats"
mkdir -p "$OUTPUT_DIR"

# Function to download a single image
download_cat() {
    local i=$1
    local filename="$OUTPUT_DIR/cat${i}.jpg"
    curl -s -o "$filename" "https://cataas.com/cat"
    echo "Downloaded $filename"
}

# Main loop for downloading in parallel batches
for ((i=1; i<=TOTAL_IMAGES; i++)); do
    download_cat "$i" &

    # Wait after every BATCH_SIZE jobs
    if (( i % BATCH_SIZE == 0 )); then
        wait
    fi
done

# Wait for any remaining background jobs
wait

echo "Downloaded $TOTAL_IMAGES cat images to $OUTPUT_DIR"
