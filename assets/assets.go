package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.png
var assetsFS embed.FS

var (
	BoyWalkImage      *ebiten.Image
	BoyBikeImage      *ebiten.Image
	BoyWalkFrames     []*ebiten.Image
	BoyBikeFrames     []*ebiten.Image
	BikeShopImg       *ebiten.Image
	IceCreamVanImg    *ebiten.Image
	RecyclingBinImg   *ebiten.Image
	TreeImage         *ebiten.Image
	IceCreamImage     *ebiten.Image
	PlasticBottleImg  *ebiten.Image
	RedBusImage       *ebiten.Image
	RedPhoneBoxImg    *ebiten.Image
	LondonPanoramaImg *ebiten.Image
)

// LoadAssets loads and decodes all embedded PNG images into Ebitengine images.
func LoadAssets() {
	BoyWalkImage, _ = loadImage("boy_walk.png")
	BoyBikeImage, _ = loadImage("boy_bike.png")

	w1, err1 := loadImage("boy_walk_1.png")
	w2, err2 := loadImage("boy_walk_2.png")
	if err1 == nil && err2 == nil {
		BoyWalkFrames = []*ebiten.Image{w1, w2}
	} else if BoyWalkImage != nil {
		BoyWalkFrames = []*ebiten.Image{BoyWalkImage}
	}

	b1, err3 := loadImage("boy_bike_1.png")
	b2, err4 := loadImage("boy_bike_2.png")
	if err3 == nil && err4 == nil {
		BoyBikeFrames = []*ebiten.Image{b1, b2}
	} else if BoyBikeImage != nil {
		BoyBikeFrames = []*ebiten.Image{BoyBikeImage}
	}

	BikeShopImg, _ = loadImage("bike_shop.png")
	IceCreamVanImg, _ = loadImage("ice_cream_van.png")
	RecyclingBinImg, _ = loadImage("recycling_bin.png")
	TreeImage, _ = loadImage("tree.png")
	IceCreamImage, _ = loadImage("ice_cream.png")
	PlasticBottleImg, _ = loadImage("plastic_bottle.png")
	RedBusImage, _ = loadImage("red_bus.png")
	RedPhoneBoxImg, _ = loadImage("red_phone_box.png")
	LondonPanoramaImg, _ = loadImage("london_panorama.png")
}

func loadImage(filename string) (*ebiten.Image, error) {
	data, err := assetsFS.ReadFile(filename)
	if err != nil {
		log.Printf("Warning: failed to read %s: %v", filename, err)
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Warning: failed to decode %s: %v", filename, err)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
