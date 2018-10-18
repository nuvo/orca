#!/bin/sh

set -e

echo "Getting stable environment"
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml

echo "Deploying dynamic environment"
orca deploy env --name $DST_NS -c charts.yaml --kube-context $DST_KUBE_CONTEXT --override $CHART_NAME=$CHART_VERSION -x
