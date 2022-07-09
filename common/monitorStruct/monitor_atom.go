package monitorStruct

import (
	"fmt"
	"market_server/common/datastruct"
	"market_server/common/util"
	"sync"
)

type MonitorData struct {
	Symbol   string
	DataType string
}

type MonitorConfig struct {
	RateParam    float64
	InitDeadLine int64
	CheckSecs    int64
}

// 时间转换为毫秒计算
type MonitorAtom struct {
	first_time       int64
	last_update_time int64
	data_count       int64

	sum_time int64
	max_time int64
	ave_time int64
	lst_time int64

	RateParam    float64
	InitDeadLine int64

	Symbol   string
	DataType string
	MetaInfo string

	InvalidInfo string

	update_mutex sync.Mutex
}

func (m *MonitorAtom) String() string {
	return fmt.Sprintf("f:%s, l:%s, lst: %dus, max: %dus, ave: %d us",
		util.TimeStrFromInt(m.first_time),
		util.TimeStrFromInt(m.last_update_time),
		m.lst_time/datastruct.NANO_PER_MICR,
		m.max_time/datastruct.NANO_PER_MICR,
		m.ave_time/datastruct.NANO_PER_MICR)
}

func NewMonitorAtom(symbol string, DataType string, meta_info string, rate_param float64, init_dead_line int64) *MonitorAtom {
	return &MonitorAtom{
		RateParam:    rate_param,
		InitDeadLine: init_dead_line * datastruct.NANO_PER_SECS,

		Symbol:   symbol,
		DataType: DataType,
		MetaInfo: meta_info,

		data_count:       0,
		first_time:       0,
		last_update_time: 0,
		sum_time:         0,
		max_time:         0,
		ave_time:         0,
		lst_time:         0,
	}
}

func (m *MonitorAtom) Update() {
	m.update_mutex.Lock()
	defer m.update_mutex.Unlock()

	cur_time := util.UTCNanoTime()

	if m.first_time == 0 {
		m.first_time = cur_time
	}

	if m.last_update_time == 0 {
		m.last_update_time = cur_time
		return
	}

	delta_time := cur_time - m.last_update_time
	// logx.Slowf("%s, %s, cur_time: %s, last_update_time:%s, delta_time: %d ", m.MetaInfo, m.Symbol,
	// 	util.TimeStrFromInt(cur_time), util.TimeStrFromInt(m.last_update_time), delta_time)

	m.last_update_time = cur_time

	m.lst_time = delta_time

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

	m.InvalidInfo = fmt.Sprintf("%s.%s, f:%s, l:%s, max: %dms, ave: %d ms;\ndelta: %d, time_limit: %d",
		m.DataType, m.Symbol,
		util.TimeStrFromInt(m.first_time),
		util.TimeStrFromInt(m.last_update_time),
		m.max_time/datastruct.NANO_PER_MILL,
		m.ave_time/datastruct.NANO_PER_MILL,
		delta_time, m.TimeLimit())

	if delta_time > m.TimeLimit() {
		return false
	} else {
		return true
	}
}
