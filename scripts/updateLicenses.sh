#!/bin/bash

set -e
set -x

python scripts/updateLicense.py $(go list -json ./... | jq -r '.Dir + "/" + (.GoFiles | .[])')
