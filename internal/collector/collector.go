package collector

import (
	"SysTrace_Agent/internal/data/static"
)

type Collector interface {
	Collect() static.Data
}
