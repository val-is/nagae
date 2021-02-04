package nagae

import (
	"errors"
	"fmt"
)

var (
	ErrComponentPresent    = errors.New("component is already present")
	ErrComponentNotPresent = errors.New("component is not present")

	ErrScenePresent = errors.New("scene is already present")
)

// ComponentType is an enum for ENGINE components. this defines what type of (default) component something is
type ComponentType uint16

const (
	ComponentTypeTransform ComponentType = iota
	ComponentTypeGraphical
	ComponentTypePhysics

	ComponentTypeSprite
	ComponentTypeSpriteAnimated

	ComponentTypeCustom
)

// ComponentSystem is an enum for ENGINE components. this defines what system uses the object
type ComponentSystem uint16

const (
	ComponentSystemDefault ComponentSystem = iota
	ComponentSystemTransform
	ComponentSystemGraphical
	ComponentSystemPhysics

	// ComponentSystemCustom is designed to be used for things like scripts
	// there can be INFINITE of the same type. be careful
	// (logically this is ok because these types are reserved for engine components)
	// (cont., any others aren't talking to the engine systems)
	ComponentSystemCustom
)

// ComponentList is a bitmask containing info on what ENGINE components are present
type ComponentList uint64

func (c ComponentList) AddComponent(other ComponentSystem) ComponentList {
	if other == ComponentSystemCustom {
		// we don't keep track of custom components that are added
		return c
	}
	return ComponentList(uint64(c) | (1 << uint16(other)))
}

func (c ComponentList) RemoveComponent(other ComponentSystem) ComponentList {
	return ComponentList(uint64(c) & (0 << uint16(other)))
}

func (c ComponentList) CheckComponent(other ComponentSystem) bool {
	return uint64(c)&(1<<uint16(other)) != 0
}

// ComponentId is a string identifier for components
type ComponentId string

// NOTE -- generating unique component ids is discouraged because of only allowing unique components per actor
var componentsCreated int = 0

func GenComponentId(baseId string) ComponentId {
	id := fmt.Sprintf("%s %d", baseId, componentsCreated)
	componentsCreated++
	return ComponentId(id)
}

// SceneId is just another identifier for a scene
type SceneId string

// ActorId is the same
type ActorId string
