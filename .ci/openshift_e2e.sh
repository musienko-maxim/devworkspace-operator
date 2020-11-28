#!/bin/bash
#
# Copyright (c) 2012-2020 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#

#!/usr/bin/env bash
# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u
# print each command before executing it
set -x

#trap 'Catch_Finish $?' EXIT SIGINT

# Catch the finish of the job and write logs in artifacts.
function Catch_Finish() {
    # grab devworkspace-controller namespace events after running e2e
    getDevWorkspaceOperatorLogs
}

# ENV used by PROW ci
export CI="openshift"
export ARTIFACTS_DIR="/tmp/artifacts"
export NAMESPACE="devworkspace-controller"
export TERMINAL_USER_SECRET_NAME="terminal-user"
export TERMINAL_USER_LOGIN="developer"
export TERMINAL_USER_PASSWORD="developer"
export PATH_TO_HTPASSWD_FILE='/tmp/users.htpasswd'
export PATH_TO_OAUTH_CR_YAML_FILE='/tmp/htpasswdProvider.yaml'
export KUBERNETES_API_ENDPOINT=$(oc whoami --show-server)

# Component is defined in Openshift CI job configuration. See: https://github.com/openshift/release/blob/master/ci-operator/config/devfile/devworkspace-operator/devfile-devworkspace-operator-master__v4.yaml#L8
export CI_COMPONENT="devworkspace-operator"

# DEVWORKSPACE_OPERATOR env var exposed by Openshift CI in e2e test pod. More info about how images are builded in Openshift CI: https://github.com/openshift/ci-tools/blob/master/TEMPLATES.md#parameters-available-to-templates
# Dependencies environment are defined here: https://github.com/openshift/release/blob/master/ci-operator/config/devfile/devworkspace-operator/devfile-devworkspace-operator-master__v5.yaml#L36-L38

#export IMG=${DEVWORKSPACE_OPERATOR}

# Pod created by openshift ci don't have user. Using this envs should avoid errors with git user.
export GIT_COMMITTER_NAME="CI BOT"
export GIT_COMMITTER_EMAIL="ci_bot@notused.com"

# Function to get all logs and events from devworkspace operator deployments
function getDevWorkspaceOperatorLogs() {
    mkdir -p ${ARTIFACTS_DIR}/devworkspace-operator
    cd ${ARTIFACTS_DIR}/devworkspace-operator
    for POD in $(oc get pods -o name -n ${NAMESPACE}); do
       for CONTAINER in $(oc get -n ${NAMESPACE} ${POD} -o jsonpath="{.spec.containers[*].name}"); do
            echo ""
            echo "<=========================Getting logs from $POD==================================>"
            echo ""
            oc logs ${POD} -c ${CONTAINER} -n ${NAMESPACE} | tee $(echo ${POD}-${CONTAINER}.log | sed 's|pod/||g')
        done
    done
    echo "======== oc get events ========"
    oc get events -n ${NAMESPACE}| tee get_events.log
}

function generateHtpasswd() {
   htpasswd -c -B -b ${PATH_TO_HTPASSWD_FILE} terminal-test terminal
}
function generateHtpasswdProviderYaml() {
    echo  "apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
spec:
  identityProviders:
  - name: htpasswd
    mappingMethod: claim
    type: HTPasswd
    htpasswd:
      fileData:
        name: ${TERMINAL_USER_SECRET_NAME}" > ${PATH_TO_OAUTH_CR_YAML_FILE}
}

function addUserToCluster (){
  oc create secret generic ${TERMINAL_USER_SECRET_NAME} \
    --from-file=htpasswd=${PATH_TO_HTPASSWD_FILE} \
    --namespace openshift-config \
    --dry-run=client \
    --output yaml | oc apply -f -
  #need timeout for applying changes into cluster properly
  sleep 4
  oc apply -f ${PATH_TO_OAUTH_CR_YAML_FILE}
  oc adm policy add-cluster-role-to-user admin ${TERMINAL_USER_LOGIN}
}

function checkLogin (){
  CURRENT_TIME=$(date +%s)
  ENDTIME=$(($CURRENT_TIME + 300))
  while [ $(date +%s) -lt $ENDTIME ]; do
      if KUBECONFIG='/tmp/checkloginconfig' oc login\
      -u ${TERMINAL_USER_LOGIN}\
      -p ${TERMINAL_USER_PASSWORD} ${KUBERNETES_API_ENDPOINT}\
      "--insecure-skip-tls-verify=true"; then
          break
      fi
      sleep 10
  done
}


# Check if operator-sdk is installed and if not install operator sdk in $GOPATH/bin dir
if ! hash operator-sdk 2>/dev/null; then
    mkdir -p $GOPATH/bin
    export PATH="$PATH:$(pwd):$GOPATH/bin"
    OPERATOR_SDK_VERSION=v0.17.0

    curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu

    chmod +x operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu && \
        cp operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu $GOPATH/bin/operator-sdk && \
        rm operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
fi

# For some reason go on PROW force usage vendor folder
# This workaround is here until we don't figure out cause
generateHtpasswd
#generateHtpasswdProviderYaml
#addUserToCluster
checkLogin
go mod tidy
go mod vendor
make install
make test_e2e
make uninstall
