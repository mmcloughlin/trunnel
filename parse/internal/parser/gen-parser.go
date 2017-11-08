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
			pos:  position{line: 55, col: 1, offset: 1408},
			expr: &actionExpr{
				pos: position{line: 55, col: 16, offset: 1423},
				run: (*parser).callonDeclaration1,
				expr: &seqExpr{
					pos: position{line: 55, col: 16, offset: 1423},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 55, col: 16, offset: 1423},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 55, col: 18, offset: 1425},
							label: "d",
							expr: &choiceExpr{
								pos: position{line: 55, col: 21, offset: 1428},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 55, col: 21, offset: 1428},
										name: "ConstDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 40, offset: 1447},
										name: "ContextDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 61, offset: 1468},
										name: "StructDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 81, offset: 1488},
										name: "PragmaDeclaration",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 100, offset: 1507},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "ConstDeclaration",
			pos:  position{line: 61, col: 1, offset: 1572},
			expr: &actionExpr{
				pos: position{line: 61, col: 21, offset: 1592},
				run: (*parser).callonConstDeclaration1,
				expr: &seqExpr{
					pos: position{line: 61, col: 21, offset: 1592},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 61, col: 21, offset: 1592},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 29, offset: 1600},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 32, offset: 1603},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 34, offset: 1605},
								name: "ConstIdentifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 50, offset: 1621},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 52, offset: 1623},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 56, offset: 1627},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 58, offset: 1629},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 60, offset: 1631},
								name: "IntLiteral",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 71, offset: 1642},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 73, offset: 1644},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContextDeclaration",
			pos:  position{line: 73, col: 1, offset: 1905},
			expr: &actionExpr{
				pos: position{line: 73, col: 23, offset: 1927},
				run: (*parser).callonContextDeclaration1,
				expr: &seqExpr{
					pos: position{line: 73, col: 23, offset: 1927},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 73, col: 23, offset: 1927},
							val:        "context",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 33, offset: 1937},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 73, col: 36, offset: 1940},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 38, offset: 1942},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 49, offset: 1953},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 73, col: 52, offset: 1956},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 56, offset: 1960},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 73, col: 58, offset: 1962},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 73, col: 61, offset: 1965},
								expr: &ruleRefExpr{
									pos:  position{line: 73, col: 61, offset: 1965},
									name: "ContextMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 76, offset: 1980},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 73, col: 78, offset: 1982},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 82, offset: 1986},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 73, col: 84, offset: 1988},
							expr: &litMatcher{
								pos:        position{line: 73, col: 84, offset: 1988},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "ContextMember",
			pos:  position{line: 84, col: 1, offset: 2192},
			expr: &actionExpr{
				pos: position{line: 84, col: 18, offset: 2209},
				run: (*parser).callonContextMember1,
				expr: &seqExpr{
					pos: position{line: 84, col: 18, offset: 2209},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 84, col: 18, offset: 2209},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 84, col: 20, offset: 2211},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 84, col: 22, offset: 2213},
								name: "IntType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 30, offset: 2221},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 84, col: 33, offset: 2224},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 84, col: 35, offset: 2226},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 46, offset: 2237},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 84, col: 48, offset: 2239},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 52, offset: 2243},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "StructDeclaration",
			pos:  position{line: 94, col: 1, offset: 2437},
			expr: &actionExpr{
				pos: position{line: 94, col: 22, offset: 2458},
				run: (*parser).callonStructDeclaration1,
				expr: &seqExpr{
					pos: position{line: 94, col: 22, offset: 2458},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 94, col: 22, offset: 2458},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 94, col: 24, offset: 2460},
								name: "StructDecl",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 94, col: 35, offset: 2471},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 94, col: 37, offset: 2473},
							expr: &litMatcher{
								pos:        position{line: 94, col: 37, offset: 2473},
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
			pos:  position{line: 100, col: 1, offset: 2575},
			expr: &actionExpr{
				pos: position{line: 100, col: 15, offset: 2589},
				run: (*parser).callonStructDecl1,
				expr: &seqExpr{
					pos: position{line: 100, col: 15, offset: 2589},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 100, col: 15, offset: 2589},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 100, col: 17, offset: 2591},
								name: "StructIdentifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 34, offset: 2608},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 100, col: 37, offset: 2611},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 100, col: 41, offset: 2615},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 100, col: 44, offset: 2618},
								expr: &ruleRefExpr{
									pos:  position{line: 100, col: 44, offset: 2618},
									name: "StructMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 58, offset: 2632},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 100, col: 60, offset: 2634},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 100, col: 62, offset: 2636},
								expr: &ruleRefExpr{
									pos:  position{line: 100, col: 62, offset: 2636},
									name: "StructEnding",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 76, offset: 2650},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 100, col: 78, offset: 2652},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructIdentifier",
			pos:  position{line: 114, col: 1, offset: 2892},
			expr: &actionExpr{
				pos: position{line: 114, col: 21, offset: 2912},
				run: (*parser).callonStructIdentifier1,
				expr: &seqExpr{
					pos: position{line: 114, col: 21, offset: 2912},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 114, col: 21, offset: 2912},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 114, col: 30, offset: 2921},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 114, col: 33, offset: 2924},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 114, col: 35, offset: 2926},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "PragmaDeclaration",
			pos:  position{line: 118, col: 1, offset: 2958},
			expr: &actionExpr{
				pos: position{line: 118, col: 22, offset: 2979},
				run: (*parser).callonPragmaDeclaration1,
				expr: &seqExpr{
					pos: position{line: 118, col: 22, offset: 2979},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 118, col: 22, offset: 2979},
							val:        "trunnel",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 32, offset: 2989},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 35, offset: 2992},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 37, offset: 2994},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 48, offset: 3005},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 51, offset: 3008},
							label: "opts",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 56, offset: 3013},
								name: "IdentifierList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 71, offset: 3028},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 118, col: 73, offset: 3030},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructMember",
			pos:  position{line: 137, col: 1, offset: 3439},
			expr: &actionExpr{
				pos: position{line: 137, col: 17, offset: 3455},
				run: (*parser).callonStructMember1,
				expr: &seqExpr{
					pos: position{line: 137, col: 17, offset: 3455},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 137, col: 17, offset: 3455},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 137, col: 19, offset: 3457},
							label: "m",
							expr: &choiceExpr{
								pos: position{line: 137, col: 22, offset: 3460},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 137, col: 22, offset: 3460},
										name: "SMArray",
									},
									&ruleRefExpr{
										pos:  position{line: 137, col: 32, offset: 3470},
										name: "SMInteger",
									},
									&ruleRefExpr{
										pos:  position{line: 137, col: 44, offset: 3482},
										name: "SMString",
									},
									&ruleRefExpr{
										pos:  position{line: 137, col: 55, offset: 3493},
										name: "SMStruct",
									},
									&ruleRefExpr{
										pos:  position{line: 137, col: 66, offset: 3504},
										name: "SMUnion",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 137, col: 75, offset: 3513},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 137, col: 77, offset: 3515},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructEnding",
			pos:  position{line: 145, col: 1, offset: 3624},
			expr: &actionExpr{
				pos: position{line: 145, col: 17, offset: 3640},
				run: (*parser).callonStructEnding1,
				expr: &seqExpr{
					pos: position{line: 145, col: 17, offset: 3640},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 145, col: 17, offset: 3640},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 145, col: 19, offset: 3642},
							label: "e",
							expr: &choiceExpr{
								pos: position{line: 145, col: 22, offset: 3645},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 145, col: 22, offset: 3645},
										name: "SMRemainder",
									},
									&ruleRefExpr{
										pos:  position{line: 145, col: 36, offset: 3659},
										name: "StructEOS",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 145, col: 47, offset: 3670},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 145, col: 49, offset: 3672},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructEOS",
			pos:  position{line: 149, col: 1, offset: 3697},
			expr: &actionExpr{
				pos: position{line: 149, col: 14, offset: 3710},
				run: (*parser).callonStructEOS1,
				expr: &litMatcher{
					pos:        position{line: 149, col: 14, offset: 3710},
					val:        "eos",
					ignoreCase: false,
				},
			},
		},
		{
			name: "SMArray",
			pos:  position{line: 156, col: 1, offset: 3803},
			expr: &actionExpr{
				pos: position{line: 156, col: 12, offset: 3814},
				run: (*parser).callonSMArray1,
				expr: &labeledExpr{
					pos:   position{line: 156, col: 12, offset: 3814},
					label: "a",
					expr: &choiceExpr{
						pos: position{line: 156, col: 15, offset: 3817},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 156, col: 15, offset: 3817},
								name: "SMFixedArray",
							},
							&ruleRefExpr{
								pos:  position{line: 156, col: 30, offset: 3832},
								name: "SMVarArray",
							},
						},
					},
				},
			},
		},
		{
			name: "SMFixedArray",
			pos:  position{line: 162, col: 1, offset: 3912},
			expr: &actionExpr{
				pos: position{line: 162, col: 17, offset: 3928},
				run: (*parser).callonSMFixedArray1,
				expr: &seqExpr{
					pos: position{line: 162, col: 17, offset: 3928},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 162, col: 17, offset: 3928},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 19, offset: 3930},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 29, offset: 3940},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 162, col: 32, offset: 3943},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 34, offset: 3945},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 45, offset: 3956},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 162, col: 47, offset: 3958},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 162, col: 51, offset: 3962},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 53, offset: 3964},
								name: "Integer",
							},
						},
						&litMatcher{
							pos:        position{line: 162, col: 61, offset: 3972},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SMVarArray",
			pos:  position{line: 173, col: 1, offset: 4192},
			expr: &actionExpr{
				pos: position{line: 173, col: 15, offset: 4206},
				run: (*parser).callonSMVarArray1,
				expr: &seqExpr{
					pos: position{line: 173, col: 15, offset: 4206},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 173, col: 15, offset: 4206},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 173, col: 17, offset: 4208},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 173, col: 27, offset: 4218},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 173, col: 30, offset: 4221},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 173, col: 32, offset: 4223},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 173, col: 43, offset: 4234},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 173, col: 45, offset: 4236},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 173, col: 49, offset: 4240},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 173, col: 51, offset: 4242},
								name: "LengthConstraint",
							},
						},
						&litMatcher{
							pos:        position{line: 173, col: 68, offset: 4259},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "LengthConstraint",
			pos:  position{line: 181, col: 1, offset: 4400},
			expr: &actionExpr{
				pos: position{line: 181, col: 21, offset: 4420},
				run: (*parser).callonLengthConstraint1,
				expr: &labeledExpr{
					pos:   position{line: 181, col: 21, offset: 4420},
					label: "l",
					expr: &choiceExpr{
						pos: position{line: 181, col: 24, offset: 4423},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 181, col: 24, offset: 4423},
								name: "Leftover",
							},
							&ruleRefExpr{
								pos:  position{line: 181, col: 35, offset: 4434},
								name: "IDRef",
							},
						},
					},
				},
			},
		},
		{
			name: "Leftover",
			pos:  position{line: 185, col: 1, offset: 4462},
			expr: &actionExpr{
				pos: position{line: 185, col: 13, offset: 4474},
				run: (*parser).callonLeftover1,
				expr: &seqExpr{
					pos: position{line: 185, col: 13, offset: 4474},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 185, col: 13, offset: 4474},
							val:        "..-",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 185, col: 19, offset: 4480},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 185, col: 21, offset: 4482},
								name: "Integer",
							},
						},
					},
				},
			},
		},
		{
			name: "SMRemainder",
			pos:  position{line: 191, col: 1, offset: 4612},
			expr: &actionExpr{
				pos: position{line: 191, col: 16, offset: 4627},
				run: (*parser).callonSMRemainder1,
				expr: &seqExpr{
					pos: position{line: 191, col: 16, offset: 4627},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 191, col: 16, offset: 4627},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 191, col: 18, offset: 4629},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 191, col: 28, offset: 4639},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 191, col: 31, offset: 4642},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 191, col: 33, offset: 4644},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 191, col: 44, offset: 4655},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 191, col: 46, offset: 4657},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 191, col: 50, offset: 4661},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 191, col: 52, offset: 4663},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ArrayBase",
			pos:  position{line: 204, col: 1, offset: 4894},
			expr: &actionExpr{
				pos: position{line: 204, col: 14, offset: 4907},
				run: (*parser).callonArrayBase1,
				expr: &labeledExpr{
					pos:   position{line: 204, col: 14, offset: 4907},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 204, col: 17, offset: 4910},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 204, col: 17, offset: 4910},
								name: "IntType",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 27, offset: 4920},
								name: "CharType",
							},
							&ruleRefExpr{
								pos:  position{line: 204, col: 38, offset: 4931},
								name: "StructRef",
							},
						},
					},
				},
			},
		},
		{
			name: "SMInteger",
			pos:  position{line: 211, col: 1, offset: 5011},
			expr: &actionExpr{
				pos: position{line: 211, col: 14, offset: 5024},
				run: (*parser).callonSMInteger1,
				expr: &seqExpr{
					pos: position{line: 211, col: 14, offset: 5024},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 211, col: 14, offset: 5024},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 16, offset: 5026},
								name: "IntType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 24, offset: 5034},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 26, offset: 5036},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 211, col: 28, offset: 5038},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 211, col: 39, offset: 5049},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 211, col: 41, offset: 5051},
							label: "cst",
							expr: &zeroOrOneExpr{
								pos: position{line: 211, col: 45, offset: 5055},
								expr: &ruleRefExpr{
									pos:  position{line: 211, col: 45, offset: 5055},
									name: "IntConstraint",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SMString",
			pos:  position{line: 224, col: 1, offset: 5265},
			expr: &actionExpr{
				pos: position{line: 224, col: 13, offset: 5277},
				run: (*parser).callonSMString1,
				expr: &seqExpr{
					pos: position{line: 224, col: 13, offset: 5277},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 224, col: 13, offset: 5277},
							val:        "nulterm",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 224, col: 23, offset: 5287},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 224, col: 26, offset: 5290},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 224, col: 28, offset: 5292},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "SMStruct",
			pos:  position{line: 231, col: 1, offset: 5421},
			expr: &actionExpr{
				pos: position{line: 231, col: 13, offset: 5433},
				run: (*parser).callonSMStruct1,
				expr: &seqExpr{
					pos: position{line: 231, col: 13, offset: 5433},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 231, col: 13, offset: 5433},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 231, col: 15, offset: 5435},
								name: "StructRef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 231, col: 25, offset: 5445},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 231, col: 28, offset: 5448},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 231, col: 30, offset: 5450},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "SMUnion",
			pos:  position{line: 240, col: 1, offset: 5622},
			expr: &actionExpr{
				pos: position{line: 240, col: 12, offset: 5633},
				run: (*parser).callonSMUnion1,
				expr: &seqExpr{
					pos: position{line: 240, col: 12, offset: 5633},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 240, col: 12, offset: 5633},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 20, offset: 5641},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 240, col: 23, offset: 5644},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 240, col: 25, offset: 5646},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 36, offset: 5657},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 240, col: 38, offset: 5659},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 240, col: 42, offset: 5663},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 240, col: 44, offset: 5665},
								name: "IDRef",
							},
						},
						&litMatcher{
							pos:        position{line: 240, col: 50, offset: 5671},
							val:        "]",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 54, offset: 5675},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 240, col: 56, offset: 5677},
							label: "l",
							expr: &zeroOrOneExpr{
								pos: position{line: 240, col: 58, offset: 5679},
								expr: &ruleRefExpr{
									pos:  position{line: 240, col: 58, offset: 5679},
									name: "UnionLength",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 71, offset: 5692},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 240, col: 73, offset: 5694},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 77, offset: 5698},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 240, col: 79, offset: 5700},
							label: "cs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 240, col: 82, offset: 5703},
								expr: &ruleRefExpr{
									pos:  position{line: 240, col: 82, offset: 5703},
									name: "UnionMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 95, offset: 5716},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 240, col: 97, offset: 5718},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "UnionLength",
			pos:  position{line: 252, col: 1, offset: 5894},
			expr: &actionExpr{
				pos: position{line: 252, col: 16, offset: 5909},
				run: (*parser).callonUnionLength1,
				expr: &seqExpr{
					pos: position{line: 252, col: 16, offset: 5909},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 252, col: 16, offset: 5909},
							val:        "with",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 252, col: 23, offset: 5916},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 252, col: 26, offset: 5919},
							val:        "length",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 252, col: 35, offset: 5928},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 252, col: 38, offset: 5931},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 252, col: 40, offset: 5933},
								name: "LengthConstraint",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionMember",
			pos:  position{line: 258, col: 1, offset: 6046},
			expr: &actionExpr{
				pos: position{line: 258, col: 16, offset: 6061},
				run: (*parser).callonUnionMember1,
				expr: &seqExpr{
					pos: position{line: 258, col: 16, offset: 6061},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 258, col: 16, offset: 6061},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 258, col: 18, offset: 6063},
							label: "cse",
							expr: &ruleRefExpr{
								pos:  position{line: 258, col: 22, offset: 6067},
								name: "UnionCase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 258, col: 32, offset: 6077},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 258, col: 34, offset: 6079},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 258, col: 38, offset: 6083},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 258, col: 40, offset: 6085},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 258, col: 43, offset: 6088},
								name: "UnionBody",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 258, col: 53, offset: 6098},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "UnionCase",
			pos:  position{line: 271, col: 1, offset: 6278},
			expr: &choiceExpr{
				pos: position{line: 271, col: 14, offset: 6291},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 271, col: 14, offset: 6291},
						run: (*parser).callonUnionCase2,
						expr: &labeledExpr{
							pos:   position{line: 271, col: 14, offset: 6291},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 271, col: 16, offset: 6293},
								name: "IntList",
							},
						},
					},
					&actionExpr{
						pos: position{line: 273, col: 5, offset: 6323},
						run: (*parser).callonUnionCase5,
						expr: &litMatcher{
							pos:        position{line: 273, col: 5, offset: 6323},
							val:        "default",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "UnionBody",
			pos:  position{line: 284, col: 1, offset: 6571},
			expr: &choiceExpr{
				pos: position{line: 284, col: 14, offset: 6584},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 284, col: 14, offset: 6584},
						run: (*parser).callonUnionBody2,
						expr: &seqExpr{
							pos: position{line: 284, col: 14, offset: 6584},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 284, col: 14, offset: 6584},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 284, col: 16, offset: 6586},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 286, col: 5, offset: 6614},
						run: (*parser).callonUnionBody6,
						expr: &seqExpr{
							pos: position{line: 286, col: 5, offset: 6614},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 286, col: 5, offset: 6614},
									val:        "fail",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 286, col: 12, offset: 6621},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 286, col: 14, offset: 6623},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 288, col: 5, offset: 6659},
						run: (*parser).callonUnionBody11,
						expr: &seqExpr{
							pos: position{line: 288, col: 5, offset: 6659},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 288, col: 5, offset: 6659},
									val:        "ignore",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 288, col: 14, offset: 6668},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 288, col: 16, offset: 6670},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 290, col: 5, offset: 6708},
						run: (*parser).callonUnionBody16,
						expr: &labeledExpr{
							pos:   position{line: 290, col: 5, offset: 6708},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 290, col: 8, offset: 6711},
								name: "UnionFields",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionFields",
			pos:  position{line: 294, col: 1, offset: 6745},
			expr: &actionExpr{
				pos: position{line: 294, col: 16, offset: 6760},
				run: (*parser).callonUnionFields1,
				expr: &seqExpr{
					pos: position{line: 294, col: 16, offset: 6760},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 294, col: 16, offset: 6760},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 294, col: 19, offset: 6763},
								expr: &ruleRefExpr{
									pos:  position{line: 294, col: 19, offset: 6763},
									name: "UnionField",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 294, col: 31, offset: 6775},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 294, col: 33, offset: 6777},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 294, col: 35, offset: 6779},
								expr: &ruleRefExpr{
									pos:  position{line: 294, col: 35, offset: 6779},
									name: "ExtentSpec",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "UnionField",
			pos:  position{line: 311, col: 1, offset: 7108},
			expr: &actionExpr{
				pos: position{line: 311, col: 15, offset: 7122},
				run: (*parser).callonUnionField1,
				expr: &seqExpr{
					pos: position{line: 311, col: 15, offset: 7122},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 311, col: 15, offset: 7122},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 311, col: 17, offset: 7124},
							label: "m",
							expr: &choiceExpr{
								pos: position{line: 311, col: 20, offset: 7127},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 311, col: 20, offset: 7127},
										name: "SMFixedArray",
									},
									&ruleRefExpr{
										pos:  position{line: 311, col: 35, offset: 7142},
										name: "SMInteger",
									},
									&ruleRefExpr{
										pos:  position{line: 311, col: 47, offset: 7154},
										name: "SMString",
									},
									&ruleRefExpr{
										pos:  position{line: 311, col: 58, offset: 7165},
										name: "SMStruct",
									},
									&ruleRefExpr{
										pos:  position{line: 311, col: 69, offset: 7176},
										name: "SMVarArray",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 81, offset: 7188},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 311, col: 83, offset: 7190},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExtentSpec",
			pos:  position{line: 319, col: 1, offset: 7302},
			expr: &choiceExpr{
				pos: position{line: 319, col: 15, offset: 7316},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 319, col: 15, offset: 7316},
						run: (*parser).callonExtentSpec2,
						expr: &seqExpr{
							pos: position{line: 319, col: 15, offset: 7316},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 319, col: 15, offset: 7316},
									val:        "...",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 319, col: 21, offset: 7322},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 319, col: 23, offset: 7324},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 321, col: 5, offset: 7362},
						run: (*parser).callonExtentSpec7,
						expr: &seqExpr{
							pos: position{line: 321, col: 5, offset: 7362},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 321, col: 5, offset: 7362},
									label: "r",
									expr: &ruleRefExpr{
										pos:  position{line: 321, col: 7, offset: 7364},
										name: "SMRemainder",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 321, col: 19, offset: 7376},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 321, col: 21, offset: 7378},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "StructRef",
			pos:  position{line: 325, col: 1, offset: 7403},
			expr: &choiceExpr{
				pos: position{line: 325, col: 14, offset: 7416},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 325, col: 14, offset: 7416},
						run: (*parser).callonStructRef2,
						expr: &labeledExpr{
							pos:   position{line: 325, col: 14, offset: 7416},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 325, col: 16, offset: 7418},
								name: "StructDecl",
							},
						},
					},
					&actionExpr{
						pos: position{line: 328, col: 5, offset: 7523},
						run: (*parser).callonStructRef5,
						expr: &labeledExpr{
							pos:   position{line: 328, col: 5, offset: 7523},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 328, col: 7, offset: 7525},
								name: "StructIdentifier",
							},
						},
					},
				},
			},
		},
		{
			name: "CharType",
			pos:  position{line: 332, col: 1, offset: 7594},
			expr: &actionExpr{
				pos: position{line: 332, col: 13, offset: 7606},
				run: (*parser).callonCharType1,
				expr: &litMatcher{
					pos:        position{line: 332, col: 13, offset: 7606},
					val:        "char",
					ignoreCase: false,
				},
			},
		},
		{
			name: "IntType",
			pos:  position{line: 341, col: 1, offset: 7728},
			expr: &actionExpr{
				pos: position{line: 341, col: 12, offset: 7739},
				run: (*parser).callonIntType1,
				expr: &seqExpr{
					pos: position{line: 341, col: 12, offset: 7739},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 341, col: 12, offset: 7739},
							val:        "u",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 341, col: 16, offset: 7743},
							label: "b",
							expr: &choiceExpr{
								pos: position{line: 341, col: 19, offset: 7746},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 341, col: 19, offset: 7746},
										val:        "8",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 341, col: 25, offset: 7752},
										val:        "16",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 341, col: 32, offset: 7759},
										val:        "32",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 341, col: 39, offset: 7766},
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
			pos:  position{line: 349, col: 1, offset: 7924},
			expr: &actionExpr{
				pos: position{line: 349, col: 18, offset: 7941},
				run: (*parser).callonIntConstraint1,
				expr: &seqExpr{
					pos: position{line: 349, col: 18, offset: 7941},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 349, col: 18, offset: 7941},
							val:        "IN",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 349, col: 23, offset: 7946},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 349, col: 26, offset: 7949},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 349, col: 30, offset: 7953},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 349, col: 32, offset: 7955},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 349, col: 34, offset: 7957},
								name: "IntList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 349, col: 42, offset: 7965},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 349, col: 44, offset: 7967},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "IntList",
			pos:  position{line: 356, col: 1, offset: 8063},
			expr: &actionExpr{
				pos: position{line: 356, col: 12, offset: 8074},
				run: (*parser).callonIntList1,
				expr: &seqExpr{
					pos: position{line: 356, col: 12, offset: 8074},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 356, col: 12, offset: 8074},
							label: "m",
							expr: &ruleRefExpr{
								pos:  position{line: 356, col: 14, offset: 8076},
								name: "IntListMember",
							},
						},
						&labeledExpr{
							pos:   position{line: 356, col: 28, offset: 8090},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 356, col: 31, offset: 8093},
								expr: &seqExpr{
									pos: position{line: 356, col: 32, offset: 8094},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 356, col: 32, offset: 8094},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 356, col: 34, offset: 8096},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 356, col: 38, offset: 8100},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 356, col: 40, offset: 8102},
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
			pos:  position{line: 367, col: 1, offset: 8392},
			expr: &actionExpr{
				pos: position{line: 367, col: 18, offset: 8409},
				run: (*parser).callonIntListMember1,
				expr: &seqExpr{
					pos: position{line: 367, col: 18, offset: 8409},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 367, col: 18, offset: 8409},
							label: "lo",
							expr: &ruleRefExpr{
								pos:  position{line: 367, col: 21, offset: 8412},
								name: "Integer",
							},
						},
						&labeledExpr{
							pos:   position{line: 367, col: 29, offset: 8420},
							label: "hi",
							expr: &zeroOrOneExpr{
								pos: position{line: 367, col: 32, offset: 8423},
								expr: &seqExpr{
									pos: position{line: 367, col: 34, offset: 8425},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 367, col: 34, offset: 8425},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 367, col: 36, offset: 8427},
											val:        "..",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 367, col: 41, offset: 8432},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 367, col: 43, offset: 8434},
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
			pos:  position{line: 380, col: 1, offset: 8639},
			expr: &actionExpr{
				pos: position{line: 380, col: 12, offset: 8650},
				run: (*parser).callonInteger1,
				expr: &labeledExpr{
					pos:   position{line: 380, col: 12, offset: 8650},
					label: "i",
					expr: &choiceExpr{
						pos: position{line: 380, col: 15, offset: 8653},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 380, col: 15, offset: 8653},
								name: "IntegerConstRef",
							},
							&ruleRefExpr{
								pos:  position{line: 380, col: 33, offset: 8671},
								name: "IntegerLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "IntegerConstRef",
			pos:  position{line: 384, col: 1, offset: 8708},
			expr: &actionExpr{
				pos: position{line: 384, col: 20, offset: 8727},
				run: (*parser).callonIntegerConstRef1,
				expr: &labeledExpr{
					pos:   position{line: 384, col: 20, offset: 8727},
					label: "n",
					expr: &ruleRefExpr{
						pos:  position{line: 384, col: 22, offset: 8729},
						name: "ConstIdentifier",
					},
				},
			},
		},
		{
			name: "IntegerLiteral",
			pos:  position{line: 388, col: 1, offset: 8803},
			expr: &actionExpr{
				pos: position{line: 388, col: 19, offset: 8821},
				run: (*parser).callonIntegerLiteral1,
				expr: &labeledExpr{
					pos:   position{line: 388, col: 19, offset: 8821},
					label: "v",
					expr: &ruleRefExpr{
						pos:  position{line: 388, col: 21, offset: 8823},
						name: "IntLiteral",
					},
				},
			},
		},
		{
			name: "IDRef",
			pos:  position{line: 397, col: 1, offset: 8947},
			expr: &actionExpr{
				pos: position{line: 397, col: 10, offset: 8956},
				run: (*parser).callonIDRef1,
				expr: &seqExpr{
					pos: position{line: 397, col: 10, offset: 8956},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 397, col: 10, offset: 8956},
							label: "s",
							expr: &zeroOrOneExpr{
								pos: position{line: 397, col: 12, offset: 8958},
								expr: &seqExpr{
									pos: position{line: 397, col: 13, offset: 8959},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 397, col: 13, offset: 8959},
											name: "Identifier",
										},
										&litMatcher{
											pos:        position{line: 397, col: 24, offset: 8970},
											val:        ".",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 397, col: 30, offset: 8976},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 397, col: 32, offset: 8978},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierList",
			pos:  position{line: 407, col: 1, offset: 9119},
			expr: &actionExpr{
				pos: position{line: 407, col: 19, offset: 9137},
				run: (*parser).callonIdentifierList1,
				expr: &seqExpr{
					pos: position{line: 407, col: 19, offset: 9137},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 407, col: 19, offset: 9137},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 407, col: 21, offset: 9139},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 407, col: 32, offset: 9150},
							label: "ns",
							expr: &zeroOrMoreExpr{
								pos: position{line: 407, col: 35, offset: 9153},
								expr: &seqExpr{
									pos: position{line: 407, col: 36, offset: 9154},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 407, col: 36, offset: 9154},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 407, col: 38, offset: 9156},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 407, col: 42, offset: 9160},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 407, col: 44, offset: 9162},
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
			pos:  position{line: 415, col: 1, offset: 9326},
			expr: &actionExpr{
				pos: position{line: 415, col: 15, offset: 9340},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 415, col: 15, offset: 9340},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 415, col: 15, offset: 9340},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 415, col: 25, offset: 9350},
							expr: &charClassMatcher{
								pos:        position{line: 415, col: 25, offset: 9350},
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
			pos:  position{line: 419, col: 1, offset: 9398},
			expr: &actionExpr{
				pos: position{line: 419, col: 20, offset: 9417},
				run: (*parser).callonConstIdentifier1,
				expr: &seqExpr{
					pos: position{line: 419, col: 20, offset: 9417},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 419, col: 20, offset: 9417},
							val:        "[A-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 419, col: 27, offset: 9424},
							expr: &charClassMatcher{
								pos:        position{line: 419, col: 27, offset: 9424},
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
			pos:  position{line: 425, col: 1, offset: 9484},
			expr: &actionExpr{
				pos: position{line: 425, col: 15, offset: 9498},
				run: (*parser).callonIntLiteral1,
				expr: &choiceExpr{
					pos: position{line: 425, col: 16, offset: 9499},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 425, col: 16, offset: 9499},
							name: "HexLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 29, offset: 9512},
							name: "OctalLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 425, col: 44, offset: 9527},
							name: "DecimalLiteral",
						},
					},
				},
			},
		},
		{
			name: "DecimalLiteral",
			pos:  position{line: 429, col: 1, offset: 9597},
			expr: &seqExpr{
				pos: position{line: 429, col: 19, offset: 9615},
				exprs: []interface{}{
					&oneOrMoreExpr{
						pos: position{line: 429, col: 19, offset: 9615},
						expr: &charClassMatcher{
							pos:        position{line: 429, col: 19, offset: 9615},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&zeroOrMoreExpr{
						pos: position{line: 429, col: 26, offset: 9622},
						expr: &charClassMatcher{
							pos:        position{line: 429, col: 26, offset: 9622},
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
			pos:  position{line: 431, col: 1, offset: 9630},
			expr: &seqExpr{
				pos: position{line: 431, col: 15, offset: 9644},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 431, col: 15, offset: 9644},
						val:        "0x",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 431, col: 20, offset: 9649},
						expr: &charClassMatcher{
							pos:        position{line: 431, col: 20, offset: 9649},
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
			pos:  position{line: 433, col: 1, offset: 9663},
			expr: &seqExpr{
				pos: position{line: 433, col: 17, offset: 9679},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 433, col: 17, offset: 9679},
						val:        "0",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 433, col: 21, offset: 9683},
						expr: &charClassMatcher{
							pos:        position{line: 433, col: 21, offset: 9683},
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
			name: "SourceChar",
			pos:  position{line: 437, col: 1, offset: 9704},
			expr: &anyMatcher{
				line: 437, col: 15, offset: 9718,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 438, col: 1, offset: 9720},
			expr: &choiceExpr{
				pos: position{line: 438, col: 12, offset: 9731},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 438, col: 12, offset: 9731},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 438, col: 31, offset: 9750},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 439, col: 1, offset: 9768},
			expr: &seqExpr{
				pos: position{line: 439, col: 21, offset: 9788},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 439, col: 21, offset: 9788},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 439, col: 26, offset: 9793},
						expr: &seqExpr{
							pos: position{line: 439, col: 28, offset: 9795},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 439, col: 28, offset: 9795},
									expr: &litMatcher{
										pos:        position{line: 439, col: 29, offset: 9796},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 439, col: 34, offset: 9801},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 439, col: 48, offset: 9815},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 440, col: 1, offset: 9820},
			expr: &seqExpr{
				pos: position{line: 440, col: 22, offset: 9841},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 440, col: 22, offset: 9841},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 440, col: 27, offset: 9846},
						expr: &seqExpr{
							pos: position{line: 440, col: 29, offset: 9848},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 440, col: 29, offset: 9848},
									expr: &ruleRefExpr{
										pos:  position{line: 440, col: 30, offset: 9849},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 440, col: 34, offset: 9853},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 444, col: 1, offset: 9890},
			expr: &oneOrMoreExpr{
				pos: position{line: 444, col: 7, offset: 9896},
				expr: &ruleRefExpr{
					pos:  position{line: 444, col: 7, offset: 9896},
					name: "Skip",
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 445, col: 1, offset: 9902},
			expr: &zeroOrMoreExpr{
				pos: position{line: 445, col: 6, offset: 9907},
				expr: &ruleRefExpr{
					pos:  position{line: 445, col: 6, offset: 9907},
					name: "Skip",
				},
			},
		},
		{
			name: "Skip",
			pos:  position{line: 447, col: 1, offset: 9914},
			expr: &choiceExpr{
				pos: position{line: 447, col: 10, offset: 9923},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 447, col: 10, offset: 9923},
						name: "Whitespace",
					},
					&ruleRefExpr{
						pos:  position{line: 447, col: 23, offset: 9936},
						name: "EOL",
					},
					&ruleRefExpr{
						pos:  position{line: 447, col: 29, offset: 9942},
						name: "Comment",
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 449, col: 1, offset: 9952},
			expr: &charClassMatcher{
				pos:        position{line: 449, col: 15, offset: 9966},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 450, col: 1, offset: 9974},
			expr: &litMatcher{
				pos:        position{line: 450, col: 8, offset: 9981},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 451, col: 1, offset: 9986},
			expr: &notExpr{
				pos: position{line: 451, col: 8, offset: 9993},
				expr: &anyMatcher{
					line: 451, col: 9, offset: 9994,
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
		case *ast.Context:
			f.Contexts = append(f.Contexts, d)
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

func (c *current) onContextDeclaration1(n, ms interface{}) (interface{}, error) {
	m := make([]*ast.IntegerMember, 0)
	for _, i := range ms.([]interface{}) {
		m = append(m, i.(*ast.IntegerMember))
	}
	return &ast.Context{
		Name:    n.(string),
		Members: m,
	}, nil
}

func (p *parser) callonContextDeclaration1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextDeclaration1(stack["n"], stack["ms"])
}

func (c *current) onContextMember1(t, n interface{}) (interface{}, error) {
	return &ast.IntegerMember{
		Type: t.(*ast.IntType),
		Name: n.(string),
	}, nil
}

func (p *parser) callonContextMember1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextMember1(stack["t"], stack["n"])
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
	if e != nil {
		m = append(m, e.(ast.Member))
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

func (c *current) onStructEnding1(e interface{}) (interface{}, error) {
	return e, nil
}

func (p *parser) callonStructEnding1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructEnding1(stack["e"])
}

func (c *current) onStructEOS1() (interface{}, error) {
	return &ast.EOS{}, nil
}

func (p *parser) callonStructEOS1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructEOS1()
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

func (c *current) onSMVarArray1(b, n, l interface{}) (interface{}, error) {
	return &ast.VarArrayMember{
		Base:       b.(ast.ArrayBase),
		Name:       n.(string),
		Constraint: l.(ast.LengthConstraint),
	}, nil
}

func (p *parser) callonSMVarArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMVarArray1(stack["b"], stack["n"], stack["l"])
}

