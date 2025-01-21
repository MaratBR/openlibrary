#!/bin/bash

# Directory containing the input files
input_dir="web/frontend/embed-assets/cover"

# Ensure ffmpeg is installed
if ! command -v ffmpeg &> /dev/null
then
    echo "ffmpeg could not be found. Please install it to use this script."
    exit 1
fi

# Process each file matching the glob pattern
for input_file in "$input_dir"/*.full.jpg; do
  # Ensure the file exists
  if [[ ! -f "$input_file" ]]; then
    echo "No files matching the pattern found."
    continue
  fi

  # Extract the base name without extension
  base_name="${input_file%.full.jpg}"

  # Create the converted files
  echo "Processing $input_file..."

  # h200.jpg
  ffmpeg -i "$input_file" -vf scale=-1:200 "${base_name}.h200.jpg" -y

  # h200.webp
  ffmpeg -i "$input_file" -vf scale=-1:200 -qscale:v 80 "${base_name}.h200.webp" -y

  # h300.jpg
  ffmpeg -i "$input_file" -vf scale=-1:300 "${base_name}.h300.jpg" -y

  # h300.webp
  ffmpeg -i "$input_file" -vf scale=-1:300 -qscale:v 80 "${base_name}.h300.webp" -y

done

echo "Conversion process completed."