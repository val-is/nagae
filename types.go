package nagae

import "errors"

var (
	ErrComponentPresent    = errors.New("component is already present")
	ErrComponentNotPresent = errors.New("component is not present")
)

// ComponentType is an enum for ENGINE components
type ComponentType uint16

const (
	ComponentTypeTransform ComponentType = iota
	ComponentTypeGraphical
	ComponentTypePhysics

	// ComponentTypeCustom is designed to be used for things like scripts
	// there can be INFINITE of the same type. be careful
	// (logically this is ok because these types are reserved for engine components)
	// (cont., any others aren't talking to the engine systems)
	ComponentTypeCustom
)

// ComponentList is a bitmask containing info on what ENGINE components are present
type ComponentList uint64

func (c ComponentList) AddComponent(other ComponentType) ComponentList {
	if other == ComponentTypeCustom {
		// we don't keep track of custom components that are added
		return c
	}
	return ComponentList(uint64(c) & (1 << uint16(other)))
}

func (c ComponentList) RemoveComponent(other ComponentType) ComponentList {
	return ComponentList(uint64(c) & (0 << uint16(other)))
}

func (c ComponentList) CheckComponent(other ComponentType) bool {
	return uint64(c)&(1<<uint16(other)) == 1
}

// ComponentId is a string identifier for components
type ComponentId string
