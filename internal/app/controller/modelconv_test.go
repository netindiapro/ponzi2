package controller

import (
	"testing"
	"time"

	"github.com/btmura/ponzi2/internal/app/model"
	"github.com/btmura/ponzi2/internal/stock/iex"
	"github.com/google/go-cmp/cmp"
)

func TestModelIntradayChart(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *iex.Chart
		want  *model.Chart
	}{
		{
			input: &iex.Chart{
				ChartPoints: []*iex.ChartPoint{
					{
						Date:   time.Date(2018, time.September, 18, 15, 57, 0, 0, time.UTC),
						Open:   218.44,
						High:   218.49,
						Low:    218.37,
						Close:  218.49,
						Volume: 2607,
					},
				},
			},
			want: &model.Chart{
				Interval: model.Intraday,
				TradingSessionSeries: &model.TradingSessionSeries{
					TradingSessions: []*model.TradingSession{
						{
							Date:   time.Date(2018, time.September, 18, 15, 57, 0, 0, time.UTC),
							Open:   218.44,
							High:   218.49,
							Low:    218.37,
							Close:  218.49,
							Volume: 2607,
						},
					},
				},
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			got := modelIntradayChart(tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestExponentialMovingAverages(t *testing.T) {
	for _, tt := range []struct {
		desc                 string
		inputTradingSessions []*model.TradingSession
		inputNumDays         int
		want                 []*model.MovingAverageValue
	}{
		{
			inputTradingSessions: []*model.TradingSession{
				{
					Date:  time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC),
					Close: 100,
				},
				{
					Date:  time.Date(2018, time.September, 2, 0, 0, 0, 0, time.UTC),
					Close: 200,
				},
				{
					Date:  time.Date(2018, time.September, 3, 0, 0, 0, 0, time.UTC),
					Close: 300,
				},
			},
			inputNumDays: 2,
			want: []*model.MovingAverageValue{
				{
					Date:  time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC),
					Value: 0,
				},
				{
					Date:  time.Date(2018, time.September, 2, 0, 0, 0, 0, time.UTC),
					Value: 0,
				},
				{
					Date:  time.Date(2018, time.September, 3, 0, 0, 0, 0, time.UTC),
					Value: 250,
				},
			},
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			got := modelExponentialMovingAverages(tt.inputTradingSessions, tt.inputNumDays)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("diff (-want, +got)\n%s", diff)
			}
		})
	}
}
