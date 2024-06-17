#!/bin/bash

# Ensure correct number of arguments
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <directory_path> <json_file> <json_property>"
    exit 1
fi

# Arguments
DIRECTORY_PATH=$1
JSON_FILE=$2
JSON_PROPERTY=$3

# Get the filename from the specified directory
FILENAME=$(ls "$DIRECTORY_PATH" | head -n 1)

# Check if the file exists
if [ -z "$FILENAME" ]; then
    echo "No files found in the directory: $DIRECTORY_PATH"
    exit 1
fi

# Create an array from the JSON property path
IFS='.' read -r -a path_array <<<"$JSON_PROPERTY"

# Construct the jq argument for setting the property
jq_arg=""
for elem in "${path_array[@]}"; do
    jq_arg="$jq_arg.\"$elem\""
done

# Extract the current path from the JSON file
current_path=$(jq -r "$jq_arg" "$JSON_FILE")
if [ -z "$current_path" ]; then
    echo "Property $JSON_PROPERTY does not exist in $JSON_FILE"
    exit 1
fi

# Extract the directory path (excluding the filename) from the current path
dir_path=$(dirname "$current_path")

# Construct the updated path with the new filename
updated_path="$dir_path/$FILENAME"

# Replace the filename part in the JSON property
jq --arg updated_path "$updated_path" "$jq_arg = \$updated_path" "$JSON_FILE" >tmp.json && mv tmp.json "$JSON_FILE"

# Check if jq command was successful
if [ $? -eq 0 ]; then
    echo "Successfully updated $JSON_PROPERTY in $JSON_FILE with value $updated_path"
else
    echo "Failed to update $JSON_FILE"
    exit 1
fi
