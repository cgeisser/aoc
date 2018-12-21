package main

import (
	"bufio"
	"fmt"
	"os"
)

type inst struct {
	op    string
	runop opfunc
	arg   args
}

type args [3]int

type opfunc func(arg args, r *reg)

func getRunner(op string) opfunc {
	switch op {
	case "addi":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] + a[1] }
	case "addr":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] + r[a[1]] }
	case "muli":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] * a[1] }
	case "mulr":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] * r[a[1]] }
	case "bani":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] & a[1] }
	case "banr":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] & r[a[1]] }
	case "bori":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] | a[1] }
	case "borr":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] | r[a[1]] }
	case "setr":
		return func(a args, r *reg) { r[a[2]] = r[a[0]] }
	case "seti":
		return func(a args, r *reg) { r[a[2]] = a[0] }
	case "gtir":
		return func(a args, r *reg) {
			if a[0] > r[a[1]] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	case "gtri":
		return func(a args, r *reg) {
			if r[a[0]] > a[1] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	case "gtrr":
		return func(a args, r *reg) {
			if r[a[0]] > r[a[1]] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	case "eqir":
		return func(a args, r *reg) {
			if a[0] == r[a[1]] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	case "eqri":
		return func(a args, r *reg) {
			if r[a[0]] == a[1] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	case "eqrr":
		return func(a args, r *reg) {
			if r[a[0]] == r[a[1]] {
				r[a[2]] = 1
			} else {
				r[a[2]] = 0
			}
		}
	default:
		panic("bad opcode")
	}
}

func parseinst(s string) (inst, error) {
	i := inst{}
	_, ok := fmt.Sscanf(s, "%s %d %d %d",
		&i.op, &i.arg[0], &i.arg[1], &i.arg[2])
	i.runop = getRunner(i.op)
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

	mininst := -1
	for zero := 0; ; {
		var r reg = reg{zero, 0, 0, 0, 0, 0}
		pc := r[ip]
		instructions := 0
		for pc < len(program) {
			op := &program[pc]
			r[ip] = pc
			op.runOp(&r)
			pc = r[ip]
			pc++
			instructions++
		}
		fmt.Print("zero ", zero)
		fmt.Println("end ", pc, r)
		fmt.Println("ran ", instructions)
		if mininst == -1 || instructions < mininst {
			mininst = instructions
			fmt.Println("new minimum! ", mininst)
		}
		break
	}
}

type reg [6]int

func (i inst) runOp(r *reg) {
	if i.runop != nil {
		i.runop(i.arg, r)
		return
	}
	fmt.Println("going the slow way")
	//	fmt.Printf("  running %v %v %v\n", op, arg, r)
	if i.op[:2] == "gt" || i.op[:2] == "eq" {
		i.compOp(r)
		return
	}

	if i.op[:3] == "set" {
		out := &r[i.arg[2]]
		switch i.op[3] {
		case 'r':
			*out = r[i.arg[0]]
		case 'i':
			*out = i.arg[0]
		}
		return
	}

	var x, y int
	switch i.op[3] {
	case 'r':
		x, y = r[i.arg[0]], r[i.arg[1]]
	case 'i':
		x, y = r[i.arg[0]], i.arg[1]
	default:
		panic("bad register value")
	}
	out := &r[i.arg[2]]
	switch i.op[:3] {
	case "add":
		*out = x + y
	case "mul":
		*out = x * y
	case "ban":
		*out = x & y
	case "bor":
		*out = x | y
	default:
		fmt.Println(i.op[:3])
		panic("bad opcode")
	}
}

func (i inst) compOp(r *reg) {
	var x, y int
	switch i.op[2:4] {
	case "rr":
		x, y = r[i.arg[0]], r[i.arg[1]]
	case "ri":
		x, y = r[i.arg[0]], i.arg[1]
	case "ir":
		x, y = i.arg[0], r[i.arg[1]]
	default:
		panic("bad register value")
	}

	var val bool
	switch i.op[:2] {
	case "gt":
		val = (x > y)
	case "eq":
		val = (x == y)
	}
	if val {
		r[i.arg[2]] = 1
	} else {
		r[i.arg[2]] = 0
	}
}
