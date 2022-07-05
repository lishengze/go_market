package monitorStruct

import (
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync"
)

// 时间转换为毫秒计算
type MonitorAtom struct {
	first_time       int64
	last_update_time int64
	data_count       int64

	sum_time int64
	max_time int64
	ave_time int64

	RateParam    float64
	InitDeadLine int64
	Symbol       string

	update_mutex sync.Mutex
}

func NewMonitorAtom(symbol string, rate_param float64, init_dead_line int64) *MonitorAtom {
	return &MonitorAtom{
		RateParam:        rate_param,
		InitDeadLine:     init_dead_line,
		Symbol:           symbol,
		data_count:       0,
		first_time:       0,
		last_update_time: 0,
		sum_time:         0,
		max_time:         0,
		ave_time:         0,
	}
}

func (m *MonitorAtom) Update() {
	m.update_mutex.Lock()
	defer m.update_mutex.Unlock()

	cur_time := util.UTCNanoTime()

	if m.first_time == 0 {
		m.first_time = cur_time
	}

	delta_time := cur_time - m.last_update_time
	m.last_update_time = cur_time

	m.sum_time += delta_time
	m.data_count++
	m.ave_time = m.sum_time / m.data_count

	if delta_time > m.max_time {
		m.max_time = delta_time
	}

	if cur_time-m.first_time > datastruct.NANO_PER_DAY*14 {
		m.sum_time = 0
		m.data_count = 0
		m.ave_time = 0
		m.first_time = 0
	}
}

func (m *MonitorAtom) TimeLimit() int64 {
	cur_time := util.UTCNanoTime()

	static_time_limit := m.max_time * int64(m.RateParam)

	if static_time_limit > m.InitDeadLine {
		return static_time_limit
	}

	if cur_time-m.first_time > datastruct.NANO_PER_DAY {
		return static_time_limit
	} else {
		return m.InitDeadLine
	}
}

func (m *MonitorAtom) IsAlive() bool {
	m.update_mutex.Lock()
	defer m.update_mutex.Unlock()

	cur_time := util.UTCNanoTime()
	delta_time := cur_time - m.last_update_time

	if delta_time > m.TimeLimit() {
		return false
	} else {
		return true
	}
}
