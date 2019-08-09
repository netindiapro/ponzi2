package controller

import (
	"sort"
	"time"

	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/errors"
	"github.com/btmura/ponzi2/internal/stock/iex"
)

// maxDataWeeks is maximum number of weeks of data to retain.
const maxDataWeeks = 12 /* months */ * 4 /* weeks = 1 year */

// Stochastic parameters.
const (
	k = 10
	d = 3
)

func modelOneDayChart(chart *iex.Chart) (*model.Chart, error) {
	// TODO(btmura): remove duplication with modelTradingSessions
	var ts []*model.TradingSession
	for _, p := range chart.ChartPoints {
		ts = append(ts, &model.TradingSession{
			Date:          p.Date,
			Open:          p.Open,
			High:          p.High,
			Low:           p.Low,
			Close:         p.Close,
			Volume:        p.Volume,
			Change:        p.Change,
			PercentChange: p.ChangePercent,
		})
	}
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Date.Before(ts[j].Date)
	})

	return &model.Chart{
		Range: model.OneDay,
		TradingSessionSeries: &model.TradingSessionSeries{
			TradingSessions: ts,
		},
	}, nil
}

func modelOneYearChart(quote *iex.Quote, chart *iex.Chart) (*model.Chart, error) {
	ds := modelTradingSessions(quote, chart)
	ws := weeklyModelTradingSessions(ds)

	m25 := modelMovingAverages(ds, 25)
	m50 := modelMovingAverages(ds, 50)
	m200 := modelMovingAverages(ds, 200)

	dsto := modelStochastics(ds)
	wsto := modelStochastics(ws)

	if len(ws) > maxDataWeeks {
		start := ws[len(ws)-maxDataWeeks:][0].Date

		trimmedTradingSessions := func(vs []*model.TradingSession) []*model.TradingSession {
			for i, v := range vs {
				if v.Date == start {
					return vs[i:]
				}
			}
			return vs
		}

		trimmedMovingAverages := func(vs []*model.MovingAverage) []*model.MovingAverage {
			for i, v := range vs {
				if v.Date == start {
					return vs[i:]
				}
			}
			return vs
		}

		trimmedStochastics := func(vs []*model.Stochastic) []*model.Stochastic {
			for i, v := range vs {
				if v.Date == start {
					return vs[i:]
				}
			}
			return vs
		}

		ds = trimmedTradingSessions(ds)
		m25 = trimmedMovingAverages(m25)
		m50 = trimmedMovingAverages(m50)
		m200 = trimmedMovingAverages(m200)
		dsto = trimmedStochastics(dsto)
		wsto = trimmedStochastics(wsto)
	}

	return &model.Chart{
		Range:                  model.OneYear,
		TradingSessionSeries:   &model.TradingSessionSeries{TradingSessions: ds},
		MovingAverageSeries25:  &model.MovingAverageSeries{MovingAverages: m25},
		MovingAverageSeries50:  &model.MovingAverageSeries{MovingAverages: m50},
		MovingAverageSeries200: &model.MovingAverageSeries{MovingAverages: m200},
		DailyStochasticSeries:  &model.StochasticSeries{Stochastics: dsto},
		WeeklyStochasticSeries: &model.StochasticSeries{Stochastics: wsto},
	}, nil
}

func modelQuote(q *iex.Quote) (*model.Quote, error) {
	if q == nil {
		return nil, errors.Errorf("missing quote")
	}

	src, err := modelSource(q.LatestSource)
	if err != nil {
		return nil, err
	}

	return &model.Quote{
		CompanyName:   q.CompanyName,
		LatestPrice:   q.LatestPrice,
		LatestSource:  src,
		LatestTime:    q.LatestTime,
		LatestUpdate:  q.LatestUpdate,
		LatestVolume:  q.LatestVolume,
		Open:          q.Open,
		High:          q.High,
		Low:           q.Low,
		Close:         q.Close,
		Change:        q.Change,
		ChangePercent: q.ChangePercent,
	}, nil
}

func modelSource(src iex.Source) (model.Source, error) {
	switch src {
	case iex.SourceUnspecified:
		return model.SourceUnspecified, nil
	case iex.IEXRealTimePrice:
		return model.IEXRealTimePrice, nil
	case iex.FifteenMinuteDelayedPrice:
		return model.FifteenMinuteDelayedPrice, nil
	case iex.Close:
		return model.Close, nil
	case iex.PreviousClose:
		return model.PreviousClose, nil
	default:
		return 0, errors.Errorf("unrecognized iex source: %v", src)
	}
}

