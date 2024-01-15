package gfxcore

import "io/fs"
import "errors"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"

type Graphics struct {
	Tiles [][]*ebiten.Image // listed by type and with variations
	BackPatternMask *ebiten.Image
	BackPatternMaskX32 *ebiten.Image
	BackParticles [4]*ebiten.Image // NE, NW, SE, SW
	MenuMask *ebiten.Image
	BackLightingSmall *ebiten.Image
	BackLightingBig *ebiten.Image

	CarrotInUseMask *ebiten.Image
	CarrotSelector *ebiten.Image
	CarrotNone *ebiten.Image
	CarrotOrange *ebiten.Image
	CarrotYellow *ebiten.Image
	CarrotPurple *ebiten.Image
}

func maskToRGBA(mask []byte) []byte{
	out := make([]byte, len(mask)*4)
	var i int
	for _, value := range mask {
		out[i + 0] = 255*value
		out[i + 1] = 255*value
		out[i + 2] = 255*value
		out[i + 3] = 255*value
		i += 4
	}
	return out
}

func New(filesys fs.FS) (*Graphics, error) {
	var err error

	// generate some custom graphics
	backPatternMask := ebiten.NewImage(8, 8)
	backPatternMask.WritePixels(maskToRGBA([]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 1, 0,
		0, 1, 1, 1, 0, 0, 1, 0,
		0, 0, 1, 0, 0, 0, 1, 0,
	}))
	backPatternMaskX32 := ebiten.NewImage(8*32, 8)
	var opts ebiten.DrawImageOptions
	for i := 0; i < 32; i++ {
		backPatternMaskX32.DrawImage(backPatternMask, &opts)
		opts.GeoM.Translate(8, 0)
	}
	
	particleNE := ebiten.NewImage(3, 3)
	particleNW := ebiten.NewImage(3, 3)
	particleSE := ebiten.NewImage(3, 3)
	particleSW := ebiten.NewImage(3, 3)
	particleNE.WritePixels(maskToRGBA([]byte{ 0, 1, 1, /**/ 1, 0, 1, /**/ 0, 1, 0 }))
	particleNW.WritePixels(maskToRGBA([]byte{ 1, 1, 0, /**/ 1, 0, 1, /**/ 0, 1, 0 }))
	particleSE.WritePixels(maskToRGBA([]byte{ 0, 1, 0, /**/ 1, 0, 1, /**/ 0, 1, 1 }))
	particleSW.WritePixels(maskToRGBA([]byte{ 0, 1, 0, /**/ 1, 0, 1, /**/ 1, 1, 0 }))

	menuMask := ebiten.NewImage(9, 10)
	menuMask.WritePixels(maskToRGBA([]byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 1, 0, 1, 1, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 1, 0, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	}))

	// load tiles
	var tiles = make([][]*ebiten.Image, tcsts.TileTypeMax)
	
	const LayerMainPath = "assets/graphics/tiles/layer_main/"
	tiles[tcsts.MainGround], err = loadTileVariants(filesys, LayerMainPath + "ground_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGroundRaiser], err = loadTileVariants(filesys, LayerMainPath + "ground_raiser_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGroundSide], err = loadTileVariants(filesys, LayerMainPath + "ground_side_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGroundCorner], err = loadTileVariants(filesys, LayerMainPath + "ground_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGroundMark], err = loadTileVariants(filesys, LayerMainPath + "ground_mark_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGroundMarkCorner], err = loadTileVariants(filesys, LayerMainPath + "ground_mark_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.MainSinglePlatform], err = loadTileVariants(filesys, LayerMainPath + "platform_single_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGrassSide], err = loadTileVariants(filesys, LayerMainPath + "grass_side_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGrassSideFull], err = loadTileVariants(filesys, LayerMainPath + "grass_side_full_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGrassCorner], err = loadTileVariants(filesys, LayerMainPath + "grass_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.MainGrassCornerFull], err = loadTileVariants(filesys, LayerMainPath + "grass_corner_full_")
	if err != nil { return nil, err }

	tiles[tcsts.MainOrangePlatSingle], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_single_")
	if err != nil { return nil, err }
	tiles[tcsts.MainOrangePlatSingleFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_single_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainOrangePlatLeft], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_left_")
	if err != nil { return nil, err }
	tiles[tcsts.MainOrangePlatLeftFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_left_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainOrangePlatRight], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_right_")
	if err != nil { return nil, err }
	tiles[tcsts.MainOrangePlatRightFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_orange_plat_right_fill_")
	if err != nil { return nil, err }

	tiles[tcsts.MainYellowPlatSingle], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_single_")
	if err != nil { return nil, err }
	tiles[tcsts.MainYellowPlatSingleFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_single_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainYellowPlatLeft], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_left_")
	if err != nil { return nil, err }
	tiles[tcsts.MainYellowPlatLeftFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_left_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainYellowPlatRight], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_right_")
	if err != nil { return nil, err }
	tiles[tcsts.MainYellowPlatRightFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_yellow_plat_right_fill_")
	if err != nil { return nil, err }

	tiles[tcsts.MainPurplePlatSingle], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_single_")
	if err != nil { return nil, err }
	tiles[tcsts.MainPurplePlatSingleFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_single_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainPurplePlatLeft], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_left_")
	if err != nil { return nil, err }
	tiles[tcsts.MainPurplePlatLeftFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_left_fill_")
	if err != nil { return nil, err }
	tiles[tcsts.MainPurplePlatRight], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_right_")
	if err != nil { return nil, err }
	tiles[tcsts.MainPurplePlatRightFill], err = loadTileVariants(filesys, LayerMainPath + "carrot_purple_plat_right_fill_")
	if err != nil { return nil, err }

	const LayerBackPath = "assets/graphics/tiles/layer_back/"
	tiles[tcsts.BackGround], err = loadTileVariants(filesys, LayerBackPath + "ground_")
	if err != nil { return nil, err }
	tiles[tcsts.BackGroundSide], err = loadTileVariants(filesys, LayerBackPath + "ground_side_")
	if err != nil { return nil, err }
	tiles[tcsts.BackGroundCorner], err = loadTileVariants(filesys, LayerBackPath + "ground_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.BackGroundMark], err = loadTileVariants(filesys, LayerBackPath + "ground_mark_")
	if err != nil { return nil, err }
	tiles[tcsts.BackGroundMarkCorner], err = loadTileVariants(filesys, LayerBackPath + "ground_mark_corner_")
	if err != nil { return nil, err }

	const LayerFrontPath = "assets/graphics/tiles/layer_front/"
	tiles[tcsts.FrontGround], err = loadTileVariants(filesys, LayerFrontPath + "ground_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGroundRaiser], err = loadTileVariants(filesys, LayerFrontPath + "ground_raiser_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGroundSide], err = loadTileVariants(filesys, LayerFrontPath + "ground_side_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGroundCorner], err = loadTileVariants(filesys, LayerFrontPath + "ground_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGroundMark], err = loadTileVariants(filesys, LayerFrontPath + "ground_mark_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGroundMarkCorner], err = loadTileVariants(filesys, LayerFrontPath + "ground_mark_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontSinglePlatform], err = loadTileVariants(filesys, LayerFrontPath + "platform_single_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGrassSide], err = loadTileVariants(filesys, LayerFrontPath + "grass_side_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGrassSideFull], err = loadTileVariants(filesys, LayerFrontPath + "grass_side_full_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGrassCorner], err = loadTileVariants(filesys, LayerFrontPath + "grass_corner_")
	if err != nil { return nil, err }
	tiles[tcsts.FrontGrassCornerFull], err = loadTileVariants(filesys, LayerFrontPath + "grass_corner_full_")
	if err != nil { return nil, err }

	const LayerSpecialPath = "assets/graphics/tiles/layer_special/"
	tiles[tcsts.RaceGoal], err = loadTileVariants(filesys, LayerSpecialPath + "race_goal_")
	if err != nil { return nil, err }
	tiles[tcsts.StartPoint], err = loadTileVariants(filesys, LayerSpecialPath + "start_point_")
	if err != nil { return nil, err }
	tiles[tcsts.CarrotOrange], err = loadTileVariants(filesys, LayerSpecialPath + "carrot_orange_")
	if err != nil { return nil, err }
	tiles[tcsts.CarrotYellow], err = loadTileVariants(filesys, LayerSpecialPath + "carrot_yellow_")
	if err != nil { return nil, err }
	tiles[tcsts.CarrotPurple], err = loadTileVariants(filesys, LayerSpecialPath + "carrot_purple_")
	if err != nil { return nil, err }
	tiles[tcsts.CarrotMissing], err = loadTileVariants(filesys, LayerSpecialPath + "carrot_missing_")
	if err != nil { return nil, err }

	tiles[tcsts.TransferUp], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_up_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferUpA], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_upA_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferUpB], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_upB_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferUpC], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_upC_")
	if err != nil { return nil, err }

	tiles[tcsts.TransferDown], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_down_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferDownA], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_downA_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferDownB], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_downB_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferDownC], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_downC_")
	if err != nil { return nil, err }

	tiles[tcsts.TransferRight], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_right_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferRightA], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_rightA_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferRightB], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_rightB_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferRightC], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_rightC_")
	if err != nil { return nil, err }

	tiles[tcsts.TransferLeft], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_left_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferLeftA], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_leftA_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferLeftB], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_leftB_")
	if err != nil { return nil, err }
	tiles[tcsts.TransferLeftC], err = loadTileVariants(filesys, LayerSpecialPath + "transfer_leftC_")
	if err != nil { return nil, err }

	// load other assets
	backLightingSmall, err := loadImage(filesys, "assets/graphics/environment/back_lighting_small.png")
	if err != nil { return nil, err }
	backLightingBig, err := loadImage(filesys, "assets/graphics/environment/back_lighting_big.png")
	if err != nil { return nil, err }

	const UIGraphicsPath = "assets/graphics/ui/"
	carrotInUseMask, err := loadImage(filesys, UIGraphicsPath + "carrot_in_use_mask.png")
	if err != nil { return nil, err }
	carrotSelector, err := loadImage(filesys, UIGraphicsPath + "carrot_selector.png")
	if err != nil { return nil, err }
	carrotNone, err := loadImage(filesys, UIGraphicsPath + "carrot_none.png")
	if err != nil { return nil, err }
	carrotOrange, err := loadImage(filesys, UIGraphicsPath + "carrot_orange.png")
	if err != nil { return nil, err }
	carrotYellow, err := loadImage(filesys, UIGraphicsPath + "carrot_yellow.png")
	if err != nil { return nil, err }
	carrotPurple, err := loadImage(filesys, UIGraphicsPath + "carrot_purple.png")
	if err != nil { return nil, err }

	// create graphics struct
	gfx := &Graphics{
		Tiles: tiles,
		BackPatternMask: backPatternMask,
		BackPatternMaskX32: backPatternMaskX32,
		BackParticles: [4]*ebiten.Image{ particleNE, particleNW, particleSE, particleSW },
		MenuMask: menuMask,
		BackLightingSmall: backLightingSmall,
		BackLightingBig: backLightingBig,
		CarrotInUseMask: carrotInUseMask,
		CarrotSelector: carrotSelector,
		CarrotNone: carrotNone,
		CarrotOrange: carrotOrange,
		CarrotYellow: carrotYellow,
		CarrotPurple: carrotPurple,
	}
	return gfx, nil
}

func loadImage(filesys fs.FS, filepath string) (*ebiten.Image, error) {
	file, err := filesys.Open(filepath)
	if err != nil { return nil, err }
	img, err := png.Decode(file)
	if err != nil { return nil, err }
	return ebiten.NewImageFromImage(img), nil
}

func loadTileVariants(filesys fs.FS, basePath string) ([]*ebiten.Image, error) {
	var list []*ebiten.Image
	for c := 'A'; c <= 'Z'; c++ {
		file, err := filesys.Open(basePath + string(c) + ".png")
		if err != nil {
			if isNotExist(err) {
				if c == 'A' {
					err = errors.New("no tiles found for '" + basePath + "' pattern")
				} else {
					err = nil // we already got something
				}
			}
			return list, err
		}
		img, err := png.Decode(file)
		if err != nil { return list, err }
		list = append(list, ebiten.NewImageFromImage(img))
	}
	
	return list, nil
}

func isNotExist(err error) bool {
	if err == fs.ErrNotExist { return true }
	pathErr, isPathErr := err.(*fs.PathError)
	if !isPathErr { return false }
	return pathErr.Err == fs.ErrNotExist
}
