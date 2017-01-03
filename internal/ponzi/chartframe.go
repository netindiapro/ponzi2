package ponzi

import (
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type chartFrame struct {
	propText       *dynamicText
	borderVAO      uint32
	borderCount    int32
	separatorVAO   uint32
	separatorCount int32
}

func createChartFrame(propText *dynamicText) *chartFrame {
	borderVAO, borderCount := createChartBorderVAO()
	separatorVAO, separatorCount := createChartSeparatorVAO()
	return &chartFrame{
		propText:       propText,
		borderVAO:      borderVAO,
		borderCount:    borderCount,
		separatorVAO:   separatorVAO,
		separatorCount: separatorCount,
	}
}

func createChartBorderVAO() (uint32, int32) {
	vertices := []float32{
		-1, 1,
		-1, -1,
		1, -1,
		1, 1,
	}

	colors := []float32{
		0, 0.25, 0.5,
		0, 0.25, 0.5,
		0, 0.25, 0.5,
		0, 0.25, 0.5,
	}

	indices := []uint16{
		0, 1,
		1, 2,
		2, 3,
		3, 0,
	}

	vbo := createArrayBuffer(vertices)
	cbo := createArrayBuffer(colors)
	ibo := createElementArrayBuffer(indices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(positionLocation)
		gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, cbo)
		gl.EnableVertexAttribArray(colorLocation)
		gl.VertexAttribPointer(colorLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	}
	gl.BindVertexArray(0)

	return vao, int32(len(indices))
}

func createChartSeparatorVAO() (uint32, int32) {
	vertices := []float32{
		-1, 0,
		1, 0,
	}

	colors := []float32{
		0, 0.25, 0.5,
		0, 0.25, 0.5,
	}

	indices := []uint16{
		0, 1,
	}

	vbo := createArrayBuffer(vertices)
	cbo := createArrayBuffer(colors)
	ibo := createElementArrayBuffer(indices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(positionLocation)
		gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, cbo)
		gl.EnableVertexAttribArray(colorLocation)
		gl.VertexAttribPointer(colorLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	}
	gl.BindVertexArray(0)

	return vao, int32(len(indices))
}

func (f *chartFrame) render(stock *modelStock, r image.Rectangle) []image.Rectangle {
	// Start rendering from the top left. Track position with point.
	c := image.Pt(r.Min.X, r.Max.Y)

	//
	// Render the frame around the chart.
	//

	setModelMatrixRectangle(r)
	gl.Uniform1f(colorMixAmountLocation, 1)

	gl.BindVertexArray(f.borderVAO)
	gl.DrawElements(gl.LINES, f.borderCount, gl.UNSIGNED_SHORT, gl.Ptr(nil))
	gl.BindVertexArray(0)

	//
	// Render the symbol and its quote.
	//

	const padding = 10
	s := f.propText.measure(stock.symbol)
	c.Y -= padding + s.Y
	{
		c := c
		c.X += padding
		c = c.Add(f.propText.render(stock.symbol, c))
		c = c.Add(f.propText.render(formatQuote(stock.quote), c))
	}

	//
	// Render the separator below the symbol and quote.
	//

	r.Max.Y = c.Y
	gl.Uniform1f(colorMixAmountLocation, 1)

	rects := sliceRectangle(r, 0.13, 0.13, 0.13, 0.6)
	for _, r := range rects {
		setModelMatrixRectangle(image.Rect(r.Min.X, r.Max.Y, r.Max.X, r.Max.Y))
		gl.BindVertexArray(f.separatorVAO)
		gl.DrawElements(gl.LINES, f.separatorCount, gl.UNSIGNED_SHORT, gl.Ptr(nil))
		gl.BindVertexArray(0)
	}

	return rects
}

func (f *chartFrame) close() {
	gl.DeleteVertexArrays(1, &f.borderVAO)
	gl.DeleteVertexArrays(1, &f.separatorVAO)
}

// createChartBackgroundVAO creates a 2-triangle VAO for a square centered around (0, 0).
func createChartBackgroundVAO(ulColor, urColor, blColor, brColor [3]float32) (uint32, int32) {
	vertices := []float32{
		-1, 1, // UL - 0
		1, 1, // UR - 1
		-1, -1, // BL - 2
		1, -1, // BR - 3
	}

	colors := []float32{
		ulColor[0], ulColor[1], ulColor[2], // UL
		urColor[0], urColor[1], urColor[2], // UR
		blColor[0], blColor[1], blColor[2], // BL
		brColor[0], brColor[1], brColor[2], // BR
	}

	indices := []uint16{
		0, 2, 1,
		1, 2, 3,
	}

	vbo := createArrayBuffer(vertices)
	cbo := createArrayBuffer(colors)
	ibo := createElementArrayBuffer(indices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	{
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(positionLocation)
		gl.VertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, cbo)
		gl.EnableVertexAttribArray(colorLocation)
		gl.VertexAttribPointer(colorLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	}
	gl.BindVertexArray(0)

	return vao, int32(len(indices))
}
