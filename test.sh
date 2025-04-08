#!/bin/bash

# Run 10 instances of "go run client/*.go" in parallel
for i in {1..100}; do
  go run client/*.go &
done

# Wait for all background processes to finish
wait

