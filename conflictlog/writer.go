package conflictlog

import (
	"io"
)

type Writer interface {
	io.Closer
	SetMeta(key string, val []byte, dataType uint8) (err error)
	SetMetaObj(key string, obj interface{}) (err error)
	Bucket() string
}
