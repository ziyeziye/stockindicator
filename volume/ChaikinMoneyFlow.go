package volume

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Chaikin Money Flow (CMF) 蔡金资金流量是用于在一段时间内衡量资金流量的技术分析指标。
// 资金流量(Marc Chaikin创立的一个概念)是用于衡量单一期间证券的买卖压力的指标。
// CMF在用户指定的回溯期内对资金流量进行加总。 任何回溯期都可使用，但最受欢迎的设定是20或21天。
//
// 1、一般而言，CMF大于零，市场处于牛市，CMF小于零，市场处于熊市。
// 2、CMF大于零（或小于零）的时间长短也值得注意。停留时间越长，趋势越可靠。
// 3、CMF可以结合趋势线及支撑线、阻力线突破进行分析。
// 4、CMF与价格之间的背离具有重要意义，通常预示着行情即将转变。
//
// Money Flow Multiplier = ((Closing - Low) - (High - Closing)) / (High - Low)
// Money Flow Volume = Money Flow Multiplier * Volume
// Chaikin Money Flow = Sum(20, Money Flow Volume) / Sum(20, Volume)
type ChaikinMoneyFlow struct {
	Name   string
	Period int
	data   []ChaikinMoneyFlowData
	kline  utils.Klines
}

// ChaikinMoneyFlowData
type ChaikinMoneyFlowData struct {
	Time  time.Time
	Value float64
}

// NewChaikinMoneyFlow new Func
func NewChaikinMoneyFlow(list utils.Klines, period int) *ChaikinMoneyFlow {
	m := &ChaikinMoneyFlow{
		Name:   "ChaikinMoneyFlow",
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultChaikinMoneyFlow new Func
func NewDefaultChaikinMoneyFlow(list utils.Klines) *ChaikinMoneyFlow {
	return NewChaikinMoneyFlow(list, 20)
}

// Calculation Func
func (e *ChaikinMoneyFlow) Calculation() *ChaikinMoneyFlow {

	period := e.Period
	var high, low, closing, volume []float64
	for _, v := range e.kline {
		high = append(high, v.High)
		low = append(low, v.Low)
		closing = append(closing, v.Close)
		volume = append(volume, v.Volume)
	}

	moneyFlowMultiplier := utils.Divide(
		utils.Subtract(utils.Subtract(closing, low), utils.Subtract(high, closing)),
		utils.Subtract(high, low))

	moneyFlowVolume := utils.Multiply(moneyFlowMultiplier, volume)

	cmf := utils.Divide(
		utils.Sum(period, moneyFlowVolume),
		utils.Sum(period, volume))

	for i := 0; i < len(cmf); i++ {
		e.data = append(e.data, ChaikinMoneyFlowData{
			Time:  e.kline[i].Time,
			Value: cmf[i],
		})
	}

	return e
}

// AnalysisSide Func
// func (e *ChaikinMoneyFlow) AnalysisSide() utils.SideData {
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
func (e *ChaikinMoneyFlow) GetData() []ChaikinMoneyFlowData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}
