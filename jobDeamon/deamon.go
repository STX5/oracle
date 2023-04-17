package jobdeamon

import clientv3 "go.etcd.io/etcd/client/v3"

// if a ETCD lock expires and its key
// still exist, means the worker failed.
// JobDeamon should update this key's version,
// so other wokers will receive watch event and do it again.
// JobDeamon should be deployed with ETCD member node
type JobDeamon struct {
	ETCDClient *clientv3.Client
}