func modelRange(r iex.Range) (model.Range, error) {
	switch r {
	case iex.RangeUnspecified:
		return model.RangeUnspecified, nil
	case iex.OneDay:
		return model.OneDay, nil
	case iex.TwoYears:
		return model.OneYear, nil
	default:
		return 0, errors.Errorf("unrecognized iex range: %v", r)
	}
}

func modelTradingSessions(quote *iex.Quote, chart *iex.Chart) []*model.TradingSession {
	var ts []*model.TradingSession

	for _, p := range chart.ChartPoints {
		ts = append(ts, &model.TradingSession{
			Date:          p.Date,
			Open:          p.Open,
			High:          p.High,
			Low:           p.Low,
			Close:         p.Close,
			Volume:        p.Volume,
			Change:        p.Change,
			PercentChange: p.ChangePercent,
		})
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Date.Before(ts[j].Date)
	})

	// Add a trading session for the current quote if we do not have data
	// for today's trading session, so that the chart includes the latest quote.

	if quote == nil {
		return ts
	}

	q := quote

	// Real-time quotes won't have OHLC set, but they will have a latest price.
	// Fake OHLC so something shows up on the chart by using the latest price.
	// TODO(btmura): considering using a different color for a fake ohlc
	o, h, l, c := q.Open, q.High, q.Low, q.Close
	if o == 0 && h == 0 && l == 0 && c == 0 {
		o = q.LatestPrice - q.Change
		c = q.LatestPrice
		l = (o + c) / 2
		h = (o + c) / 2
	}

	t := &model.TradingSession{
		Date:          q.LatestTime,
		Open:          o,
		High:          h,
		Low:           l,
		Close:         c,
		Volume:        q.LatestVolume,
		Change:        q.Change,
		PercentChange: q.ChangePercent,
	}

	if len(ts) == 0 {
		return []*model.TradingSession{t}
	}

	clean := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	if clean(t.Date) == clean(ts[len(ts)-1].Date) {
		return ts
	}

	return append(ts, t)
}

func weeklyModelTradingSessions(ds []*model.TradingSession) (ws []*model.TradingSession) {
	for _, p := range ds {
		diffWeek := ws == nil
		if !diffWeek {
			_, week := p.Date.ISOWeek()
			_, prevWeek := ws[len(ws)-1].Date.ISOWeek()
			diffWeek = week != prevWeek
		}

		if diffWeek {
			pcopy := *p
			ws = append(ws, &pcopy)
		} else {
			ls := ws[len(ws)-1]
			if ls.High < p.High {
				ls.High = p.High
			}
			if ls.Low > p.Low {
				ls.Low = p.Low
			}
			ls.Close = p.Close
			ls.Volume += p.Volume
		}
	}
	return ws
}

func modelMovingAverages(ts []*model.TradingSession, n int) []*model.MovingAverage {
	average := func(i, n int) (avg float32) {
		if i+1-n < 0 {
			return 0 // Not enough data
		}
		var sum float32
		for j := 0; j < n; j++ {
			sum += ts[i-j].Close
		}
		return sum / float32(n)
	}

	var ms []*model.MovingAverage
	for i := range ts {
		ms = append(ms, &model.MovingAverage{
			Date:  ts[i].Date,
			Value: average(i, n),
		})
	}
	return ms
}

func modelStochastics(ts []*model.TradingSession) []*model.Stochastic {
	// Calculate fast %K for stochastics.
	fastK := make([]float32, len(ts))
	for i := range ts {
		if i+1 < k {
			continue
		}

		highestHigh, lowestLow := ts[i].High, ts[i].Low
		for j := 0; j < k; j++ {
			if highestHigh < ts[i-j].High {
				highestHigh = ts[i-j].High
			}
			if lowestLow > ts[i-j].Low {
				lowestLow = ts[i-j].Low
			}
		}
		fastK[i] = (ts[i].Close - lowestLow) / (highestHigh - lowestLow)
	}

	// Setup slice to hold stochastics.
	var ms []*model.Stochastic
	for i := range ts {
		ms = append(ms, &model.Stochastic{Date: ts[i].Date})
	}

	// Calculate fast %D (slow %K) for stochastics.
	for i := range ts {
		if i+1 < k+d {
			continue
		}
		ms[i].K = (fastK[i] + fastK[i-1] + fastK[i-2]) / 3
	}

	// Calculate slow %D for stochastics.
	for i := range ts {
		if i+1 < k+d+d {
			continue
		}
		ms[i].D = (ms[i].K + ms[i-1].K + ms[i-2].K) / 3
	}

	return ms
}
