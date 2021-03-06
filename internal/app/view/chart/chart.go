package chart

import (
	"image"
	"math"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/btmura/ponzi2/internal/app/gfx"
	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/app/view"
	"github.com/btmura/ponzi2/internal/app/view/animation"
	"github.com/btmura/ponzi2/internal/app/view/rect"
	"github.com/btmura/ponzi2/internal/app/view/status"
	"github.com/btmura/ponzi2/internal/app/view/text"
	"github.com/btmura/ponzi2/internal/app/view/vao"
	"github.com/btmura/ponzi2/internal/logger"
)

// Embed resources into the application. Get esc from github.com/mjibson/esc.
//go:generate esc -o bindata.go -pkg chart -include ".*(ply|png|ttf)" -modtime 1337 -private data

const (
	chartRounding       = 10
	chartSectionPadding = 5
	chartTextPadding    = 20
	chartVolumePercent  = 0.25
)

var (
	chartSymbolQuoteTextRenderer = gfx.NewTextRenderer(goregular.TTF, 24)
	chartQuotePrinter            = func(q *model.Quote) string { return status.Join(status.PriceChange(q), status.SourceUpdate(q)) }
)

const axisLabelPadding = 4

var (
	axisLabelBubble       = rect.NewBubble(6)
	axisLabelTextRenderer = gfx.NewTextRenderer(goregular.TTF, 12)
)

var (
	cursorHorizLine = vao.HorizLine(view.LightGray, view.LightGray)
	cursorVertLine  = vao.VertLine(view.LightGray, view.LightGray)
)

var movingAverageColors = map[model.Interval]map[int]view.Color{
	model.Daily: {
		8:   view.Purple,
		21:  view.Green,
		50:  view.Red,
		200: view.White,
	},
	model.Weekly: {
		10: view.Red,
		40: view.White,
	},
}

// PriceStyle is visual style of the chart's prices.
type PriceStyle int

// PriceStyle values.
//go:generate stringer -type=PriceStyle
const (
	PriceStyleUnspecified PriceStyle = iota
	Bar
	Candlestick
)

// ZoomChange specifies whether the user has zoomed in or not.
type ZoomChange int

// ZoomChange values.
//go:generate stringer -type=ZoomChange
const (
	ZoomChangeUnspecified ZoomChange = iota
	ZoomIn
	ZoomOut
)

// Chart shows a stock chart for a single stock.
type Chart struct {
	// frameBubble is the border around the entire chart.
	frameBubble *rect.Bubble

	// header renders the header with the symbol, quote, and buttons.
	header *header

	price         *price
	priceLevel    *priceLevel
	priceCursor   *priceCursor
	priceTimeline *timeline

	movingAverages []*movingAverage

	volume         *volume
	volumeLevel    *volumeLevel
	volumeCursor   *volumeCursor
	volumeTimeline *timeline

	timelineAxis   *timelineAxis
	timelineCursor *timelineCursor

	legend *legend

	// loadingTextBox renders the loading text shown when loading from a fresh state.
	loadingTextBox *text.Box

	// errorTextBox renders the error text.
	errorTextBox *text.Box

	// loading is true when data for the stock is being retrieved.
	loading bool

	// hasStockUpdated is true if the stock has reported data.
	hasStockUpdated bool

	// hasError is true there was a loading issue.
	hasError bool

	// fadeIn fades in the data after it loads.
	fadeIn *animation.Animation

	// showMovingAverages is whether to render the moving averages.
	showMovingAverages bool

	// bounds is the rect with global coords that should be drawn within.
	bounds image.Rectangle

	// bodyBounds is the bounds below the header.
	bodyBounds image.Rectangle

	// sectionDividers are bounds of the sections inside the body to render dividers.
	sectionDividers []image.Rectangle

	// zoomChangeCallback is fired when the zoom is changed. Nil if no callback registered.
	zoomChangeCallback func(zoomChange ZoomChange)
}

