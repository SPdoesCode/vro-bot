#!/bin/sh

set -e

echo "Building..."

go build -o vro-bot bot.go

echo "Done!"

echo "To run use ./vro-bot"

exit 0
