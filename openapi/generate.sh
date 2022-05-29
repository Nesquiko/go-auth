#!/bin/bash

echo "Generating OpenApi specification code with Oapi-codegen"

if ! command -v oapi-codegen &> /dev/null
then
	echo "Oapi-codegen not installed."
	echo "Run 'go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest'"
	exit
fi

echo -n "Generating chi-server code..."
oapi-codegen --config ./openapi/server.cfg.yaml ./openapi/openapi.yaml

if [ $? -eq 0 ]; then
	echo -e " \x1b[32;1mSUCCESS\x1b[0m"
else
	echo -e " \x1b[31;1mFAILED\x1b[0m"
	exit
fi

echo -n "Generating types..."
oapi-codegen --config ./openapi/types.cfg.yaml ./openapi/openapi.yaml

if [ $? -eq 0 ]; then
	echo -e " \x1b[32;1mSUCCESS\x1b[0m"
else
	echo -e " \x1b[31;1mFAILED\x1b[0m"
	exit
fi