// NewChart creates a new Chart.
func NewChart(priceStyle PriceStyle) *Chart {
	if priceStyle == PriceStyleUnspecified {
		logger.Error("unspecified price style")
		return nil
	}

	return &Chart{
		frameBubble: rect.NewBubble(chartRounding),
		header: newHeader(&headerArgs{
			SymbolQuoteTextRenderer: chartSymbolQuoteTextRenderer,
			QuotePrinter:            chartQuotePrinter,
			ShowBarButton:           true,
			ShowCandlestickButton:   true,
			ShowRefreshButton:       true,
			ShowAddButton:           true,
			Rounding:                chartRounding,
			Padding:                 chartSectionPadding,
		}),

		price:         newPrice(priceStyle),
		priceLevel:    newPriceLevel(),
		priceCursor:   new(priceCursor),
		priceTimeline: newTimeline(view.TransparentLightGray, view.LightGray, view.TransparentGray, view.Gray),

		volume:         newVolume(priceStyle),
		volumeLevel:    newVolumeLevel(),
		volumeCursor:   new(volumeCursor),
		volumeTimeline: newTimeline(view.LightGray, view.TransparentLightGray, view.Gray, view.TransparentGray),

		timelineAxis:   new(timelineAxis),
		timelineCursor: new(timelineCursor),

		legend: newLegend(),

		loadingTextBox: text.NewBox(chartSymbolQuoteTextRenderer, "LOADING...", text.Padding(chartTextPadding)),
		errorTextBox:   text.NewBox(chartSymbolQuoteTextRenderer, "ERROR", text.Color(view.Orange), text.Padding(chartTextPadding)),
		loading:        true,
		fadeIn:         animation.New(1 * view.FPS),
	}
}

// SetPriceStyle sets the chart's price style.
func (ch *Chart) SetPriceStyle(newPriceStyle PriceStyle) {
	if newPriceStyle == PriceStyleUnspecified {
		logger.Error("unspecified price style")
		return
	}
	ch.price.SetStyle(newPriceStyle)
	ch.volume.SetStyle(newPriceStyle)
}

// SetLoading toggles the Chart's loading indicator.
func (ch *Chart) SetLoading(loading bool) {
	ch.loading = loading
	ch.header.SetLoading(loading)
}

// SetErrorMessage sets or resets an error message on the chart.
// An empty error message clears any previously set error messages.
func (ch *Chart) SetErrorMessage(errorMessage string) {
	ch.hasError = errorMessage != ""
	if errorMessage != "" {
		ch.errorTextBox.SetText(errorMessage)
	}
	ch.header.SetErrorMessage(errorMessage)
}

// Data is argument to SetData.
type Data struct {
	// Symbol is a required stock symbol.
	Symbol string

	// Quote is an optional quote. Nil when no data has been received yet.
	Quote *model.Quote

	// Chart is optional chart data. Nil when data hasn't been received yet.
	Chart *model.Chart
}

// SetData sets the data to be shown on the chart.
func (ch *Chart) SetData(data Data) {
	if !ch.hasStockUpdated && data.Chart != nil {
		ch.fadeIn.Start()
	}
	ch.hasStockUpdated = data.Chart != nil

	ch.header.SetData(data)

	dc := data.Chart
	if dc == nil {
		return
	}

	switch dc.Interval {
	case model.Intraday:
		ch.showMovingAverages = false
	case model.Daily, model.Weekly:
		ch.showMovingAverages = true
	default:
		logger.Errorf("bad interval: %v", dc.Interval)
		return
	}

	ts := dc.TradingSessionSeries

	ch.price.SetData(priceData{ts})
	ch.priceLevel.SetData(priceLevelData{ts})
	ch.priceCursor.SetData(priceCursorData{ts})
	ch.priceTimeline.SetData(timelineData{dc.Interval, ts})

	if ch.showMovingAverages {
		for _, ma := range ch.movingAverages {
			ma.Close()
		}

		ch.movingAverages = nil
		for _, ma := range dc.MovingAverageSeriesSet {
			m := newMovingAverage(movingAverageColors[dc.Interval][ma.Intervals])
			m.SetData(movingAverageData{ts, ma})
			ch.movingAverages = append(ch.movingAverages, m)
		}
	}

	ch.volume.SetData(volumeData{ts, dc.AverageVolumeSeries})
	ch.volumeLevel.SetData(volumeLevelData{ts})
	ch.volumeCursor.SetData(volumeCursorData{ts})
	ch.volumeTimeline.SetData(timelineData{dc.Interval, ts})

	ch.timelineAxis.SetData(timelineAxisData{dc.Interval, ts})
	ch.timelineCursor.SetData(timelineCursorData{dc.Interval, ts})

	ch.legend.SetData(legendData{dc.Interval, ts, dc.MovingAverageSeriesSet})
}

