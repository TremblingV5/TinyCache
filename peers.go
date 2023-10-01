package tinycache

import "github.com/TremblingV5/TinyCache/pb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerHandle, ok bool)
}

// PeerHandle is the interface that must be implemented by a peer.
type PeerHandle interface {
	Get(in *pb.GetKeyRequest, out *pb.GetKeyResponse) error
	Set(in *pb.SetKeyRequest, out *pb.SetKeyResponse) error
	Delete(in *pb.DeleteKeyRequest, out *pb.DeleteKeyRequest) error
	DeleteBucket(in *pb.DeleteBucketRequest, out *pb.DeleteBucketResponse) error
}
