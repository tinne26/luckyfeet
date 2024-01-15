//go:build windows

package autoconfig

func Apply() {
	DetectResizable()
	PreferOpenGL()
	DetectScreenshotKey()
	DetectWindowed()
}
