package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.x[rd] = cpu.xint(cpu.x[rs1] + imm)

	case 0b_001: // slli
		if imm < cpu.xlen {
			cpu.x[rd] = cpu.xint(cpu.x[rs1] << imm)
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_010: // slti
		if cpu.x[rs1] < imm {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.x[rs1]) < cpu.xuint(imm) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_100: // xori
		cpu.x[rd] = cpu.x[rs1] ^ imm

	case 0b_101:
		if imm < cpu.xlen { // srli
			cpu.x[rd] = cpu.xint(int(cpu.xuint(cpu.x[rs1]) >> cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.xlen { // srai
			cpu.x[rd] = cpu.x[rs1] >> imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_110: // ori
		cpu.x[rd] = cpu.x[rs1] | imm

	case 0b_111: // andi
		cpu.x[rd] = cpu.x[rs1] & imm
	}
}

func (cpu *CPU) execComputeR(f7, rs2, rs1, f3, rd int) {
	if f7 == 1 {
		cpu.execComputeM(rs2, rs1, f3, rd)
		return
	}

	op := bit(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000: // add
		cpu.x[rd] = cpu.xint(cpu.x[rs1] + cpu.x[rs2])

	case 0b_1_000: // sub
		cpu.x[rd] = cpu.xint(cpu.x[rs1] - cpu.x[rs2])

	case 0b_001: // sll
		cpu.x[rd] = cpu.xint(cpu.x[rs1] << (cpu.x[rs2] & (cpu.xlen - 1)))

	case 0b_010: // slt
		if cpu.x[rs1] < cpu.x[rs2] {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.x[rs1]) < cpu.xuint(cpu.x[rs2]) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_100: // xor
		cpu.x[rd] = cpu.x[rs1] ^ cpu.x[rs2]

	case 0b_101: // srl
		cpu.x[rd] = cpu.xint(int(cpu.xuint(cpu.x[rs1]) >> cpu.xuint(cpu.x[rs2]&(cpu.xlen-1))))

	case 0b_1_101: // sra
		cpu.x[rd] = cpu.x[rs1] >> (cpu.x[rs2] & (cpu.xlen - 1))

	case 0b_110: // or
		cpu.x[rd] = cpu.x[rs1] | cpu.x[rs2]

	case 0b_111: // and
		cpu.x[rd] = cpu.x[rs1] & cpu.x[rs2]

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
