package data

import (
	"sync"
)

//TagContainer can be set to a object, then check if contains a tag
type TagContainer interface {
	Contains(names ...string) bool
	Add(names ...string)
	Remove(names ...string)
}

type TagContainerImpl struct {
	sync.Mutex
	Tags map[string]bool
}

func (t TagContainerImpl) Contains(names ...string) bool {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		_, ok := t.Tags[name]
		if !ok {
			return false
		}
	}

	return true
}

func (t TagContainerImpl) Add(names ...string) {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		t.Tags[name] = true
	}
}

func (t TagContainerImpl) Remove(names ...string) {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		delete(t.Tags, name)
	}
}
