package main

import (
	"fmt"

	"alde.nu/emulator6502/cpu"
	"alde.nu/emulator6502/memory"
)

func main() {
	memory := &memory.Memory{}
	c := &cpu.CPU{}
	c.Reset(memory)
	memory.Data[0xFFFC] = cpu.INS_JSR
	memory.Data[0xFFFD] = 0x42
	memory.Data[0xFFFE] = 0x42
	memory.Data[0x4242] = cpu.INS_LDA_IM
	memory.Data[0x4243] = 0x84

	c.Execute(9, memory)
	fmt.Printf("%#+v\n", c)

	fmt.Println("cpu register A should be 0x84")
	if c.A == 0x84 {
		fmt.Println(" 😀 and it is!")
	} else {
		fmt.Println(" ☹️ something went wrong ")
	}

}