func (ch *Chart) SetBounds(bounds image.Rectangle) {
	ch.bounds = bounds
}

func (ch *Chart) ProcessInput(input *view.Input) {
	ch.frameBubble.SetBounds(ch.bounds)

	ch.header.SetBounds(ch.bounds)
	r, _ := ch.header.ProcessInput(input)

	ch.bodyBounds = r
	ch.loadingTextBox.SetBounds(r)
	ch.errorTextBox.SetBounds(r)

	// Calculate percentage needed for each section.
	timeLabelsPercent := float32(ch.timelineAxis.MaxLabelSize.Y+chartSectionPadding*2) / float32(r.Dy())

	// Divide up the rectangle into sections.
	rects := rect.Slice(r, timeLabelsPercent, chartVolumePercent)

	pr, vr, tr := rects[2], rects[1], rects[0]

	ch.sectionDividers = []image.Rectangle{vr, tr}

	// Pad all the rects.
	pr = pr.Inset(chartSectionPadding)
	vr = vr.Inset(chartSectionPadding)
	tr = tr.Inset(chartSectionPadding)

	// Create separate rects for each section's labels shown on the right.
	plr, vlr := pr, vr

	// Figure out width to trim off on the right of each rect for the labels.
	maxWidth := ch.priceLevel.MaxLabelSize.X
	if w := ch.volumeLevel.MaxLabelSize.X; w > maxWidth {
		maxWidth = w
	}

	// Set left side of label rects.
	plr.Min.X = pr.Max.X - maxWidth
	vlr.Min.X = vr.Max.X - maxWidth

	// Trim off the label rects from the main rects.
	pr.Max.X = plr.Min.X - chartSectionPadding
	vr.Max.X = vlr.Min.X - chartSectionPadding

	// Time labels and its cursors labels overlap and use the same rect.
	tr.Max.X = plr.Min.X
	tlr := tr

	ch.price.SetBounds(pr)
	ch.priceLevel.SetBounds(pr, plr)
	ch.priceCursor.SetBounds(pr, plr)
	ch.priceTimeline.SetBounds(pr)

	for _, ma := range ch.movingAverages {
		ma.SetBounds(pr)
	}

	ch.volume.SetBounds(vr)
	ch.volumeLevel.SetBounds(vr, vlr)
	ch.volumeCursor.SetBounds(vr, vlr)
	ch.volumeTimeline.SetBounds(vr)

	ch.timelineAxis.SetBounds(tr)
	ch.timelineCursor.SetBounds(tr, tlr)

	ch.legend.SetBounds(pr)

	ch.priceCursor.ProcessInput(input)
	ch.volumeCursor.ProcessInput(input)
	ch.timelineCursor.ProcessInput(input)
	ch.legend.ProcessInput(input)

	if input.MouseScrolled.In(ch.bounds) && ch.zoomChangeCallback != nil {
		zoomChange := ZoomChangeUnspecified
		switch input.MouseScrolled.Direction {
		case view.ScrollDown:
			zoomChange = ZoomOut
		case view.ScrollUp:
			zoomChange = ZoomIn
		default:
			logger.Error("mouse scroll event missing direction")
			return
		}

		input.AddFiredCallback(func() {
			if ch.zoomChangeCallback != nil {
				ch.zoomChangeCallback(zoomChange)
			}
		})
	}
}

// Update updates the Chart.
func (ch *Chart) Update() (dirty bool) {
	if ch.header.Update() {
		dirty = true
	}
	if ch.price.Update() {
		dirty = true
	}
	if ch.volume.Update() {
		dirty = true
	}
	if ch.legend.Update() {
		dirty = true
	}
	if ch.loadingTextBox.Update() {
		dirty = true
	}
	if ch.errorTextBox.Update() {
		dirty = true
	}
	if ch.fadeIn.Update() {
		dirty = true
	}
	return dirty
}

