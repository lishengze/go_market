package monitorStruct

import "time"

type MonitorChannel struct {
	DepthChan chan *MonitorAtom
	TradeChan chan *MonitorAtom
	KlineChan chan *MonitorAtom
}

func NewMonitorChannel() *MonitorChannel {
	return &MonitorChannel{
		DepthChan: make(chan *MonitorAtom),
		TradeChan: make(chan *MonitorAtom),
		KlineChan: make(chan *MonitorAtom),
	}
}

type MonitorMarketData struct {
	depth_cache_map map[string](*MonitorAtom)
	trade_cache_map map[string](*MonitorAtom)
	kline_cache_map map[string](*MonitorAtom)

	RateParam    float64
	InitDeadLine int64
	CheckSecs    int64
	MonitorChan  *MonitorChannel
}

func NewMonitorMarketData(rate_param float64, init_dead_line int64, check_secs int64, monitor_channel *MonitorChannel) *MonitorMarketData {
	return &MonitorMarketData{
		depth_cache_map: make(map[string]*MonitorAtom),
		trade_cache_map: make(map[string]*MonitorAtom),
		kline_cache_map: make(map[string]*MonitorAtom),
		RateParam:       rate_param,
		InitDeadLine:    init_dead_line,
		CheckSecs:       check_secs,
		MonitorChan:     monitor_channel,
	}
}

func (m *MonitorMarketData) StartCheck() {
	timer := time.Tick(time.Duration(m.CheckSecs) * time.Second)

	for {
		select {
		case <-timer:
			m.CheckAll()
		}
	}
}

func (m *MonitorMarketData) CheckAll() {
	for _, monitor_atom := range m.depth_cache_map {
		if !monitor_atom.IsAlive() {
			m.MonitorChan.DepthChan <- monitor_atom
		}
	}

	for _, monitor_atom := range m.trade_cache_map {
		if !monitor_atom.IsAlive() {
			m.MonitorChan.TradeChan <- monitor_atom
		}
	}

	for _, monitor_atom := range m.kline_cache_map {
		if !monitor_atom.IsAlive() {
			m.MonitorChan.KlineChan <- monitor_atom
		}
	}
}

func (m *MonitorMarketData) UpdateDepth(symbol string) {
	if _, ok := m.depth_cache_map[symbol]; !ok {
		m.depth_cache_map[symbol] = NewMonitorAtom(symbol, m.RateParam, int64(m.InitDeadLine))
	}

	m.depth_cache_map[symbol].Update()
}

func (m *MonitorMarketData) UpdateTrade(symbol string) {
	if _, ok := m.trade_cache_map[symbol]; !ok {
		m.trade_cache_map[symbol] = NewMonitorAtom(symbol, m.RateParam, int64(m.InitDeadLine))
	}

	m.trade_cache_map[symbol].Update()
}

func (m *MonitorMarketData) UpdateKline(symbol string) {
	if _, ok := m.kline_cache_map[symbol]; !ok {
		m.kline_cache_map[symbol] = NewMonitorAtom(symbol, m.RateParam, int64(m.InitDeadLine))
	}

	m.kline_cache_map[symbol].Update()
}
