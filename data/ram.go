package data

type Memory struct {
	Total       uint64
	Used        uint64
	Available   uint64
	UsedPercent float64
	Model       string
	Speed       uint64
}

func (m Memory) GetTotal() uint64 {
	return m.Total
}

func (m Memory) GetUsed() uint64 {
	return m.Used
}

func (m Memory) GetAvailable() uint64 {
	return m.Available
}

func (m Memory) GetUsedPercent() float64 {
	return m.UsedPercent
}

func (m Memory) GetModel() string {
	return m.Model
}

func (m Memory) GetSpeed() uint64 {
	return m.Speed
}

func (m *Memory) SetTotal(total uint64) {
	m.Total = total
}

func (m *Memory) SetUsed(used uint64) {
	m.Used = used
}

func (m *Memory) SetAvailable(available uint64) {
	m.Available = available
}

func (m *Memory) SetUsedPercent(usedPercent float64) {
	m.UsedPercent = usedPercent
}

func (m *Memory) SetModel(model string) {
	m.Model = model
}

func (m *Memory) SetSpeed(speed uint64) {
	m.Speed = speed
}