func (c *current) onLengthConstraint1(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonLengthConstraint1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLengthConstraint1(stack["l"])
}

func (c *current) onLeftover1(i interface{}) (interface{}, error) {
	return &ast.Leftover{Num: i.(ast.Integer)}, nil
}

func (p *parser) callonLeftover1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLeftover1(stack["i"])
}

func (c *current) onSMRemainder1(b, n interface{}) (interface{}, error) {
	return &ast.VarArrayMember{
		Base:       b.(ast.ArrayBase),
		Name:       n.(string),
		Constraint: nil,
	}, nil
}

func (p *parser) callonSMRemainder1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMRemainder1(stack["b"], stack["n"])
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

func (c *current) onSMString1(n interface{}) (interface{}, error) {
	return &ast.NulTermString{Name: n.(string)}, nil
}

func (p *parser) callonSMString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMString1(stack["n"])
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

func (c *current) onSMUnion1(n, t, l, cs interface{}) (interface{}, error) {
	u := &ast.UnionMember{
		Name:  n.(string),
		Tag:   t.(*ast.IDRef),
		Cases: cs,
	}
	if l != nil {
		u.Length = l.(ast.LengthConstraint)
	}
	return u, nil
}

func (p *parser) callonSMUnion1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMUnion1(stack["n"], stack["t"], stack["l"], stack["cs"])
}

