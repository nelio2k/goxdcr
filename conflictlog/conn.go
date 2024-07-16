package conflictlog

import (
	"io"
)

type Connection interface {
	io.Closer
	SetMeta(key string, val []byte, dataType uint8, target Target) (err error)
	SetMetaObj(key string, obj interface{}, target Target) (err error)
}
