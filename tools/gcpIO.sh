#!/usr/bin/env bash
set -u

. ./clusterRunProvision.shlib
if (($? != 0)); then
	exit $?
fi

function printHelp() {
	cat <<EOF
Usage: $0 [-s<hip>] -r <host> [ -r <host2> ]... -b <bucketName> [-b <bucketName2>]... -t <threadPerBucket> -l <totalRatePerBucket>

EOF
}

declare -a hostNames
declare -a bucketNames
declare -i threadPerBucket=0
declare -i totalRatePerBucket=0
declare -i shipSpecified

while getopts "hr:b:t:l:s" opt; do
	case ${opt} in
	h)
	  printHelp
	  exit 0
		;;
	r)
		hostNames+=("$OPTARG")
		;;
  b)
		bucketNames+=("$OPTARG")
		;;
  t)
    threadPerBucket=$OPTARG
    ;;
  l)
    totalRatePerBucket=$OPTARG
    ;;
  s)
    shipSpecified=1
    ;;
	esac
done

if (( $totalRatePerBucket == 0 ));then
  printHelp
  exit 1
fi

if (( $threadPerBucket == 0 ));then
  printHelp
  exit 1
fi


# Total rate / bucket = number of hosts * rate / host
# rate / host = (numberOfThread / bucket) * (rate / thread)

declare -i totalNumOfHosts=${#hostNames[@]}
declare -i ratePerHost=$(echo "$totalRatePerBucket/$totalNumOfHosts" | bc)
declare -i ratePerThread=$(echo "$ratePerHost/$threadPerBucket" | bc)

declare -i threadPerBucketPerNode=$(echo "$threadPerBucket/$totalNumOfHosts" | bc)

pillowfightBin=/opt/couchbase/bin/cbc-pillowfight

function launchWorker {
  local host=$1
  local bucket=$2

  #gcloud compute ssh xdcr-source-1 --zone us-central1-f -- "$pillowfightBin -P Couchbase1 -u Administrator -U couchbase://localhost/$bucket --rate-limit $ratePerThread -t $threadPerBucketPerNode" > /dev/null 2>&1
  gcloud compute ssh $host --zone us-central1-f -- "$pillowfightBin -P Couchbase1 -u Administrator -U couchbase://localhost/$bucket --rate-limit $ratePerThread -t $threadPerBucketPerNode" > /dev/null 2>&1
}

function ctrl_c() {
  echo "Stopping..."
  ps -ef | grep gcloud | grep pillow | awk '{print $2}' | xargs kill
  sleep 5
  ps -ef | grep gcloud | grep pillow | awk '{print $2}' | xargs kill
  exit 0
}

echo "Given $totalNumOfHosts nodes to do work on ${#bucketNames[@]}: $threadPerBucketPerNode threads/bucket/node $ratePerThread IO_limit/thread (per bucket)"
declare -i bgCounter=0
for host in ${hostNames[@]}
do
  for bucket in ${bucketNames[@]}
  do
    echo "Launching pillowfight on $host bucket $bucket..."
    launchWorker "$host" "$bucket" &
    bgCounter=$(( $bgCounter+1 ))
  done
done

sleep 5
echo "Running load. Ctrl-C to stop..."

trap ctrl_c INT

sleep 600000000
