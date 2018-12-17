package main

import (
	"fmt"
	"os"
	"bufio"
)

type sample struct {
	before, after reg
	i inst
	matches []string
}

type inst struct {
	op int
	arg [3]int
}

func parseinst(s string) inst {
	i := inst{}
	fmt.Sscanf(s, "%d %d %d %d",
		&i.op, &i.arg[0], &i.arg[1], &i.arg[2])
	return i
}

func main() {
	f, _ := os.Open("16.data")
	input := bufio.NewScanner(f)

	samples := make([]*sample, 0)
	for input.Scan() {
		s := &sample{}
		if len(input.Text()) == 0 || input.Text()[0] != 'B' {
			break
		}
		fmt.Sscanf(input.Text(),
			"Before: [%d, %d, %d, %d]", &s.before[0],
			&s.before[1], &s.before[2], &s.before[3])
		input.Scan()
		s.i = parseinst(input.Text())
		input.Scan()
		fmt.Sscanf(input.Text(),
			"After:  [%d, %d, %d, %d]", &s.after[0],
			&s.after[1], &s.after[2], &s.after[3])
		input.Scan()
		fmt.Println(s)
		samples = append(samples, s)
		//break
	}
	threeses := 0
	for _, s := range samples {
		s.matches = opMatches(s.i, s.before, s.after)
		fmt.Println(*s, s.matches)
		if len(s.matches) >= 3 {
			threeses++
		}
	}
	fmt.Println("threeses: ", threeses)

	opmap := [16]map[string]bool{}
	
	for _, s := range samples {
		if opmap[s.i.op] == nil {
			opmap[s.i.op] = make(map[string]bool)
		}
		for _, op := range s.matches {
			opmap[s.i.op][op]=true
		}
	}

	opcode := [16]string{}
	for ; !allfound(opcode) ; {
		oneloc, op := findOne(opmap,opcode)
		opcode[oneloc] = op
		for i, v := range opmap {
			if i == oneloc {
				continue
			}
			delete(v, op)
		}
		//fmt.Println(oneloc, op, opmap)
	}
	fmt.Println("Final opcodes ", opcode)

	program := make([]inst, 0)
	for input.Scan() {
		if len(input.Text()) == 0 {
			continue
		}
		program = append(program, parseinst(input.Text()))
	}
	fmt.Println("program", program)

	var r reg
	for c, i := range program {
		fmt.Println(" ", c, i, r)
		r = runOp(opcode[i.op], i.arg, r)
	}
	fmt.Println(r)
}

func findOne(opmap [16]map[string]bool, opcode [16]string) (int, string){
	for i, v := range opmap {
		if len(v) == 1 && opcode[i] == "" {
			for k := range v {
				return i, k
			}
		}
	}
	panic("couldn't find singleton")
}
		
func allfound(opcode [16]string) bool {
	for _, v := range opcode {
		if v == ""  {
			return false
		}
	}
	return true
}

func opMatches(i inst, in, exp reg) []string {
	matches := make([]string, 0)
	for _, op := range getAllOps() {
		out := runOp(op, i.arg, in)
		if out == exp {
			matches = append(matches, op)
		}
	}
	return matches
}

func getAllOps() [16]string {
	return [16]string{ "addr", "addi", "mulr", "muli", "banr", "bani", "borr",
		"bori", "setr", "seti", "gtir", "gtri", "gtrr", "eqir",
		"eqri", "eqrr" }
}

type reg [4]int
func runOp(op string, arg [3]int, r reg) reg {
	if op[:2] == "gt" || op[:2] == "eq" {
		return compOp(op, arg, r)
	}

        var x, y int
	switch op[3] {
	case 'r': x, y = r[arg[0]], r[arg[1]]
	case 'i': x, y = r[arg[0]], arg[1]
	default: panic("bad register value")
	}
	out := &r[arg[2]]	
	switch op[:3] {
	case "add": *out = x+y
	case "mul": *out = x*y
	case "ban": *out = x & y
	case "bor": *out = x|y
	case "set":
		switch op[3] {
		case 'r': *out = r[arg[0]]
			case'i': *out = arg[0]
		}
	default:
		fmt.Println(op[:3])
		panic("bad opcode")
	}
	return r	
}

func compOp(op string, arg [3]int, r reg) reg{
	var x, y int
	switch op[2:4] {
	case "rr": x, y = r[arg[0]], r[arg[1]]
	case "ri": x, y = r[arg[0]], arg[1]
	case "ir": x, y = arg[0], r[arg[1]]
	default: panic("bad register value")
	}
	
	var val bool
	switch op[:2] {
	case "gt": val = (x > y)
		case "eq": val = (x == y)
	}
	if val {
		r[arg[2]] = 1
	} else {
		r[arg[2]] = 0
	}
	return r
}


