package cpu

import (
	"fmt"

	"alde.nu/emulator6502/memory"
)

type CPU struct {
	PC memory.Word // program counter
	SP memory.Word // stack pointer

	A, X, Y memory.Byte // registers
	PS      memory.Byte // processor status
	/// [ 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 ]
	/// [ N | V | U | B | D | I | Z | C ]

	// C carry
	// Z zero
	// I irq disable
	// D decimal mode
	// B break command
	// U unused
	// V overflow
	// N negative

	cycles uint32
}

// OpCodes
const (
	INS_LDA_IM  = 0xA9
	INS_LDA_ZP  = 0xA5
	INS_LDA_ZPX = 0xB5

	INS_JSR = 0x20
)

func (c *CPU) Reset(mem *memory.Memory) {
	c.PC = 0xFFFC
	c.SP = 0x0100

	c.A = 0
	c.X = 0
	c.Y = 0

	c.PS = 0b00000000

	mem.Initialize()
}

func (c *CPU) Fetch(mem *memory.Memory) memory.Byte {
	data := mem.Read(c.PC)
	c.PC++
	c.cycles--

	return data
}

func (c *CPU) FetchWord(mem *memory.Memory) memory.Word {
	// 6502 is little-endian
	data := mem.Read(c.PC)
	c.PC++

	data2 := memory.Word(mem.Read(c.PC)) << 8
	c.PC++

	c.cycles -= 2

	return memory.Word(uint16(data) | uint16(data2))
}

func (c *CPU) Read(address memory.Word, mem *memory.Memory) memory.Byte {
	data := mem.Read(address)
	c.cycles--
	return data
}

func (c *CPU) LDAUpdateFlags() {
	// Set Zero flag if A register is 0
	if c.A == 0 { //    NVUBDIZC
		c.PS = c.PS | 0b00000010
	}
	// Set Negative flag if bit 7 of register A is set
	if c.A&0b10000000 > 0 {
		c.PS = c.PS | 0b10000000
	}
}

func (c *CPU) Execute(cycles uint32, mem *memory.Memory) {
	c.cycles = cycles
	for c.cycles > 0 {
		ins := c.Fetch(mem)
		switch ins {
		case INS_LDA_IM:
			value := c.Fetch(mem)
			c.A = value
			c.LDAUpdateFlags()
		case INS_LDA_ZP:
			zpa := memory.Word(c.Fetch(mem))
			c.A = c.Read(zpa, mem)
			c.LDAUpdateFlags()
		case INS_LDA_ZPX:
			zpa := c.Fetch(mem)
			zpa += c.X
			c.cycles--
			c.A = c.Read(memory.Word(zpa), mem)

		case INS_JSR:
			sra := c.FetchWord(mem)
			mem.WriteWord(c.PC-1, c.SP)
			c.PC = sra
			c.cycles -= 3
		default:
			fmt.Printf("invalid instruction %+v", ins)
		}
	}
}
