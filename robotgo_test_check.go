package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	pid := robotgo.GetPid()
	fmt.Printf("PID: %v\n", pid)

	x, y, w, h := robotgo.GetBounds(int(pid))
	fmt.Printf("Bounds: %d %d %d %d\n", x, y, w, h)

	bitmap, _ := robotgo.CaptureImg(x, y, w, h)
	fmt.Printf("Bitmap type: %T\n", bitmap)
}
