#!/bin/bash

for i in {1..100}
do
  echo "Run #$i"
  go run client/*.go -drop=true
done

