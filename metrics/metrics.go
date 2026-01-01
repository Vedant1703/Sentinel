package metrics

import "sync/atomic"

var Allowed uint64
var Blocked uint64
var Errors uint64

func IncAllowed() {
	atomic.AddUint64(&Allowed, 1)
}

func IncBlocked() {
	atomic.AddUint64(&Blocked, 1)
}

func IncErrors() {
	atomic.AddUint64(&Errors, 1)
}
