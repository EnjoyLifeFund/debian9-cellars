// Do not edit. Bootstrap copy of /tmp/go-20170529-7428-4cl5ea/go/src/cmd/asm/internal/arch/mips.go

//line /tmp/go-20170529-7428-4cl5ea/go/src/cmd/asm/internal/arch/mips.go:1
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the
// MIPS (MIPS64) instruction set, to minimize its interaction
// with the core of the assembler.

package arch

import (
	"bootstrap/cmd/internal/obj"
	"bootstrap/cmd/internal/obj/mips"
)

func jumpMIPS(word string) bool {
	switch word {
	case "BEQ", "BFPF", "BFPT", "BGEZ", "BGEZAL", "BGTZ", "BLEZ", "BLTZ", "BLTZAL", "BNE", "JMP", "JAL", "CALL":
		return true
	}
	return false
}

// IsMIPSCMP reports whether the op (as defined by an mips.A* constant) is
// one of the CMP instructions that require special handling.
func IsMIPSCMP(op obj.As) bool {
	switch op {
	case mips.ACMPEQF, mips.ACMPEQD, mips.ACMPGEF, mips.ACMPGED,
		mips.ACMPGTF, mips.ACMPGTD:
		return true
	}
	return false
}

// IsMIPSMUL reports whether the op (as defined by an mips.A* constant) is
// one of the MUL/DIV/REM instructions that require special handling.
func IsMIPSMUL(op obj.As) bool {
	switch op {
	case mips.AMUL, mips.AMULU, mips.AMULV, mips.AMULVU,
		mips.ADIV, mips.ADIVU, mips.ADIVV, mips.ADIVVU,
		mips.AREM, mips.AREMU, mips.AREMV, mips.AREMVU:
		return true
	}
	return false
}

func mipsRegisterNumber(name string, n int16) (int16, bool) {
	switch name {
	case "F":
		if 0 <= n && n <= 31 {
			return mips.REG_F0 + n, true
		}
	case "FCR":
		if 0 <= n && n <= 31 {
			return mips.REG_FCR0 + n, true
		}
	case "M":
		if 0 <= n && n <= 31 {
			return mips.REG_M0 + n, true
		}
	case "R":
		if 0 <= n && n <= 31 {
			return mips.REG_R0 + n, true
		}
	}
	return 0, false
}
