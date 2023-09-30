package tinycache

import "github.com/TremblingV5/TinyCache/internal/errno"

var (
	Success = errno.NewErrNo(0, "success")

	ErrBucketNotFound       = errno.NewErrNo(1001, "bucket not found")
	ErrKeyNotFound          = errno.NewErrNo(1002, "key not found")
	ErrBucketAlreadyExisted = errno.NewErrNo(1003, "bucket already existed")
)
