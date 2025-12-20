package etcd

import "errors"

var (
	errNoLease = errors.New("no lease available")
)
