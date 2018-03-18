package model

import (
	"sort"
	"time"

	"github.com/golang/glog"

	"github.com/btmura/ponzi2/internal/stock"
)

// Model keeps track of the user's stocks.
type Model struct {
	// CurrentStock is the stock currently being viewed.
	CurrentStock *Stock

	// SavedStocks is the user's saved stocks.
	SavedStocks []*Stock
}

// Stock models a single stock.
type Stock struct {
	// Symbol is the symbol of the stock.
	Symbol string

	// DailySessions are trading sessions that span a single day.
	DailySessions []*TradingSession

	// WeeklySessions are trading sessions that span a single week.
	WeeklySessions []*TradingSession

	// LastUpdateTime is when the ModelStock was last updated.
	LastUpdateTime time.Time
}

// TradingSession models a single trading session.
type TradingSession struct {
	Date   time.Time
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume int

	Change        float32
	PercentChange float32

	K float32
	D float32

	MovingAverage25  float32
	MovingAverage50  float32
	MovingAverage200 float32
}

// StockUpdate is an update that can be applied to the model.
type StockUpdate struct {
	// Symbol is the symbol of the stock.
	Symbol string

	// DailySessions are trading sessions that span a single day.
	DailySessions []*TradingSession

	// WeeklySessions are trading sessions that span a single week.
	WeeklySessions []*TradingSession
}

// NewModel creates a new Model.
func NewModel() *Model {
	return &Model{}
}

// SetCurrentStock sets the current stock by symbol. It returns the
// corresponding ModelStock and true if the current stock changed.
func (m *Model) SetCurrentStock(symbol string) (st *Stock, changed bool) {
	if symbol == "" {
		glog.Fatal("SetCurrentStock: cannot set current stock to empty symbol")
	}

	if m.CurrentStock != nil && m.CurrentStock.Symbol == symbol {
		return m.CurrentStock, false
	}

	if m.CurrentStock = m.stock(symbol); m.CurrentStock == nil {
		m.CurrentStock = &Stock{Symbol: symbol}
	}
	return m.CurrentStock, true
}

// AddSavedStock adds the stock by symbol. It returns the corresponding
// ModelStock and true if the stock was newly added.
func (m *Model) AddSavedStock(symbol string) (st *Stock, added bool) {
	if symbol == "" {
		glog.Fatal("AddSavedStock: cannot add empty symbol")
	}

	for _, st := range m.SavedStocks {
		if st.Symbol == symbol {
			return st, false
		}
	}

	if st = m.stock(symbol); st == nil {
		st = &Stock{Symbol: symbol}
	}
	m.SavedStocks = append(m.SavedStocks, st)
	return st, true
}

