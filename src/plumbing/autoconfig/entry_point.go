//go:build !windows

package autoconfig

func Apply() {
	DetectResizable()
	DetectScreenshotKey()
	DetectWindowed()
}
