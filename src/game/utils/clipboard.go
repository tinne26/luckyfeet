//go:build !wasm

// === IMPORTANT NOTICE ====================================================
// This file has been shamelessly stolen from https://github.com/ketMix/retromancer.
// The original license is GPL-3.0.
//
// If you are an ebitengine enthusiast and you haven't already, go check
// anything that kettek and liqMix create; there's a reason why they are 
// always on the first places of any ebitengine game jam.
// =========================================================================

package utils

import "golang.design/x/clipboard"

func init() {
	if err := clipboard.Init(); err != nil {
		panic(err)
	}
}

func ReadClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func WriteClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
