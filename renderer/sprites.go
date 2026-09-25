package renderer

import (
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadSprite(fsys fs.FS, path string) (*ebiten.Image, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func LoadSprites(fsys fs.FS, paths map[string]string) (map[string]*ebiten.Image, error) {
	sprites := make(map[string]*ebiten.Image, len(paths))
	for name, path := range paths {
		sprite, err := LoadSprite(fsys, path)
		if err != nil {
			return nil, err
		}
		sprites[name] = sprite
	}
	return sprites, nil
}
