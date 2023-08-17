package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gocv.io/x/gocv"
)

// ObjectDetector interface for different detection models
type ObjectDetector interface {
	Detect(img gocv.Mat) ([]image.Rectangle, error)
}

// SimpleFaceDetector implements ObjectDetector using Haar Cascades
type SimpleFaceDetector struct {
	classifier gocv.CascadeClassifier	
}

func NewSimpleFaceDetector(xmlPath string) (*SimpleFaceDetector, error) {
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(xmlPath) {
		return nil, fmt.Errorf("error reading cascade file: %s", xmlPath)
	}
	return &SimpleFaceDetector{classifier: classifier}, nil
}

func (s *SimpleFaceDetector) Detect(img gocv.Mat) ([]image.Rectangle, error) {
	rects := s.classifier.DetectMultiScale(img)
	return rects, nil
}

// Tracker interface for different tracking algorithms
type ObjectTracker interface {
	Init(img gocv.Mat, bbox image.Rectangle) error
	Update(img gocv.Mat) (image.Rectangle, error)
}

// KCFTracker implements ObjectTracker using KCF algorithm
type KCFTracker struct {
	tracker gocv.Tracker
}

func NewKCFTracker() (*KCFTracker, error) {
	tracker := gocv.NewTrackerKCF()
	return &KCFTracker{tracker: tracker}, nil
}

func (k *KCFTracker) Init(img gocv.Mat, bbox image.Rectangle) error {
	return k.tracker.Init(img, bbox)
}

func (k *KCFTracker) Update(img gocv.Mat) (image.Rectangle, error) {
	bbox, _ := k.tracker.Update(img)
	return bbox, nil
}

// RealtimeCVPipeline orchestrates detection and tracking
type RealtimeCVPipeline struct {
	webcam *gocv.VideoCapture
	detector ObjectDetector
	tracker ObjectTracker
	window *gocv.Window
	tracking bool
	trackedBox image.Rectangle
}

func NewRealtimeCVPipeline(deviceID int, detector ObjectDetector, tracker ObjectTracker) (*RealtimeCVPipeline, error) {
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("error opening video capture device: %v", err)
	}

	window := gocv.NewWindow("Realtime CV Pipeline")

	return &RealtimeCVPipeline{
		webcam: webcam,
		detector: detector,
		tracker: tracker,
		window: window,
		tracking: false,
	},
	nil
}

func (p *RealtimeCVPipeline) Run() {
	img := gocv.NewMat()
	defer img.Close()
	defer p.webcam.Close()
	defer p.window.Close()

	fmt.Println("Starting Realtime CV Pipeline. Press ESC to quit.")
	fmt.Println("Press D to toggle detection, T to toggle tracking.")

	detectionInterval := 1 * time.Second
	lastDetectionTime := time.Now()

	for {
		if ok := p.webcam.Read(&img); !ok || img.Empty() {
			fmt.Printf("Cannot read device %d\n", p.webcam.Device)
			return
		}

		// Perform detection periodically or if not tracking
		if !p.tracking && time.Since(lastDetectionTime) > detectionInterval {
			rects, err := p.detector.Detect(img)
			if err != nil {
				log.Printf("Detection error: %v", err)
			} else if len(rects) > 0 {
				// Start tracking the first detected object
				p.trackedBox = rects[0]
				err = p.tracker.Init(img, p.trackedBox)
				if err != nil {
					log.Printf("Tracker init error: %v", err)
					p.tracking = false
				} else {
					p.tracking = true
					fmt.Println("Started tracking object.")
				}
			}
			lastDetectionTime = time.Now()
		}

		if p.tracking {
			bbox, err := p.tracker.Update(img)
			if err != nil || bbox.Dx() == 0 || bbox.Dy() == 0 {
				log.Printf("Tracking lost: %v", err)
				p.tracking = false
				fmt.Println("Tracking lost. Re-initializing detector.")
			} else {
				p.trackedBox = bbox
				gocv.Rectangle(&img, p.trackedBox, color.RGBA{0, 255, 0, 0}, 2)
			}
		}

		p.window.IMShow(img)
		key := p.window.WaitKey(1)

		if key == 27 { // ESC key
			break
		}
		// Add more key handlers for toggling detection/tracking if needed
	}
}

func main() {
	// Ensure the Haar Cascade XML file is available
	// You might need to download it from OpenCV's GitHub: https://github.com/opencv/opencv/blob/master/data/haarcascades/haarcascade_frontalface_default.xml
	faceClassifier := "haarcascade_frontalface_default.xml"
	if _, err := os.Stat(faceClassifier); os.IsNotExist(err) {
		log.Fatalf("Error: Face cascade XML file not found at %s. Please download it.", faceClassifier)
	}

	detector, err := NewSimpleFaceDetector(faceClassifier)
	if err != nil {
		log.Fatalf("Error creating face detector: %v", err)
	}

	tracker, err := NewKCFTracker()
	if err != nil {
		log.Fatalf("Error creating KCF tracker: %v", err)
	}

	pipeline, err := NewRealtimeCVPipeline(0, detector, tracker)
	if err != nil {
		log.Fatalf("Error creating CV pipeline: %v", err)
	}

	pipeline.Run()
}
// Update on 2023-08-17 00:00:00 - 38
