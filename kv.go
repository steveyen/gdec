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

func KVProtocolInit(d *D, prefix string) *D {
	d.DeclareChannel(prefix+"KVPut", KVPut{})
	d.DeclareChannel(prefix+"KVPutResponse", KVPutResponse{})
	d.DeclareChannel(prefix+"KVGet", KVGet{})
	d.DeclareChannel(prefix+"KVGetResponse", KVGetResponse{})
	return d
}

// Simple KV replica that merges the values for a key, which works for
// monotonically increasing LMap's.

func KVInit(d *D, prefix string) *D {
	KVProtocolInit(d, prefix)

	kvput := d.Channels[prefix+"KVPut"]
	kvputr := d.Channels[prefix+"KVPutResponse"]
	kvget := d.Channels[prefix+"KVGet"]
	kvgetr := d.Channels[prefix+"KVGetResponse"]

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

type KVRepl struct {
	Addr       string
	TargetAddr string
}

type KVReplPropagate struct {
	Addr  string
	KVMap *LMap
}

func ReplicatedKVInit(d *D, prefix string) *D {
	KVInit(d, prefix)

	kvrepl := d.DeclareChannel(prefix+"KVRepl", KVRepl{})
	kvreplPropagate := d.DeclareChannel(prefix+"KVReplPropagate", KVReplPropagate{})

	kvmap := d.Lattices[prefix+"kvMap"].(*LMap)

	d.Join(kvrepl, func(r *KVRepl) *KVReplPropagate {
		return &KVReplPropagate{r.TargetAddr, kvmap}
	}).IntoAsync(kvreplPropagate)

	d.Join(kvreplPropagate, func(r *KVReplPropagate) *LMap {
		return r.KVMap
	}).Into(kvmap)

	return d
}
