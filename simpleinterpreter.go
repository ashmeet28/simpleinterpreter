package main

import "fmt"

/*

111 000 ECALL

001 000 ADD dest src1 src2
001 001 SUB dest src1 src2
001 100 XOR dest src1 src2
001 110 OR  dest src1 src2
001 111 AND dest src1 src2
001 010 SRA dest src1 src2
001 011 SRL dest src1 src2
001 101 SLL dest src1 src2

010 000 LB  dest base offset
010 001 LH  dest base offset
010 010 LW  dest base offset
010 100 LBU dest base offset
010 101 LHU dest base offset

011 000 SB src base offset
011 001 SH src base offset
011 010 SW src base offset

100 001 LUI  dest imm1 imm2
100 010 LLI  dest imm1 imm2
100 011 LLIU dest imm1 imm2

101 000 BEQ  offset src1 src2
101 001 BNE  offset src1 src2
101 100 BLT  offset src1 src2
101 101 BGE  offset src1 src2
101 110 BLTU offset src1 src2
101 111 BGEU offset src1 src2

110 000 JAL  dest zero offset
110 001 JALR dest base offset

*/

func main() {
	fmt.Println("Hello")
}
