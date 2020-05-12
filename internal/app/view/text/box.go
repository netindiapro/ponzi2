package text

import (
	"image"

	"github.com/btmura/ponzi2/internal/app/gfx"
	"github.com/btmura/ponzi2/internal/app/view/color"
)

// Box draws text that is horizontally and vertically centered.
type Box struct {
	// textRenderer renders the text.
	textRenderer *gfx.TextRenderer

	// text is the text to render.
	text string

	// color is the color to render the text in.
	color color.RGBA

	// padding is the padding around the text.
	padding int

	// bounds is the rectangle with global coords that should be drawn within.
	bounds image.Rectangle

	// adjustedBounds is the adjusted bounds shrink-wrapped around the text.
	adjustedBounds image.Rectangle

	// dirty is true if the text or bounds has changed and requires an update.
	dirty bool
}

// Option is an option to pass to New.
type Option func(c *Box)

// Color returns an option to set the text color.
func Color(color color.RGBA) Option {
	return func(b *Box) {
		b.color = color
	}
}

// Padding returns an option to set the padding.
func Padding(padding int) Option {
	return func(b *Box) {
		b.padding = padding
	}
}

// NewBox creates a new Box.
func NewBox(textRenderer *gfx.TextRenderer, text string, opts ...Option) *Box {
	b := &Box{
		textRenderer: textRenderer,
		text:         text,
		color:        color.RGBA{1, 1, 1, 1},
		dirty:        true,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

// Text returns the text that will be shown in the box.
func (b *Box) Text() string {
	return b.text
}

// SetText sets the text to show in the box.
func (b *Box) SetText(text string) {
	if b.text == text {
		return
	}
	b.text = text
	b.dirty = true
}

// SetBounds sets the bounds with global coordinates to draw within.
func (b *Box) SetBounds(bounds image.Rectangle) {
	if b.bounds == bounds {
		return
	}
	b.bounds = bounds
	b.dirty = true
}

// Update updates the state by one frame and returns true if another update is needed for animation.
func (b *Box) Update() (dirty bool) {
	if !b.dirty {
		return false
	}

	textSize := b.textRenderer.Measure(b.text)
	if textSize.X > b.bounds.Dx()-b.padding {
		textSize.X = b.bounds.Dx() - b.padding
	}

	textPt := image.Point{
		X: b.bounds.Min.X + b.bounds.Dx()/2 - textSize.X/2,
		Y: b.bounds.Min.Y + b.bounds.Dy()/2 - textSize.Y/2,
	}
	b.adjustedBounds = image.Rectangle{
		Min: textPt,
		Max: textPt.Add(textSize),
	}

	b.dirty = false

	return false
}

// Render renders the current state to the screen.
func (b *Box) Render(float32) {
	if b.text == "" {
		return
	}

	if b.bounds.Empty() {
		return
	}

	b.textRenderer.Render(b.text, b.adjustedBounds.Min, gfx.TextColor(b.color), gfx.TextRenderMaxWidth(b.adjustedBounds.Dx()))
}
