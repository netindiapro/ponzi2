package view

import (
	"fmt"
	"image"
	"math"
	"time"

	"gitlab.com/btmura/ponzi2/internal/app/model"
)

// longTime is a time that takes the most display width for measuring purposes.
var longTime = time.Date(2019, time.December, 31, 23, 59, 0, 0, time.UTC)

// chartTimeLabels renders the time labels for a single stock.
type chartTimeLabels struct {
	// renderable is whether the ChartTimeLabels can be rendered.
	renderable bool

	// MaxLabelSize is the maximum label size useful for rendering measurements.
	MaxLabelSize image.Point

	// labels bundle rendering measurements for time labels.
	labels []chartTimeLabel

	// dates are session dates shown for the cursor.
	dates []time.Time
}

func newChartTimeLabels() *chartTimeLabels {
	return &chartTimeLabels{}
}

func (ch *chartTimeLabels) SetData(r model.Range, ts *model.TradingSessionSeries) error {
	// Reset everything.
	ch.Close()

	// Bail out if there is no data yet.
	if ts == nil {
		return nil
	}

	txt, err := chartTimeLabelText(r, longTime)
	if err != nil {
		return err
	}
	ch.MaxLabelSize = chartAxisLabelTextRenderer.Measure(txt)

	labels, err := makeChartTimeLabels(r, ts.TradingSessions)
	if err != nil {
		return err
	}
	ch.labels = labels

	ch.dates = nil
	for _, s := range ts.TradingSessions {
		ch.dates = append(ch.dates, s.Date)
	}

	ch.renderable = true

	return nil
}

func (ch *chartTimeLabels) Render(r image.Rectangle) {
	if !ch.renderable {
		return
	}

	for _, l := range ch.labels {
		tp := image.Point{
			X: r.Min.X + int(float32(r.Dx())*l.percent) - l.size.X/2,
			Y: r.Min.Y + r.Dy()/2 - l.size.Y/2,
		}
		chartAxisLabelTextRenderer.Render(l.text, tp, white)
	}
}

func (ch *chartTimeLabels) RenderCursorLabels(mainRect, labelRect image.Rectangle, mousePos image.Point) {
	if !ch.renderable {
		return
	}

	if mousePos.X < mainRect.Min.X || mousePos.X > mainRect.Max.X {
		return
	}

	l := chartTimeLabel{
		percent: float32(mousePos.X-mainRect.Min.X) / float32(mainRect.Dx()),
	}

	i := int(math.Floor(float64(len(ch.dates))*float64(l.percent) + 0.5))
	if i >= len(ch.dates) {
		i = len(ch.dates) - 1
	}
	l.text = ch.dates[i].Format("1/2/06")
	l.size = chartAxisLabelTextRenderer.Measure(l.text)

	tp := image.Point{
		X: mousePos.X - l.size.X/2,
		Y: labelRect.Min.Y + labelRect.Dy()/2 - l.size.Y/2,
	}

	renderBubble(tp, l.size, chartAxisLabelBubbleSpec)
	chartAxisLabelTextRenderer.Render(l.text, tp, white)
}

func (ch *chartTimeLabels) Close() {
	ch.renderable = false
}

type chartTimeLabel struct {
	percent float32
	text    string
	size    image.Point
}

func chartTimeLabelText(r model.Range, t time.Time) (string, error) {
	switch r {
	case model.OneDay:
		return t.Format("3:04"), nil
	case model.TwoYears:
		return t.Format("Jan"), nil
	default:
		return "", fmt.Errorf("bad range: %v", r)
	}
}

func makeChartTimeLabels(r model.Range, ts []*model.TradingSession) ([]chartTimeLabel, error) {
	var ls []chartTimeLabel

	for i := range ts {
		// Skip if we can't check the previous value.
		if i == 0 {
			continue
		}

		// Skip if the values being printed aren't changing.
		switch r {
		case model.OneDay:
			prev := ts[i-1].Date.Hour()
			curr := ts[i].Date.Hour()
			if prev == curr {
				continue
			}

		case model.TwoYears:
			pm := ts[i-1].Date.Month()
			m := ts[i].Date.Month()
			if pm == m {
				continue
			}

		default:
			return nil, fmt.Errorf("bad range: %v", r)
		}

		// Generate the label text and its position.

		txt, err := chartTimeLabelText(r, ts[i].Date)
		if err != nil {
			return nil, err
		}

		ls = append(ls, chartTimeLabel{
			percent: float32(i) / float32(len(ts)),
			text:    txt,
			size:    chartAxisLabelTextRenderer.Measure(txt),
		})
	}

	return ls, nil
}
