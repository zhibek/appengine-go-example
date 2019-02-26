#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Set root directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

# Run dev server
echo "Running dev server"
dev_appserver.py ${DIR}/frontend/app.yaml ${DIR}/api/app.yaml --skip_sdk_update_check --log_level=info --dev_appserver_log_level=info --support_datastore_emulator --require_indexes --datastore_consistency_policy consistent --datastore_path ${DATASTORE_PATH} --host 0.0.0.0 --admin_host 0.0.0.0 --port 8080 --admin_port 8000
