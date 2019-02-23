#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Set root directory
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

# Check GAE_KEY is set
if [ -z "${GAE_KEY}" ]; then
  echo "GAE_KEY must be set to a string containing a Google Cloud authentication key"
  exit 1
fi

# Init Google Cloud base path (script will exit if this path does not exist)
source /google-cloud-sdk/path.bash.inc

# Setup Google App Engine service account for authentication & login to Google Cloud
GAE_KEY=$(eval "echo $GAE_KEY")
echo "$GAE_KEY" > key.json
gcloud auth activate-service-account --key-file=key.json

# Check states (passed in via git commit messages)
# [SKIP-DEPLOY]
DEPLOY_STATE="true"
if [ `git log -1 --pretty=%B | grep -ic "\[\s*skip-deploy\s*\]"` -ge 1 ]; then
    DEPLOY_STATE="false"
fi

# [PROMOTE]
PROMOTE_STATE="false"
if [ `git log -1 --pretty=%B | grep -ic "\[\s*promote\s*\]"` -ge 1 ]; then
    PROMOTE_STATE="true"
fi

# Skip deploys when DEPLOY_STATE is false
if [ "$DEPLOY_STATE" = "false" ]; then
    echo "Skipping deploy"
    exit 0
fi

# Logic for production
if [ $1 = "master" ]; then
    echo "Initializing production deploy"

    # Current timestamp is passed as the second parameter
    VERSION=$2
    # Create tag for version
    TAG_MESSAGE=`git log -1  --pretty=%B | sed -n 's/Merge pull request .*\/\([^ ]*\).*/\1/p'`
    echo "Creating tag '$VERSION' for feature '$TAG_MESSAGE'"
    git -c user.name='Drone CI' -c user.email='<>' tag -a $VERSION -m "$TAG_MESSAGE"
    git -c user.name='Drone CI' -c user.email='<>' push origin $VERSION
fi

# Check if we should automatically promote new version
# Default behaviour is no-promote
# If "[PROMOTE]" is present in commit message, then promote
PROMOTE_FLAG="--no-promote"
if [ "$PROMOTE_STATE" = "true" ]; then
    PROMOTE_FLAG="--promote"
fi
echo "Promote state set to '$PROMOTE_FLAG'"

SERVICES="${DIR}/frontend/app.yaml ${DIR}/api/app.yaml"
gcloud app deploy $SERVICES --project $PROJECT_ID --version $VERSION $PROMOTE_FLAG
