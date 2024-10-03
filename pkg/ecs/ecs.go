package ecs

import (
	"container/list"
	"fmt"
	"reflect"

	"github.com/kubil6y/go_game_engine/internal/type_registry"
	"github.com/kubil6y/go_game_engine/internal/utils"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/kubil6y/go_game_engine/pkg/set"
)

type Entity struct {
	ID int
}

type Component interface {
	GetID() int
	fmt.Stringer
}

func NewEntity(id int) Entity {
	return Entity{ID: id}
}

func (e Entity) GetID() int {
	return e.ID
}

type Registry struct {
	numEntities int
	// [index = entity id]
	entityComponentSignatures []*bitset.Bitset32
	// [index = component id] [index = entity id]
	componentPools        []*[]Component
	systems               map[reflect.Type]*System
	entitiesToBeAdded     set.Set[Entity]
	entitiesToBeKilled    set.Set[Entity]
	freeIDs               *list.List
	logger                *logger.Logger
	componentTypeRegistry *type_registry.TypeRegistry
}

func NewRegistry(maxComponentCount int, logger *logger.Logger) *Registry {
	return &Registry{
		numEntities:               0,
		entityComponentSignatures: make([]*bitset.Bitset32, 10),
		componentPools:            make([]*[]Component, 10),
		systems:                   make(map[reflect.Type]*System),
		entitiesToBeAdded:         set.New[Entity](),
		entitiesToBeKilled:        set.New[Entity](),
		freeIDs:                   list.New(),
		logger:                    logger,
		componentTypeRegistry:     type_registry.New(maxComponentCount),
	}
}

func (r *Registry) CreateEntity() Entity {
	var entityID int
	if r.freeIDs.Len() == 0 {
		r.numEntities++
		entityID = r.numEntities
		if entityID >= len(r.entityComponentSignatures) {
			utils.ResizeArray(r.entityComponentSignatures, entityID+1)
			for i := len(r.entityComponentSignatures); i <= entityID; i++ {
				r.entityComponentSignatures[i] = bitset.NewBitset32()
			}
		}
	} else {
		frontElement := r.freeIDs.Front()
		entityID = frontElement.Value.(int)
		r.freeIDs.Remove(frontElement)
	}
	entity := NewEntity(entityID)
	r.entitiesToBeAdded.Add(entity)
	r.logger.Info(fmt.Sprintf("Entity created with id = %d", entityID), nil)
	return entity
}

func (r *Registry) KillEntity(entity Entity) {
	r.logger.Info(fmt.Sprintf("Entity killed with id = %d", entity.GetID()), nil)
	r.entitiesToBeKilled.Add(entity)
}

func (r *Registry) AddComponent(entity Entity, component Component) error {
	entityID := entity.GetID()
	componentID, err := r.componentTypeRegistry.Register(component)
	if err != nil {
		switch err {
		case type_registry.ErrNilItem:
			panic("can not register null item")
		case type_registry.ErrMaxItemsExceeded:
			panic("too many types registered!")
		default:
			return nil
		}
	}

	if componentID >= len(r.componentPools) {
		newSize := componentID + 1
		r.componentPools = utils.ResizeArray(r.componentPools, newSize)
	}

	if r.componentPools[componentID] == nil {
		newComponentPool := make([]Component, r.numEntities)
		r.componentPools[componentID] = &newComponentPool
	}

	componentPool := r.componentPools[componentID]
	if entityID >= len(*componentPool) {
		newSize := entityID + 1 // Resize to at least accommodate the new entityID
		*componentPool = utils.ResizeArray(*componentPool, newSize)
	}
	(*componentPool)[entityID] = component
	r.logger.Info(fmt.Sprintf("%s registered with id: %d", component, componentID), nil)
	return nil
}

func (r *Registry) GetComponentTypeRegistry() *type_registry.TypeRegistry {
	return r.componentTypeRegistry
}
