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
	kvput := d.DeclareChannel(prefix+"KVPut", KVPut{})
	kvputr := d.DeclareChannel(prefix+"KVPutResponse", KVPutResponse{})
	kvget := d.DeclareChannel(prefix+"KVGet", KVGet{})
	kvgetr := d.DeclareChannel(prefix+"KVGetResponse", KVGetResponse{})

	kvmap := d.DeclareLMap(prefix + "kvMap")

	d.Join(kvput, func(k *KVPut) *KVPutResponse {
		return &KVPutResponse{k.ReqId, k.ClientAddr, d.Addr}
	}).IntoAsync(kvputr)

	d.Join(kvget, func(k *KVGet) *KVGetResponse {
		return &KVGetResponse{k.ReqId, k.ClientAddr, d.Addr, k.Key,
			kvmap.At(k.Key)}
	}).IntoAsync(kvgetr)

	d.Join(kvput, func(k *KVPut) (interface{}, Lattice) {
		return k.Key, k.Val
	}).Into(kvmap)

	return d
}
