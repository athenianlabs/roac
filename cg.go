package main

import "fmt"

func write(s string) {
	OutFile.WriteString(s)
}

func writef(s string, args ...interface{}) {
	write(fmt.Sprintf(s, args...))
}

// List of available registers
// and their names
var freereg = [4]int{}
var reglist = [4]string{"%r8", "%r9", "%r10", "%r11"}
var breglist = [4]string{"%r8b", "%r9b", "%r10b", "%r11b"}

// Set all registers as available
func freeall_registers() {
	for i := range freereg {
		freereg[i] = 1
	}
}

// Allocate a free register. Return the number of
// the register. Die if no available registers.
func alloc_register() int {
	for i := 0; i < 4; i++ {
		if freereg[i] == 1 {
			freereg[i] = 0
			return i
		}
	}
	fatal("out of registers\n")
	return 0
}

// Return a register to the list of available registers.
// Check to see if it's not already there.
func free_register(reg int) {
	if freereg[reg] != 0 {
		fatal("Error trying to free register %d\n", reg)
	}
	freereg[reg] = 1
}

// Print out the assembly preamble
func cgpreamble() {
	freeall_registers()
	write("\t.text\n")
	write(".LC0:\n")
	write("\t.string\t\"%d\\n\"\n")
	write("printint:\n")
	write("\tpushq\t%rbp\n")
	write("\tmovq\t%rsp, %rbp\n")
	write("\tsubq\t$16, %rsp\n")
	write("\tmovl\t%edi, -4(%rbp)\n")
	write("\tmovl\t-4(%rbp), %eax\n")
	write("\tmovl\t%eax, %esi\n")
	write("\tleaq	.LC0(%rip), %rdi\n")
	write("\tmovl	$0, %eax\n")
	write("\tcall	printf@PLT\n")
	write("\tnop\n")
	write("\tleave\n")
	write("\tret\n")
	write("\n")
}

// Print out the assembly postamble
func cgpostamble() {
	write("\tmovl	$0, %eax\n")
	write("\tpopq	%rbp\n")
	write("\tret\n")
}

// Load an integer literal value into a register.
// Return the number of the register
func cgloadint(value int) int {
	// Get a new register
	r := alloc_register()
	// Print out the code to initialise it
	writef("\tmovq\t$%d, %s\n", value, reglist[r])
	return r
}

// Add two registers together and return
// the number of the register with the result
func cgadd(r1, r2 int) int {
	writef("\taddq\t%s, %s\n", reglist[r1], reglist[r2])
	free_register(r1)
	return r2
}

// Subtract the second register from the first and
// return the number of the register with the result
func cgsub(r1, r2 int) int {
	writef("\tsubq\t%s, %s\n", reglist[r2], reglist[r1])
	free_register(r2)
	return r1
}

// Multiply two registers together and return
// the number of the register with the result
func cgmul(r1, r2 int) int {
	writef("\timulq\t%s, %s\n", reglist[r1], reglist[r2])
	free_register(r1)
	return r2
}

// Divide the first register by the second and
// return the number of the register with the result
func cgdiv(r1, r2 int) int {
	writef("\tmovq\t%s,%%rax\n", reglist[r1])
	write("\tcqo\n")
	writef("\tidivq\t%s\n", reglist[r2])
	writef("\tmovq\t%%rax,%s\n", reglist[r1])
	free_register(r2)
	return r1
}

// Call printint() with the given register
func cgprintint(r int) {
	writef("\tmovq\t%s, %%rdi\n", reglist[r])
	write("\tcall\tprintint\n")
	free_register(r)
}

// Load a value from a variable into a register.
// Return the number of the register
func cgloadglob(sym *Symbol) int {
	// Get a new register
	r := alloc_register()
	// Print out the code to initialize it
	if sym.t == NodeInt {
		writef("\tmovq\t%s(%%rip), %s\n", sym.name, reglist[r])
	} else {
		writef("\tmovzbq\t%s(%%rip), %s\n", sym.name, reglist[r])
	}
	return r
}

// Store a register's value into a variable
func cgstorglob(r int, sym *Symbol) int {
	if sym.t == NodeInt {
		writef("\tmovq\t%s, %s(%%rip)\n", reglist[r], sym.name)
	} else {
		writef("\tmovb\t%s, %s(%%rip)\n", breglist[r], sym.name)
	}
	return r
}

// Generate a global symbol
func cgglobsym(sym *Symbol) {
	if sym.t == NodeInt {
		writef("\t.comm\t%s,8,8\n", sym.name)
	} else {
		writef("\t.comm\t%s,1,1\n", sym.name)
	}
}

// Widen the value in the register from the old
// to the new type, and return a register with
// this new value
func cgwiden(reg int, oldtype, newtype NodeType) int {
	// Nothing to do
	return reg
}

// List of comparison instructions,
// in AST order: A_EQ, A_NE, A_LT, A_GT, A_LE, A_GE
var cmplist = map[OpType]string{
	OpEqual:              "sete",
	OpNotEqual:           "setne",
	OpLessThan:           "setl",
	OpGreaterThan:        "setg",
	OpLessThanOrEqual:    "setle",
	OpGreaterThanOrEqual: "setge",
}

// Compare two registers and set if true.
func cgcompare_and_set(ASTop OpType, r1, r2 int) int {
	// Check the range of the AST operation
	op, ok := cmplist[ASTop]
	if !ok {
		fatal("dad AST Op in cgcompare_and_set()\n")
	}
	writef("\tcmpq\t%s, %s\n", reglist[r2], reglist[r1])
	writef("\t%s\t%s\n", op, breglist[r2])
	writef("\tmovzbq\t%s, %s\n", breglist[r2], reglist[r2])
	free_register(r1)
	return (r2)
}

// List of inverted jump instructions,
// in AST order: A_EQ, A_NE, A_LT, A_GT, A_LE, A_GE
var invertedcmplist = map[OpType]string{
	OpEqual:              "jne",
	OpNotEqual:           "je",
	OpLessThan:           "jge",
	OpGreaterThan:        "jle",
	OpLessThanOrEqual:    "jg",
	OpGreaterThanOrEqual: "jl",
}

// Compare two registers and jump if false.
func cgcompare_and_jump(ASTop OpType, r1, r2, label int) int {
	// Check the range of the AST operation
	op, ok := invertedcmplist[ASTop]
	if !ok {
		fatal("bad AST Op in cgcompare_and_jump()\n")
	}
	writef("\tcmpq\t%s, %s\n", reglist[r2], reglist[r1])
	writef("\t%s\tL%d\n", op, label)
	freeall_registers()
	return NoReg
}

// Generate a label
func cglabel(l int) {
	writef("L%d:\n", l)
}

// Generate a jump to a label
func cgjump(l int) {
	writef("\tjmp\tL%d\n", l)
}

// Print out a function preamble
func cgfuncpreamble(name string) {
	write("\t.text\n")
	writef("\t.globl\t%s\n", name)
	writef("\t.type\t%s, @function\n", name)
	writef("%s:\n", name)
	write("\tpushq\t%rbp\n")
	write("\tmovq\t%rsp, %rbp\n")
}

// Print out a function postamble
func cgfuncpostamble() {
	write("\tmovl $0, %eax\n")
	write("\tpopq     %rbp\n")
	write("\tret\n")
}
