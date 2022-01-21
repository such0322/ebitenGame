package assets

import (
	"bytes"
	"ebitenGame/mygame3/assets/fonts"
	"image/png"

	"ebitenGame/mygame3/assets/images"
	"ebitenGame/mygame3/spritetools"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
)

type Characters struct {
	Pacman *ebiten.Image
	Ghost1 *ebiten.Image
	Ghost2 *ebiten.Image
	Ghost3 *ebiten.Image
	Ghost4 *ebiten.Image
}

type Powers struct {
	Life          *ebiten.Image
	Invincibility *ebiten.Image
}

type Walls struct {
	ActiveCorner   *ebiten.Image
	ActiveSide     *ebiten.Image
	InActiveCorner *ebiten.Image
	InActiveSide   *ebiten.Image
}

type Assets struct {
	ArcadeFont *truetype.Font
	Skin       *ebiten.Image
	Characters *Characters
	Powers     *Powers
	Walls      *Walls
}

func LoadAssets() (*Assets, error) {
	skin, err := loadSkin()
	if err != nil {
		return nil, err
	}

	powers, err := loadPowers()
	if err != nil {
		return nil, err
	}

	characters, err := loadCharacters()
	if err != nil {
		return nil, err
	}

	font, err := loadArcadeFont()
	if err != nil {
		return nil, err
	}

	walls, err := loadWalls()
	if err != nil {
		return nil, err
	}
	return &Assets{
		Skin:       skin,
		Powers:     powers,
		Characters: characters,
		ArcadeFont: font,
		Walls:      walls,
	}, nil
}

func loadWalls() (*Walls, error) {
	wImage, err := png.Decode(bytes.NewReader(images.WallsPng))
	if err != nil {
		return nil, err
	}
	walls := ebiten.NewImageFromImage(wImage)

	inactiveCorner := spritetools.GetSprite(12, 12, 0, 0, walls)
	inactiveSide := spritetools.GetSprite(40, 12, 12, 0, walls)
	activeCorner := spritetools.GetSprite(12, 12, 52, 0, walls)
	activeSide := spritetools.GetSprite(40, 12, 64, 0, walls)
	return &Walls{
		ActiveCorner:   activeCorner,
		ActiveSide:     activeSide,
		InActiveCorner: inactiveCorner,
		InActiveSide:   inactiveSide,
	}, err
}

func loadArcadeFont() (*truetype.Font, error) {
	return truetype.Parse(fonts.ArcadeTTF)
}

func loadCharacters() (*Characters, error) {
	cImage, err := png.Decode(bytes.NewReader(images.CharactersPng))
	if err != nil {
		return nil, err
	}

	characters := ebiten.NewImageFromImage(cImage)

	pacman := spritetools.GetSprite(61, 64, 0, 0, characters)
	ghost1 := spritetools.GetSprite(56, 64, 66, 0, characters)
	ghost2 := spritetools.GetSprite(56, 64, 125, 0, characters)
	ghost3 := spritetools.GetSprite(56, 64, 185, 0, characters)
	ghost4 := spritetools.GetSprite(56, 64, 244, 0, characters)
	return &Characters{
		pacman,
		ghost1,
		ghost2,
		ghost3,
		ghost4,
	}, nil
}

func loadPowers() (*Powers, error) {
	pImage, err := png.Decode(bytes.NewReader(images.PowersPng))
	if err != nil {
		return nil, err
	}
	powers := ebiten.NewImageFromImage(pImage)
	life := spritetools.GetSprite(64, 64, 0, 0, powers)
	invinc := spritetools.GetSprite(64, 64, 67, 0, powers)
	return &Powers{
		life,
		invinc,
	}, nil
}

func loadSkin() (*ebiten.Image, error) {
	sImage, err := png.Decode(bytes.NewReader(images.SkinPng))
	if err != nil {
		return nil, err
	}
	skin := ebiten.NewImageFromImage(sImage)
	return skin, nil
}
