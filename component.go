package nagae

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

var componentsCreated int = 0

func GenComponentId(baseId string) ComponentId {
	id := fmt.Sprintf("%s %d", baseId, componentsCreated)
	componentsCreated++
	return ComponentId(id)
}

// base component. type case to other things to use
type Component interface {
	Init() error
	Update(dt float64) error // for base components, the system controls this

	Type() ComponentType
	Id() ComponentId
	Parent() Actor
}

type componentImpl struct {
	cType      ComponentType
	id         ComponentId
	boundActor Actor
}

func (c componentImpl) Init() error             { return nil }
func (c componentImpl) Update(dt float64) error { return nil }

func (c componentImpl) Type() ComponentType { return c.cType }
func (c componentImpl) Id() ComponentId     { return c.id }
func (c componentImpl) Parent() Actor       { return c.boundActor }

func NewComponent(cType ComponentType, baseId string, parent Actor) (Component, error) {
	return &componentImpl{
		cType:      cType,
		id:         GenComponentId(baseId),
		boundActor: parent,
	}, nil
}

// graphical component. found and called when rendering
type ComponentGraphical interface {
	Component
	ToDraw() *ebiten.Image
	DrawOrder() int
}

type componentGraphicalImpl struct {
	componentImpl
	drawOrderPos int
}

func (c componentGraphicalImpl) ToDraw() *ebiten.Image { return nil }
func (c componentGraphicalImpl) DrawOrder() int        { return c.drawOrderPos }

func NewComponentGraphical(baseId string, parent Actor, drawOrderPos int) (Component, error) {
	baseComponent, err := NewComponent(ComponentTypeGraphical, baseId, parent)
	if err != nil {
		return nil, err
	}
	return &componentGraphicalImpl{
		componentImpl: baseComponent.(componentImpl),
		drawOrderPos:  drawOrderPos,
	}, nil
}

// transform interface. stores position info
type ComponentTransform interface {
	Component

	Position() Vec2
	SetPosition(newPos Vec2)
	Translate(delta Vec2)

	Scale() Vec2 // by default, 100px = 1 unit
	SetScale(newScale Vec2)
	ScaleBy(percent float64)
	ScaleTo(percent float64)

	Rotation() float64
	SetRotation(newRotation float64)
}

type componentTransformImpl struct {
	componentImpl
	pos      Vec2
	scale    Vec2
	rotation float64
}

func (c componentTransformImpl) Position() Vec2           { return c.pos }
func (c *componentTransformImpl) SetPosition(newPos Vec2) { c.pos = newPos }
func (c *componentTransformImpl) Translate(delta Vec2)    { c.pos.Translate(delta) }

func (c componentTransformImpl) Scale() Vec2              { return c.scale }
func (c *componentTransformImpl) SetScale(newScale Vec2)  { c.scale = newScale }
func (c *componentTransformImpl) ScaleBy(percent float64) { c.scale.MultScalar(percent) }
func (c *componentTransformImpl) ScaleTo(percent float64) { c.scale = Vec2{1, 1} }

func (c componentTransformImpl) Rotation() float64                { return c.rotation }
func (c *componentTransformImpl) SetRotation(newRotation float64) { c.rotation = newRotation }

func NewComponentTransform(parent Actor) (Component, error) {
	baseComponent, err := NewComponent(ComponentTypeTransform, "transform", parent)
	if err != nil {
		return nil, err
	}
	return &componentTransformImpl{
		componentImpl: baseComponent.(componentImpl),
	}, nil
}

// stores velocity, mass, and applies accelerations
type ComponentPhysics interface {
	Component

	Mass() float64
	SetMass(newMass float64) error

	Velocity() Vec2
	SetVelocity(newVel Vec2)

	SetFriction(friction Vec2)
	SetGravity(accel Vec2)

	Accelerate(acceleration Vec2)
	ApplyForce(force Vec2)
}

type componentPhysicsImpl struct {
	componentImpl

	mass     float64
	velocity Vec2
	friction Vec2
	gravity  Vec2

	frameAcceleration Vec2
}

func (c componentPhysicsImpl) Mass() float64 { return c.mass }
func (c *componentPhysicsImpl) SetMass(newMass float64) error {
	if newMass <= 0 {
		return fmt.Errorf("mass (%f) may not be less than or equal to zero", newMass)
	}
	c.mass = newMass
	return nil
}

func (c componentPhysicsImpl) Velocity() Vec2           { return c.velocity }
func (c *componentPhysicsImpl) SetVelocity(newVel Vec2) { c.velocity = newVel }

func (c *componentPhysicsImpl) SetFriction(friction Vec2) { c.friction = friction }
func (c *componentPhysicsImpl) SetGravity(accel Vec2)     { c.gravity = accel }

func (c *componentPhysicsImpl) Accelerate(acceleration Vec2) {
	c.frameAcceleration.Translate(acceleration)
}
func (c *componentPhysicsImpl) ApplyForce(force Vec2) {
	force.MultScalar(1 / c.Mass())
	c.Accelerate(force)
}

func NewComponentPhysics(parent Actor) (Component, error) {
	baseComponent, err := NewComponent(ComponentTypePhysics, "physics", parent)
	if err != nil {
		return nil, err
	}
	return &componentPhysicsImpl{
		componentImpl: baseComponent.(componentImpl),
		mass:          1,
	}, nil
}
