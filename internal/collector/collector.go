package collector

import "SysTrace_Agent/internal/data"

type Collector interface {
	Collect() data.Data
}
