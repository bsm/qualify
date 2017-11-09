// +build amd64,!appengine
// Copyright (c) 2014 Will Fitzgerald. All rights reserved.

TEXT ·hasAsm(SB),4,$0-1
MOVQ $1, AX
CPUID
SHRQ $23, CX
ANDQ $1, CX
MOVB CX, ret+0(FP)
RET

#define POPCNTQ_DX_DX BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0xd2

TEXT ·popcntSliceAsm(SB),4,$0-32
XORQ  AX, AX
MOVQ  s+0(FP), SI
MOVQ  s_len+8(FP), CX
TESTQ CX, CX
JZ    popcntSliceEnd
popcntSliceLoop:
BYTE $0xf3; BYTE $0x48; BYTE $0x0f; BYTE $0xb8; BYTE $0x16 // POPCNTQ (SI), DX
ADDQ  DX, AX
ADDQ  $8, SI
LOOP  popcntSliceLoop
popcntSliceEnd:
MOVQ  AX, ret+24(FP)
RET
