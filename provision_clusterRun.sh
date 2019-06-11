#!/bin/bash

# globals
DEFAULT_ADMIN="Administrator"
DEFAULT_PW="wewewe"
declare -a NODEPORTS
NODEPORTS=(9000 9001)
NODENAMES=("C1" "C2")
BUCKETNAMES=("B1" "B2")
RAMQUOTA=(100 100)

# Tests and see if "cluster_run -n x" has been run
function testForClusterRun {
	for port in ${NODEPORTS[@]}
		do
		curl -u $DEFAULT_ADMIN:$DEFAULT_PW -X GET http://localhost:$port/nodes/self/controller/settings > /dev/null 2>&1
		if (( $? != 0 ));then
			echo "Node $port not found. Skipping cluster_run XDCR init"
			return 1
		fi
	done
}

function setupCluster {
	for ((i=0; i < ${#NODEPORTS[@]}; i++))
	do
		echo "Setting up Node port ${NODEPORTS[$i]} name ${NODENAMES[$i]}"
		curl -u $DEFAULT_ADMIN:$DEFAULT_PW -X POST http://localhost:${NODEPORTS[$i]}/nodes/self/controller/settings > /dev/null 2>&1
		curl -X POST http://localhost:${NODEPORTS[$i]}/node/controller/rename -d hostname=127.0.0.1 > /dev/null 2>&1
		curl -X POST http://localhost:${NODEPORTS[$i]}/node/controller/setupServices -d services=kv > /dev/null 2>&1
		curl -u $DEFAULT_ADMIN:$DEFAULT_PW -v -X POST http://localhost:${NODEPORTS[$i]}/settings/web -d password=$DEFAULT_PW -d username=$DEFAULT_ADMIN -d port=${NODEPORTS[$i]} > /dev/null 2>&1

	done
}

function setupBuckets {
	for ((i=0; i < ${#NODEPORTS[@]}; i++))
	do
		echo "Setting up bucket ${BUCKETNAMES[$i]}"
curl -u $DEFAULT_ADMIN:$DEFAULT_PW -X POST http://localhost:${NODEPORTS[$i]}/pools/default/buckets -d ramQuotaMB=${RAMQUOTA[$i]} -d name=${BUCKETNAMES[$i]} > /dev/null 2>&1
	done
}

# Takes 2 arguments:
# 1- the ID number of node as source
# 2- The ID number of node as target
function createRemoteClusterReference {
	local sourceID=$1
	local targetID=$2

	echo "Creating remote cluster reference from ${NODENAMES[$sourceID]} to ${NODENAMES[$targetID]}"
	curl -X POST -u $DEFAULT_ADMIN:$DEFAULT_PW http://127.0.0.1:${NODEPORTS[$sourceId]}/pools/default/remoteClusters -d name=${NODENAMES[$targetID]} -d hostname=127.0.0.1:${NODEPORTS[$targetID]} -d username=$DEFAULT_ADMIN -d password=$DEFAULT_PW
	echo ""
}

function createReplicationInternal {
	local sourceID=$1
	local targetID=$2

	replicationID=`curl -X POST -u $DEFAULT_ADMIN:$DEFAULT_PW http://127.0.0.1:${NODEPORTS[$sourceId]}/controller/createReplication -d fromBucket=${BUCKETNAMES[$sourceID]} -d toCluster=${NODENAMES[$targetID]} -d toBucket=${BUCKETNAMES[$targetID]} -d replicationType=continuous -d checkpointInterval=60 -d statsInterval=500`
	echo "$replicationID"
}

function checkJQ {
	which jq > /dev/null 2>&1
	echo $?
}

function getKeyUsingJQ {
	local arg=$1
	local key=$2
	if [[ -z "$arg" ]];then
		return ""
	fi

	echo $arg | jq $key
}

function createReplication {
	local sourceID=$1
	local targetID=$2

	echo "Creating replication from ${NODENAMES[$sourceID]} to ${NODENAMES[$targetID]}"
	repId=`createReplicationInternal $sourceID $targetID`
	if (( `checkJQ` == 0 ));then
		repIdReal=`getKeyUsingJQ $repId '.id'`
		restFriendlyReplID=`echo $repIdReal | sed 's|/|%2F|g'`
		echo "ReplicationID: $repIdReal Rest-FriendlyID: $restFriendlyReplID"
	else
		echo "Captured replicationID: $repId"
	fi
}

#MAIN
testForClusterRun 
if (( $? == 0 ));then
	setupCluster 
	setupBuckets
	sleep 1
	createRemoteClusterReference 0 1
	createReplication 0 1
fi
