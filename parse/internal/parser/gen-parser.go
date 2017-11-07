package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mmcloughlin/trunnel/ast"
)

const lingeringDeclarationsKey = "lingering_declarations"

func (c *current) addLingeringDeclaration(d interface{}) {
	var ds []interface{}
	if e, ok := c.state[lingeringDeclarationsKey]; ok {
		ds = e.([]interface{})
	}
	c.state[lingeringDeclarationsKey] = append(ds, d)
}

func (c *current) getLingeringDeclarations() []interface{} {
	if ds, ok := c.state[lingeringDeclarationsKey]; ok {
		return ds.([]interface{})
	}
	return []interface{}{}
}

var g = &grammar{
	rules: []*rule{
		{
			name: "File",
			pos:  position{line: 28, col: 1, offset: 587},
			expr: &actionExpr{
				pos: position{line: 28, col: 9, offset: 595},
				run: (*parser).callonFile1,
				expr: &seqExpr{
					pos: position{line: 28, col: 9, offset: 595},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 28, col: 9, offset: 595},
							label: "ds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 28, col: 12, offset: 598},
								expr: &ruleRefExpr{
									pos:  position{line: 28, col: 12, offset: 598},
									name: "Declaration",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 25, offset: 611},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Declaration",
			pos:  position{line: 53, col: 1, offset: 1340},
			expr: &actionExpr{
				pos: position{line: 53, col: 16, offset: 1355},
				run: (*parser).callonDeclaration1,
				expr: &seqExpr{
					pos: position{line: 53, col: 16, offset: 1355},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 53, col: 16, offset: 1355},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 53, col: 18, offset: 1357},
							label: "d",
							expr: &choiceExpr{
								pos: position{line: 53, col: 21, offset: 1360},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 53, col: 21, offset: 1360},
										name: "ConstDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 53, col: 40, offset: 1379},
										name: "StructDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 53, col: 60, offset: 1399},
										name: "PragmaDeclaration",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 53, col: 79, offset: 1418},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "ConstDeclaration",
			pos:  position{line: 59, col: 1, offset: 1483},
			expr: &actionExpr{
				pos: position{line: 59, col: 21, offset: 1503},
				run: (*parser).callonConstDeclaration1,
				expr: &seqExpr{
					pos: position{line: 59, col: 21, offset: 1503},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 59, col: 21, offset: 1503},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 59, col: 29, offset: 1511},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 59, col: 32, offset: 1514},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 34, offset: 1516},
								name: "ConstIdentifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 59, col: 50, offset: 1532},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 59, col: 52, offset: 1534},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 59, col: 56, offset: 1538},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 59, col: 58, offset: 1540},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 60, offset: 1542},
								name: "IntLiteral",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 59, col: 71, offset: 1553},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 59, col: 73, offset: 1555},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructDeclaration",
			pos:  position{line: 69, col: 1, offset: 1740},
			expr: &actionExpr{
				pos: position{line: 69, col: 22, offset: 1761},
				run: (*parser).callonStructDeclaration1,
				expr: &seqExpr{
					pos: position{line: 69, col: 22, offset: 1761},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 69, col: 22, offset: 1761},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 24, offset: 1763},
								name: "StructDecl",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 35, offset: 1774},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 69, col: 37, offset: 1776},
							expr: &litMatcher{
								pos:        position{line: 69, col: 37, offset: 1776},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "StructDecl",
			pos:  position{line: 75, col: 1, offset: 1878},
			expr: &actionExpr{
				pos: position{line: 75, col: 15, offset: 1892},
				run: (*parser).callonStructDecl1,
				expr: &seqExpr{
					pos: position{line: 75, col: 15, offset: 1892},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 75, col: 15, offset: 1892},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 75, col: 17, offset: 1894},
								name: "StructIdentifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 34, offset: 1911},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 75, col: 37, offset: 1914},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 75, col: 41, offset: 1918},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 75, col: 44, offset: 1921},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 44, offset: 1921},
									name: "StructMember",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 75, col: 58, offset: 1935},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 75, col: 60, offset: 1937},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 60, offset: 1937},
									name: "StructEnding",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 75, col: 74, offset: 1951},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructIdentifier",
			pos:  position{line: 86, col: 1, offset: 2137},
			expr: &actionExpr{
				pos: position{line: 86, col: 21, offset: 2157},
				run: (*parser).callonStructIdentifier1,
				expr: &seqExpr{
					pos: position{line: 86, col: 21, offset: 2157},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 86, col: 21, offset: 2157},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 30, offset: 2166},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 33, offset: 2169},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 86, col: 35, offset: 2171},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "PragmaDeclaration",
			pos:  position{line: 90, col: 1, offset: 2203},
			expr: &actionExpr{
				pos: position{line: 90, col: 22, offset: 2224},
				run: (*parser).callonPragmaDeclaration1,
				expr: &seqExpr{
					pos: position{line: 90, col: 22, offset: 2224},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 90, col: 22, offset: 2224},
							val:        "trunnel",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 32, offset: 2234},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 90, col: 35, offset: 2237},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 37, offset: 2239},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 48, offset: 2250},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 90, col: 51, offset: 2253},
							label: "opts",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 56, offset: 2258},
								name: "IdentifierList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 71, offset: 2273},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 90, col: 73, offset: 2275},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructMember",
			pos:  position{line: 109, col: 1, offset: 2688},
			expr: &actionExpr{
				pos: position{line: 109, col: 17, offset: 2704},
				run: (*parser).callonStructMember1,
				expr: &seqExpr{
					pos: position{line: 109, col: 17, offset: 2704},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 109, col: 17, offset: 2704},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 109, col: 19, offset: 2706},
							label: "m",
							expr: &choiceExpr{
								pos: position{line: 109, col: 22, offset: 2709},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 109, col: 22, offset: 2709},
										name: "SMArray",
									},
									&ruleRefExpr{
										pos:  position{line: 109, col: 32, offset: 2719},
										name: "SMInteger",
									},
									&ruleRefExpr{
										pos:  position{line: 109, col: 44, offset: 2731},
										name: "SMStruct",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 109, col: 54, offset: 2741},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "StructEnding",
			pos:  position{line: 117, col: 1, offset: 2848},
			expr: &litMatcher{
				pos:        position{line: 117, col: 17, offset: 2864},
				val:        "eos",
				ignoreCase: false,
			},
		},
		{
			name: "SMArray",
			pos:  position{line: 122, col: 1, offset: 2932},
			expr: &actionExpr{
				pos: position{line: 122, col: 12, offset: 2943},
				run: (*parser).callonSMArray1,
				expr: &seqExpr{
					pos: position{line: 122, col: 12, offset: 2943},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 122, col: 12, offset: 2943},
							label: "a",
							expr: &ruleRefExpr{
								pos:  position{line: 122, col: 14, offset: 2945},
								name: "SMFixedArray",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 27, offset: 2958},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 122, col: 29, offset: 2960},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SMFixedArray",
			pos:  position{line: 128, col: 1, offset: 3032},
			expr: &actionExpr{
				pos: position{line: 128, col: 17, offset: 3048},
				run: (*parser).callonSMFixedArray1,
				expr: &seqExpr{
					pos: position{line: 128, col: 17, offset: 3048},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 128, col: 17, offset: 3048},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 128, col: 19, offset: 3050},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 128, col: 29, offset: 3060},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 128, col: 32, offset: 3063},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 128, col: 34, offset: 3065},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 128, col: 45, offset: 3076},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 128, col: 47, offset: 3078},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 128, col: 51, offset: 3082},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 128, col: 53, offset: 3084},
								name: "Integer",
							},
						},
						&litMatcher{
							pos:        position{line: 128, col: 61, offset: 3092},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ArrayBase",
			pos:  position{line: 141, col: 1, offset: 3331},
			expr: &actionExpr{
				pos: position{line: 141, col: 14, offset: 3344},
				run: (*parser).callonArrayBase1,
				expr: &labeledExpr{
					pos:   position{line: 141, col: 14, offset: 3344},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 141, col: 17, offset: 3347},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 141, col: 17, offset: 3347},
								name: "IntType",
							},
							&ruleRefExpr{
								pos:  position{line: 141, col: 27, offset: 3357},
								name: "CharType",
							},
							&ruleRefExpr{
								pos:  position{line: 141, col: 38, offset: 3368},
								name: "StructRef",
							},
						},
					},
				},
			},
		},
		{
			name: "SMInteger",
			pos:  position{line: 148, col: 1, offset: 3448},
			expr: &actionExpr{
				pos: position{line: 148, col: 14, offset: 3461},
				run: (*parser).callonSMInteger1,
				expr: &seqExpr{
					pos: position{line: 148, col: 14, offset: 3461},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 148, col: 14, offset: 3461},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 148, col: 16, offset: 3463},
								name: "IntType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 148, col: 24, offset: 3471},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 148, col: 26, offset: 3473},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 148, col: 28, offset: 3475},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 148, col: 39, offset: 3486},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 148, col: 41, offset: 3488},
							label: "cst",
							expr: &zeroOrOneExpr{
								pos: position{line: 148, col: 45, offset: 3492},
								expr: &ruleRefExpr{
									pos:  position{line: 148, col: 45, offset: 3492},
									name: "IntConstraint",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 148, col: 60, offset: 3507},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 148, col: 62, offset: 3509},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SMStruct",
			pos:  position{line: 162, col: 1, offset: 3741},
			expr: &actionExpr{
				pos: position{line: 162, col: 13, offset: 3753},
				run: (*parser).callonSMStruct1,
				expr: &seqExpr{
					pos: position{line: 162, col: 13, offset: 3753},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 162, col: 13, offset: 3753},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 15, offset: 3755},
								name: "StructRef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 25, offset: 3765},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 162, col: 28, offset: 3768},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 30, offset: 3770},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 41, offset: 3781},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 162, col: 43, offset: 3783},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructRef",
			pos:  position{line: 169, col: 1, offset: 3880},
			expr: &choiceExpr{
				pos: position{line: 169, col: 14, offset: 3893},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 169, col: 14, offset: 3893},
						run: (*parser).callonStructRef2,
						expr: &labeledExpr{
							pos:   position{line: 169, col: 14, offset: 3893},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 16, offset: 3895},
								name: "StructDecl",
							},
						},
					},
					&actionExpr{
						pos: position{line: 172, col: 5, offset: 4000},
						run: (*parser).callonStructRef5,
						expr: &labeledExpr{
							pos:   position{line: 172, col: 5, offset: 4000},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 172, col: 7, offset: 4002},
								name: "StructIdentifier",
							},
						},
					},
				},
			},
		},
		{
			name: "CharType",
			pos:  position{line: 176, col: 1, offset: 4071},
			expr: &actionExpr{
				pos: position{line: 176, col: 13, offset: 4083},
				run: (*parser).callonCharType1,
				expr: &litMatcher{
					pos:        position{line: 176, col: 13, offset: 4083},
					val:        "char",
					ignoreCase: false,
				},
			},
		},
		{
			name: "IntType",
			pos:  position{line: 185, col: 1, offset: 4205},
			expr: &actionExpr{
				pos: position{line: 185, col: 12, offset: 4216},
				run: (*parser).callonIntType1,
				expr: &seqExpr{
					pos: position{line: 185, col: 12, offset: 4216},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 185, col: 12, offset: 4216},
							val:        "u",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 185, col: 16, offset: 4220},
							label: "b",
							expr: &choiceExpr{
								pos: position{line: 185, col: 19, offset: 4223},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 185, col: 19, offset: 4223},
										val:        "8",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 185, col: 25, offset: 4229},
										val:        "16",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 185, col: 32, offset: 4236},
										val:        "32",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 185, col: 39, offset: 4243},
										val:        "64",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IntConstraint",
			pos:  position{line: 193, col: 1, offset: 4401},
			expr: &actionExpr{
				pos: position{line: 193, col: 18, offset: 4418},
				run: (*parser).callonIntConstraint1,
				expr: &seqExpr{
					pos: position{line: 193, col: 18, offset: 4418},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 193, col: 18, offset: 4418},
							val:        "IN",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 193, col: 23, offset: 4423},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 193, col: 26, offset: 4426},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 193, col: 30, offset: 4430},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 193, col: 32, offset: 4432},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 193, col: 34, offset: 4434},
								name: "IntList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 193, col: 42, offset: 4442},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 193, col: 44, offset: 4444},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "IntList",
			pos:  position{line: 200, col: 1, offset: 4540},
			expr: &actionExpr{
				pos: position{line: 200, col: 12, offset: 4551},
				run: (*parser).callonIntList1,
				expr: &seqExpr{
					pos: position{line: 200, col: 12, offset: 4551},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 200, col: 12, offset: 4551},
							label: "m",
							expr: &ruleRefExpr{
								pos:  position{line: 200, col: 14, offset: 4553},
								name: "IntListMember",
							},
						},
						&labeledExpr{
							pos:   position{line: 200, col: 28, offset: 4567},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 200, col: 31, offset: 4570},
								expr: &seqExpr{
									pos: position{line: 200, col: 32, offset: 4571},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 200, col: 32, offset: 4571},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 200, col: 34, offset: 4573},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 200, col: 38, offset: 4577},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 200, col: 40, offset: 4579},
											name: "IntListMember",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IntListMember",
			pos:  position{line: 211, col: 1, offset: 4869},
			expr: &actionExpr{
				pos: position{line: 211, col: 18, offset: 4886},
				run: (*parser).callonIntListMember1,
				expr: &seqExpr{
					pos: position{line: 211, col: 18, offset: 4886},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 211, col: 18, offset: 4886},
							label: "lo",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 21, offset: 4889},
								name: "Integer",
							},
						},
						&labeledExpr{
							pos:   position{line: 211, col: 29, offset: 4897},
							label: "hi",
							expr: &zeroOrOneExpr{
								pos: position{line: 211, col: 32, offset: 4900},
								expr: &seqExpr{
									pos: position{line: 211, col: 34, offset: 4902},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 211, col: 34, offset: 4902},
											val:        "..",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 211, col: 39, offset: 4907},
											name: "Integer",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 224, col: 1, offset: 5112},
			expr: &actionExpr{
				pos: position{line: 224, col: 12, offset: 5123},
				run: (*parser).callonInteger1,
				expr: &labeledExpr{
					pos:   position{line: 224, col: 12, offset: 5123},
					label: "i",
					expr: &choiceExpr{
						pos: position{line: 224, col: 15, offset: 5126},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 224, col: 15, offset: 5126},
								name: "IntegerConstRef",
							},
							&ruleRefExpr{
								pos:  position{line: 224, col: 33, offset: 5144},
								name: "IntegerLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "IntegerConstRef",
			pos:  position{line: 228, col: 1, offset: 5181},
			expr: &actionExpr{
				pos: position{line: 228, col: 20, offset: 5200},
				run: (*parser).callonIntegerConstRef1,
				expr: &labeledExpr{
					pos:   position{line: 228, col: 20, offset: 5200},
					label: "n",
					expr: &ruleRefExpr{
						pos:  position{line: 228, col: 22, offset: 5202},
						name: "ConstIdentifier",
					},
				},
			},
		},
		{
			name: "IntegerLiteral",
			pos:  position{line: 232, col: 1, offset: 5276},
			expr: &actionExpr{
				pos: position{line: 232, col: 19, offset: 5294},
				run: (*parser).callonIntegerLiteral1,
				expr: &labeledExpr{
					pos:   position{line: 232, col: 19, offset: 5294},
					label: "v",
					expr: &ruleRefExpr{
						pos:  position{line: 232, col: 21, offset: 5296},
						name: "IntLiteral",
					},
				},
			},
		},
		{
			name: "IdentifierList",
			pos:  position{line: 238, col: 1, offset: 5380},
			expr: &actionExpr{
				pos: position{line: 238, col: 19, offset: 5398},
				run: (*parser).callonIdentifierList1,
				expr: &seqExpr{
					pos: position{line: 238, col: 19, offset: 5398},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 238, col: 19, offset: 5398},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 238, col: 21, offset: 5400},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 238, col: 32, offset: 5411},
							label: "ns",
							expr: &zeroOrMoreExpr{
								pos: position{line: 238, col: 35, offset: 5414},
								expr: &seqExpr{
									pos: position{line: 238, col: 36, offset: 5415},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 238, col: 36, offset: 5415},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 238, col: 38, offset: 5417},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 238, col: 42, offset: 5421},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 238, col: 44, offset: 5423},
											name: "Identifier",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 246, col: 1, offset: 5587},
			expr: &actionExpr{
				pos: position{line: 246, col: 15, offset: 5601},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 246, col: 15, offset: 5601},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 246, col: 15, offset: 5601},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 246, col: 25, offset: 5611},
							expr: &charClassMatcher{
								pos:        position{line: 246, col: 25, offset: 5611},
								val:        "[a-zA-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "ConstIdentifier",
			pos:  position{line: 250, col: 1, offset: 5659},
			expr: &actionExpr{
				pos: position{line: 250, col: 20, offset: 5678},
				run: (*parser).callonConstIdentifier1,
				expr: &seqExpr{
					pos: position{line: 250, col: 20, offset: 5678},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 250, col: 20, offset: 5678},
							val:        "[A-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 250, col: 27, offset: 5685},
							expr: &charClassMatcher{
								pos:        position{line: 250, col: 27, offset: 5685},
								val:        "[A-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "IntLiteral",
			pos:  position{line: 256, col: 1, offset: 5745},
			expr: &actionExpr{
				pos: position{line: 256, col: 15, offset: 5759},
				run: (*parser).callonIntLiteral1,
				expr: &choiceExpr{
					pos: position{line: 256, col: 16, offset: 5760},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 256, col: 16, offset: 5760},
							name: "HexLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 256, col: 29, offset: 5773},
							name: "OctalLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 256, col: 44, offset: 5788},
							name: "DecimalLiteral",
						},
					},
				},
			},
		},
		{
			name: "DecimalLiteral",
			pos:  position{line: 260, col: 1, offset: 5858},
			expr: &seqExpr{
				pos: position{line: 260, col: 19, offset: 5876},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 260, col: 19, offset: 5876},
						expr: &charClassMatcher{
							pos:        position{line: 260, col: 19, offset: 5876},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 260, col: 26, offset: 5883},
						expr: &charClassMatcher{
							pos:        position{line: 260, col: 26, offset: 5883},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name: "HexLiteral",
			pos:  position{line: 262, col: 1, offset: 5891},
			expr: &seqExpr{
				pos: position{line: 262, col: 15, offset: 5905},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 262, col: 15, offset: 5905},
						val:        "0x",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 262, col: 20, offset: 5910},
						expr: &charClassMatcher{
							pos:        position{line: 262, col: 20, offset: 5910},
							val:        "[0-9a-fA-F]",
							ranges:     []rune{'0', '9', 'a', 'f', 'A', 'F'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name: "OctalLiteral",
			pos:  position{line: 264, col: 1, offset: 5924},
			expr: &seqExpr{
				pos: position{line: 264, col: 17, offset: 5940},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 264, col: 17, offset: 5940},
						val:        "0",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 264, col: 21, offset: 5944},
						expr: &charClassMatcher{
							pos:        position{line: 264, col: 21, offset: 5944},
							val:        "[0-7]",
							ranges:     []rune{'0', '7'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"optional whitespace\"",
			pos:         position{line: 268, col: 1, offset: 5974},
			expr: &zeroOrMoreExpr{
				pos: position{line: 268, col: 28, offset: 6001},
				expr: &charClassMatcher{
					pos:        position{line: 268, col: 28, offset: 6001},
					val:        "[ \\t\\n]",
					chars:      []rune{' ', '\t', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name:        "__",
			displayName: "\"whitespace\"",
			pos:         position{line: 269, col: 1, offset: 6010},
			expr: &oneOrMoreExpr{
				pos: position{line: 269, col: 20, offset: 6029},
				expr: &charClassMatcher{
					pos:        position{line: 269, col: 20, offset: 6029},
					val:        "[ \\t\\n]",
					chars:      []rune{' ', '\t', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 271, col: 1, offset: 6039},
			expr: &litMatcher{
				pos:        position{line: 271, col: 8, offset: 6046},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 272, col: 1, offset: 6051},
			expr: &notExpr{
				pos: position{line: 272, col: 8, offset: 6058},
				expr: &anyMatcher{
					line: 272, col: 9, offset: 6059,
				},
			},
		},
	},
}

func (c *current) onFile1(ds interface{}) (interface{}, error) {
	f := &ast.File{}
	lingering := c.getLingeringDeclarations()
	decls := append(ds.([]interface{}), lingering...)
	for _, i := range decls {
		switch d := i.(type) {
		case *ast.Constant:
			f.Constants = append(f.Constants, d)
		case *ast.Struct:
			f.Structs = append(f.Structs, d)
		case *ast.Pragma:
			f.Pragmas = append(f.Pragmas, d)
		default:
			return nil, errors.New("unknown declaration")
		}
	}
	return f, nil
}

func (p *parser) callonFile1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFile1(stack["ds"])
}

func (c *current) onDeclaration1(d interface{}) (interface{}, error) {
	return d, nil
}

func (p *parser) callonDeclaration1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDeclaration1(stack["d"])
}

func (c *current) onConstDeclaration1(n, v interface{}) (interface{}, error) {
	return &ast.Constant{
		Name:  n.(string),
		Value: v.(int64),
	}, nil
}

func (p *parser) callonConstDeclaration1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstDeclaration1(stack["n"], stack["v"])
}

func (c *current) onStructDeclaration1(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonStructDeclaration1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructDeclaration1(stack["s"])
}

func (c *current) onStructDecl1(n, ms, e interface{}) (interface{}, error) {
	m := make([]ast.Member, 0)
	for _, i := range ms.([]interface{}) {
		m = append(m, i.(ast.Member))
	}
	return &ast.Struct{
		Name:    n.(string),
		Members: m,
	}, nil
}

func (p *parser) callonStructDecl1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructDecl1(stack["n"], stack["ms"], stack["e"])
}

func (c *current) onStructIdentifier1(n interface{}) (interface{}, error) {
	return n, nil
}

func (p *parser) callonStructIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructIdentifier1(stack["n"])
}

func (c *current) onPragmaDeclaration1(n, opts interface{}) (interface{}, error) {
	return &ast.Pragma{
		Type:    n.(string),
		Options: opts.([]string),
	}, nil
}

func (p *parser) callonPragmaDeclaration1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPragmaDeclaration1(stack["n"], stack["opts"])
}

