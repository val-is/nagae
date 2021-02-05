package nagae

import (
	"sort"

	"github.com/hajimehoshi/ebiten"
)

type System interface {
	Init() error
	Update(dt float64) error
}

type systemImpl struct {
	attachedScene *Scene
}

func (s systemImpl) Init() error             { return nil }
func (s systemImpl) Update(dt float64) error { return nil }

// PhysicsSystem handles updating physics objects (purely forces and velocities)
type PhysicsSystem interface {
	System
}

type physicsSystemImpl struct {
	systemImpl
}

func NewPhysicsSystem(scene *Scene) PhysicsSystem {
	return &physicsSystemImpl{
		systemImpl: systemImpl{
			attachedScene: scene,
		},
	}
}

func (p *physicsSystemImpl) Update(dt float64) error {
	for _, actor := range p.attachedScene.actors {
		physicsComp, present := actor.GetComponentBySystemType(ComponentSystemPhysics)
		if !present {
			continue
		}
		transformComp, present := actor.GetComponentBySystemType(ComponentSystemTransform)
		if !present {
			continue
		}

		// update velocity based on acceleration
		physicsCompImpl := physicsComp.(*componentPhysicsImpl)
		transformCompImpl := transformComp.(*componentTransformImpl)

		physicsCompImpl.frameAcceleration.MultScalar(dt)
		physicsCompImpl.velocity.Translate(physicsCompImpl.frameAcceleration)
		physicsCompImpl.frameAcceleration = Vec2{0, 0}

		// NOTE we don't care about collisions or friction in this system.
		// those will be handled by other systems that then apply forces to our physics body

		// update position based on velocity
		vel := physicsCompImpl.velocity
		vel.MultScalar(dt)
		transformCompImpl.pos.Translate(vel)
	}
	return nil
}

// GraphicsSystem handles drawing all components to the screen -- updating animation controllers is handled by the overarching system
type GraphicsSystem interface {
	System
	Draw(screen *ebiten.Image) error
}

type graphicsSystemImpl struct {
	systemImpl
}

func NewGraphicsSystem(scene *Scene) GraphicsSystem {
	return &graphicsSystemImpl{
		systemImpl: systemImpl{
			attachedScene: scene,
		},
	}
}

func (g *graphicsSystemImpl) Draw(screen *ebiten.Image) error {
	// TODO optimize. HACKY CURSED CODE
	// NOTE ALSO NEVER ROTATE SPRITES
	drawOrders := make(map[int][]DrawCall)
	drawLayers := make([]int, 0)
	for _, actor := range g.attachedScene.actors {
		graphicalComp, present := actor.GetComponentBySystemType(ComponentSystemGraphical)
		if !present {
			continue
		}
		transformComp, present := actor.GetComponentBySystemType(ComponentSystemTransform)
		if !present {
			continue
		}

		transformCompImpl := transformComp.(ComponentTransform)

		transPos := transformCompImpl.Position()
		transSize := transformCompImpl.Scale()
		transRot := transformCompImpl.Rotation()

		graphicalBaseCompImpl := graphicalComp.(ComponentGraphicalBase)
		drawCall := graphicalBaseCompImpl.Draw
		order := graphicalBaseCompImpl.DrawOrder()

		if !graphicalBaseCompImpl.Raw() {
			graphicalCompImpl := graphicalComp.(ComponentGraphical)

			img := graphicalCompImpl.ToDraw()
			if img == nil {
				continue
			}

			// calculate position relative to transform
			relativePos := graphicalCompImpl.RelativePos()
			// relativePos.MultScalar(0.01)
			relativeSize := graphicalCompImpl.Size()
			// relativeSize.MultScalar(0.01)
			relativeRot := graphicalCompImpl.Rotation()

			relativeSize.MultVec(transSize)

			s := relativeSize
			s.MultScalar(-0.5)
			relativePos.Translate(s)

			relativePos.Rotate(transRot)
			relativePos.Translate(transPos)
			relativeRot += transRot

			drawCall = GetDrawCall(img, relativePos.X, relativePos.Y, relativeSize.X, relativeSize.Y, relativeRot)
		}

		if calls, present := drawOrders[order]; present {
			drawOrders[order] = append(calls, drawCall)
		} else {
			drawLayers = append(drawLayers, order)
			drawOrders[order] = []DrawCall{drawCall}
		}
	}
	sort.Ints(drawLayers)
	for _, layer := range drawLayers {
		for _, order := range drawOrders[layer] {
			if err := order(screen); err != nil {
				return err
			}
		}
	}
	return nil
}
