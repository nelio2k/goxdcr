#!/usr/bin/env bash
set -u

. ./clusterRunProvision.shlib
if (($? != 0)); then
	exit $?
fi

CLUSTER_NAME_PORT_MAP=(["C1"]=8091)
CLUSTER_NAME_XDCR_PORT_MAP=(["C1"]=9998)

cluster1BucketsArr=("source_A" "source_B" "source_C" "source_D" "source_E" "source_F" "source_G" "source_H" "source_I" "source_J" "source_K" "source_L" "source_M" "source_N" "source_O" "source_P")
CLUSTER_NAME_BUCKET_MAP=(["C1"]=${cluster1BucketsArr[@]})

declare -A DefaultBucketReplProperties=(["replicationType"]="continuous" ["checkpointInterval"]=60 ["statsInterval"]=500)

cluster2BucketsArr=("dest_A" "dest_B" "dest_C" "dest_D" "dest_E" "dest_F" "dest_G" "dest_H" "dest_I" "dest_J" "dest_K" "dest_L" "dest_M" "dest_N" "dest_O" "dest_P")

for bucket in ${cluster1BucketsArr[@]};do
  for bucket2 in ${cluster2BucketsArr[@]};do
    insertPropertyIntoBucketReplPropertyMap "C1" "$bucket" "ship" "$bucket2" DefaultBucketReplProperties
  done
done

exportProvisionedConfig
