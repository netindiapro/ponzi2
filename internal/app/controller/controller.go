// Package controller contains code for the controller in the MVC pattern.
package controller

import (
	"context"

	"github.com/btmura/ponzi2/internal/app/config"
	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/app/view/chart"
	"github.com/btmura/ponzi2/internal/app/view/ui"
	"github.com/btmura/ponzi2/internal/errors"
	"github.com/btmura/ponzi2/internal/logger"
	"github.com/btmura/ponzi2/internal/stock/iex"
)

// Controller runs the program in a "game loop".
type Controller struct {
	// model is the data that the Controller connects to the View.
	model *model.Model

	// ui is the UI that the Controller updates.
	ui *ui.UI

	// chartRange is the current data range to use for charts.
	chartRange model.Range

	// chartThumbRange is the current data range to use for thumbnails.
	chartThumbRange model.Range

	// chartPriceStyle is the current price style for charts and thumbnails.
	chartPriceStyle chart.PriceStyle

	// stockRefresher offers methods to refresh one or many stocks.
	stockRefresher *stockRefresher

	// configSaver offers methods to save configs in the background.
	configSaver *configSaver

	// eventController offers methods to queue and process events in the main loop.
	eventController *eventController
}

// iexClientInterface is implemented by clients in the iex package to get stock data.
type iexClientInterface interface {
	GetQuotes(ctx context.Context, req *iex.GetQuotesRequest) ([]*iex.Quote, error)
	GetCharts(ctx context.Context, req *iex.GetChartsRequest) ([]*iex.Chart, error)
}

// New creates a new Controller.
func New(iexClient iexClientInterface, token string) *Controller {
	c := &Controller{
		model:           model.New(),
		ui:              ui.New(),
		chartRange:      model.OneYear,
		chartThumbRange: model.OneYear,
		chartPriceStyle: chart.Bar,
		configSaver:     newConfigSaver(),
	}
	c.eventController = newEventController(c)
	c.stockRefresher = newStockRefresher(iexClient, token, c.eventController)
	return c
}

// RunLoop runs the loop until the user exits the app.
func (c *Controller) RunLoop() error {
	ctx := context.Background()

	cleanup, err := c.ui.Init(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// Load the config and setup the initial UI.
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.CurrentStock != nil {
		if s := cfg.CurrentStock.Symbol; s != "" {
			if err := c.setChart(ctx, s); err != nil {
				return err
			}
		}
	}

	for _, cs := range cfg.Stocks {
		if s := cs.Symbol; s != "" {
			if err := c.addChartThumb(ctx, s); err != nil {
				return err
			}
		}
	}

	if priceStyle := cfg.Settings.ChartSettings.PriceStyle; priceStyle != chart.PriceStyleUnspecified {
		c.setChartPriceStyle(priceStyle)
	}

	c.ui.SetInputSymbolSubmittedCallback(func(symbol string) {
		if err := c.setChart(ctx, symbol); err != nil {
			logger.Errorf("setChart: %v", err)
		}
	})

	c.ui.SetSidebarChangeCallback(func(symbols []string) {
		if err := c.setSidebarSymbols(symbols); err != nil {
			logger.Errorf("setSidebarSymbols: %v", err)
		}
	})

	c.ui.SetChartZoomChangeCallback(func(zoomChange chart.ZoomChange) {
		r := nextRange(c.chartRange, zoomChange)
		if c.chartRange == r {
			return
		}

		c.chartRange = r
		if err := c.refreshCurrentStock(ctx); err != nil {
			logger.Errorf("refreshCurrentStock: %v", err)
		}
	})

	c.ui.SetChartPriceStyleButtonClickCallback(func(newPriceStyle chart.PriceStyle) {
		if newPriceStyle == chart.PriceStyleUnspecified {
			logger.Error("unspecified price style")
			return
		}
		c.setChartPriceStyle(newPriceStyle)
	})

	c.ui.SetChartRefreshButtonClickCallback(func(symbol string) {
		if err := c.refreshAllStocks(ctx); err != nil {
			logger.Errorf("refreshAllStocks: %v", err)
		}
	})

	c.ui.SetChartAddButtonClickCallback(func(symbol string) {
		if err := c.addChartThumb(ctx, symbol); err != nil {
			logger.Errorf("addChartThumb: %v", err)
		}
	})

	c.ui.SetThumbRemoveButtonClickCallback(func(symbol string) {
		if err := c.removeChartThumb(symbol); err != nil {
			logger.Errorf("removeChartThumb: %v", err)
		}
	})

	c.ui.SetThumbClickCallback(func(symbol string) {
		if err := c.setChart(ctx, symbol); err != nil {
			logger.Errorf("setChart: %v", err)
		}
	})

	// Process stock refreshes and config changes in the background until the program ends.
	go c.stockRefresher.refreshLoop()
	go c.configSaver.saveLoop()

	defer func() {
		c.stockRefresher.stop()
		c.configSaver.stop()
	}()

	c.stockRefresher.start()
	c.configSaver.start()

	// Fire requests to get data for the entire UI.
	if err := c.refreshAllStocks(ctx); err != nil {
		return err
	}

	return c.ui.RunLoop(ctx, c.eventController.process)
}

func (c *Controller) setChart(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.Errorf("missing symbol")
	}

	changed, err := c.model.SetCurrentSymbol(symbol)
	if err != nil {
		return err
	}

	// If the stock is already the current one, just refresh it.
	if !changed {
		return c.refreshCurrentStock(ctx)
	}

	data, err := c.chartData(symbol, c.chartRange)
	if err != nil {
		return err
	}

	if err := c.ui.SetChart(symbol, data, c.chartPriceStyle); err != nil {
		return err
	}

	if err := c.refreshCurrentStock(ctx); err != nil {
		return err
	}

	c.configSaver.save(c.makeConfig())

	return nil
}

