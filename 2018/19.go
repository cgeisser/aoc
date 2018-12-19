package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
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

	// crash logging
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	var tombstone ringbuf
	var stop bool
	go func() {
		<-c
		stop = true
		fmt.Println(tombstone)
	}()

	var r reg = reg{1, 0, 0, 0, 0, 0}
	pc := r[ip]
	instructions := 0
	for pc < len(program) && !stop {
		op := &program[pc]
		r[ip] = pc
		savr := r
		op.runOp(&r)
		tombstone.put(fmt.Sprintf("ip=%d %v %v %v\n", pc, savr, op, r))
		pc = r[ip]
		pc++
		instructions++
	}
	fmt.Println("end ", pc, r)
	fmt.Println("ran ", instructions)
}

type reg [6]int

func (i inst) runOp(r *reg) {
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

type ringbuf struct {
	buf   [10]string
	index int
}

func (r *ringbuf) put(s string) {
	r.buf[r.index] = s
	r.index = (r.index + 1) % len(r.buf)
}

func (r ringbuf) String() string {
	var buf bytes.Buffer
	buf.WriteString("crash dump:\n")
	for i := 0; i < len(r.buf); i++ {
		buf.WriteString(r.buf[(i+r.index)%len(r.buf)])
	}
	return buf.String()
}
