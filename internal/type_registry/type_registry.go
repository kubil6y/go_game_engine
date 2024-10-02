package type_registry

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrMaxItemsExceeded = errors.New("too many types registered!")
	ErrTypeNotFound     = errors.New("type not found")
	ErrNilItem          = errors.New("can not register nil item!")
)

type TypeRegistry struct {
	maxItems int
	nextID   int
	typeIDs  map[reflect.Type]int
	mu       sync.Mutex
}

func New(maxItems int) *TypeRegistry {
	return &TypeRegistry{
		maxItems: maxItems,
		nextID:   0,
		typeIDs:  make(map[reflect.Type]int),
	}
}

func (r *TypeRegistry) Size() int {
    return len(r.typeIDs)
}

func (r * TypeRegistry) GetTypeIDs() map[reflect.Type]int {
    return r.typeIDs
}

func (r *TypeRegistry) Register(item any) (int, error) {
	if item == nil {
		return -1, ErrNilItem
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.nextID >= r.maxItems {
		return -1, ErrMaxItemsExceeded
	}
	itemType := reflect.TypeOf(item)
	itemTypeID, exists := r.typeIDs[itemType]
	if exists {
		return itemTypeID, nil
	}
	id := r.nextID
	r.typeIDs[itemType] = id
	r.nextID++
	return id, nil
}

func (r *TypeRegistry) Get(item any) (int, error) {
	if item == nil {
		return -1, ErrNilItem
	}
	itemType := reflect.TypeOf(item)
	r.mu.Lock()
	defer r.mu.Unlock()
	itemTypeID, exists := r.typeIDs[itemType]
	if !exists {
		return -1, ErrTypeNotFound
	}
	return itemTypeID, nil
}
