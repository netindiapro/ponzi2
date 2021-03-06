package chart

import (
	"fmt"
	"image"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/btmura/ponzi2/internal/app/gfx"
	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/app/view"
	"github.com/btmura/ponzi2/internal/app/view/rect"
)

const (
	legendBubbleMargin   = 10
	legendBubbleRounding = 10
	legendTablePadding   = 10
	legendFontSize       = 20
)

var (
	legendTextRenderer = gfx.NewTextRenderer(goregular.TTF, legendFontSize)

	// legendGeometricShapeRenderer renderers geometric shapes.
	// https://www.fileformat.info/info/unicode/font/dejavu_sans_mono/blockview.htm?block=geometric_shapes
	legendGeometricShapeRenderer = gfx.NewTextRenderer(_escFSMustByte(false, "/data/DejaVuSans.ttf"), legendFontSize)
)

// legend is a bubble that shows a trading session's stats where the mouse cursor is.
type legend struct {
	// data is the data necessary to render.
	data legendData
	// bounds is the bounds to draw the legend within.
	bounds image.Rectangle
	// mousePos is the current mouse position. Nil for no mouse input.
	mousePos *view.MousePosition

	// table renders the information of a single trading session.
	table legendTable

	// needUpdate is true if the table and tableBubble need updating.
	needUpdate bool
	// renderable is true if there is something to render.
	renderable bool
}

type legendTable struct {
	bubble  *rect.Bubble
	rows    [][3]legendCell
	columns [3]legendColumn
}

func (l *legendTable) SetBounds(bounds image.Rectangle) {
	l.bubble.SetBounds(bounds)
}

func (l *legendTable) Render(fudge float32) {
	l.bubble.Render(fudge)

	lowerLeft := l.bubble.Bounds().Inset(legendTablePadding).Min
	for i := len(l.rows) - 1; i >= 0; i-- {
		{
			row := l.rows[i]
			pt := lowerLeft

			row[0].Render(pt)
			pt.X += l.columns[0].maxWidth + legendTablePadding

			row[1].Render(pt)
			pt.X += l.columns[1].maxWidth + legendTablePadding

			row[2].Render(pt)
		}
		lowerLeft.Y += legendTextRenderer.LineHeight()
	}
}

type legendCell struct {
	renderer *gfx.TextRenderer
	text     string
	color    view.Color
	size     image.Point
}

func (l *legendCell) Render(pt image.Point) {
	l.renderer.Render(l.text, pt, gfx.TextColor(l.color))
}

type legendColumn struct {
	maxWidth int
}

func newLegend() *legend {
	return &legend{}
}

type legendData struct {
	Interval               model.Interval
	TradingSessionSeries   *model.TradingSessionSeries
	MovingAverageSeriesSet []*model.MovingAverageSeries
}

func (l *legend) SetData(data legendData) {
	l.data = data
	l.needUpdate = true
}

func (l *legend) SetBounds(bounds image.Rectangle) {
	if l.bounds == bounds {
		return
	}
	l.bounds = bounds
	l.needUpdate = true
}

func (l *legend) ProcessInput(input *view.Input) {
	if l.mousePos == input.MousePos {
		return
	}
	l.mousePos = input.MousePos
	l.needUpdate = true
}

