#!/bin/bash
set -e

# Set variables
APP_NAME="tt"
INSTALL_PATH="$GOPATH/bin"

# Build the application
# -o flag specifies the output path for the binary
echo "Building the application..."
go build -o $INSTALL_PATH/$APP_NAME ./cmd/tt

# Output success message
echo "$APP_NAME has been installed to $INSTALL_PATH"
echo "You can now run the command 'tt' in your terminal."
