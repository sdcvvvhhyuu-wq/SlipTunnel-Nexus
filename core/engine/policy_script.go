package engine

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

// =====================================================================
// LIGHTWEIGHT RUNTIME POLICY VIRTUAL MACHINE (WASM/LUA ABSTRACTION)
// =====================================================================
// Executes dynamic routing rules and endpoint mutations without rebuilding the binary.

type OpCode int

const (
	OpPush OpCode = iota
	OpPop
	OpEq
	OpGt
	OpJmpIf
	OpSetRoute
	OpHalt
)

type Instruction struct {
	Op      OpCode
	Operand string
}

type PolicyVM struct {
	mu           sync.Mutex
	stack        []int
	instructions []Instruction
	pc           int
	ActiveRoute  string
}

func NewPolicyVM() *PolicyVM {
	return &PolicyVM{
		stack:       make([]int, 0),
		ActiveRoute: "DEFAULT",
	}
}

// LoadScript parses a custom assembly-like script into VM instructions.
// Example Script: "PUSH 150; PUSH 200; GT; JMPIF 1; SETROUTE CDN_RELAY; HALT"
func (vm *PolicyVM) LoadScript(script string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.instructions = []Instruction{}
	vm.pc = 0

	tokens := strings.Split(script, ";")
	for _, token := range tokens {
		parts := strings.Fields(strings.TrimSpace(token))
		if len(parts) == 0 {
			continue
		}

		inst := Instruction{Operand: ""}
		if len(parts) > 1 {
			inst.Operand = parts[1]
		}

		switch strings.ToUpper(parts[0]) {
		case "PUSH":
			inst.Op = OpPush
		case "POP":
			inst.Op = OpPop
		case "EQ":
			inst.Op = OpEq
		case "GT":
			inst.Op = OpGt
		case "JMPIF":
			inst.Op = OpJmpIf
		case "SETROUTE":
			inst.Op = OpSetRoute
		case "HALT":
			inst.Op = OpHalt
		default:
			return errors.New("unknown opcode: " + parts[0])
		}
		vm.instructions = append(vm.instructions, inst)
	}
	return nil
}

func (vm *PolicyVM) Execute() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	for vm.pc < len(vm.instructions) {
		inst := vm.instructions[vm.pc]
		switch inst.Op {
		case OpPush:
			val, _ := strconv.Atoi(inst.Operand)
			vm.stack = append(vm.stack, val)
		case OpPop:
			if len(vm.stack) > 0 {
				vm.stack = vm.stack[:len(vm.stack)-1]
			}
		case OpEq:
			if len(vm.stack) >= 2 {
				a := vm.stack[len(vm.stack)-1]
				b := vm.stack[len(vm.stack)-2]
				vm.stack = vm.stack[:len(vm.stack)-2]
				if a == b {
					vm.stack = append(vm.stack, 1)
				} else {
					vm.stack = append(vm.stack, 0)
				}
			}
		case OpGt:
			if len(vm.stack) >= 2 {
				a := vm.stack[len(vm.stack)-1]
				b := vm.stack[len(vm.stack)-2]
				vm.stack = vm.stack[:len(vm.stack)-2]
				if b > a {
					vm.stack = append(vm.stack, 1)
				} else {
					vm.stack = append(vm.stack, 0)
				}
			}
		case OpJmpIf:
			if len(vm.stack) > 0 {
				cond := vm.stack[len(vm.stack)-1]
				vm.stack = vm.stack[:len(vm.stack)-1]
				if cond == 1 {
					offset, _ := strconv.Atoi(inst.Operand)
					vm.pc += offset
					continue
				}
			}
		case OpSetRoute:
			vm.ActiveRoute = inst.Operand
		case OpHalt:
			return nil
		}
		vm.pc++
	}
	return nil
}
