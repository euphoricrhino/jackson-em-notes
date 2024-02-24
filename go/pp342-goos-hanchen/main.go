package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/draw"
	"math"
	"math/cmplx"

	fieldrenderer "github.com/euphoricrhino/jackson-em-notes/go/pkg/field-renderer"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	// We assume wavelength lambda=1.
	k = 2.0 * math.Pi

	// Gaussian cutoff at e^{-cutoff}.
	cutoff = 15.0

	kappaSamples = 1000
)

type waveParams struct {
	amplitude complex128
	kx        complex128
	kz        complex128
}

var (
	heatmap = flag.String("heatmap", "", "heatmap file")
	output  = flag.String("output", "", "output file")
	gamma   = flag.Float64("gamma", 1.0, "gamma correction")
	width   = flag.Int("width", 800, "output width")
	height  = flag.Int("height", 800, "output height")

	beta           = flag.Float64("beam-width-in-lambdas", 0.0, "transverse distribution of electric field is Gaussian ~ exp(-x^2/(beta*lambda)^2)")
	refrIdx        = flag.Float64("refr-idx", 0.0, "relative refraction index n'/n")
	perpPol        = flag.Bool("perp-pol", true, "whether to consider perpendicular or parallel polarization")
	widthInLambdas = flag.Float64("width-in-lambdas", 0, "how many lambdas does the image's width represent")
)

func main() {
	flag.Parse()

	kappaLimit := 2.0 * math.Sqrt(cutoff) / *beta
	dkappa := 2.0 * kappaLimit / float64(kappaSamples)

	const frames = 180
	deltaAng := math.Pi / 2.0 / frames

	imgWidth, imgHeight := float64(*width), float64(*height)
	centerPixelX, centerPixelY := imgWidth/2.0, imgHeight/3.0
	for f := 0; f <= frames; f++ {
		incAng := float64(f) * deltaAng
		// Decompose incident beam into many plane waves, each with slightly different wave vector and amplitude, by Fourier transform.
		incWp := constructIncidentWaveParams(kappaLimit, dkappa, *beta, incAng)
		reflWp, transWp := reflectAndTransmit(incWp, *perpPol, complex(*refrIdx, 0))

		pixelsPerWavelength := imgWidth / *widthInLambdas

		field := func(screenX, screenY int) float64 {
			pixelX := float64(screenX) - centerPixelX
			pixelY := centerPixelY - float64(screenY)
			x, z := complex(pixelX/pixelsPerWavelength, 0), complex(pixelY/pixelsPerWavelength, 0)

			amp := complex(0, 0)
			if real(z) > 0.0 {
				for _, w := range transWp {
					amp += w.amplitude * cmplx.Exp((w.kx*x+w.kz*z)*1i)
				}
			} else {
				for _, w := range incWp {
					amp += w.amplitude * cmplx.Exp((w.kx*x+w.kz*z)*1i)
				}
				for _, w := range reflWp {
					amp += w.amplitude * cmplx.Exp((w.kx*x+w.kz*z)*1i)
				}
			}
			amp *= complex(dkappa, 0)
			abs := cmplx.Abs(amp)
			return abs * abs
		}

		postEdit := func(img draw.Image) {
			gc := draw2dimg.NewGraphicContext(img)
			gc.SetLineWidth(1)

			// Draw incident optical center line.
			drawRay(gc, imgWidth, imgHeight, centerPixelX, centerPixelY, math.Pi/2+incAng, color.RGBA{0, 0, 0xff, 0xff})

			// Draw unshifted reflected optical center line.
			drawRay(gc, imgWidth, imgHeight, centerPixelX, centerPixelY, math.Pi/2-incAng, color.RGBA{0, 0, 0xff, 0xff})

			cosi, sini := math.Cos(incAng), math.Sin(incAng)

			draw2d.SetFontFolder("/Users/xni/Library/Fonts")
			draw2d.SetFontNamer(func(_ draw2d.FontData) string { return "MonoLisaVariableNormal.ttf" })

			text := fmt.Sprintf("n'/n=%.02f", *refrIdx)
			if *refrIdx < 1.0 {
				text += fmt.Sprintf(", i0=%.02f°", math.Asin(*refrIdx)*180.0/math.Pi)
			}
			text += fmt.Sprintf(", i=%.1f°", incAng*180.0/math.Pi)
			textColor := color.RGBA{0, 0xcc, 0xcc, 0xff}
			gc.SetFillColor(textColor)
			gc.SetStrokeColor(textColor)
			gc.SetDPI(288)
			gc.SetFontSize(3.5)
			gc.FillStringAt(text, 20.0, 20.0)
			text = fmt.Sprintf("image-W=%.1fλ, beam-W=%.1fλ", *widthInLambdas, *beta)
			gc.FillStringAt(text, 20.0, 40.0)
			sini2 := sini * sini
			n2 := *refrIdx * *refrIdx
			if sini > *refrIdx {
				ghs := sini / math.Pi / math.Sqrt(sini2-n2)
				symbol := ""
				if *perpPol {
					symbol = "perp"
				} else {
					symbol = "para"
					ghs *= n2 / (sini2 - (1.0-sini2)*n2)
				}
				text = fmt.Sprintf("theoretical D(%v)=%.02fλ", symbol, ghs)
				gc.FillStringAt(text, 20.0, 60.0)
				// Find the center optical line of reflected field.
				// Only do it when incident angle is not close to right angle.
				if incAng < 85.0/180.0*math.Pi {
					// Measure the field in the range of 2 beam widths and find the peak position.
					lo, hi := -2.0*(*beta)/cosi, 2.0*(*beta)/cosi
					// We want to keep the error or measurement within 0.01 lambdas.
					delta := 0.01
					x0 := 0.0
					maxv := 0.0
					for x := lo; x <= hi; x += delta {
						amp := complex(0, 0)
						for _, w := range reflWp {
							amp += w.amplitude * cmplx.Exp(w.kx*complex(x, 0)*1i)
						}
						abs := cmplx.Abs(amp)
						abs2 := abs * abs
						if maxv < abs2 {
							maxv = abs2
							x0 = x
						}
					}
					text = fmt.Sprintf("measured D=%.02fλ", x0*cosi)
					gc.FillStringAt(text, 20.0, 80.0)
					// Draw the measured optical center line for reflected field.
					drawRay(gc, imgWidth, imgHeight, centerPixelX+x0*pixelsPerWavelength, centerPixelY, math.Pi/2-incAng, color.RGBA{0xff, 0, 0, 0xff})
				}
			}

			// Draw the boundary interface.
			gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			gc.MoveTo(0.0, centerPixelY)
			gc.LineTo(imgWidth, centerPixelY)
			gc.Stroke()
		}

		if err := fieldrenderer.Run(fieldrenderer.Options{
			HeatMapFile: *heatmap,
			OutputFile:  fmt.Sprintf("%s-%04d.png", *output, f),
			Gamma:       *gamma,
			Width:       *width,
			Height:      *height,
			Field:       field,
			PostEdit:    postEdit,
		}); err != nil {
			panic(err)
		}
		fmt.Printf("frame %04d done\n", f)
	}
}

