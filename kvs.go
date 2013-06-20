package gdec

type KVPut struct {
	ReqId      int64
	Addr       string
	ClientAddr string
	Key        string
	Val        Lattice
}
type KVPutResponse struct {
	ReqId       int64
	Addr        string
	ReplicaAddr string
}
type KVGet struct {
	ReqId      int64
	Addr       string
	ClientAddr string
	Key        string
}
type KVGetResponse struct {
	ReqId       int64
	Addr        string
	ReplicaAddr string
	Key         string
	Val         Lattice
}
