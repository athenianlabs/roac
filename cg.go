package main

import "fmt"

const NoReg = -1

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
	OutFile.WriteString("\t.text\n")
	OutFile.WriteString(".LC0:\n")
	OutFile.WriteString("\t.string\t\"%d\\n\"\n")
	OutFile.WriteString("printint:\n")
	OutFile.WriteString("\tpushq\t%rbp\n")
	OutFile.WriteString("\tmovq\t%rsp, %rbp\n")
	OutFile.WriteString("\tsubq\t$16, %rsp\n")
	OutFile.WriteString("\tmovl\t%edi, -4(%rbp)\n")
	OutFile.WriteString("\tmovl\t-4(%rbp), %eax\n")
	OutFile.WriteString("\tmovl\t%eax, %esi\n")
	OutFile.WriteString("\tleaq	.LC0(%rip), %rdi\n")
	OutFile.WriteString("\tmovl	$0, %eax\n")
	OutFile.WriteString("\tcall	printf@PLT\n")
	OutFile.WriteString("\tnop\n")
	OutFile.WriteString("\tleave\n")
	OutFile.WriteString("\tret\n")
	OutFile.WriteString("\n")
	OutFile.WriteString("\t.globl\tmain\n")
	OutFile.WriteString("\t.type\tmain, @function\n")
	OutFile.WriteString("main:\n")
	OutFile.WriteString("\tpushq\t%rbp\n")
	OutFile.WriteString("\tmovq	%rsp, %rbp\n")
}

// Print out the assembly postamble
func cgpostamble() {
	OutFile.WriteString("\tmovl	$0, %eax\n")
	OutFile.WriteString("\tpopq	%rbp\n")
	OutFile.WriteString("\tret\n")
}

// Load an integer literal value into a register.
// Return the number of the register
func cgloadint(value int) int {
	// Get a new register
	r := alloc_register()
	// Print out the code to initialise it
	OutFile.WriteString(fmt.Sprintf("\tmovq\t$%d, %s\n", value, reglist[r]))
	return r
}

// Add two registers together and return
// the number of the register with the result
func cgadd(r1, r2 int) int {
	OutFile.WriteString(fmt.Sprintf("\taddq\t%s, %s\n", reglist[r1], reglist[r2]))
	free_register(r1)
	return r2
}

// Subtract the second register from the first and
// return the number of the register with the result
func cgsub(r1, r2 int) int {
	OutFile.WriteString(fmt.Sprintf("\tsubq\t%s, %s\n", reglist[r2], reglist[r1]))
	free_register(r2)
	return r1
}

// Multiply two registers together and return
// the number of the register with the result
func cgmul(r1, r2 int) int {
	OutFile.WriteString(fmt.Sprintf("\timulq\t%s, %s\n", reglist[r1], reglist[r2]))
	free_register(r1)
	return r2
}

// Divide the first register by the second and
// return the number of the register with the result
func cgdiv(r1, r2 int) int {
	OutFile.WriteString(fmt.Sprintf("\tmovq\t%s,%%rax\n", reglist[r1]))
	OutFile.WriteString("\tcqo\n")
	OutFile.WriteString(fmt.Sprintf("\tidivq\t%s\n", reglist[r2]))
	OutFile.WriteString(fmt.Sprintf("\tmovq\t%%rax,%s\n", reglist[r1]))
	free_register(r2)
	return r1
}

// Call printint() with the given register
func cgprintint(r int) {
	OutFile.WriteString(fmt.Sprintf("\tmovq\t%s, %%rdi\n", reglist[r]))
	OutFile.WriteString("\tcall\tprintint\n")
	free_register(r)
}

// Load a value from a variable into a register.
// Return the number of the register
func cgloadglob(ident string) int {
	// Get a new register
	r := alloc_register()
	// Print out the code to initialise it
	OutFile.WriteString(fmt.Sprintf("\tmovq\t%s(%%rip), %s\n", ident, reglist[r]))
	return r
}

// Store a register's value into a variable
func cgstorglob(r int, ident string) int {
	OutFile.WriteString(fmt.Sprintf("\tmovq\t%s, %s(%%rip)\n", reglist[r], ident))
	return r
}

// Generate a global symbol
func cgglobsym(sym string) {
	OutFile.WriteString(fmt.Sprintf("\t.comm\t%s,8,8\n", sym))
}

// List of comparison instructions,
// in AST order: A_EQ, A_NE, A_LT, A_GT, A_LE, A_GE
var cmplist = map[NodeType]string{
	NodeEqual:              "sete",
	NodeNotEqual:           "setne",
	NodeLessThan:           "setl",
	NodeGreaterThan:        "setg",
	NodeLessThanOrEqual:    "setle",
	NodeGreaterThanOrEqual: "setge",
}

// Compare two registers and set if true.
func cgcompare_and_set(ASTop NodeType, r1, r2 int) int {
	// Check the range of the AST operation
	op, ok := cmplist[ASTop]
	if !ok {
		fatal("dad AST Op in cgcompare_and_set()\n")
	}
	OutFile.WriteString(fmt.Sprintf("\tcmpq\t%s, %s\n", reglist[r2], reglist[r1]))
	OutFile.WriteString(fmt.Sprintf("\t%s\t%s\n", op, breglist[r2]))
	OutFile.WriteString(fmt.Sprintf("\tmovzbq\t%s, %s\n", breglist[r2], reglist[r2]))
	free_register(r1)
	return (r2)
}

// List of inverted jump instructions,
// in AST order: A_EQ, A_NE, A_LT, A_GT, A_LE, A_GE
var invertedcmplist = map[NodeType]string{
	NodeEqual:              "jne",
	NodeNotEqual:           "je",
	NodeLessThan:           "jge",
	NodeGreaterThan:        "jle",
	NodeLessThanOrEqual:    "jg",
	NodeGreaterThanOrEqual: "jl",
}

// Compare two registers and jump if false.
func cgcompare_and_jump(ASTop NodeType, r1, r2, label int) int {
	// Check the range of the AST operation
	op, ok := invertedcmplist[ASTop]
	if !ok {
		fatal("bad AST Op in cgcompare_and_jump()\n")
	}
	OutFile.WriteString(fmt.Sprintf("\tcmpq\t%s, %s\n", reglist[r2], reglist[r1]))
	OutFile.WriteString(fmt.Sprintf("\t%s\tL%d\n", op, label))
	freeall_registers()
	return NoReg
}

// Generate a label
func cglabel(l int) {
	OutFile.WriteString(fmt.Sprintf("L%d:\n", l))
}

// Generate a jump to a label
func cgjump(l int) {
	OutFile.WriteString(fmt.Sprintf("\tjmp\tL%d\n", l))
}