func (c *Controller) addChartThumb(ctx context.Context, symbol string) error {
	if symbol == "" {
		return errors.Errorf("missing symbol")
	}

	added, err := c.model.AddSidebarSymbol(symbol)
	if err != nil {
		return err
	}

	// If the stock is already added, just refresh it.
	if !added {
		return c.stockRefresher.refreshOne(ctx, symbol, c.chartThumbRange)
	}

	data, err := c.chartData(symbol, c.chartThumbRange)
	if err != nil {
		return err
	}

	if err := c.ui.AddChartThumb(symbol, data, c.chartPriceStyle); err != nil {
		return err
	}

	if err := c.stockRefresher.refreshOne(ctx, symbol, c.chartThumbRange); err != nil {
		return err
	}

	c.configSaver.save(c.makeConfig())

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

	if err := c.ui.RemoveChartThumb(symbol); err != nil {
		return err
	}

	c.configSaver.save(c.makeConfig())

	return nil
}

func (c *Controller) setSidebarSymbols(symbols []string) error {
	for _, s := range symbols {
		if err := model.ValidateSymbol(s); err != nil {
			return err
		}
	}

	if err := c.model.SetSidebarSymbols(symbols); err != nil {
		return err
	}

	c.configSaver.save(c.makeConfig())

	return nil
}

func (c *Controller) setChartPriceStyle(newPriceStyle chart.PriceStyle) {
	if newPriceStyle == chart.PriceStyleUnspecified {
		logger.Error("unspecified price style")
		return
	}

	if newPriceStyle == c.chartPriceStyle {
		return
	}

	c.chartPriceStyle = newPriceStyle
	c.ui.SetChartPriceStyle(newPriceStyle)
	c.configSaver.save(c.makeConfig())
}

func (c *Controller) chartData(symbol string, dataRange model.Range) (chart.Data, error) {
	if symbol == "" {
		return chart.Data{}, errors.Errorf("missing symbol")
	}

	data := chart.Data{Symbol: symbol}

	st, err := c.model.Stock(symbol)
	if err != nil {
		return data, err
	}

	if st == nil {
		return data, nil
	}

	for _, ch := range st.Charts {
		if ch.Range == dataRange {
			data.Quote = st.Quote
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

// onStockRefreshStarted implements the eventHandler interface.
func (c *Controller) onStockRefreshStarted(symbol string, dataRange model.Range) error {
	return c.ui.SetLoading(symbol, dataRange)
}

// onStockUpdate implements the eventHandler interface.
func (c *Controller) onStockUpdate(symbol string, q *model.Quote, ch *model.Chart) error {
	if err := model.ValidateSymbol(symbol); err != nil {
		return err
	}

	if q != nil {
		if err := model.ValidateQuote(q); err != nil {
			return err
		}

		if err := c.model.UpdateStockQuote(symbol, q); err != nil {
			return err
		}
	}

	if ch != nil {
		if err := model.ValidateChart(ch); err != nil {
			return err
		}

		if err := c.model.UpdateStockChart(symbol, ch); err != nil {
			return err
		}
	}

	if q != nil || ch != nil {
		data, err := c.chartData(symbol, c.chartRange)
		if err != nil {
			return err
		}

		return c.ui.SetData(symbol, data)
	}

	return nil
}

// onStockUpdateError implements the eventHandler interface.
func (c *Controller) onStockUpdateError(symbol string, updateErr error) error {
	logger.Errorf("stock update for %s failed: %v\n", symbol, updateErr)
	return c.ui.SetError(symbol, updateErr)
}

// onRefreshAllStocksRequest implements the eventHandler interface.
func (c *Controller) onRefreshAllStocksRequest(ctx context.Context) error {
	return c.refreshAllStocks(ctx)
}

// onEventAdded implements the eventHandler interface.
func (c *Controller) onEventAdded() {
	c.ui.WakeLoop()
}

func nextRange(r model.Range, zoomChange chart.ZoomChange) model.Range {
	// zoomRanges are the ranges from most zoomed out to most zoomed in.
	var zoomRanges = []model.Range{
		model.OneYear,
		model.OneDay,
	}

	// Find the current zoom range.
	i := 0
	for j := range zoomRanges {
		if zoomRanges[j] == r {
			i = j
		}
	}

	// Adjust the zoom one increment.
	switch zoomChange {
	case chart.ZoomIn:
		if i+1 < len(zoomRanges) {
			i++
		}
	case chart.ZoomOut:
		if i-1 >= 0 {
			i--
		}
	}

	return zoomRanges[i]
}

func (c *Controller) makeConfig() *config.Config {
	cfg := &config.Config{}
	if s := c.model.CurrentSymbol(); s != "" {
		cfg.CurrentStock = &config.Stock{Symbol: s}
	}
	for _, s := range c.model.SidebarSymbols() {
		cfg.Stocks = append(cfg.Stocks, &config.Stock{Symbol: s})
	}
	cfg.Settings.ChartSettings.PriceStyle = c.chartPriceStyle
	return cfg
}
