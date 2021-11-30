#!/usr/bin/env bash

go build -i -o ./tools ./src/tools/main.go

sudo chown $USER:$USER ./tools
sudo chmod 755 ./tools