func (c *current) onStructMember1(m interface{}) (interface{}, error) {
	return m, nil
}

func (p *parser) callonStructMember1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructMember1(stack["m"])
}

func (c *current) onSMArray1(a interface{}) (interface{}, error) {
	return a, nil
}

func (p *parser) callonSMArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMArray1(stack["a"])
}

func (c *current) onSMFixedArray1(b, n, s interface{}) (interface{}, error) {
	return &ast.FixedArrayMember{
		Base: b.(ast.ArrayBase),
		Name: n.(string),
		Size: s.(ast.Integer),
	}, nil
}

func (p *parser) callonSMFixedArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMFixedArray1(stack["b"], stack["n"], stack["s"])
}

func (c *current) onArrayBase1(t interface{}) (interface{}, error) {
	return t, nil
}

func (p *parser) callonArrayBase1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArrayBase1(stack["t"])
}

func (c *current) onSMInteger1(t, n, cst interface{}) (interface{}, error) {
	m := &ast.IntegerMember{
		Type: t.(*ast.IntType),
		Name: n.(string),
	}
	if cst != nil {
		m.Constraint = cst.(*ast.IntegerList)
	}
	return m, nil
}

func (p *parser) callonSMInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMInteger1(stack["t"], stack["n"], stack["cst"])
}

