package example_test

import (
	"context"
	"log"
	"os"
	"time"

	"pipelined.dev/audio"
	"pipelined.dev/mp3"
	"pipelined.dev/pipe"
)

// This example demonstrates how to cut a clip from .mp3 file
// and save the result to new .mp3 file.
func Example_6() {
	// open source wav file.
	mp3File, err := os.Open("_testdata/sample.mp3")
	if err != nil {
		log.Fatalf("failed to open mp3 file: %v", err)
	}
	defer mp3File.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out6.mp3")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// asset to keep mp3 signal.
	a := &audio.Asset{}
	// build pipe with a single line.
	p, err := pipe.New(
		&pipe.Line{
			// mp3 pump.
			Pump: &mp3.Pump{
				Reader: mp3File,
			},
			// asset sink.
			Sinks: pipe.Sinks(a),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind asset pipeline: %v", err)
	}
	defer p.Close()
	// run the line.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute asset pipeline: %v", err)
	}

	// get the sample rate of the asset.
	// it will be needed to convert duration.
	sampleRate := a.SampleRate()
	// cut the clip that starts at 1 second and lasts 2.5 seconds.
	c := a.Clip(
		sampleRate.SamplesIn(time.Millisecond*1000),
		sampleRate.SamplesIn(time.Millisecond*2500),
	)
	// build pipe with a single line.
	p, err = pipe.New(
		&pipe.Line{
			// mp3 pump.
			Pump: c,
			// asset sink.
			Sinks: pipe.Sinks(
				// mp3 sink
				&mp3.Sink{
					Writer:      outputFile,
					BitRateMode: mp3.CBR(320),
					ChannelMode: mp3.JointStereo,
				},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind output pipeline: %v", err)
	}
	defer p.Close()
	// run the line.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute output pipeline: %v", err)
	}
}
