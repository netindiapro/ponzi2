package ponzi

import (
	"image"
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type chartPrices struct {
	stock      *modelStock
	stickLines *vao
	stickRects *vao
}

func createChartPrices(stock *modelStock) *chartPrices {
	return &chartPrices{
		stock: stock,
	}
}

func (p *chartPrices) update() {
	if p.stock.dailySessions == nil {
		return
	}

	if p.stickLines != nil {
		return
	}

	// Find the min and max prices for the y-axis.
	min := float32(math.MaxFloat32)
	max := float32(0)
	for _, s := range p.stock.dailySessions {
		if s.low < min {
			min = s.low
		}
		if s.high > max {
			max = s.high
		}
	}

	// Calculate vertices and indices for the candlesticks.
	var vertices []float32
	var colors []float32
	var lineIndices []uint16
	var triangleIndices []uint16

	stickWidth := 2.0 / float32(len(p.stock.dailySessions)) // (-1 to 1) on X-axis
	leftX := -1.0 + stickWidth*0.1
	midX := -1.0 + stickWidth*0.5
	rightX := -1.0 + stickWidth*0.9

	calcY := func(value float32) float32 {
		return 2*(value-min)/(max-min) - 1
	}

	for i, s := range p.stock.dailySessions {
		// Figure out Y coordinates of the key levels.
		lowY, highY, openY, closeY := calcY(s.low), calcY(s.high), calcY(s.open), calcY(s.close)

		// Figure out the top and bottom of the candlestick.
		topY, botY := openY, closeY
		if openY < closeY {
			topY, botY = closeY, openY
		}

		// Add the vertices needed to create the candlestick.
		vertices = append(vertices,
			midX, highY, // 0
			midX, topY, // 1
			midX, lowY, // 2
			midX, botY, // 3
			leftX, topY, // 4 - Upper left of box
			rightX, topY, // 5 - Upper right of box
			leftX, botY, // 6 - Bottom left of box
			rightX, botY, // 7 - Bottom right of box
		)

		// Add the colors corresponding to the vertices.
		var c [3]float32
		switch {
		case s.close > s.open:
			c = green
		case s.close < s.open:
			c = red
		default:
			c = yellow
		}

		colors = append(colors,
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
			c[0], c[1], c[2],
		)

		// idx is function to refer to the vertices above.
		idx := func(j uint16) uint16 {
			return uint16(i)*8 + j
		}

		// Add the vertex indices to render the candlestick.
		lineIndices = append(lineIndices,
			// Top and bottom lines around the box.
			idx(0), idx(1),
			idx(2), idx(3),
		)

		if s.close > s.open {
			// Use lines for open candlestick on higher closes.
			lineIndices = append(lineIndices,
				idx(4), idx(5),
				idx(6), idx(7),
				idx(4), idx(6),
				idx(5), idx(7),
			)
		} else {
			// Use triangles for filled candlestick on lower closes.
			triangleIndices = append(triangleIndices,
				idx(4), idx(6), idx(5),
				idx(5), idx(6), idx(7),
			)
		}

		// Move the X coordinates one stick over.
		leftX += stickWidth
		midX += stickWidth
		rightX += stickWidth
	}

	p.stickLines = createVAO(gl.LINES, vertices, colors, lineIndices)
	p.stickRects = createVAO(gl.TRIANGLES, vertices, colors, triangleIndices)
}

func (p *chartPrices) render(r image.Rectangle) {
	if p == nil {
		return
	}
	setModelMatrixRectangle(r)
	p.stickLines.render()
	p.stickRects.render()
}

func (p *chartPrices) close() {
	if p == nil {
		return
	}
	p.stickLines.close()
	p.stickRects.close()
}
