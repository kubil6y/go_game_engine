package ecs

import (
	"container/list"
	"fmt"
	"math"
	"reflect"

	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/logger"
)

// NOTE:
// - Claude.AI helped a lot with reflect stuff lul
// - Resizing slices +1 is the original code im just converting to cpp to go

type SystemTypeID int
type ComponentTypeID int

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
	entityComponentSignatures []bitset.Bitset32
	// [index = component id] [index = entity id]
	componentPools     map[ComponentTypeID]any
	systems            map[SystemTypeID]System
	entitiesToBeAdded  []Entity
	entitiesToBeKilled []Entity
	freeIDs            *list.List
	logger             *logger.Logger
}

func NewRegistry(maxComponentCount int, logger *logger.Logger) *Registry {
	return &Registry{
		numEntities:               0,
		entityComponentSignatures: make([]bitset.Bitset32, 10),
		componentPools:            make(map[ComponentTypeID]any),
		systems:                   make(map[SystemTypeID]System),
		entitiesToBeAdded:         make([]Entity, 0),
		entitiesToBeKilled:        make([]Entity, 0),
		freeIDs:                   list.New(),
		logger:                    logger,
	}
}

// ENTITY MANAGEMENT ////////////////////
func (r *Registry) CreateEntity() Entity {
	var entityID int
	if r.freeIDs.Len() == 0 {
		r.numEntities++
		entityID = r.numEntities
		if entityID >= len(r.entityComponentSignatures) {
			// WARNING newSize := entityID + 1 // This is insane but thats the code in pikuma.com
			newSize := int(float32(len(r.entityComponentSignatures)) * 1.5)
			r.logger.Info(fmt.Sprintf("resize entityComponentSignatures %d -> %d", len(r.entityComponentSignatures), newSize), nil)
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

// COMPONENT MANAGEMENT ////////////////////
func (r *Registry) AddComponent(entity Entity, componentID ComponentTypeID, component interface{}) error {
	entityID := entity.GetID()

	// Check if the component pool exists, if not, create it
	if _, exists := r.componentPools[componentID]; !exists {
		initialCapacity := max(r.numEntities, 8) // Start with at least 8 elements
		r.componentPools[componentID] = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(component)), 0, initialCapacity).Interface()
	}

	// Get the component pool and ensure it's large enough
	pool := reflect.ValueOf(r.componentPools[componentID])
	if entityID >= pool.Cap() {
		newCapacity := max(entityID+1, int(math.Ceil(float64(pool.Cap())*1.5)))
		newPool := reflect.MakeSlice(pool.Type(), pool.Len(), newCapacity)
		reflect.Copy(newPool, pool)
		pool = newPool
		r.componentPools[componentID] = pool.Interface()
	}

	// Extend the slice if necessary
	if entityID >= pool.Len() {
		pool = pool.Slice(0, entityID+1)
		r.componentPools[componentID] = pool.Interface()
	}

	// Add the component to the pool
	pool.Index(entityID).Set(reflect.ValueOf(component))

	// Update the entity's component signature
	r.entityComponentSignatures[entityID].Set(int(componentID))

	return nil
}

// Helper function to find the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r *Registry) RemoveComponent(entity Entity, componentID ComponentTypeID) {
	entityID := entity.GetID()
	r.entityComponentSignatures[entityID].Clear(int(componentID))
}

func (r *Registry) HasComponent(entity Entity, componentID ComponentTypeID) bool {
	entityID := entity.GetID()
	signature := r.entityComponentSignatures[entityID]
	return signature.IsSet(int(componentID))
}

func (r *Registry) GetComponent(entity Entity, componentID ComponentTypeID) interface{} {
	pool := reflect.ValueOf(r.componentPools[componentID])
	return pool.Index(entity.GetID()).Interface()
}

func (r *Registry) GetComponentPtr(entity Entity, componentID ComponentTypeID) interface{} {
	pool := reflect.ValueOf(r.componentPools[componentID])
	return pool.Index(entity.GetID()).Addr().Interface()
}

// SYSTEM MANAGEMENT ////////////////////
func (r *Registry) AddSystem(systemID SystemTypeID, system System) {
	_, exists := r.systems[systemID]
	if !exists {
		r.systems[systemID] = system
		r.logger.Info(fmt.Sprintf("%s{%d} is registered", system.GetName(), systemID), nil)
	}
}

func (r *Registry) RemoveSystem(systemID SystemTypeID) {
	delete(r.systems, systemID)
}

func (r *Registry) GetSystem(systemID SystemTypeID) System {
	return r.systems[systemID]
}

func (r *Registry) HasSystem(systemID SystemTypeID) bool {
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