func (l *legend) Update() (dirty bool) {
	if !l.needUpdate {
		return false
	}

	defer func() { l.needUpdate = false }()

	if l.data.TradingSessionSeries == nil {
		l.renderable = false
		return true
	}

	tss := l.data.TradingSessionSeries.TradingSessions
	if len(tss) == 0 {
		l.renderable = false
		return true
	}

	i, ts := len(tss)-1, tss[len(tss)-1]
	if l.mousePos.WithinX(l.bounds) {
		i, ts = tradingSessionAtX(tss, l.bounds, l.mousePos.X)
	}

	curr := ts
	prev := curr
	if i > 0 {
		prev = tss[i-1]
	}

	formatFloat := func(value float32) string {
		return fmt.Sprintf("%.2f", value)
	}

	formatChange := func(change float32) string {
		return fmt.Sprintf("%+.2f", change)
	}

	formatPercentChange := func(percentChange float32) string {
		return fmt.Sprintf("%+.2f%%", percentChange)
	}

	var empty legendCell

	text := func(text string) legendCell {
		return legendCell{
			renderer: legendTextRenderer,
			text:     text,
			color:    view.White,
			size:     legendTextRenderer.Measure(text),
		}
	}

	symbol := func(text string, color view.Color) legendCell {
		return legendCell{
			renderer: legendGeometricShapeRenderer,
			text:     text,
			color:    color,
			size:     legendGeometricShapeRenderer.Measure(text),
		}
	}

	whiteArrow := func(change float32) legendCell {
		switch {
		case change > 0:
			return symbol("△", view.White)
		case change < 0:
			return symbol("▽", view.White)
		default:
			return empty
		}
	}

	colorArrow := func(change float32) legendCell {
		switch {
		case change > 0:
			return symbol("▲", view.Green)
		case change < 0:
			return symbol("▼", view.Red)
		default:
			return empty
		}
	}

	rows := [][3]legendCell{
		{empty, text(curr.Date.Format("1/2/06")), empty},
		{empty, empty, empty},
		{
			whiteArrow(curr.Open - prev.Open),
			text("Open"),
			text(formatFloat(curr.Open)),
		},
		{
			whiteArrow(curr.High - prev.High),
			text("High"),
			text(formatFloat(curr.High)),
		},
		{
			whiteArrow(curr.Low - prev.Low),
			text("Low"),
			text(formatFloat(curr.Low)),
		},
		{
			whiteArrow(curr.Close - prev.Close),
			text("Close"),
			text(formatFloat(curr.Close)),
		},
		{empty, empty, empty},
		{
			colorArrow(curr.Change),
			text("Change"),
			text(formatChange(curr.Change)),
		},
		{empty, empty, text(formatPercentChange(curr.PercentChange))},
	}

	if len(l.data.MovingAverageSeriesSet) != 0 {
		rows = append(rows, [3]legendCell{empty, empty, empty})
	}

	for _, ma := range l.data.MovingAverageSeriesSet {
		typeLabel := "?"
		switch ma.Type {
		case model.Simple:
			typeLabel = "SMA"
		case model.Exponential:
			typeLabel = "EMA"
		}

		rows = append(rows, [3]legendCell{
			symbol("◼", movingAverageColors[l.data.Interval][ma.Intervals]),
			text(fmt.Sprintf("%s %d", typeLabel, ma.Intervals)),
			text(formatFloat(ma.Values[i].Value)),
		})
	}

	if curr.Volume != 0 {
		dv := curr.Volume - prev.Volume
		rows = append(rows,
			[3]legendCell{empty, empty, empty},
			[3]legendCell{
				whiteArrow(float32(dv)),
				text("Volume"),
				text(volumeText(curr.Volume)),
			},
			[3]legendCell{
				empty,
				empty,
				text(volumeChangeText(dv)),
			},
			[3]legendCell{
				empty,
				empty,
				text(formatPercentChange(curr.VolumePercentChange)),
			},
		)
	}

	columns := [3]legendColumn{}
	for i := range rows {
		for j := range columns {
			if w := rows[i][j].size.X; w > columns[j].maxWidth {
				columns[j].maxWidth = w
			}
		}
	}

	tableBounds := image.Rect(
		0,
		0,
		legendTablePadding+columns[0].maxWidth+
			legendTablePadding+columns[1].maxWidth+
			legendTablePadding+columns[2].maxWidth+legendTablePadding,
		legendTablePadding+len(rows)*legendTextRenderer.LineHeight()+legendTablePadding,
	)

	// Move the table to the lower left.
	bounds := l.bounds.Inset(legendBubbleMargin)
	tableBounds = tableBounds.Add(bounds.Min)

	// Move the table to the right if the mouse is in the bounds.
	if l.mousePos.In(tableBounds) {
		tableBounds = tableBounds.Add(image.Pt(bounds.Dx()-tableBounds.Dx(), 0))
	}

	l.table = legendTable{
		bubble:  rect.NewBubble(legendBubbleRounding),
		rows:    rows,
		columns: columns,
	}
	l.table.SetBounds(tableBounds)
	l.renderable = true

	return true
}

func (l *legend) Render(fudge float32) {
	if !l.renderable {
		return
	}

	l.table.Render(fudge)
}

func (l *legend) Close() {
	l.renderable = false
}
