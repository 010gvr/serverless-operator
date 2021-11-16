#!/usr/bin/env bash

set -Eeuo pipefail

root="$(dirname "${BASH_SOURCE[0]}")/../.."

# Source the main vars file to get the serving/eventing version to be used.
# shellcheck disable=SC1091,SC1090
source "$root/hack/lib/__sources__.bash"

kafka_channel_files=(channel-consolidated)
kafka_source_files=(source)

function download_kafka {
  component=$1
  subdir=$2
  version=$3
  shift
  shift
  shift

  files=("$@")

  component_dir="$root/knative-operator/deploy/resources/knativekafka"
  target_dir="${component_dir}"

  for (( i=0; i<${#files[@]}; i++ ));
  do
    index=$(( i+1 ))
    file="${files[$i]}.yaml"
    target_file="$target_dir/$subdir/$index-$file"
    url="https://github.com/knative-sandbox/$component/releases/download/$version/$file"

    wget --no-check-certificate "$url" -O "$target_file"

    # Break all image references so we know our overrides work correctly.
    yaml.break_image_references "$target_file"
  done
}

download_kafka eventing-kafka channel "$KNATIVE_EVENTING_KAFKA_VERSION" "${kafka_channel_files[@]}"
download_kafka eventing-kafka source "$KNATIVE_EVENTING_KAFKA_VERSION" "${kafka_source_files[@]}"

# For 1.17 we still skip HPA
git apply "$root/knative-operator/hack/001-eventing-kafka-remove_hpa.patch"

# SRVKE-919: Change the minavailable pdb for kafka-webhook to 0
git apply "$root/knative-operator/hack/007-eventing-kafka-patch-pdb.patch"

# NOTE: With upstream 0.27 (1.0) this patch is not needed:
# The kafka-ch-controller requires DELETE on deployment in OpenShift
git apply "$root/knative-operator/hack/002-eventing-kafka-ctor-role.patch"
