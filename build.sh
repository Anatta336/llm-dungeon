#!/bin/bash
set -e

for file in ./models/*; do
    ollama create "$(basename "$file")" -f "$file"
done

go build -o ./bin/dungeon ./main.go
