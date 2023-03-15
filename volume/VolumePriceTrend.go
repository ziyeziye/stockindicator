package volume

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

// The Volume Price Trend (VPT) provides a correlation between the
// volume and the price.
//
// VPT = Previous VPT + (Volume * (Current Closing - Previous Closing) / Previous Closing)
type VolumePriceTrend struct {
	Name   string
	Period int
	data   []VolumePriceTrendData
	kline  utils.Klines
}

// VolumePriceTrendData
type VolumePriceTrendData struct {
	Time  time.Time
	Value float64
}

// NewVolumePriceTrend new Func
func NewVolumePriceTrend(list utils.Klines, period int) *VolumePriceTrend {
	m := &VolumePriceTrend{
		Name:   "VolumePriceTrend",
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultVolumePriceTrend new Func
func NewDefaultVolumePriceTrend(list utils.Klines) *VolumePriceTrend {
	return NewVolumePriceTrend(list, 14)
}

// Calculation Func
func (e *VolumePriceTrend) Calculation() *VolumePriceTrend {

	period := e.Period
	var closing, volume []float64
	for _, v := range e.kline {
		closing = append(closing, v.Close)
		volume = append(volume, v.Volume)
	}

	previousClosing := utils.ShiftRightAndFillBy(period, closing[0], closing)
	vpt := utils.Multiply(volume, utils.Divide(utils.Subtract(closing, previousClosing), previousClosing))
	vals := utils.Sum(len(vpt), vpt)

	for i := 0; i < len(vals); i++ {
		e.data = append(e.data, VolumePriceTrendData{
			Time:  e.kline[i].Time,
			Value: vals[i],
		})
	}

	return e
}

// AnalysisSide Func
// func (e *VolumePriceTrend) AnalysisSide() utils.SideData {
// 	sides := make([]utils.Side, len(e.kline))

// 	if len(e.data) == 0 {
// 		e = e.Calculation()
// 	}

// 	for i, v := range e.data {
// 		if i < 1 {
// 			continue
// 		}

// 		var prevItem = e.data[i-1]

// 		if v.Value < 10 && prevItem.Value > 10 {
// 			sides[i] = utils.Buy
// 		} else if v.Value > 90 && prevItem.Value < 90 {
// 			sides[i] = utils.Sell
// 		} else {
// 			sides[i] = utils.Hold
// 		}
// 	}
// 	return utils.SideData{
// 		Name: e.Name,
// 		Data: sides,
// 	}
// }

// GetData Func
func (e *VolumePriceTrend) GetData() []VolumePriceTrendData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}
