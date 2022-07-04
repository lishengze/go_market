package monitorStruct

import (
	"market_server/common/datastruct"
	"market_server/common/util"
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
}

func NewMonitorAtom(rate_param float64, init_dead_line int64, data_cout int) *MonitorAtom {
	return &MonitorAtom{
		RateParam:        rate_param,
		InitDeadLine:     init_dead_line,
		data_count:       0,
		first_time:       0,
		last_update_time: 0,
		sum_time:         0,
		max_time:         0,
		ave_time:         0,
	}
}

func (m *MonitorAtom) SetCurTime(cur_time int64) {

	m.first_time = cur_time

	cur_time = cur_time / datastruct.NANO_PER_MILL
	delta_time := cur_time - m.last_update_time

	m.last_update_time = delta_time

	m.sum_time += delta_time
	m.data_count++

	m.ave_time = m.sum_time / m.data_count
	if delta_time > m.max_time {
		m.max_time = delta_time
	}

}

func (m *MonitorAtom) TimeLimit() int64 {
	cur_time := util.UTCNanoTime() / datastruct.NANO_PER_MILL
	if cur_time-m.first_time > datastruct.SECS_PER_DAY*datastruct.MILL_PER_SECS {

	} else {

	}

	return 0
}

func (m *MonitorAtom) IsAlive() bool {
	cur_time := util.UTCNanoTime() / datastruct.NANO_PER_MILL
	delta_time := cur_time - m.last_update_time

	if delta_time > m.TimeLimit() {
		return false
	} else {
		return true
	}
}
