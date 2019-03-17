// Package controller contains code for the controller in the MVC pattern.
package controller

import (
	"context"
	"log"
	"time"

	"github.com/golang/glog"

	"github.com/btmura/ponzi2/internal/app/config"
	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/app/view"
	"github.com/btmura/ponzi2/internal/status"
	"github.com/btmura/ponzi2/internal/stock/iex"
)

// loc is the timezone to use when parsing dates.
var loc = mustLoadLocation("America/New_York")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Fatalf("time.LoadLocation(%s) failed: %v", name, err)
	}
	return loc
}

// zoomRanges are the ranges from most zoomed out to most zoomed in.
var zoomRanges = []model.Range{
	model.OneYear,
	model.OneDay,
}

// Controller runs the program in a "game loop".
type Controller struct {
	// model is the data that the Controller connects to the View.
	model *model.Model

	// view is the UI that the Controller updates.
	view *view.View

	// title controls the title bar.
	title *view.Title

	// symbolToChartMap maps symbol to Chart. Only one entry right now.
	symbolToChartMap map[string]*view.Chart

	// symbolToChartThumbMap maps symbol to ChartThumbnail.
	symbolToChartThumbMap map[string]*view.ChartThumb

	// chartRange is the current data range to use for Charts.
	chartRange model.Range

	// chartThumbRange is the current data range to use for ChartThumbnails.
	chartThumbRange model.Range

	// stockRefresher handles refreshing stock data.
	stockRefresher *stockRefresher

	// eventController offers methods to manage events like stock updates.
	eventController *eventController

	// configController controls loading and saving configs.
	configController *configController
}

// New creates a new Controller.
func New(iexClient *iex.Client) *Controller {
	c := &Controller{
		model:                 model.New(),
		view:                  view.New(),
		title:                 view.NewTitle(),
		symbolToChartMap:      map[string]*view.Chart{},
		symbolToChartThumbMap: map[string]*view.ChartThumb{},
		chartRange:            model.OneYear,
		chartThumbRange:       model.OneYear,
		configController:      newConfigController(),
	}
	c.eventController = newEventController(c)
	c.stockRefresher = newStockRefresher(iexClient, c.eventController)
	return c
}

// RunLoop runs the loop until the user exits the app.
func (c *Controller) RunLoop() error {
	ctx := context.Background()

	cleanup, err := c.view.Init(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// Load the config and setup the initial UI.
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if s := cfg.GetCurrentStock().GetSymbol(); s != "" {
		c.setChart(ctx, s)
	}

	for _, cs := range cfg.GetStocks() {
		if s := cs.GetSymbol(); s != "" {
			if err := c.addChartThumb(ctx, s); err != nil {
				return err
			}
		}
	}

	c.view.SetTitle(c.title)

	c.view.SetInputSymbolSubmittedCallback(func(symbol string) {
		c.setChart(ctx, symbol)
	})

	c.view.SetChartZoomChangeCallback(func(zoomChange view.ZoomChange) {
		// Find the current zoom range.
		i := 0
		for j := range zoomRanges {
			if zoomRanges[j] == c.chartRange {
				i = j
			}
		}

		// Adjust the zoom one increment.
		switch zoomChange {
		case view.ZoomIn:
			if i+1 < len(zoomRanges) {
				i++
			}
		case view.ZoomOut:
			if i-1 >= 0 {
				i--
			}
		}

		// Ignore if no change in zoom.
		if c.chartRange == zoomRanges[i] {
			return
		}

		// Set zoom and refresh the current stock.
		c.chartRange = zoomRanges[i]

		if err := c.refreshCurrentStock(ctx); err != nil {
			glog.Fatalf("TODO(btmura): remove log fatal, refreshStocks: %v", err)
		}
	})

	// Process config changes in the background until the program ends.
	go c.configController.saveLoop()

	// Refresh stocks during market hours.
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for t := range ticker.C {
			n := time.Now()
			open := time.Date(n.Year(), n.Month(), n.Day(), 9, 30, 0, 0, loc)
			close := time.Date(n.Year(), n.Month(), n.Day(), 16, 0, 0, 0, loc)

			if t.Before(open) || t.After(close) {
				glog.V(2).Infof("ignoring refresh ticker at %v", t.Format("1/2/2006 3:04:05 PM"))
				continue
			}

			c.eventController.addEventLocked(event{refreshAllStocks: true})
		}
	}()

	defer func() {
		ticker.Stop()
		c.stockRefresher.stop()
		c.configController.stop()
	}()

	c.stockRefresher.start()
	c.configController.start()

	// Fire requests to get data for the entire UI.
	if err := c.refreshAllStocks(ctx); err != nil {
		return err
	}

	return c.view.RunLoop(ctx, c.eventController.process)
}

