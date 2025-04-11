#!/bin/bash

for i in {1..10}; do
    echo "Starting run #$i in background"
    go run ./client/*.go &
done

# Wait for all background jobs to finish
wait
echo "All client runs completed."


