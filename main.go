package main

import (
	"fmt"
	"image"
	"image/color"
	"gocv.io/x/gocv"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting Realtime CV Pipeline...")

	// Simulate opening a webcam
	// webcam, err := gocv.OpenVideoCapture(0)
	// if err != nil {
	// 	fmt.Printf("Error opening video capture: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer webcam.Close()

	// Simulate a video frame
	img := gocv.NewMatWithSize(480, 640, gocv.MatTypeCV8U)
	defer img.Close()

	window := gocv.NewWindow("Realtime CV")
	defer window.Close()

	fmt.Println("Processing frames...")
	for i := 0; i < 5; i++ { // Process 5 dummy frames
		// In a real scenario, webcam.Read(&img) would be here
		// Simulate some processing
		gocv.Rectangle(&img, image.Rect(50+i*10, 50+i*10, 100+i*10, 100+i*10), color.RGBA{255, 0, 0, 0}, 2)
		window.IMShow(img)
		window.WaitKey(100) // Display for 100ms
		fmt.Printf("Processed frame %d\n", i+1)
	}

	fmt.Println("Realtime CV Pipeline finished.")
}
