package in

import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/input"

func LoadAndConfigure(kbgp *input.KBGP) error {
	// keyboard configuration
	kbConfig := kbgp.Keyboard().Config()
	kbConfig.SetDirTriggers(input.DirKeysArrowsAndWASD())
	kbConfig.MapTriggerActionToKey(ActionJump, ebiten.KeySpace)
	kbConfig.MapTriggerActionToKey(ActionUseCarrot, ebiten.KeyO)
	kbConfig.MapTriggerActionToKey(ActionPrevCarrot, ebiten.KeyI)
	kbConfig.MapTriggerActionToKey(ActionNextCarrot, ebiten.KeyP)
	kbConfig.MapTriggerActionToKey(ActionFullscreen, ebiten.KeyF)
	kbConfig.MapTriggerActionToKey(ActionToggleFPS, ebiten.KeyDigit1)
	kbConfig.MapTriggerActionToKey(ActionModKey, ebiten.KeyAltLeft)
	var multi input.MultiKey
	multi.AddNormalKey(ebiten.KeyTab)
	multi.AddBreakerKey(ebiten.KeyAltLeft)
	kbConfig.MapTriggerAction(ActionMenu, multi)
	kbConfig.MapTriggerActionToKey(ActionMenuBrowserAlt, ebiten.KeyM)
	kbConfig.MapTriggerAction(ActionBack, input.NewKeyList(ebiten.KeyBackspace, ebiten.KeyEscape))
	kbConfig.MapTriggerAction(ActionConfirm,
		input.KeyList([]ebiten.Key{ebiten.KeyEnter, ebiten.KeySpace}),
	)
	kbConfig.MapTriggerAction(ActionNextTile, input.NewMultiKey(ebiten.KeyD, ebiten.KeyAltLeft))
	kbConfig.MapTriggerAction(ActionPrevTile, input.NewMultiKey(ebiten.KeyA, ebiten.KeyAltLeft))
	kbConfig.MapTriggerAction(ActionNextTileGroup, input.NewMultiKey(ebiten.KeyS, ebiten.KeyAltLeft))
	kbConfig.MapTriggerAction(ActionPrevTileGroup, input.NewMultiKey(ebiten.KeyW, ebiten.KeyAltLeft))
	kbConfig.MapTriggerActionToKey(ActionTileVariation, ebiten.KeyV)
	kbConfig.MapTriggerActionToKey(ActionRotateTileRight, ebiten.KeyE)
	kbConfig.MapTriggerActionToKey(ActionRotateTileLeft, ebiten.KeyQ)

	// gamepad configuration
	gpConfig := kbgp.Gamepad().Config()
	gpConfig.SetDirButtons(input.DirButtonsDPad())
	gpConfig.MapTriggerActionToButton(ActionJump, input.GamepadButtonBottom)
	gpConfig.MapTriggerActionToButton(ActionFullscreen, input.GamepadSelect)
	gpConfig.MapTriggerActionToButton(ActionMenu, input.GamepadStart)
	gpConfig.MapTriggerActionToButton(ActionConfirm, input.GamepadButtonBottom)
	gpConfig.MapTriggerActionToButton(ActionBack, input.GamepadButtonRight)
	gpConfig.MapTriggerActionToButton(ActionModKey, input.GamepadButtonLeft)
	gpConfig.MapTriggerActionToButton(ActionUseCarrot, input.GamepadButtonLeft)
	gpConfig.MapTriggerActionToButton(ActionPrevCarrot, input.GamepadShoulderLeft)
	gpConfig.MapTriggerActionToButton(ActionNextCarrot, input.GamepadShoulderRight)

	var multiPrevTile input.MultiButton
	multiPrevTile.AddBreaker(input.GamepadButtonLeft)
	multiPrevTile.Add(input.GamepadShoulderLeft)
	gpConfig.MapTriggerAction(ActionPrevTile, multiPrevTile)
	var multiNextTile input.MultiButton
	multiNextTile.AddBreaker(input.GamepadButtonLeft)
	multiNextTile.Add(input.GamepadShoulderRight)
	gpConfig.MapTriggerAction(ActionNextTile, multiNextTile)
	var multiPrevTileGroup input.MultiButton
	multiPrevTileGroup.Add(input.GamepadButtonLeft)
	multiPrevTileGroup.Add(input.GamepadShoulderLeft)
	gpConfig.MapTriggerAction(ActionPrevTileGroup, multiPrevTileGroup)
	var multiNextTileGroup input.MultiButton
	multiNextTileGroup.Add(input.GamepadButtonLeft)
	multiNextTileGroup.Add(input.GamepadShoulderRight)
	gpConfig.MapTriggerAction(ActionNextTileGroup, multiNextTileGroup)
	gpConfig.MapTriggerActionToButton(ActionTileVariation, input.GamepadButtonTop)
	gpConfig.MapTriggerActionToButton(ActionRotateTileRight, input.GamepadTriggerRight)
	gpConfig.MapTriggerActionToButton(ActionRotateTileLeft, input.GamepadTriggerLeft)
	
	return nil
}

func LoadAndConfigureGamepad(gamepad *input.Gamepad, filesys fs.FS) error {
	
	return nil
}
