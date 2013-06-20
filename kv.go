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

func KVInit(d *D, prefix string) *D {
	kvput  := d.DeclareChannel(prefix + "kvput", KVPut{})
	kvputr := d.DeclareChannel(prefix + "kvputresponse", KVPutResponse{})
	kvget  := d.DeclareChannel(prefix + "kvget", KVGet{})
	kvgetr := d.DeclareChannel(prefix + "kvgetresponse", KVGetResponse{})
	kvmap := d.DeclareLMap(prefix + "kvmap")

	kvmap.JoinUpdate(kvput,
		func(k *KVPut) (interface{}, Lattice) { return k.Key, k.Val })

	kvputr.JoinUpdateAsync(kvput,
		func(k *KVPut) *KVPutResponse {
			return &KVPutResponse{k.ReqId, k.ClientAddr, d.Addr}
		})

	kvgetr.JoinUpdateAsync(kvget,
		func(k *KVGet) *KVGetResponse {
			return &KVGetResponse{k.ReqId, k.ClientAddr, d.Addr, k.Key,
				kvmap.At(k.Key)}
		})

	return d
}