func (c *current) onUnionLength1(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonUnionLength1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionLength1(stack["l"])
}

func (c *current) onUnionMember1(cse, fs interface{}) (interface{}, error) {
	uc := &ast.UnionCase{
		Fields: fs,
	}
	if cse != nil {
		uc.Case = cse.(*ast.IntegerList)
	}
	return uc, nil
}

func (p *parser) callonUnionMember1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionMember1(stack["cse"], stack["fs"])
}

func (c *current) onUnionCase2(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonUnionCase2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionCase2(stack["l"])
}

func (c *current) onUnionCase5() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonUnionCase5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionCase5()
}

func (c *current) onUnionBody2() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonUnionBody2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionBody2()
}

func (c *current) onUnionBody6() (interface{}, error) {
	return &ast.Fail{}, nil
}

func (p *parser) callonUnionBody6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionBody6()
}

func (c *current) onUnionBody11() (interface{}, error) {
	return &ast.Ignore{}, nil
}

func (p *parser) callonUnionBody11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionBody11()
}

func (c *current) onUnionBody16(fs interface{}) (interface{}, error) {
	return fs, nil
}

func (p *parser) callonUnionBody16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionBody16(stack["fs"])
}

func (c *current) onUnionFields1(ms, e interface{}) (interface{}, error) {
	fs := []ast.Member{}
	for _, i := range ms.([]interface{}) {
		fs = append(fs, i.(ast.Member))
	}
	if e != nil {
		fs = append(fs, e)
	}
	return fs, nil
}

func (p *parser) callonUnionFields1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionFields1(stack["ms"], stack["e"])
}

func (c *current) onUnionField1(m interface{}) (interface{}, error) {
	return m, nil
}

func (p *parser) callonUnionField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnionField1(stack["m"])
}

func (c *current) onExtentSpec2() (interface{}, error) {
	return &ast.Ignore{}, nil
}

func (p *parser) callonExtentSpec2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExtentSpec2()
}

func (c *current) onExtentSpec7(r interface{}) (interface{}, error) {
	return r, nil
}

func (p *parser) callonExtentSpec7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExtentSpec7(stack["r"])
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
		r.High = hi.([]interface{})[3].(ast.Integer)
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

func (c *current) onIDRef1(s, n interface{}) (interface{}, error) {
	r := &ast.IDRef{
		Name: n.(string),
	}
	if s != nil {
		r.Scope = s.([]interface{})[0].(string)
	}
	return r, nil
}

func (p *parser) callonIDRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIDRef1(stack["s"], stack["n"])
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
