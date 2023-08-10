package memory

type Byte uint8
type Word uint16

const MAX_MEMORY = 1024 * 64

type Memory struct {
	Data [MAX_MEMORY]Byte
}

func (m *Memory) Read(address Word) Byte {
	return m.Data[address]
}
func (m *Memory) WriteWord(address Word, value Word) {
	m.Data[address] = Byte(value & 0xFF)
	m.Data[address+1] = Byte(value >> 8)

}

func (m *Memory) Initialize() {
	for mem := 0; mem < MAX_MEMORY; mem++ {
		m.Data[mem] = 0
	}
}