// Render renders the Chart.
func (ch *Chart) Render(fudge float32) {
	ch.frameBubble.Render(fudge)
	ch.header.Render(fudge)
	rect.RenderLineAtTop(ch.bodyBounds)

	// Only show messages if no prior data to show.
	if !ch.hasStockUpdated {
		if ch.loading {
			ch.loadingTextBox.Render(fudge)
			return
		}

		if ch.hasError {
			ch.errorTextBox.Render(fudge)
			return
		}
	}

	old := gfx.Alpha()
	gfx.SetAlpha(old * ch.fadeIn.Value(fudge))
	defer gfx.SetAlpha(old)

	// Render the dividers between the sections.
	for _, r := range ch.sectionDividers {
		rect.RenderLineAtTop(r)
	}

	ch.priceTimeline.Render(fudge)
	ch.priceLevel.Render(fudge)
	ch.price.Render(fudge)
	if ch.showMovingAverages {
		for _, ma := range ch.movingAverages {
			ma.Render(fudge)
		}
	}
	ch.priceCursor.Render(fudge)

	ch.volumeTimeline.Render(fudge)
	ch.volumeLevel.Render(fudge)
	ch.volume.Render(fudge)
	ch.volumeCursor.Render(fudge)

	ch.timelineAxis.Render(fudge)
	ch.timelineCursor.Render(fudge)

	ch.legend.Render(fudge)
}

// SetBarButtonClickCallback sets the callback for bar button clicks.
func (ch *Chart) SetBarButtonClickCallback(cb func()) {
	ch.header.SetBarButtonClickCallback(cb)
}

// SetCandlestickButtonClickCallback sets the callback for the candlestick button clicks.
func (ch *Chart) SetCandlestickButtonClickCallback(cb func()) {
	ch.header.SetCandlestickButtonClickCallback(cb)
}

// SetRefreshButtonClickCallback sets the callback for refresh button clicks.
func (ch *Chart) SetRefreshButtonClickCallback(cb func()) {
	ch.header.SetRefreshButtonClickCallback(cb)
}

// SetAddButtonClickCallback sets the callback for add button clicks.
func (ch *Chart) SetAddButtonClickCallback(cb func()) {
	ch.header.SetAddButtonClickCallback(cb)
}

// SetZoomChangeCallback sets the callback for zoom changes.
func (ch *Chart) SetZoomChangeCallback(cb func(zoomChange ZoomChange)) {
	ch.zoomChangeCallback = cb
}

// Close frees the resources backing the chart.
func (ch *Chart) Close() {
	ch.header.Close()
	ch.price.Close()
	ch.priceLevel.Close()
	ch.priceCursor.Close()
	ch.priceTimeline.Close()
	for _, ma := range ch.movingAverages {
		ma.Close()
	}
	ch.movingAverages = nil
	ch.volume.Close()
	ch.volumeLevel.Close()
	ch.volumeCursor.Close()
	ch.volumeTimeline.Close()
	ch.timelineAxis.Close()
	ch.timelineCursor.Close()
	ch.legend.Close()
	ch.zoomChangeCallback = nil
}

func renderCursorLines(r image.Rectangle, mousePos *view.MousePosition) {
	if mousePos.In(r) {
		gfx.SetModelMatrixRect(image.Rect(r.Min.X, mousePos.Y, r.Max.X, mousePos.Y))
		cursorHorizLine.Render()
	}

	if mousePos.WithinX(r) {
		gfx.SetModelMatrixRect(image.Rect(mousePos.X, r.Min.Y, mousePos.X, r.Max.Y))
		cursorVertLine.Render()
	}
}

func tradingSessionAtX(ts []*model.TradingSession, r image.Rectangle, x int) (int, *model.TradingSession) {
	dx := x - r.Min.X
	xPercent := float64(dx) / float64(r.Dx())
	i := int(math.Floor(float64(len(ts)) * xPercent))
	if i >= len(ts) {
		i = len(ts) - 1
	}
	return i, ts[i]
}
