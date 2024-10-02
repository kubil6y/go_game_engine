package ecs

import (
	"container/list"
	"reflect"

	"github.com/kubil6y/go_game_engine/internal/utils"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/set"
)

var (
	entityIDGenerator    = utils.CreateIDGenerator()
	componentIDGenerator = utils.CreateIDGenerator()
)

type Entity int

type Component struct {
	ID int
}

type Register struct {
	numEntities               int
    // [index = entity id]
	entityComponentSignatures []bitset.Bitset32
    // [index = component id] [index = entity id]
	componentPools            [][]Component
	systems                   map[reflect.Type]*System
	entitiesToBeAdded         set.Set[Entity]
	entitiesToBeKilled        set.Set[Entity]
	freeIDs                   *list.List
}

func (r *Register) AddComponent(componentTypeId int) {
}
