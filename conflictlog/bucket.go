package conflictlog

// BucketInfo, VBucketServerMap are the structs for the bucket topology json returned
// in the response body when memcached returns NOT_MY_BUCKET.

type BucketInfo struct {
	VBucketServerMap VBucketServerMap `json:"vBucketServerMap"`
}

// VBucketServerMap is the a mapping of vbuckets to nodes.
type VBucketServerMap struct {
	NumReplicas int      `json:"numReplicas"`
	ServerList  []string `json:"serverList"`
	VBucketMap  [][]int  `json:"vBucketMap"`
}

func (v *VBucketServerMap) GetAddrByVB(vbno uint16, replicaNum int) (addr string) {
	idx := v.VBucketMap[int(vbno)][replicaNum]
	return v.ServerList[idx]
}
