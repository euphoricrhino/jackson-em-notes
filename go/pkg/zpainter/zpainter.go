package zpainter

import (
	"image"
	"sort"
	"sync"
)

// Represents a rasterized pixel's color together with its z depth.
type zColor struct {
	// Color premultiplied with rasterizer span alpha.
	r, g, b, a uint32
	// Depth.
	z float64
}

type sortByZ []*zColor

func (zh sortByZ) Len() int      { return len(zh) }
func (zh sortByZ) Swap(i, j int) { zh[i], zh[j] = zh[j], zh[i] }

// We are looking from -z to +z, so a greater z value needs to be painted first.
func (zh sortByZ) Less(i, j int) bool { return zh[i].z > zh[j].z }

// ZPainter enables strokes to be rendered with an order sorted by the corresponding z depth.
// ZPainter.Shard(i) implements the Painter interface from github.com/llgcode/draw2d/draw2dimg and can be used with draw2dimg.NewGraphicContextWithPainter for stroking in parallel.
// ZPainter.Commit() will merge and sort the depth info at each pixel and paint the pixel in the sorted z order.
type ZPainter struct {
	img *image.RGBA

	width  int
	height int

	shards []*zShard
}

func NewZPainter(img *image.RGBA, shardCnt int) *ZPainter {
	b := img.Bounds()
	zp := &ZPainter{
		img:    img,
		width:  b.Max.X - b.Min.X,
		height: b.Max.Y - b.Min.Y,
		shards: make([]*zShard, shardCnt),
	}
	for i := range zp.shards {
		zp.shards[i] = &zShard{
			parent: zp,
			zbuf:   make([][]*zColor, zp.width*zp.height),
		}
	}
	return zp
}

func (zp *ZPainter) Shard(i int) *zShard {
	return zp.shards[i]
}

func (zp *ZPainter) Commit(workers int) {
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(wk int) {
			for i := 0; i < zp.width*zp.height; i++ {
				if i%workers != wk {
					continue
				}
				// Merge and sort z-buffers from all concurrent shards.
				l := 0
				for _, sh := range zp.shards {
					l += len(sh.zbuf[i])
				}
				sorted := make([]*zColor, 0, l)
				for _, sh := range zp.shards {
					sorted = append(sorted, sh.zbuf[i]...)
				}
				sort.Sort(sortByZ(sorted))

				x, y := i%zp.width, i/zp.width
				idx := y*zp.img.Stride + x*4

				// Paint the pixels from far to near, multiplying the alpha at each layer in order.
				for _, zc := range sorted {
					const m = 1<<16 - 1
					dr := uint32(zp.img.Pix[idx+0])
					dg := uint32(zp.img.Pix[idx+1])
					db := uint32(zp.img.Pix[idx+2])
					da := uint32(zp.img.Pix[idx+3])
					a := (m - (zc.a / m)) * 0x101
					zp.img.Pix[idx+0] = uint8((dr*a + zc.r) / m >> 8)
					zp.img.Pix[idx+1] = uint8((dg*a + zc.g) / m >> 8)
					zp.img.Pix[idx+2] = uint8((db*a + zc.b) / m >> 8)
					zp.img.Pix[idx+3] = uint8((da*a + zc.a) / m >> 8)
				}
			}
			wg.Done()
		}(w)
	}

	wg.Wait()
}