func constructIncidentWaveParams(kappaLimit, dkappa, beta, incAng float64) []waveParams {
	wp := make([]waveParams, kappaSamples+1)
	cosi, sini := math.Cos(incAng), math.Sin(incAng)
	for s := 0; s <= kappaSamples; s++ {
		kappa := -kappaLimit + float64(s)*dkappa
		v := kappa * beta / 2
		wp[s].amplitude = complex(beta/(2.0*math.Sqrt(math.Pi))*math.Exp(-v*v), 0)
		wp[s].kx = complex(k*sini-kappa*cosi, 0)
		wp[s].kz = complex(k*cosi+kappa*sini, 0)
	}
	return wp
}

func reflectAndTransmit(incWp []waveParams, perp bool, refrIdx complex128) ([]waveParams, []waveParams) {
	reflWp := make([]waveParams, len(incWp))
	transWp := make([]waveParams, len(incWp))
	sini02 := refrIdx * refrIdx
	for s, w := range incWp {
		// Incident wave vectors are always real.
		sini := math.Sin(math.Atan(real(w.kx) / real(w.kz)))
		sini2 := complex(sini*sini, 0)
		cosi := cmplx.Sqrt(1.0 - sini2)
		v := cmplx.Sqrt(sini02 - sini2)
		if perp {
			// See Jackson eq (7.39).
			reflWp[s].amplitude = w.amplitude * (cosi - v) / (cosi + v)
			transWp[s].amplitude = w.amplitude * 2.0 * cosi / (cosi + v)
		} else {
			// See Jackson eq (7.41).
			reflWp[s].amplitude = w.amplitude * (sini02*cosi - v) / (sini02*cosi + v)
			transWp[s].amplitude = w.amplitude * 2.0 * refrIdx * cosi / (sini02*cosi + v)
		}
		reflWp[s].kx, reflWp[s].kz = w.kx, -w.kz
		transWp[s].kx, transWp[s].kz = w.kx, cmplx.Sqrt(w.kx*w.kx+w.kz*w.kz)*v
	}
	return reflWp, transWp
}

func drawRay(gc draw2d.GraphicContext, width, height, fromX, fromY, angle float64, color color.RGBA) {
	gc.SetStrokeColor(color)
	gc.MoveTo(fromX, fromY)
	// Use a ray length that's guaranteed to reach beyond the boundary.
	dist := math.Sqrt(width*width + height*height)
	gc.LineTo(fromX+dist*math.Cos(angle), fromY+dist*math.Sin(angle))
	gc.Stroke()
}
