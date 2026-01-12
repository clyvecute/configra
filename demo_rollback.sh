#!/bin/bash
# demo_rollback.sh - Simulates a simplified workflow

echo "1. Starting clean..."
# (In a real script we would reset DB here if possible, but let's assume we proceed)

echo "2. Pushing Config V1 (Valid)"
echo '{ "app_name": "MyApp", "max_retries": 3, "mode": "release" }' > config_v1.json
# We reuse the existing schema.json
go run ./cmd/cli push -file config_v1.json -project 1

echo "3. Pushing Config V2 (Update)"
echo '{ "app_name": "MyApp", "max_retries": 5, "mode": "debug" }' > config_v2.json
go run ./cmd/cli push -file config_v2.json -project 1

echo "4. Attempting Invalid Push (Should Fail)"
echo '{ "app_name": "MyApp", "max_retries": 100, "mode": "debug" }' > config_bad.json
go run ./cmd/cli push -file config_bad.json -project 1
# This should print a validation error

echo "5. Rolling back to Version 1..."
# We assume the key is "feature_flags" as hardcoded in our CLI push
go run ./cmd/cli rollback -project 1 -key feature_flags -version 1

echo "Done! Check database to see Version 3 is a copy of Version 1."
