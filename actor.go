package nagae

type Actor struct {
	parentScene Scene

	componentMask ComponentList
	components    map[ComponentId]Component
}

func (a Actor) GetComponentByType(componentType ComponentType) (Component, bool) {
	if !a.componentMask.CheckComponent(componentType) {
		return nil, false
	}
	for _, component := range a.components {
		if component.Type() == componentType {
			return component, true
		}
	}
	return nil, false
}

func (a Actor) GetComponentById(componentId ComponentId) (Component, bool) {
	if component, present := a.components[componentId]; !present {
		return nil, false
	} else {
		return component, true
	}
}

func (a *Actor) AddComponent(component Component) error {
	if a.componentMask.CheckComponent(component.Type()) {
		return ErrComponentPresent
	}
	a.components[component.Id()] = component
	a.componentMask.AddComponent(component.Type())
	return nil
}

func (a *Actor) RemoveComponentByType(componentType ComponentType) error {
	if !a.componentMask.CheckComponent(componentType) {
		return ErrComponentNotPresent
	}
	if component, present := a.GetComponentByType(componentType); !present {
		return ErrComponentNotPresent
	} else {
		delete(a.components, component.Id())
		return nil
	}
}

func (a *Actor) RemoveComponentById(componentId ComponentId) error {
	if component, present := a.GetComponentById(componentId); !present {
		return ErrComponentNotPresent
	} else {
		delete(a.components, component.Id())
		return nil
	}
}
