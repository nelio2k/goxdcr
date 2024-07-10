package conflictlog

import (
	"container/list"
	"sync"
	"time"
)

type writerPool struct {
	writers        map[string]*writerList
	poolLock       sync.Mutex
	createWriterFn func(bucketName string) (Writer, error)
	timeout        time.Duration
}

type writerList struct {
	listLock sync.Mutex
	connList *list.List
}

func (wl *writerList) pop() Writer {
	wl.listLock.Lock()
	defer wl.listLock.Unlock()

	ele := wl.connList.Front()
	if ele == nil {
		return nil
	}
	w, _ := ele.Value.(Writer)
	wl.connList.Remove(ele)

	return w
}

func (wl *writerList) push(w Writer) {
	wl.listLock.Lock()
	wl.connList.PushBack(w)
	wl.listLock.Unlock()
}

func newWriterPool(createWriterFn func(bucketName string) (Writer, error)) *writerPool {
	return &writerPool{
		writers:        map[string]*writerList{},
		poolLock:       sync.Mutex{},
		createWriterFn: createWriterFn,
	}
}

func (pool *writerPool) getList(bucketName string) *writerList {
	pool.poolLock.Lock()
	wlist, ok := pool.writers[bucketName]
	if !ok {
		wlist = &writerList{
			listLock: sync.Mutex{},
			connList: list.New(),
		}
		pool.writers[bucketName] = wlist
	}
	defer pool.poolLock.Unlock()

	return wlist
}

func (pool *writerPool) get(bucketName string) (w Writer, err error) {
	wlist := pool.getList(bucketName)

	w = wlist.pop()
	if w != nil {
		return
	}

	w, err = pool.createWriterFn(bucketName)
	if err != nil {
		return
	}

	return
}

func (pool *writerPool) release(w Writer) {
	wlist := pool.getList(w.Bucket())
	wlist.push(w)
}