// RemoveSavedStock removes the stock by symbol and returns true if removed.
func (m *Model) RemoveSavedStock(symbol string) (removed bool) {
	if symbol == "" {
		glog.Fatal("RemovedSavedStock: cannot remove empty symbol")
	}

	for i, st := range m.SavedStocks {
		if st.Symbol == symbol {
			m.SavedStocks = append(m.SavedStocks[:i], m.SavedStocks[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateStock updates a stock with the update if it is in the model.
func (m *Model) UpdateStock(update *StockUpdate) (st *Stock, updated bool) {
	if st = m.stock(update.Symbol); st == nil {
		return nil, false
	}
	st.DailySessions = update.DailySessions
	st.WeeklySessions = update.WeeklySessions
	st.LastUpdateTime = time.Now()
	return st, true
}

func (m *Model) stock(symbol string) *Stock {
	if m.CurrentStock != nil && m.CurrentStock.Symbol == symbol {
		return m.CurrentStock
	}

	for _, st := range m.SavedStocks {
		if st.Symbol == symbol {
			return st
		}
	}

	return nil
}

// Price returns the most recent price or 0 if no data.
func (m *Stock) Price() float32 {
	if len(m.DailySessions) == 0 {
		return 0
	}
	return m.DailySessions[len(m.DailySessions)-1].Close
}

// Change returns the most recent change or 0 if no data.
func (m *Stock) Change() float32 {
	if len(m.DailySessions) == 0 {
		return 0
	}
	return m.DailySessions[len(m.DailySessions)-1].Change
}

// PercentChange returns the most recent percent change or 0 if no data.
func (m *Stock) PercentChange() float32 {
	if len(m.DailySessions) == 0 {
		return 0
	}
	return m.DailySessions[len(m.DailySessions)-1].PercentChange
}

// Date returns the most recent date or zero time if no data.
func (m *Stock) Date() time.Time {
	if len(m.DailySessions) == 0 {
		return time.Time{}
	}
	return m.DailySessions[len(m.DailySessions)-1].Date
}

// NewStockUpdate returns a StockUpdate that can be applied to the model.
func NewStockUpdate(symbol string, ts []*stock.TradingSession) *StockUpdate {
	ds := dailySessions(ts)
	ws := weeklySessions(ds)

	fillChangeValues(ds)
	fillChangeValues(ws)

	fillStochastics(ds)
	fillStochastics(ws)

	fillMovingAverages(ds)
	fillMovingAverages(ws)

	ds, ws = trimSessions(ds, ws)

	return &StockUpdate{
		Symbol:         symbol,
		DailySessions:  ds,
		WeeklySessions: ws,
	}
}

func dailySessions(ts []*stock.TradingSession) (ds []*TradingSession) {
	for _, s := range ts {
		ds = append(ds, &TradingSession{
			Date:   s.Date,
			Open:   s.Open,
			High:   s.High,
			Low:    s.Low,
			Close:  s.Close,
			Volume: s.Volume,
		})
	}
	sort.Slice(ds, func(i, j int) bool {
		return ds[i].Date.Before(ds[j].Date)
	})
	return ds
}

func weeklySessions(ds []*TradingSession) (ws []*TradingSession) {
	for _, s := range ds {
		diffWeek := ws == nil
		if !diffWeek {
			_, week := s.Date.ISOWeek()
			_, prevWeek := ws[len(ws)-1].Date.ISOWeek()
			diffWeek = week != prevWeek
		}

		if diffWeek {
			sc := *s
			ws = append(ws, &sc)
		} else {
			ls := ws[len(ws)-1]
			if ls.High < s.High {
				ls.High = s.High
			}
			if ls.Low > s.Low {
				ls.Low = s.Low
			}
			ls.Close = s.Close
			ls.Volume += s.Volume
		}
	}
	return ws
}

func fillChangeValues(ss []*TradingSession) {
	for i := range ss {
		if i > 0 {
			ss[i].Change = ss[i].Close - ss[i-1].Close
			ss[i].PercentChange = ss[i].Change / ss[i-1].Close
		}
	}
}

func fillStochastics(ss []*TradingSession) {
	const (
		k = 10
		d = 3
	)

	// Calculate fast %K for stochastics.
	fastK := make([]float32, len(ss))
	for i := range ss {
		if i+1 < k {
			continue
		}

		highestHigh, lowestLow := ss[i].High, ss[i].Low
		for j := 0; j < k; j++ {
			if highestHigh < ss[i-j].High {
				highestHigh = ss[i-j].High
			}
			if lowestLow > ss[i-j].Low {
				lowestLow = ss[i-j].Low
			}
		}
		fastK[i] = (ss[i].Close - lowestLow) / (highestHigh - lowestLow)
	}

	// Calculate fast %D (slow %K) for stochastics.
	for i := range ss {
		if i+1 < k+d {
			continue
		}
		ss[i].K = (fastK[i] + fastK[i-1] + fastK[i-2]) / 3
	}

	// Calculate slow %D for stochastics.
	for i := range ss {
		if i+1 < k+d+d {
			continue
		}
		ss[i].D = (ss[i].K + ss[i-1].K + ss[i-2].K) / 3
	}
}

func fillMovingAverages(ss []*TradingSession) {
	average := func(i, n int) (avg float32) {
		if i+1-n < 0 {
			return 0 // Not enough data
		}
		var sum float32
		for j := 0; j < n; j++ {
			sum += ss[i-j].Close
		}
		return sum / float32(n)
	}

	for i := range ss {
		ss[i].MovingAverage25 = average(i, 25)
		ss[i].MovingAverage50 = average(i, 50)
		ss[i].MovingAverage200 = average(i, 200)
	}
}

func trimSessions(ds, ws []*TradingSession) (trimDs, trimWs []*TradingSession) {
	const sixMonthWeeks = 4 /* weeks */ * 6 /* months */
	if len(ws) >= sixMonthWeeks {
		ws = ws[len(ws)-sixMonthWeeks:]
		for i := range ds {
			if ds[i].Date == ws[0].Date {
				ds = ds[i:]
				return ds, ws
			}
		}
	}
	return ds, ws
}
