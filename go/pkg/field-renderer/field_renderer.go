package fieldrenderer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/euphoricrhino/jackson-em-notes/go/pkg/heatmap"
)

// Options represents options to run the field renderer.
type Options struct {
	HeatMapFile string
	OutputFile  string
	// Gamma correction to be applied to heatmap.
	Gamma  float64
	Width  int
	Height int
	// Field function for pixel (x,y) ranging from 0 to (Width|Height)-1. Return math.NaN to indicate divergence.
	Field func(x, y int) float64
	// Function to edit the generated image after all the field pixels have rendered.
	PostEdit func(img draw.Image)
}

// Run runs the field renderer with the given options.
func Run(opts Options) error {
	hm, err := heatmap.Load(opts.HeatMapFile, opts.Gamma)
	if err != nil {
		return err
	}

	data := make([]float64, opts.Width*opts.Height)
	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workers)
	cnt := int32(0)
	for w := 0; w < workers; w++ {
		go func(i int) {
			defer wg.Done()
			for x := 0; x < opts.Width; x++ {
				if x%workers != i {
					continue
				}
				for y := 0; y < opts.Height; y++ {
					data[y*opts.Width+x] = opts.Field(x, y)
					atomic.AddInt32(&cnt, 1)
				}
			}
		}(w)
	}
	// Progress counter.
	counterDone := make(chan struct{})
	go func() {
		erase := strings.Repeat(" ", 80)
		nextMark := 1.0
		for {
			doneCnt := int(atomic.LoadInt32(&cnt))
			if doneCnt == opts.Width*opts.Height {
				fmt.Printf("\r%v\rrendering complete\n", erase)
				close(counterDone)
				return
			}
			progress := float64(doneCnt) / float64(opts.Width*opts.Height) * 100.0
			if progress >= nextMark {
				fmt.Printf("\r%v\rrendering... %.2f%% done", erase, progress)
				nextMark = math.Ceil(progress)
			}
			runtime.Gosched()
		}
	}()
	wg.Wait()
	<-counterDone

	// Normalize the data.
	max, min := math.NaN(), math.NaN()
	for i := 0; i < len(data); i++ {
		v := data[i]
		if !math.IsNaN(v) {
			if math.IsNaN(max) || max < v {
				max = v
			}
			if math.IsNaN(min) || min > v {
				min = v
			}
		}
	}
	if !math.IsNaN(max) {
		// Not all pixels are NaN.
		spread := max - min
		if spread > 0.0 {
			// Normalize into range [0,1].
			for i := range data {
				if math.IsNaN(data[i]) {
					data[i] = 0.0
				} else {
					data[i] = (data[i] - min) / spread
				}
			}
		} else {
			// Corner case: ~onstant field.
			for i := range data {
				data[i] = 0.5
			}
		}
	} else {
		// All NaNs.
		for i := range data {
			data[i] = 0.0
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	for x := 0; x < opts.Width; x++ {
		for y := 0; y < opts.Height; y++ {
			pixel := data[y*opts.Width+x]
			pos := int(pixel * float64(len(hm)-1))
			r, g, b, a := hm[pos].RGBA()
			img.SetRGBA64(x, y, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
		}
	}
	if opts.PostEdit != nil {
		opts.PostEdit(img)
	}
	out, err := os.Create(opts.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%v': %v", opts.OutputFile, err)
	}
	defer out.Close()
	if err := png.Encode(out, img); err != nil {
		return fmt.Errorf("failed to encode to PNG: %v", err)
	}
	return nil
}
