package info

import "github.com/tinne26/luckyfeet/src/lib/text"

var ControlsKB = []string{
	"MENU: " + string(text.KeyTAB),
	"MOVEMENT: WASD",
	"JUMP: SPACEBAR",
	"",
	"SELECT CARROT: " + string(text.KeyI) + " AND " + string(text.KeyP),
	"USE CARROT: " + string(text.KeyO),
}
var ControlsGP = []string{
	"MENU: START",
	"MOVEMENT: D-PAD",
	"CONFIRM/JUMP: BOTTOM BUTTON " + string(text.GpBtBottom),
	"CANCEL/BACK: RIGHT BUTTON " + string(text.GpBtRight),
	"",
	"SELECT CARROT: L/R SHOULDERS " + string(text.GpShoulders),
	"USE CARROT: LEFT BUTTON " + string(text.GpBtLeft),
}

var EditorControlsKB = []string{
	"MENU: " + string(text.KeyTAB),
	"MOVEMENT: WASD",
	"ADD/REMOVE BLOCK: ENTER/BACKSPACE",
	"CHANGE BLOCK/GROUP: ALT + MOVEMENT",
	"ROTATE BLOCK: Q/E",
	"MIRROR BLOCK: ALT + Q/E",
	"BLOCK VARIATION: V",
	"MODIFY SPECIAL BLOCK: ALT + ENTER",
}
var EditorControlsGP = []string{
	"MENU: START",
	"MOVEMENT: D-PAD",
	"ADD/REMOVE BLOCK: " + string(text.GpBtBottom) + "/" + string(text.GpBtRight) + " BUTTONS",
	"CHANGE BLOCK: L/R SHOULDERS " + string(text.GpShoulders),
	"CHANGE BLOCK GROUP: LEFT BUTTON " + string(text.GpBtLeft) + " + L/R SHOULDERS " + string(text.GpShoulders),
	"ROTATE BLOCK: L/R TRIGGERS " + string(text.GpTriggers),
	"MIRROR BLOCK: LEFT BUTTON "  + string(text.GpBtLeft) + " + L/R TRIGGERS " + string(text.GpTriggers),
	"BLOCK VARIATION: UP BUTTON " + string(text.GpBtTop),
	"MODIFY SPECIAL BLOCK: LEFT BUTTON " + string(text.GpBtLeft) + " + BOTTOM BUTTON " + string(text.GpBtBottom),
}
