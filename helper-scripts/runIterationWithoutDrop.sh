#!/bin/bash

for i in {1..100}
do
  echo "Run #$i"
  go run client/*.go -link="https://i.ytimg.com/vi/2DjGg77iz-A/sddefault.jpg"
done
