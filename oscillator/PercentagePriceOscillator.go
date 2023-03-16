package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// Percentage Price Oscillator (PPO). 百分比价格振荡器,是一个动量技术指标，它表示动量方向作为振荡器的迹象
// 虽然它确实与一些更受欢迎的振荡器（如 MACD）有一些相似之处, 这也是相当奇特的，因为它使用价格的百分比变化来计算动量，而不是绝对价格.
// 正 PPO 线表示看涨趋势, 而负 PPO 线表示看跌趋势.
// 还可以根据 PPO 线和信号线的交互方式确定动量. 每当 PPO 线越过信号线上方时，动量是看涨的, 当PPPO线越过信号线下方时，看跌.
//
// PPO = ((EMA(fastPeriod, prices) - EMA(slowPeriod, prices)) / EMA(longPeriod, prices)) * 100
// Signal = EMA(9, PPO)
// Histogram = PPO - Signal
type PercentagePriceOscillator struct {
	Name         string
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
	data         []PercentagePriceOscillatorData
	kline        utils.Klines
}

// PercentagePriceOscillatorData
type PercentagePriceOscillatorData struct {
	Time      time.Time
	PPO       float64
	Signal    float64
	Histogram float64
}

// NewPercentagePriceOscillator new Func
func NewPercentagePriceOscillator(list utils.Klines, fastPeriod, slowPeriod, signalPeriod int) *PercentagePriceOscillator {
	m := &PercentagePriceOscillator{
		Name:         fmt.Sprintf("PercentagePriceOscillator%d-%d-%d", fastPeriod, slowPeriod, signalPeriod),
		kline:        list,
		FastPeriod:   fastPeriod,
		SlowPeriod:   slowPeriod,
		SignalPeriod: signalPeriod,
	}
	return m
}

// NewDefaultPercentagePriceOscillator new Func
func NewDefaultPercentagePriceOscillator(list utils.Klines) *PercentagePriceOscillator {
	return NewPercentagePriceOscillator(list, 12, 26, 9)
}

// Calculation Func
func (e *PercentagePriceOscillator) Calculation() *PercentagePriceOscillator {

	var price []float64
	var fastPeriod = e.FastPeriod
	var slowPeriod = e.SlowPeriod
	var signalPeriod = e.SignalPeriod
	for _, v := range e.kline {
		price = append(price, v.Close)
	}

	fastEma := trend.NewEma(utils.CloseArrayToKline(price), fastPeriod).GetValues()
	slowEma := trend.NewEma(utils.CloseArrayToKline(price), slowPeriod).GetValues()

	ppo := utils.MultiplyBy(utils.Divide(utils.Subtract(fastEma, slowEma), slowEma), 100)
	signal := trend.NewEma(utils.CloseArrayToKline(ppo), signalPeriod).GetValues()
	histogram := utils.Subtract(ppo, signal)

	for i := 0; i < len(ppo); i++ {
		e.data = append(e.data, PercentagePriceOscillatorData{
			Time:      e.kline[i].Time,
			PPO:       ppo[i],
			Signal:    signal[i],
			Histogram: histogram[i],
		})
	}
	return e
}

// AnalysisSide Func
// 正 PPO 线表示看涨趋势, 而负 PPO 线表示看跌趋势.
// 还可以根据 PPO 线和信号线的交互方式确定动量. 每当 PPO 线越过信号线上方时，动量是看涨的, 当PPPO线越过信号线下方时，看跌.
func (e *PercentagePriceOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if v.PPO > v.Signal && prevItem.PPO < prevItem.Signal {
			sides[i] = utils.Buy
		} else if v.PPO < v.Signal && prevItem.PPO > prevItem.Signal {
			sides[i] = utils.Sell
		} else {
			sides[i] = utils.Hold
		}
	}
	return utils.SideData{
		Name: e.Name,
		Data: sides,
	}
}

// GetData Func
func (e *PercentagePriceOscillator) GetData() []PercentagePriceOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}
