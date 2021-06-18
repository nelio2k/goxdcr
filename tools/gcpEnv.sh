#!/usr/bin/env bash
set -u

. ./importExporter.shlib
if (($? != 0)); then
	exit $?
fi

function insertPropertyIntoBucketReplPropertyMap {
	local sourceCluster=$1
	local sourceBucket=$2
	local targetCluster=$3
	local targetBucket=$4
	local -n incomingMap=$5

	BUCKET_REPL_PROPERTIES_MAP=()
	for key in ${!incomingMap[@]}; do
		BUCKET_REPL_PROPERTIES_MAP["${sourceCluster},${sourceBucket},${targetCluster},${targetBucket},${key}"]=${incomingMap[$key]}
	done
	export BUCKET_REPL_PROPERTIES_MAP
}

declare -A BUCKET_REPL_PROPERTIES_MAP
declare -A CLUSTER_NAME_PORT_MAP
declare -A CLUSTER_NAME_XDCR_PORT_MAP
declare -A CLUSTER_NAME_BUCKET_MAP
declare -a cluster1BucketsArr
declare -a cluster2BucketsArr

declare -A BUCKET_NAME_RAMQUOTA_MAP
declare -A BUCKET_NAME_SCOPE_MAP
declare -A SCOPE_NAME_COLLECTION_MAP
declare -A BUCKET_PROPERTIES_OUTPUT_MAP
declare -A BUCKET_REPL_PROPERTIES_OUTPUT_MAP
declare -A BUCKET_REPL_EXPORT_MAP
declare -A CLUSTER_NAME_BUCKET_DONE_MAP
declare -A BUCKET_NAME_SCOPE_DONE_MAP
declare -A SCOPE_NAME_COLLECTION_DONE_MAP
DEFAULT_ADMIN="Administrator"
DEFAULT_PW="Couchbase1"

CLUSTER_NAME_PORT_MAP=(["C1"]=8091)
CLUSTER_NAME_XDCR_PORT_MAP=(["C1"]=9998)

cluster1BucketsArr=("source_A" "source_B" "source_C" "source_D" "source_E" "source_F" "source_G" "source_H" "source_I" "source_J" "source_K" "source_L" "source_M" "source_N" "source_O" "source_P")
CLUSTER_NAME_BUCKET_MAP=(["C1"]=${cluster1BucketsArr[@]})

declare -A DefaultBucketReplProperties=(["replicationType"]="continuous" ["checkpointInterval"]=60 ["statsInterval"]=500)

cluster2BucketsArr=("dest_A" "dest_B" "dest_C" "dest_D" "dest_E" "dest_F" "dest_G" "dest_H" "dest_I" "dest_J" "dest_K" "dest_L" "dest_M" "dest_N" "dest_O" "dest_P")

function getRemoteClusterUuid {
	uuid=$(curl -sX GET -u Administrator:Couchbase1 http://127.0.0.1:8091/pools/default/remoteClusters | jq '.[0]' | jq '.uuid' | sed 's/"//g')
	echo "$uuid"
}

function getReplRestID {
	local sourceBucket=$1
	local targetBucket=$2
	local remClusterId
	remClusterId=$(getRemoteClusterUuid)

	local restID="$remClusterId/$sourceBucket/$targetBucket"

	local restFriendlyReplID=$(echo $restID | sed 's|/|%2F|g')

	echo "$restFriendlyReplID"
}

for bucket in ${cluster1BucketsArr[@]}; do
	for bucket2 in ${cluster2BucketsArr[@]}; do
		restID=$(getReplRestID "$bucket" "$bucket2")
		insertBucketReplIntoExportMap "C1" "$bucket" "ship" "$bucket2" "$restID"
		insertPropertyIntoBucketReplPropertyMap "C1" "$bucket" "ship" "$bucket2" DefaultBucketReplProperties
	done
done

exportProvisionedConfig