func (c *current) onSMStruct1(s, n interface{}) (interface{}, error) {
	return &ast.StructMember{
		Name: n.(string),
		Ref:  s.(*ast.StructRef),
	}, nil
}

func (p *parser) callonSMStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMStruct1(stack["s"], stack["n"])
}

func (c *current) onStructRef2(s interface{}) (interface{}, error) {
	c.addLingeringDeclaration(s)
	return &ast.StructRef{Name: s.(*ast.Struct).Name}, nil
}

func (p *parser) callonStructRef2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructRef2(stack["s"])
}

func (c *current) onStructRef5(n interface{}) (interface{}, error) {
	return &ast.StructRef{Name: n.(string)}, nil
}

func (p *parser) callonStructRef5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructRef5(stack["n"])
}

func (c *current) onCharType1() (interface{}, error) {
	return &ast.CharType{}, nil
}

func (p *parser) callonCharType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharType1()
}

func (c *current) onIntType1(b interface{}) (interface{}, error) {
	s, err := strconv.Atoi(string(b.([]byte)))
	return &ast.IntType{Size: s}, err
}

func (p *parser) callonIntType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntType1(stack["b"])
}

func (c *current) onIntConstraint1(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonIntConstraint1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntConstraint1(stack["l"])
}

