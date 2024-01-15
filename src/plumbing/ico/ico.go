package ico

import "runtime"
import "image"
import "image/png"
import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"

func LoadAndSetWindowIcons(filesys fs.FS) error {
	if runtime.GOOS == "js" { return nil } // no window icons on the web

	file, err := filesys.Open("assets/graphics/ico/16x16.png")
	if err != nil { return err }
	ico16, err := png.Decode(file)
	if err != nil { return err }
	file, err = filesys.Open("assets/graphics/ico/32x32.png")
	if err != nil { return err }
	ico32, err := png.Decode(file)
	if err != nil { return err }
	file, err = filesys.Open("assets/graphics/ico/48x48.png")
	if err != nil { return err }
	ico48, err := png.Decode(file)
	if err != nil { return err }

	ebiten.SetWindowIcon([]image.Image{ico16, ico32, ico48})
	return nil
}

