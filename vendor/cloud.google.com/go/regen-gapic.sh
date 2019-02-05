#!/bin/bash

# This script generates all GAPIC clients in this repo.
# One-time setup:
#   cd path/to/googleapis # https://github.com/googleapis/googleapis
#   virtualenv env
#   . env/bin/activate
#   pip install googleapis-artman
#   deactivate
#
# Regenerate:
#   cd path/to/googleapis
#   . env/bin/activate
#   $GOPATH/src/cloud.google.com/go/regen-gapic.sh
#   deactivate
#
# Being in googleapis directory is important;
# that's where we find YAML files and where artman puts the "artman-genfiles" directory.
#
# NOTE: This script does not generate the "raw" gRPC client found in google.golang.org/genproto.
# To do that, use the regen.sh script in the genproto repo instead.

set -ex

APIS=(
google/api/expr/artman_cel.yaml
google/iam/artman_iam_admin.yaml
google/cloud/asset/artman_cloudasset_v1beta1.yaml
google/iam/credentials/artman_iamcredentials_v1.yaml
google/cloud/bigquery/datatransfer/artman_bigquerydatatransfer.yaml
google/cloud/bigquery/storage/artman_bigquerystorage_v1beta1.yaml
google/cloud/dataproc/artman_dataproc_v1.yaml
google/cloud/dataproc/artman_dataproc_v1beta2.yaml
google/cloud/dialogflow/artman_dialogflow_v2.yaml
google/cloud/irm/artman_irm_v1alpha2.yaml
google/cloud/kms/artman_cloudkms.yaml
google/cloud/language/artman_language_v1.yaml
google/cloud/language/artman_language_v1beta2.yaml
google/cloud/oslogin/artman_oslogin_v1.yaml
google/cloud/oslogin/artman_oslogin_v1beta.yaml
google/cloud/redis/artman_redis_v1beta1.yaml
google/cloud/redis/artman_redis_v1.yaml
google/cloud/scheduler/artman_cloudscheduler_v1beta1.yaml
google/cloud/securitycenter/artman_securitycenter_v1beta1.yaml
google/cloud/speech/artman_speech_v1.yaml
google/cloud/speech/artman_speech_v1p1beta1.yaml
google/cloud/talent/artman_jobs_v4beta1.yaml
google/cloud/tasks/artman_cloudtasks_v2beta2.yaml
google/cloud/tasks/artman_cloudtasks_v2beta3.yaml
google/cloud/texttospeech/artman_texttospeech_v1.yaml
google/cloud/videointelligence/artman_videointelligence_v1.yaml
google/cloud/videointelligence/artman_videointelligence_v1beta1.yaml
google/cloud/videointelligence/artman_videointelligence_v1beta2.yaml
google/cloud/vision/artman_vision_v1.yaml
google/cloud/vision/artman_vision_v1p1beta1.yaml
google/devtools/artman_clouddebugger.yaml
google/devtools/clouderrorreporting/artman_errorreporting.yaml
google/devtools/cloudtrace/artman_cloudtrace_v1.yaml
google/devtools/cloudtrace/artman_cloudtrace_v2.yaml
google/devtools/containeranalysis/artman_containeranalysis_v1beta1.yaml
google/firestore/artman_firestore.yaml
google/logging/artman_logging.yaml
google/longrunning/artman_longrunning.yaml
google/monitoring/artman_monitoring.yaml
google/privacy/dlp/artman_dlp_v2.yaml
google/pubsub/artman_pubsub.yaml
google/spanner/admin/database/artman_spanner_admin_database.yaml
google/spanner/admin/instance/artman_spanner_admin_instance.yaml
google/spanner/artman_spanner.yaml
)

for api in "${APIS[@]}"; do
  rm -rf artman-genfiles/*
  artman --config "$api" generate go_gapic
  cp -r artman-genfiles/gapi-*/cloud.google.com/go/* $GOPATH/src/cloud.google.com/go/
done

pushd $GOPATH/src/cloud.google.com/go/
  gofmt -s -d -l -w . && goimports -w .

  # NOTE(pongad): `sed -i` doesn't work on Macs, because -i option needs an argument.
  # `-i ''` doesn't work on GNU, since the empty string is treated as a file name.
  # So we just create the backup and delete it after.
  ver=$(date +%Y%m%d)
  git ls-files -mo | while read modified; do
    dir=${modified%/*.*}
    find . -path "*/$dir/doc.go" -exec sed -i.backup -e "s/^const versionClient.*/const versionClient = \"$ver\"/" '{}' +
  done
popd


HASMANUAL=(
errorreporting/apiv1beta1
firestore/apiv1beta1
firestore/apiv1
logging/apiv2
longrunning/autogen
pubsub/apiv1
spanner/apiv1
trace/apiv1
)
for dir in "${HASMANUAL[@]}"; do
	find "$GOPATH/src/cloud.google.com/go/$dir" -name '*.go' -exec sed -i.backup -e 's/setGoogleClientInfo/SetGoogleClientInfo/g' '{}' '+'
done

find $GOPATH/src/cloud.google.com/go/ -name '*.backup' -delete
