package sexp

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tests = map[string]struct {
	input string
	node  *Node
	err   error
}{
	"success 01": {
		input: `(module
  (func $main
    (result i32)
   i32.const 5
   i32.const 20
   i32.add
   i32.const 4
   i32.sub
   i32.const 3
   i32.mul
   i32.const 7
   i32.div_s
   return
     drop 
    i32.const 0
  )
  (export "main" (func $main))
)`,
		// (
		node: &Node{
			Type: NodeCell,
			Car: &Node{ // module
				Type:  NodeSymbol,
				Value: "module",
			},
			Cdr: &Node{
				Type: NodeCell,
				Car: &Node{ // (
					Type: NodeCell,
					Car: &Node{ // func
						Type:  NodeSymbol,
						Value: "func",
					},
					Cdr: &Node{
						Type: NodeCell,
						Car: &Node{ // $main
							Type:  NodeSymbol,
							Value: "$main",
						},
						Cdr: &Node{
							Type: NodeCell,
							Car: &Node{ // (
								Type: NodeCell,
								Car: &Node{
									Type:  NodeSymbol,
									Value: "result",
								},
								Cdr: &Node{
									Type: NodeCell,
									Car: &Node{
										Type:  NodeSymbol,
										Value: "i32",
									},
								},
							},
							Cdr: &Node{
								Type: NodeCell,
								Car: &Node{ // i32.const
									Type:  NodeSymbol,
									Value: "i32.const",
								},
								Cdr: &Node{
									Type: NodeCell,
									Car: &Node{ // 5
										Type:  NodeInt,
										Value: int64(5),
									},
									Cdr: &Node{
										Type: NodeCell,
										Car: &Node{ // i32.const
											Type:  NodeSymbol,
											Value: "i32.const",
										},
										Cdr: &Node{
											Type: NodeCell,
											Car: &Node{ // 20
												Type:  NodeInt,
												Value: int64(20),
											},
											Cdr: &Node{
												Type: NodeCell,
												Car: &Node{ // i32.add
													Type:  NodeSymbol,
													Value: "i32.add",
												},
												Cdr: &Node{
													Type: NodeCell,
													Car: &Node{ // i32.const
														Type:  NodeSymbol,
														Value: "i32.const",
													},
													Cdr: &Node{
														Type: NodeCell,
														Car: &Node{ // 4
															Type:  NodeInt,
															Value: int64(4),
														},
														Cdr: &Node{
															Type: NodeCell,
															Car: &Node{ // i32.sub
																Type:  NodeSymbol,
																Value: "i32.sub",
															},
															Cdr: &Node{
																Type: NodeCell,
																Car: &Node{ // i32.const
																	Type:  NodeSymbol,
																	Value: "i32.const",
																},
																Cdr: &Node{
																	Type: NodeCell,
																	Car: &Node{ // 3
																		Type:  NodeInt,
																		Value: int64(3),
																	},
																	Cdr: &Node{
																		Type: NodeCell,
																		Car: &Node{ // i32.mul
																			Type:  NodeSymbol,
																			Value: "i32.mul",
																		},
																		Cdr: &Node{
																			Type: NodeCell,
																			Car: &Node{ // i32.const
																				Type:  NodeSymbol,
																				Value: "i32.const",
																			},
																			Cdr: &Node{
																				Type: NodeCell,
																				Car: &Node{ // 7
																					Type:  NodeInt,
																					Value: int64(7),
																				},
																				Cdr: &Node{
																					Type: NodeCell,
																					Car: &Node{ // i32.div_s
																						Type:  NodeSymbol,
																						Value: "i32.div_s",
																					},
																					Cdr: &Node{
																						Type: NodeCell,
																						Car: &Node{ // return
																							Type:  NodeSymbol,
																							Value: "return",
																						},
																						Cdr: &Node{
																							Type: NodeCell,
																							Car: &Node{ // drop
																								Type:  NodeSymbol,
																								Value: "drop",
																							},
																							Cdr: &Node{
																								Type: NodeCell,
																								Car: &Node{ // i32.const
																									Type:  NodeSymbol,
																									Value: "i32.const",
																								},
																								Cdr: &Node{
																									Type: NodeCell,
																									Car: &Node{ // 0
																										Type:  NodeInt,
																										Value: int64(0),
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Cdr: &Node{
					Type: NodeCell,
					Car: &Node{ // (
						Type: NodeCell,
						Car: &Node{ // export
							Type:  NodeSymbol,
							Value: "export",
						},
						Cdr: &Node{
							Type: NodeCell,
							Car: &Node{ // "main"
								Type:  NodeString,
								Value: "main",
							},
							Cdr: &Node{
								Type: NodeCell,
								Car: &Node{ // (
									Type: NodeCell,
									Car: &Node{ // func
										Type:  NodeSymbol,
										Value: "func",
									},
									Cdr: &Node{
										Type: NodeCell,
										Car: &Node{ // $main
											Type:  NodeSymbol,
											Value: "$main",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		err: nil,
	},
	"empty": {
		input: "",
		node:  nil,
		err:   nil,
	},
	"space only": {
		input: " \t\v\r\n",
		node:  nil,
		err:   nil,
	},
	"empty list": {
		input: "()",
		node: &Node{
			Type: NodeCell,
		},
		err: nil,
	},
	"unexpected EOF 01": {
		input: "(aaa bbb",
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"unexpected EOF 02": {
		input: `(   `,
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"unexpected EOF 03": {
		input: `(   aaa`,
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"unexpected EOF 04": {
		input: `(aaa "bbb ccc`,
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"unexpected EOF 05": {
		input: `(aaa "bbb ccc\`,
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"unexpected EOF 06": {
		input: `(aaa 1234`,
		node:  nil,
		err:   io.ErrUnexpectedEOF,
	},
	"invalid format 01": {
		input: `(aaa 1\234)`,
		node:  nil,
		err:   ErrInvalidFormat,
	},
	"invalid format 02": {
		input: `(aaa \1234`,
		node:  nil,
		err:   ErrInvalidFormat,
	},
}

func Test_Parser(t *testing.T) {
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := New(strings.NewReader(tt.input))
			node, err := p.Parse()
			if err != tt.err {
				t.Errorf("Parser.Parse(): err: want: %v, got:%v", tt.err, err)
			}

			if diff := cmp.Diff(node, tt.node); diff != "" {
				t.Errorf("Parser.Parse(), differs: (-got +want)\n%s", diff)
			}
		})
	}
}
