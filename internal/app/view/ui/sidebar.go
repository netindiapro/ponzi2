package ui

import (
	"image"

	"github.com/btmura/ponzi2/internal/app/view/rect"

	"github.com/btmura/ponzi2/internal/app/view"
	"github.com/btmura/ponzi2/internal/app/view/chart"
)

var (
	chartThumbSize         = image.Pt(155, 105)
	chartThumbRenderOffset = image.Pt(0, viewPadding+chartThumbSize.Y)
	sidebarScrollAmount    = chartThumbRenderOffset
)

type sidebar struct {
	// slots are slots which can have a thumbnail or be a drop site.
	slots []*sidebarSlot

	// draggedSlot if not nil is the slot with the thumbnail being dragged.
	draggedSlot *sidebarSlot

	// sidebarScrollOffset stores the Y offset accumulated from scroll events
	// that should be used to calculate the sidebar's bounds.
	sidebarScrollOffset image.Point

	// bounds is the rectangle to draw within.
	bounds image.Rectangle
}

func (s *sidebar) AddChartThumb(th *chart.Thumb) {
	s.slots = append(s.slots, newSidebarSlot(th))
}

func (s *sidebar) RemoveChartThumb(th *chart.Thumb) {
	for _, slot := range s.slots {
		if slot.thumbnail == th {
			slot.FadeOut()
			break
		}
	}
}

func (s *sidebar) SetBounds(bounds image.Rectangle) {
	s.bounds = bounds
}

func (s *sidebar) ProcessInput(input *view.Input) {
	// TODO(btmura): Change these to be constants somewhere.
	thumbSize := image.Pt(s.bounds.Dx(), chartThumbSize.Y)

	slotBounds := image.Rect(
		s.bounds.Min.X, s.bounds.Max.Y-viewPadding-thumbSize.Y,
		s.bounds.Max.X, s.bounds.Max.Y-viewPadding,
	)

	if !input.MouseLeftButtonDragging {
		s.draggedSlot = nil
	}

	draggedSlotIndex := -1

	// Go down the sidebar and assign bounds to each slot and identify the dragged slot.
	for i, slot := range s.slots {
		slot.SetBounds(slotBounds)
		slot.SetThumbnailBounds(slotBounds)

		if s.draggedSlot == nil {
			dragging := input.MouseLeftButtonDragging &&
				input.MouseLeftButtonDraggingStartedPos.In(slotBounds)
			if dragging {
				s.draggedSlot = slot
			}
		}

		if s.draggedSlot == slot {
			draggedSlotIndex = i
		}

		slotBounds = slotBounds.Sub(chartThumbRenderOffset)
	}

	if s.draggedSlot != nil {
		// Float the dragged slot's thumbnail to be under the mouse cursor.
		s.draggedSlot.SetThumbnailBounds(rect.FromCenterAndSize(input.MousePos, thumbSize))

		// Move the dragged slot up or down.
		swapIndex := draggedSlotIndex - 1
		if c := rect.Center(s.draggedSlot.Bounds()); input.MousePos.Y < c.Y {
			swapIndex = draggedSlotIndex + 1
		}

		// Swap with adjacent slot if the mouse has moved over it.
		if swapIndex >= 0 && swapIndex < len(s.slots) && input.MousePos.In(s.slots[swapIndex].Bounds()) {
			i, j := draggedSlotIndex, swapIndex
			s.slots[i], s.slots[j] = s.slots[j], s.slots[i]
		}
	}

	// Forward the input such as clicks to the adjusted sidebar now.
	for _, slot := range s.slots {
		slot.ProcessInput(input)
	}
}

// Update moves the animation one step forward.
func (s *sidebar) Update() (dirty bool) {
	for i := 0; i < len(s.slots); i++ {
		slot := s.slots[i]
		if slot.Update() {
			dirty = true
		}
		if slot.DoneFadingOut() {
			s.slots = append(s.slots[:i], s.slots[i+1:]...)
			slot.Close()
			i--
		}
	}
	return dirty
}

// Render renders a frame.
func (s *sidebar) Render(fudge float32) {
	// Draw the non-dragged thumbnails first, so they appear under the dragged thumbnail.
	for _, slot := range s.slots {
		if slot != s.draggedSlot {
			slot.Render(fudge)
		}
	}

	if s.draggedSlot != nil {
		s.draggedSlot.Render(fudge)
	}
}

type sidebarSlot struct {
	// bounds is the rectangle to draw the slot within. Can be empty for collapsed slots.
	bounds image.Rectangle

	// thumbnail is an optional stock thumbnail. Could have bounds outside the slot if dragged.
	thumbnail *chart.Thumb

	// fader fades out the slot.
	fader *view.Fader
}

func newSidebarSlot(th *chart.Thumb) *sidebarSlot {
	return &sidebarSlot{
		thumbnail: th,
		fader:     view.NewFader(1 * fps),
	}
}

func (s *sidebarSlot) FadeOut() {
	s.fader.FadeOut()
}

func (s *sidebarSlot) DoneFadingOut() bool {
	return s.fader.DoneFadingOut()
}

func (s *sidebarSlot) SetBounds(bounds image.Rectangle) {
	s.bounds = bounds
}

func (s *sidebarSlot) Bounds() image.Rectangle {
	return s.bounds
}

func (s *sidebarSlot) SetThumbnailBounds(bounds image.Rectangle) {
	if s.thumbnail != nil {
		s.thumbnail.SetBounds(bounds)
	}
}

func (s *sidebarSlot) ProcessInput(input *view.Input) {
	if s.thumbnail != nil {
		s.thumbnail.ProcessInput(input)
	}
}

func (s *sidebarSlot) Update() (dirty bool) {
	if s.thumbnail != nil && s.thumbnail.Update() {
		dirty = true
	}
	if s.fader.Update() {
		dirty = true
	}
	return dirty
}

func (s *sidebarSlot) Render(fudge float32) {
	s.fader.Render(func(fudge float32) {
		if s.thumbnail != nil {
			s.thumbnail.Render(fudge)
		}
	}, fudge)
}

func (s *sidebarSlot) Close() {
	if s.thumbnail != nil {
		s.thumbnail.Close()
		s.thumbnail = nil
	}
}
