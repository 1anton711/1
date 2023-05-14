package main

// Define a new struct for the beep sound.
type beepSound struct{}

// Implement the Streamer interface for the beep sound.
func (bs *beepSound) Streamer(samples [][2]float64) (n int, ok bool) {
	return 0, false
}

// Create a new instance of the beep sound.
var sound = &beepSound{}
