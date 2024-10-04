package ecs

import (
	"container/list"
	"fmt"

	"github.com/kubil6y/go_game_engine/internal/utils"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/logger"
)

type Entity struct {
	ID int
}

type Component interface {
	GetID() (int, error)
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
	entityComponentSignatures []bitset.Bitset32
	// [index = component id] [index = entity id]
	componentPools        []*[]Component
	systems               map[int]System
	entitiesToBeAdded     []Entity
	entitiesToBeKilled    []Entity
	freeIDs               *list.List
	logger                *logger.Logger
	componentTypeRegistry *TypeRegistry
	systemTypeRegistry    *TypeRegistry
}

func NewRegistry(maxComponentCount int, logger *logger.Logger, componentTypeRegistry *TypeRegistry, systemTypeRegistry *TypeRegistry) *Registry {
	return &Registry{
		numEntities:               0,
		entityComponentSignatures: make([]bitset.Bitset32, 10),
		componentPools:            make([]*[]Component, 10),
		systems:                   make(map[int]System),
		entitiesToBeAdded:         make([]Entity, 0),
		entitiesToBeKilled:        make([]Entity, 0),
		freeIDs:                   list.New(),
		logger:                    logger,
		componentTypeRegistry:     componentTypeRegistry,
		systemTypeRegistry:        systemTypeRegistry,
	}
}

func (r *Registry) GetComponentTypeRegistry() *TypeRegistry {
	return r.componentTypeRegistry
}

// ENTITY MANAGEMENT ////////////////////
func (r *Registry) CreateEntity() Entity {
	var entityID int
	if r.freeIDs.Len() == 0 {
		r.numEntities++
		entityID = r.numEntities
		if entityID >= len(r.entityComponentSignatures) {
			newSize := entityID + 1 // This is insane but thats the code in pikuma.com
			// newSize := int(float32(len(r.entityComponentSignatures)) * 1.5)
			// r.logger.Info(fmt.Sprintf("resize entityComponentSignatures %d -> %d", len(r.entityComponentSignatures), newSize), nil)
			newSignatureSlice := make([]bitset.Bitset32, newSize)
			for i := 0; i < len(r.entityComponentSignatures); i++ {
				newSignatureSlice[i] = r.entityComponentSignatures[i]
			}
			r.entityComponentSignatures = newSignatureSlice
		}
	} else {
		frontElement := r.freeIDs.Front()
		entityID = frontElement.Value.(int)
		r.freeIDs.Remove(frontElement)
	}
	entity := NewEntity(entityID)

	// Handle entities to be added
	var exists bool
	for _, e := range r.entitiesToBeAdded {
		if e.GetID() == entity.GetID() {
			exists = true
			break
		}
	}
	if !exists {
		r.entitiesToBeAdded = append(r.entitiesToBeAdded, entity)
	}

	r.logger.Debug(fmt.Sprintf("Entity{%d} created", entityID), nil)
	return entity
}

func (r *Registry) KillEntity(entity Entity) {
	r.logger.Debug(fmt.Sprintf("Entity killed with id = %d", entity.GetID()), nil)
	// Handle entities to be added
	var exists bool
	for _, e := range r.entitiesToBeKilled {
		if e.GetID() == entity.GetID() {
			exists = true
			break
		}
	}
	if !exists {
		r.entitiesToBeKilled = append(r.entitiesToBeKilled, entity)
	}
}

func (r *Registry) RegisterComponent(component Component) error {
	_, err := r.componentTypeRegistry.Register(component)
	if err != nil {
		switch err {
		case ErrNilItem:
			panic("can not register null item")
		case ErrMaxItemsExceeded:
			panic("too many types registered!")
		default:
			return nil
		}
	}
	componentID, err := component.GetID()
	if err != nil {
		return err
	}
	r.logger.Info(fmt.Sprintf("%s{%d} is registered", component.String(), componentID), nil)
	return nil
}

// COMPONENT MANAGEMENT ////////////////////
func (r *Registry) AddComponent(entity Entity, component Component) error {
	entityID := entity.GetID()
	componentID, err := r.componentTypeRegistry.Get(component)
	if err != nil {
		return err
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

	r.entityComponentSignatures[entityID].Set(componentID)
	r.logger.Debug(fmt.Sprintf("%s{%d} added to entity{%d}", component, componentID, entityID), nil)
	return nil
}

func (r *Registry) RemoveComponent(entity Entity, component Component) {
	entityID := entity.GetID()
	componentID, err := r.componentTypeRegistry.Get(component)
	if err != nil {
		r.logger.Error(err, fmt.Sprintf("Failed to retrieve %s component", component), nil)
		return
	}
	r.entityComponentSignatures[entityID].Clear(componentID)
}

func (r *Registry) HasComponent(entity Entity, component Component) bool {
	entityID := entity.GetID()
	componentID, err := r.componentTypeRegistry.Get(component)
	if err != nil {
		r.logger.Error(err, fmt.Sprintf("Failed to retrieve %s component", component), nil)
		return false
	}
	signature := r.entityComponentSignatures[entityID]
	return signature.IsSet(componentID)
}

func (r *Registry) GetComponent(entity Entity, component Component) Component {
	componentID, err := r.componentTypeRegistry.Get(component)
	if err != nil {
		r.logger.Error(err, fmt.Sprintf("Registry failed to add [%s] to Entity{%d}", component, entity.GetID()), nil)
	}
	return (*r.componentPools[componentID])[entity.GetID()]
}

// SYSTEM MANAGEMENT ////////////////////
func (r *Registry) AddSystem(system System) {
	systemID, err := r.systemTypeRegistry.Register(system)

	if err != nil {
		r.logger.Error(err, fmt.Sprintf("could not register system: %s", system.GetName()), nil)
	}

	_, exists := r.systems[systemID]
	if !exists {
		r.systems[systemID] = system
		r.logger.Info(fmt.Sprintf("%s{%d} is registered", system.GetName(), systemID), nil)
	}
}

func (r *Registry) RemoveSystem(system System) {
	systemID, err := r.systemTypeRegistry.Get(system)
	if err != nil {
		r.logger.Error(err, fmt.Sprintf("could not get system: %s", system.GetName()), nil)
		return
	}
	delete(r.systems, systemID)
}

func (r *Registry) GetSystem(systemID int) System {
	return r.systems[systemID]
}

func (r *Registry) HasSystem(systemID int) bool {
	_, exists := r.systems[systemID]
	return exists
}

func (r *Registry) Update() {
	for _, entity := range r.entitiesToBeAdded {
		r.AddEntityToSystems(entity)
	}
	r.entitiesToBeAdded = r.entitiesToBeAdded[:0]

	for _, entity := range r.entitiesToBeKilled {
		fmt.Printf("TODO entitiesToBeKilled: %d\n", entity.GetID())
	}
	r.entitiesToBeKilled = r.entitiesToBeKilled[:0]
}

func (r *Registry) AddEntityToSystems(entity Entity) {
	entitySignature := r.entityComponentSignatures[entity.GetID()]
	for _, system := range r.systems {
		if (entitySignature.Get32() & system.GetSignature().Get32()) == system.GetSignature().Get32() {
			system.AddEntityToSystem(entity)
		}
	}
}

func (r *Registry) RemoveEntityFromSystems(entity Entity) {
	for _, system := range r.systems {
		system.RemoveEntityFromSystem(entity)
		r.logger.Debug(fmt.Sprintf("Entity{%d} %d is removed from %s", entity.GetID(), system.GetName()), nil)
	}
}