func (c *Controller) setChart(ctx context.Context, symbol string) error {
	if symbol == "" {
		return status.Error("missing symbol")
	}

	changed, err := c.model.SetCurrentSymbol(symbol)
	if err != nil {
		return err
	}

	// If the stock is already the current one, just refresh it.
	if !changed {
		return c.refreshCurrentStock(ctx)
	}

	for symbol, ch := range c.symbolToChartMap {
		delete(c.symbolToChartMap, symbol)
		ch.Close()
	}

	ch := view.NewChart()
	c.symbolToChartMap[symbol] = ch

	data, err := c.chartData(symbol, c.chartRange)
	if err != nil {
		return err
	}

	if err := c.title.SetData(data); err != nil {
		return err
	}

	if err := ch.SetData(data); err != nil {
		return err
	}

	ch.SetRefreshButtonClickCallback(func() {
		if err := c.refreshAllStocks(ctx); err != nil {
			glog.Fatalf("TODO(btmura): remove log fatal, refreshStocks: %v", err)
		}
	})
	ch.SetAddButtonClickCallback(func() {
		if err := c.addChartThumb(ctx, symbol); err != nil {
			glog.Fatalf("TODO(btmura): remove log fatal, addChartThumb: %v", err)
		}
	})

	c.view.SetChart(ch)

	if err := c.refreshCurrentStock(ctx); err != nil {
		return err
	}

	c.configController.save(c.model)

	return nil
}

func (c *Controller) addChartThumb(ctx context.Context, symbol string) error {
	if symbol == "" {
		return status.Error("missing symbol")
	}

	added, err := c.model.AddSidebarSymbol(symbol)
	if err != nil {
		return err
	}

	// If the stock is already added, just refresh it.
	if !added {
		return c.stockRefresher.refreshOne(ctx, symbol, c.chartThumbRange)
	}

	th := view.NewChartThumb()
	c.symbolToChartThumbMap[symbol] = th

	data, err := c.chartData(symbol, c.chartThumbRange)
	if err != nil {
		return err
	}

	if err := th.SetData(data); err != nil {
		return err
	}

	th.SetRemoveButtonClickCallback(func() {
		if err := c.removeChartThumb(symbol); err != nil {
			glog.Fatalf("TODO(btmura): remove log fatal, removeChartThumb: %v", err)
		}
	})
	th.SetThumbClickCallback(func() {
		if err := c.setChart(ctx, symbol); err != nil {
			glog.Fatalf("TODO(btmura): remove log fatal, setChart: %v", err)
		}
	})

	c.view.AddChartThumb(th)

	if err := c.stockRefresher.refreshOne(ctx, symbol, c.chartThumbRange); err != nil {
		return err
	}

	c.configController.save(c.model)

	return nil
}

func (c *Controller) removeChartThumb(symbol string) error {
	if symbol == "" {
		return nil
	}

	removed, err := c.model.RemoveSidebarSymbol(symbol)
	if err != nil {
		return err
	}

	if !removed {
		return nil
	}

	th := c.symbolToChartThumbMap[symbol]
	delete(c.symbolToChartThumbMap, symbol)
	th.Close()

	c.view.RemoveChartThumb(th)
	c.configController.save(c.model)

	return nil
}

