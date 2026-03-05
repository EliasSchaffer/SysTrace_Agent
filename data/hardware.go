package data

type Hardware struct {
	CPU    CPU
	MEMORY Memory
}

// Getter-Methoden
func (h Hardware) GetCPU() CPU {
	return h.CPU
}

func (h Hardware) GetMemory() Memory {
	return h.MEMORY
}

// Setter-Methoden
func (h *Hardware) SetCPU(cpu CPU) {
	h.CPU = cpu
}

func (h *Hardware) SetMemory(memory Memory) {
	h.MEMORY = memory
}