func (c *current) onIntList1(m, ms interface{}) (interface{}, error) {
	r := []*ast.IntegerRange{m.(*ast.IntegerRange)}
	for _, i := range ms.([]interface{}) {
		r = append(r, i.([]interface{})[3].(*ast.IntegerRange))
	}
	return &ast.IntegerList{Ranges: r}, nil
}

func (p *parser) callonIntList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntList1(stack["m"], stack["ms"])
}

func (c *current) onIntListMember1(lo, hi interface{}) (interface{}, error) {
	r := &ast.IntegerRange{
		Low: lo.(ast.Integer),
	}
	if hi != nil {
		r.High = hi.([]interface{})[1].(ast.Integer)
	}
	return r, nil
}

func (p *parser) callonIntListMember1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntListMember1(stack["lo"], stack["hi"])
}

func (c *current) onInteger1(i interface{}) (interface{}, error) {
	return i, nil
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1(stack["i"])
}

func (c *current) onIntegerConstRef1(n interface{}) (interface{}, error) {
	return &ast.IntegerConstRef{Name: n.(string)}, nil
}

func (p *parser) callonIntegerConstRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerConstRef1(stack["n"])
}

func (c *current) onIntegerLiteral1(v interface{}) (interface{}, error) {
	return &ast.IntegerLiteral{Value: v.(int64)}, nil
}

func (p *parser) callonIntegerLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntegerLiteral1(stack["v"])
}

func (c *current) onIdentifierList1(n, ns interface{}) (interface{}, error) {
	ids := []string{n.(string)}
	for _, i := range ns.([]interface{}) {
		ids = append(ids, i.([]interface{})[3].(string))
	}
	return ids, nil
}

func (p *parser) callonIdentifierList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifierList1(stack["n"], stack["ns"])
}

func (c *current) onIdentifier1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onConstIdentifier1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonConstIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstIdentifier1()
}

func (c *current) onIntLiteral1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 0, 64)
}

func (p *parser) callonIntLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntLiteral1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
