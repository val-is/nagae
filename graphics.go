package nagae

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

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
