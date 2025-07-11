package rv

import (
	"fmt"
	"github.com/deadsy/rvda"
	"strings"
)

var (
	cycleCount = 0
	trapCount  = 0
	trace      [][4]uint
)

func (cpu *CPU) debugStep() bool {
	cycleCount++

	pc := cpu.pc
	oldRegs := cpu.x

	opcode := cpu.innerStep()

	entry := [4]uint{uint(pc), uint(opcode)}
	for i, val := range cpu.x {
		if val != oldRegs[i] {
			entry[2], entry[3] = uint(i), uint(val)
			break
		}
	}

	trace = append(trace, entry)
	if n := 100; len(trace) == n+1 {
		copy(trace[:n], trace[1:])
		trace = trace[:n]
	}

	if cpu.isTrapped /*|| cycleCount == 1_000_000*/ {
		trapCount++

		//if trapCount == 8 {
		//	fmt.Printf("Cycle: %v, trap: %v\r\n\r\n", cycleCount, trapCount)
		//
		//	debugDump(cpu)
		//
		//	return false
		//}
	}

	return true
}

func debugDump(cpu *CPU) {
	isa, _ := rvda.New(uint(cpu.xlen), rvda.RV64gc)

	for _, entry := range trace {
		fmt.Printf("%v\r\n", disassemble(isa, entry))
	}

	fmt.Printf("\r\nmcause:%x, mepc:%x, mtvec:%x, mstatus:%x\r\n",
		uint(cpu.csr.mcause), uint(cpu.csr.mepc), uint(cpu.csr.mtvec), uint(cpu.csr.mstatus))
	fmt.Printf("scause:%x, sepc:%x, stvec:%x, satp:%x\r\n",
		uint(cpu.csr.scause), uint(cpu.csr.sepc), uint(cpu.csr.stvec), uint(cpu.csr.satp))
	fmt.Printf("priv:%v, medeleg:%x\r\n\r\n", cpu.priv, uint(cpu.csr.medeleg))

	for i := range 16 {
		fmt.Printf("% 5v:%16x      % 5v:%16x\r\n",
			regNames[i], uint(cpu.x[i]), regNames[16+i], uint(cpu.x[16+i]))
	}
}

func disassemble(isa *rvda.ISA, entry [4]uint) string {
	addr, code, reg, regVal := entry[0], entry[1], entry[2], entry[3]

	line := isa.Disassemble(addr, code).String()
	parts := strings.Split(line, "\t")
	ops := strings.Split(parts[1], " ")
	for len(ops) < 2 {
		ops = append(ops, "")
	}

	line = fmt.Sprintf("%-30v %-6v %-16v", parts[0], ops[0], ops[1])

	if entry[2] != 0 {
		line += fmt.Sprintf("// %v=%x", regNames[reg], regVal)
	}

	return line
}

var regNames = []string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}
