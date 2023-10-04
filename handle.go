package tinycache

import (
	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/transmit"
)

func GetBucket(bucket string) *base.Bucket {
	return base.GetBucket(bucket)
}

func AddBucketLocally(bucket string, maxBytes int64) {
	base.AddBucketLocally(bucket, maxBytes)
}

func RemoveBucketLocally(bucket string) {
	base.RemoveBucketLocally(bucket)
}

func Get(bucket, key string) (base.ByteView, error) {
	v, err := base.GetBucket(bucket).GetLocally(key)
	if err == nil {
		return v, nil
	}

	v1, err := transmit.Get(bucket, key)
	if err != nil {
		return base.ByteView{}, err
	}
	return base.ByteView{
		B: v1,
	}, nil
}

func Set(bucket, key string, value []byte) error {
	return transmit.Set(bucket, key, value)
}

func Del(bucket, key string) error {
	return transmit.Delete(bucket, key)
}
