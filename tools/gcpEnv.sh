#!/home/neil.huang/bash-5.0/bash
set -u

. ./clusterRunProvision.shlib
if (($? != 0)); then
	exit $?
fi

#function insertPropertyIntoBucketReplPropertyMap {
#	local sourceCluster=$1
#	local sourceBucket=$2
#	local targetCluster=$3
#	local targetBucket=$4
#	local -n incomingMap=$5
#
#	BUCKET_REPL_PROPERTIES_MAP=()
#	for key in ${!incomingMap[@]}; do
#		BUCKET_REPL_PROPERTIES_MAP["${sourceCluster},${sourceBucket},${targetCluster},${targetBucket},${key}"]=${incomingMap[$key]}
#	done
#	export BUCKET_REPL_PROPERTIES_MAP
#}

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

#CLUSTER_NAME_BUCKET_MAP=(["C1"]=${cluster1BucketsArr[@]})

declare -A DefaultBucketReplProperties=(["replicationType"]="continuous" ["checkpointInterval"]=60 ["statsInterval"]=500)


function getRemoteClusterUuid {
	uuid=$(curl -sX GET -u Administrator:Couchbase1 http://127.0.0.1:8091/pools/default/remoteClusters | jq '.[0]' | jq '.uuid' | sed 's/"//g')
	echo "$uuid"
}

function getReplRestID {
	local sourceBucket=$1
	local targetBucket=$2
	local remClusterId=$3

	local restID="$remClusterId/$sourceBucket/$targetBucket"

	local restFriendlyReplID=$(echo $restID | sed 's|/|%2F|g')

	echo "\"$restFriendlyReplID\""
}

function getListOfRepls {
  curl -sX GET -u Administrator:Couchbase1 http://127.0.0.1:9998/pools/default/replications | jq | grep \"id\" | cut -d: -f2 | sed 's/,//g' | sed 's/"//g' | cut -d/ -f2,3
}

remClusterId=$(getRemoteClusterUuid)
for oneRepl in $(getListOfRepls)
do
    bucket=$(echo "$oneRepl" | cut -d/ -f1)
    bucket2=$(echo "$oneRepl" | cut -d/ -f2)
		restID=$(getReplRestID "$bucket" "$bucket2" "$remClusterId")
		insertBucketReplIntoExportMap "C1" "$bucket" "ship" "$bucket2" "$restID"
		insertPropertyIntoBucketReplPropertyMap "C1" "$bucket" "ship" "$bucket2" DefaultBucketReplProperties
done

exportProvisionedConfig

