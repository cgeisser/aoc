package main

import (
	"bufio"
	"fmt"
	"os"
)

type inst struct {
	op  string
	arg [3]int
}

func parseinst(s string) (inst, error) {
	i := inst{}
	_, ok := fmt.Sscanf(s, "%s %d %d %d",
		&i.op, &i.arg[0], &i.arg[1], &i.arg[2])
	return i, ok
}

func main() {
	f, _ := os.Open(os.Args[1])
	input := bufio.NewScanner(f)

	program := make([]inst, 0)
	var ip int
	input.Scan()
	_, ok := fmt.Sscanf(input.Text(), "#ip %d", &ip)
	if ok != nil {
		panic(fmt.Sprintf("bad ip %v %v", ok, input.Text()))
	}
	for input.Scan() {
		if len(input.Text()) == 0 {
			continue
		}
		i, ok := parseinst(input.Text())
		if ok != nil {
			panic(fmt.Sprintf("bad line: %v %v", ok, input.Text()))
		}
		program = append(program, i)
	}
	fmt.Println("ip ", ip)
	fmt.Println("program", program)
	var r reg = reg{0, 0, 0, 0, 0, 0}
	pc := r[ip]
	for pc < len(program) {
		op := program[pc]
		r[ip] = pc
		//savr := r
		r = runOp(op.op, op.arg, r)
		//fmt.Printf("ip=%d %v %v %v\n", pc, savr, op, r)
		pc = r[ip]
		pc++
		//		break
	}
	fmt.Println("end ", pc, r)
}

type reg [6]int

func runOp(op string, arg [3]int, r reg) reg {
	//	fmt.Printf("  running %v %v %v\n", op, arg, r)
	if op[:2] == "gt" || op[:2] == "eq" {
		return compOp(op, arg, r)
	}

	if op[:3] == "set" {
		out := &r[arg[2]]
		switch op[3] {
		case 'r':
			*out = r[arg[0]]
		case 'i':
			*out = arg[0]
		}
		return r
	}

	var x, y int
	switch op[3] {
	case 'r':
		x, y = r[arg[0]], r[arg[1]]
	case 'i':
		x, y = r[arg[0]], arg[1]
	default:
		panic("bad register value")
	}
	out := &r[arg[2]]
	switch op[:3] {
	case "add":
		*out = x + y
	case "mul":
		*out = x * y
	case "ban":
		*out = x & y
	case "bor":
		*out = x | y
	default:
		fmt.Println(op[:3])
		panic("bad opcode")
	}
	return r
}

func compOp(op string, arg [3]int, r reg) reg {
	var x, y int
	switch op[2:4] {
	case "rr":
		x, y = r[arg[0]], r[arg[1]]
	case "ri":
		x, y = r[arg[0]], arg[1]
	case "ir":
		x, y = arg[0], r[arg[1]]
	default:
		panic("bad register value")
	}

	var val bool
	switch op[:2] {
	case "gt":
		val = (x > y)
	case "eq":
		val = (x == y)
	}
	if val {
		r[arg[2]] = 1
	} else {
		r[arg[2]] = 0
	}
	return r
}
