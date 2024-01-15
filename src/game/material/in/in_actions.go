package in

import "github.com/tinne26/luckyfeet/src/lib/input"

const (
	ActionFullscreen input.TriggerAction = iota + 1
	ActionToggleFPS
	ActionConfirm
	ActionBack
	ActionMenu
	
	ActionJump
	ActionUseCarrot
	ActionNextCarrot
	ActionPrevCarrot
	
	// --- editor bs ---
	ActionModKey
	ActionNextTile
	ActionPrevTile
	ActionNextTileGroup
	ActionPrevTileGroup
	ActionTileVariation
	ActionRotateTileRight
	ActionRotateTileLeft
	ActionMirrorTile

	// ...
)
