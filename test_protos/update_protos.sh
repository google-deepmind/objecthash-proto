#!/bin/bash -eu
#
# Copyright 2017 The ObjectHash-Proto Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

readonly TEST_PROTOS_DIR="$(cd "$(dirname "$0")"; pwd)"
readonly SCHEMA_DIR="${TEST_PROTOS_DIR}/schema"
readonly GENERATED_DIR="${TEST_PROTOS_DIR}/generated"
readonly LATEST_COMMIT="${TEST_PROTOS_DIR}/latest_commit.txt"
readonly LATEST_DIR="${GENERATED_DIR}/latest"

readonly TMP_GOPATH="$(mktemp -d)"
trap "rm -rf ${TMP_GOPATH}" EXIT

readonly GOPATH="${TMP_GOPATH}"
readonly GOBIN="${TMP_GOPATH}/bin"
readonly PATH="${TMP_GOPATH}/bin:${PATH}"

readonly TMP_PROTOC_PATH="$(mktemp -d)"
trap "rm -rf ${TMP_PROTOC_PATH}" EXIT

readonly PROTOC_VERSION="3.5.1"
readonly PROTOC_URL="https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
readonly PROTOC_BIN="${TMP_PROTOC_PATH}/protoc/bin/protoc"

generate_protos() {
  local schema_dir="$1"
  local output_dir="$2"

  # Make sure the directory exists and is empty.
  rm -rf "${output_dir}"
  mkdir -p "${output_dir}"

  for version in proto2 proto3; do
    mkdir -p "${output_dir}/${version}"
    "${PROTOC_BIN}" \
      --proto_path="${schema_dir}/${version}" \
      --go_out="${output_dir}/${version}" \
      "${schema_dir}/${version}"/*.proto
  done
}

install_protoc() {
  curl --location "${PROTOC_URL}" \
    --output "${TMP_PROTOC_PATH}/protoc.zip"
  unzip "${TMP_PROTOC_PATH}/protoc.zip" \
    -d "${TMP_PROTOC_PATH}/protoc"
}

clone_protoc_gen_go() {
  go get -u github.com/golang/protobuf/protoc-gen-go
}

build_protoc_gen_go() {
  rm -rf "${TMP_GOPATH}/bin"
  rm -rf "${TMP_GOPATH}/pkg"

  mkdir -p "${TMP_GOPATH}/bin"
  mkdir -p "${TMP_GOPATH}/pkg"
  (
    cd "${TMP_GOPATH}/src/github.com/golang/protobuf"
    go install ./protoc-gen-go
  )
}

get_protoc_gen_commit() {
  (
    cd "${TMP_GOPATH}/src/github.com/golang/protobuf"
    git rev-parse HEAD
  )
}

run() {
  rm -rf "${GENERATED_DIR}"
  mkdir -p "${GENERATED_DIR}"

  install_protoc
  clone_protoc_gen_go
  build_protoc_gen_go
  generate_protos "${SCHEMA_DIR}" "${LATEST_DIR}"

  local commit="$(get_protoc_gen_commit)"
  echo "${commit}" > "${LATEST_COMMIT}"
}

run