func (c *Controller) chartData(symbol string, dataRange model.Range) (*view.ChartData, error) {
	if symbol == "" {
		return nil, status.Error("missing symbol")
	}

	data := &view.ChartData{Symbol: symbol}

	st, err := c.model.Stock(symbol)
	if err != nil {
		return nil, err
	}

	if st == nil {
		return data, nil
	}

	for _, ch := range st.Charts {
		if ch.Range == dataRange {
			data.Quote = ch.Quote
			data.Chart = ch
			return data, nil
		}
	}

	return data, nil
}

func (c *Controller) refreshCurrentStock(ctx context.Context) error {
	d := new(dataRequestBuilder)
	if s := c.model.CurrentSymbol(); s != "" {
		if err := d.add([]string{s}, c.chartRange); err != nil {
			return err
		}
	}
	return c.stockRefresher.refresh(ctx, d)
}

func (c *Controller) refreshAllStocks(ctx context.Context) error {
	d := new(dataRequestBuilder)

	if s := c.model.CurrentSymbol(); s != "" {
		if err := d.add([]string{s}, c.chartRange); err != nil {
			return err
		}
	}

	if err := d.add(c.model.SidebarSymbols(), c.chartThumbRange); err != nil {
		return err
	}

	return c.stockRefresher.refresh(ctx, d)
}

// processStockRefreshStarted implements the eventProcessor interface.
func (c *Controller) processStockRefreshStarted(symbol string, dataRange model.Range) error {
	if err := model.ValidateSymbol(symbol); err != nil {
		return err
	}

	if dataRange == model.RangeUnspecified {
		return status.Error("range not set")
	}

	for s, ch := range c.symbolToChartMap {
		if s == symbol && c.chartRange == dataRange {
			ch.SetLoading(true)
			ch.SetError(false)
		}
	}

	for s, th := range c.symbolToChartThumbMap {
		if s == symbol && c.chartThumbRange == dataRange {
			th.SetLoading(true)
			th.SetError(false)
		}
	}

	return nil
}

// processStockChartUpdate implements the eventProcessor interface.
func (c *Controller) processStockChartUpdate(symbol string, ch *model.Chart) error {
	if err := model.ValidateSymbol(symbol); err != nil {
		return err
	}

	if err := model.ValidateChart(ch); err != nil {
		return err
	}

	if err := c.model.UpdateStockChart(symbol, ch); err != nil {
		return err
	}

	if ch, ok := c.symbolToChartMap[symbol]; ok {
		ch.SetLoading(false)

		data, err := c.chartData(symbol, c.chartRange)
		if err != nil {
			return err
		}

		if err := c.title.SetData(data); err != nil {
			return err
		}

		if err := ch.SetData(data); err != nil {
			return err
		}
	}

	if th, ok := c.symbolToChartThumbMap[symbol]; ok {
		th.SetLoading(false)

		data, err := c.chartData(symbol, c.chartThumbRange)
		if err != nil {
			return err
		}

		if err := th.SetData(data); err != nil {
			return err
		}
	}

	return nil
}

// processStockChartUpdateError implements the eventProcessor interface.
func (c *Controller) processStockChartUpdateError(symbol string, updateErr error) error {
	if err := model.ValidateSymbol(symbol); err != nil {
		return err
	}

	if ch, ok := c.symbolToChartMap[symbol]; ok {
		ch.SetLoading(false)
		ch.SetError(true)
	}

	if th, ok := c.symbolToChartThumbMap[symbol]; ok {
		th.SetLoading(false)
		th.SetError(true)
	}

	return nil
}

// processRefreshAllStocks implements the eventProcessor interface.
func (c *Controller) processRefreshAllStocks(ctx context.Context) error {
	return c.refreshAllStocks(ctx)
}

// notifyProcessor implements the eventProcessor interface.
func (c *Controller) notifyProcessor() {
	c.view.WakeLoop()
}
