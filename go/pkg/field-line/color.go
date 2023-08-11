package fieldline

import (
	"math"
	"math/rand"
)

// The following formula is from wikipedia: https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB_alternative
// With S=V=1.
func min(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}

func hueAsRGB(h float64) [3]float64 {
	rr := math.Mod(5.0+h*6.0, 6.0)
	gg := math.Mod(3.0+h*6.0, 6.0)
	bb := math.Mod(1.0+h*6.0, 6.0)

	r := 1.0 - math.Max(min(rr, 4.0-rr, 1.0), 0.0)
	g := 1.0 - math.Max(min(gg, 4.0-gg, 1.0), 0.0)
	b := 1.0 - math.Max(min(bb, 4.0-bb, 1.0), 0.0)

	return [3]float64{r, g, b}
}

// Pre-generated full-hue colors in RGB.
var hueRGB [][3]float64

func init() {
	for h := 0.0; h < 1; h += 1.0 / 1024.0 {
		hueRGB = append(hueRGB, hueAsRGB(h))
	}
}

func RandColor() [3]float64 {
	return hueRGB[rand.Intn(len(hueRGB))]
}
