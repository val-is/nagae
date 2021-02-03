package nagae

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func LoadImageFromPath(path string) (*ebiten.Image, error) {
	imgReader, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}
	ebitenImage, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	return ebitenImage, nil
}

type Sprite interface {
	Image() *ebiten.Image
	GetSize() (float64, float64)
}

type spriteImpl struct {
	loadedImage   *ebiten.Image
	width, height float64
}

func (s spriteImpl) Image() *ebiten.Image        { return s.loadedImage }
func (s spriteImpl) GetSize() (float64, float64) { return s.width, s.height }

func NewStaticSprite(image *ebiten.Image) Sprite {
	wInt, hInt := image.Size()
	return &spriteImpl{
		loadedImage: image,
		width:       float64(wInt),
		height:      float64(hInt),
	}
}

type AnimatedSprite interface {
	Sprite
	Active() bool
	SetActive(active bool)

	CurrentFrame() int
	NumFrames() int

	NextFrame()
	SetFrame(frameNum int) bool

	SetLooping(loop bool)

	TicksPerFrame() int
	SetTicksPerFrame(ticks int)
	SetSecondsToRun(seconds float64)
	ResetTicks()
}

type animatedSpriteImpl struct {
	loadedFrames         []*ebiten.Image
	currentFrame         int
	ticks, ticksPerFrame int
	loop                 bool
	active               bool
}

func (a animatedSpriteImpl) Active() bool                { return a.active }
func (a *animatedSpriteImpl) SetActive(active bool)      { a.active = active }
func (a animatedSpriteImpl) CurrentFrame() int           { return a.currentFrame }
func (a animatedSpriteImpl) NumFrames() int              { return len(a.loadedFrames) }
func (a *animatedSpriteImpl) SetLooping(loop bool)       { a.loop = loop }
func (a animatedSpriteImpl) TicksPerFrame() int          { return a.ticks }
func (a *animatedSpriteImpl) ResetTicks()                { a.ticks = 0 }
func (a *animatedSpriteImpl) SetTicksPerFrame(ticks int) { a.ticks = ticks }

func (a *animatedSpriteImpl) Image() *ebiten.Image {
	if !a.active {
		return nil
	}
	a.ticks++
	if a.ticks > a.ticksPerFrame {
		a.NextFrame()
	}
	return a.loadedFrames[a.CurrentFrame()]
}

func (a animatedSpriteImpl) GetSize() (float64, float64) {
	intW, intH := a.loadedFrames[a.CurrentFrame()].Size()
	return float64(intW), float64(intH)
}

func (a *animatedSpriteImpl) NextFrame() {
	if !a.active {
		return
	}
	a.ticks = 0
	a.currentFrame++
	if a.currentFrame >= len(a.loadedFrames) {
		a.currentFrame = 0
		if !a.loop {
			a.active = false
		}
	}
}

func (a *animatedSpriteImpl) SetFrame(frameNum int) bool {
	if frameNum >= a.NumFrames() {
		return false
	}
	a.currentFrame = frameNum
	return true
}

func (a *animatedSpriteImpl) SetSecondsToRun(seconds float64) {
	// 60 fps constant draw loop call speed
	framesPerSecond := float64(a.NumFrames()) / seconds
	ticksPerFrame := 60 / framesPerSecond
	a.ticksPerFrame = int(ticksPerFrame)
}

type DrawCall func(screen *ebiten.Image) error

func GetDrawCall(image *ebiten.Image, x, y, w, h, angle float64) DrawCall {
	drawOptions := ebiten.DrawImageOptions{}
	drawOptions.GeoM.Reset()
	imageW, imageH := image.Size()
	drawOptions.GeoM.Scale(w/float64(imageW), h/float64(imageH))
	drawOptions.GeoM.Rotate(2 * math.Pi * angle)
	drawOptions.GeoM.Translate(x, y)
	return func(screen *ebiten.Image) error {
		return screen.DrawImage(image, &drawOptions)
	}
}
