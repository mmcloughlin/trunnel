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
						&ruleRefExpr{
							pos:  position{line: 28, col: 9, offset: 595},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 28, col: 11, offset: 597},
							label: "ds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 28, col: 14, offset: 600},
								expr: &ruleRefExpr{
									pos:  position{line: 28, col: 14, offset: 600},
									name: "Declaration",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 27, offset: 613},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 28, col: 29, offset: 615},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Declaration",
			pos:  position{line: 55, col: 1, offset: 1412},
			expr: &actionExpr{
				pos: position{line: 55, col: 16, offset: 1427},
				run: (*parser).callonDeclaration1,
				expr: &seqExpr{
					pos: position{line: 55, col: 16, offset: 1427},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 55, col: 16, offset: 1427},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 55, col: 18, offset: 1429},
							label: "d",
							expr: &choiceExpr{
								pos: position{line: 55, col: 21, offset: 1432},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 55, col: 21, offset: 1432},
										name: "ConstDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 40, offset: 1451},
										name: "ContextDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 61, offset: 1472},
										name: "StructDeclaration",
									},
									&ruleRefExpr{
										pos:  position{line: 55, col: 81, offset: 1492},
										name: "PragmaDeclaration",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 100, offset: 1511},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "ConstDeclaration",
			pos:  position{line: 61, col: 1, offset: 1576},
			expr: &actionExpr{
				pos: position{line: 61, col: 21, offset: 1596},
				run: (*parser).callonConstDeclaration1,
				expr: &seqExpr{
					pos: position{line: 61, col: 21, offset: 1596},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 61, col: 21, offset: 1596},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 29, offset: 1604},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 32, offset: 1607},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 34, offset: 1609},
								name: "ConstIdentifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 50, offset: 1625},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 52, offset: 1627},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 56, offset: 1631},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 58, offset: 1633},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 61, col: 60, offset: 1635},
								name: "IntLiteral",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 71, offset: 1646},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 73, offset: 1648},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContextDeclaration",
			pos:  position{line: 73, col: 1, offset: 1909},
			expr: &actionExpr{
				pos: position{line: 73, col: 23, offset: 1931},
				run: (*parser).callonContextDeclaration1,
				expr: &seqExpr{
					pos: position{line: 73, col: 23, offset: 1931},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 73, col: 23, offset: 1931},
							val:        "context",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 33, offset: 1941},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 73, col: 36, offset: 1944},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 38, offset: 1946},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 49, offset: 1957},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 73, col: 52, offset: 1960},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 56, offset: 1964},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 73, col: 58, offset: 1966},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 73, col: 61, offset: 1969},
								expr: &ruleRefExpr{
									pos:  position{line: 73, col: 61, offset: 1969},
									name: "ContextMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 76, offset: 1984},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 73, col: 78, offset: 1986},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 73, col: 82, offset: 1990},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 73, col: 84, offset: 1992},
							expr: &litMatcher{
								pos:        position{line: 73, col: 84, offset: 1992},
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
			pos:  position{line: 84, col: 1, offset: 2196},
			expr: &actionExpr{
				pos: position{line: 84, col: 18, offset: 2213},
				run: (*parser).callonContextMember1,
				expr: &seqExpr{
					pos: position{line: 84, col: 18, offset: 2213},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 84, col: 18, offset: 2213},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 84, col: 20, offset: 2215},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 84, col: 22, offset: 2217},
								name: "IntType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 30, offset: 2225},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 84, col: 33, offset: 2228},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 84, col: 35, offset: 2230},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 46, offset: 2241},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 84, col: 48, offset: 2243},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 84, col: 52, offset: 2247},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "StructDeclaration",
			pos:  position{line: 94, col: 1, offset: 2441},
			expr: &actionExpr{
				pos: position{line: 94, col: 22, offset: 2462},
				run: (*parser).callonStructDeclaration1,
				expr: &seqExpr{
					pos: position{line: 94, col: 22, offset: 2462},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 94, col: 22, offset: 2462},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 94, col: 24, offset: 2464},
								name: "StructDecl",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 94, col: 35, offset: 2475},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 94, col: 37, offset: 2477},
							expr: &litMatcher{
								pos:        position{line: 94, col: 37, offset: 2477},
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
			pos:  position{line: 100, col: 1, offset: 2595},
			expr: &actionExpr{
				pos: position{line: 100, col: 15, offset: 2609},
				run: (*parser).callonStructDecl1,
				expr: &seqExpr{
					pos: position{line: 100, col: 15, offset: 2609},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 100, col: 15, offset: 2609},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 100, col: 17, offset: 2611},
								name: "StructIdentifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 100, col: 34, offset: 2628},
							label: "ctx",
							expr: &zeroOrOneExpr{
								pos: position{line: 100, col: 38, offset: 2632},
								expr: &seqExpr{
									pos: position{line: 100, col: 39, offset: 2633},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 100, col: 39, offset: 2633},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 100, col: 42, offset: 2636},
											name: "ContextRefs",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 56, offset: 2650},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 100, col: 58, offset: 2652},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 100, col: 62, offset: 2656},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 100, col: 65, offset: 2659},
								expr: &ruleRefExpr{
									pos:  position{line: 100, col: 65, offset: 2659},
									name: "StructMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 79, offset: 2673},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 100, col: 81, offset: 2675},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 100, col: 83, offset: 2677},
								expr: &ruleRefExpr{
									pos:  position{line: 100, col: 83, offset: 2677},
									name: "StructEnding",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 97, offset: 2691},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 100, col: 99, offset: 2693},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructIdentifier",
			pos:  position{line: 121, col: 1, offset: 3018},
			expr: &actionExpr{
				pos: position{line: 121, col: 21, offset: 3038},
				run: (*parser).callonStructIdentifier1,
				expr: &seqExpr{
					pos: position{line: 121, col: 21, offset: 3038},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 121, col: 21, offset: 3038},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 30, offset: 3047},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 121, col: 33, offset: 3050},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 121, col: 35, offset: 3052},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "PragmaDeclaration",
			pos:  position{line: 125, col: 1, offset: 3084},
			expr: &actionExpr{
				pos: position{line: 125, col: 22, offset: 3105},
				run: (*parser).callonPragmaDeclaration1,
				expr: &seqExpr{
					pos: position{line: 125, col: 22, offset: 3105},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 125, col: 22, offset: 3105},
							val:        "trunnel",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 125, col: 32, offset: 3115},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 125, col: 35, offset: 3118},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 125, col: 37, offset: 3120},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 125, col: 48, offset: 3131},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 125, col: 51, offset: 3134},
							label: "opts",
							expr: &ruleRefExpr{
								pos:  position{line: 125, col: 56, offset: 3139},
								name: "IdentifierList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 125, col: 71, offset: 3154},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 125, col: 73, offset: 3156},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContextRefs",
			pos:  position{line: 135, col: 1, offset: 3319},
			expr: &actionExpr{
				pos: position{line: 135, col: 16, offset: 3334},
				run: (*parser).callonContextRefs1,
				expr: &seqExpr{
					pos: position{line: 135, col: 16, offset: 3334},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 135, col: 16, offset: 3334},
							val:        "with",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 23, offset: 3341},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 135, col: 26, offset: 3344},
							val:        "context",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 36, offset: 3354},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 135, col: 39, offset: 3357},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 135, col: 42, offset: 3360},
								name: "IdentifierList",
							},
						},
					},
				},
			},
		},
		{
			name: "StructMember",
			pos:  position{line: 149, col: 1, offset: 3679},
			expr: &actionExpr{
				pos: position{line: 149, col: 17, offset: 3695},
				run: (*parser).callonStructMember1,
				expr: &seqExpr{
					pos: position{line: 149, col: 17, offset: 3695},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 149, col: 17, offset: 3695},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 149, col: 19, offset: 3697},
							label: "m",
							expr: &choiceExpr{
								pos: position{line: 149, col: 22, offset: 3700},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 149, col: 22, offset: 3700},
										name: "SMArray",
									},
									&ruleRefExpr{
										pos:  position{line: 149, col: 32, offset: 3710},
										name: "SMInteger",
									},
									&ruleRefExpr{
										pos:  position{line: 149, col: 44, offset: 3722},
										name: "SMPosition",
									},
									&ruleRefExpr{
										pos:  position{line: 149, col: 57, offset: 3735},
										name: "SMString",
									},
									&ruleRefExpr{
										pos:  position{line: 149, col: 68, offset: 3746},
										name: "SMStruct",
									},
									&ruleRefExpr{
										pos:  position{line: 149, col: 79, offset: 3757},
										name: "SMUnion",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 149, col: 88, offset: 3766},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 149, col: 90, offset: 3768},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructEnding",
			pos:  position{line: 157, col: 1, offset: 3877},
			expr: &actionExpr{
				pos: position{line: 157, col: 17, offset: 3893},
				run: (*parser).callonStructEnding1,
				expr: &seqExpr{
					pos: position{line: 157, col: 17, offset: 3893},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 157, col: 17, offset: 3893},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 157, col: 19, offset: 3895},
							label: "e",
							expr: &choiceExpr{
								pos: position{line: 157, col: 22, offset: 3898},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 157, col: 22, offset: 3898},
										name: "SMRemainder",
									},
									&ruleRefExpr{
										pos:  position{line: 157, col: 36, offset: 3912},
										name: "StructEOS",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 157, col: 47, offset: 3923},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 157, col: 49, offset: 3925},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "StructEOS",
			pos:  position{line: 161, col: 1, offset: 3950},
			expr: &actionExpr{
				pos: position{line: 161, col: 14, offset: 3963},
				run: (*parser).callonStructEOS1,
				expr: &litMatcher{
					pos:        position{line: 161, col: 14, offset: 3963},
					val:        "eos",
					ignoreCase: false,
				},
			},
		},
		{
			name: "SMArray",
			pos:  position{line: 168, col: 1, offset: 4056},
			expr: &actionExpr{
				pos: position{line: 168, col: 12, offset: 4067},
				run: (*parser).callonSMArray1,
				expr: &labeledExpr{
					pos:   position{line: 168, col: 12, offset: 4067},
					label: "a",
					expr: &choiceExpr{
						pos: position{line: 168, col: 15, offset: 4070},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 168, col: 15, offset: 4070},
								name: "SMFixedArray",
							},
							&ruleRefExpr{
								pos:  position{line: 168, col: 30, offset: 4085},
								name: "SMVarArray",
							},
						},
					},
				},
			},
		},
		{
			name: "SMFixedArray",
			pos:  position{line: 174, col: 1, offset: 4165},
			expr: &actionExpr{
				pos: position{line: 174, col: 17, offset: 4181},
				run: (*parser).callonSMFixedArray1,
				expr: &seqExpr{
					pos: position{line: 174, col: 17, offset: 4181},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 174, col: 17, offset: 4181},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 19, offset: 4183},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 174, col: 29, offset: 4193},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 174, col: 32, offset: 4196},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 34, offset: 4198},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 174, col: 45, offset: 4209},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 174, col: 47, offset: 4211},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 174, col: 51, offset: 4215},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 174, col: 53, offset: 4217},
								name: "Integer",
							},
						},
						&litMatcher{
							pos:        position{line: 174, col: 61, offset: 4225},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SMVarArray",
			pos:  position{line: 185, col: 1, offset: 4445},
			expr: &actionExpr{
				pos: position{line: 185, col: 15, offset: 4459},
				run: (*parser).callonSMVarArray1,
				expr: &seqExpr{
					pos: position{line: 185, col: 15, offset: 4459},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 185, col: 15, offset: 4459},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 185, col: 17, offset: 4461},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 185, col: 27, offset: 4471},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 185, col: 30, offset: 4474},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 185, col: 32, offset: 4476},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 185, col: 43, offset: 4487},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 185, col: 45, offset: 4489},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 185, col: 49, offset: 4493},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 185, col: 51, offset: 4495},
								name: "LengthConstraint",
							},
						},
						&litMatcher{
							pos:        position{line: 185, col: 68, offset: 4512},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "LengthConstraint",
			pos:  position{line: 193, col: 1, offset: 4653},
			expr: &actionExpr{
				pos: position{line: 193, col: 21, offset: 4673},
				run: (*parser).callonLengthConstraint1,
				expr: &labeledExpr{
					pos:   position{line: 193, col: 21, offset: 4673},
					label: "l",
					expr: &choiceExpr{
						pos: position{line: 193, col: 24, offset: 4676},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 193, col: 24, offset: 4676},
								name: "Leftover",
							},
							&ruleRefExpr{
								pos:  position{line: 193, col: 35, offset: 4687},
								name: "IDRef",
							},
						},
					},
				},
			},
		},
		{
			name: "Leftover",
			pos:  position{line: 197, col: 1, offset: 4715},
			expr: &actionExpr{
				pos: position{line: 197, col: 13, offset: 4727},
				run: (*parser).callonLeftover1,
				expr: &seqExpr{
					pos: position{line: 197, col: 13, offset: 4727},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 197, col: 13, offset: 4727},
							val:        "..-",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 197, col: 19, offset: 4733},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 197, col: 21, offset: 4735},
								name: "Integer",
							},
						},
					},
				},
			},
		},
		{
			name: "SMRemainder",
			pos:  position{line: 203, col: 1, offset: 4865},
			expr: &actionExpr{
				pos: position{line: 203, col: 16, offset: 4880},
				run: (*parser).callonSMRemainder1,
				expr: &seqExpr{
					pos: position{line: 203, col: 16, offset: 4880},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 203, col: 16, offset: 4880},
							label: "b",
							expr: &ruleRefExpr{
								pos:  position{line: 203, col: 18, offset: 4882},
								name: "ArrayBase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 28, offset: 4892},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 203, col: 31, offset: 4895},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 203, col: 33, offset: 4897},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 44, offset: 4908},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 203, col: 46, offset: 4910},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 50, offset: 4914},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 203, col: 52, offset: 4916},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ArrayBase",
			pos:  position{line: 216, col: 1, offset: 5147},
			expr: &actionExpr{
				pos: position{line: 216, col: 14, offset: 5160},
				run: (*parser).callonArrayBase1,
				expr: &labeledExpr{
					pos:   position{line: 216, col: 14, offset: 5160},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 216, col: 17, offset: 5163},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 216, col: 17, offset: 5163},
								name: "IntType",
							},
							&ruleRefExpr{
								pos:  position{line: 216, col: 27, offset: 5173},
								name: "CharType",
							},
							&ruleRefExpr{
								pos:  position{line: 216, col: 38, offset: 5184},
								name: "StructRef",
							},
						},
					},
				},
			},
		},
		{
			name: "SMInteger",
			pos:  position{line: 223, col: 1, offset: 5264},
			expr: &actionExpr{
				pos: position{line: 223, col: 14, offset: 5277},
				run: (*parser).callonSMInteger1,
				expr: &seqExpr{
					pos: position{line: 223, col: 14, offset: 5277},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 223, col: 14, offset: 5277},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 223, col: 16, offset: 5279},
								name: "IntType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 24, offset: 5287},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 223, col: 26, offset: 5289},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 223, col: 28, offset: 5291},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 39, offset: 5302},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 223, col: 41, offset: 5304},
							label: "cst",
							expr: &zeroOrOneExpr{
								pos: position{line: 223, col: 45, offset: 5308},
								expr: &ruleRefExpr{
									pos:  position{line: 223, col: 45, offset: 5308},
									name: "IntConstraint",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SMPosition",
			pos:  position{line: 236, col: 1, offset: 5519},
			expr: &actionExpr{
				pos: position{line: 236, col: 15, offset: 5533},
				run: (*parser).callonSMPosition1,
				expr: &seqExpr{
					pos: position{line: 236, col: 15, offset: 5533},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 236, col: 15, offset: 5533},
							val:        "@ptr",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 236, col: 22, offset: 5540},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 236, col: 25, offset: 5543},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 236, col: 27, offset: 5545},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "SMString",
			pos:  position{line: 244, col: 1, offset: 5640},
			expr: &actionExpr{
				pos: position{line: 244, col: 13, offset: 5652},
				run: (*parser).callonSMString1,
				expr: &seqExpr{
					pos: position{line: 244, col: 13, offset: 5652},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 244, col: 13, offset: 5652},
							val:        "nulterm",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 244, col: 23, offset: 5662},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 244, col: 26, offset: 5665},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 244, col: 28, offset: 5667},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "SMStruct",
			pos:  position{line: 251, col: 1, offset: 5796},
			expr: &actionExpr{
				pos: position{line: 251, col: 13, offset: 5808},
				run: (*parser).callonSMStruct1,
				expr: &seqExpr{
					pos: position{line: 251, col: 13, offset: 5808},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 251, col: 13, offset: 5808},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 251, col: 15, offset: 5810},
								name: "StructRef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 251, col: 25, offset: 5820},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 251, col: 28, offset: 5823},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 251, col: 30, offset: 5825},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "SMUnion",
			pos:  position{line: 260, col: 1, offset: 5997},
			expr: &actionExpr{
				pos: position{line: 260, col: 12, offset: 6008},
				run: (*parser).callonSMUnion1,
				expr: &seqExpr{
					pos: position{line: 260, col: 12, offset: 6008},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 260, col: 12, offset: 6008},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 20, offset: 6016},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 23, offset: 6019},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 25, offset: 6021},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 36, offset: 6032},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 260, col: 38, offset: 6034},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 260, col: 42, offset: 6038},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 44, offset: 6040},
								name: "IDRef",
							},
						},
						&litMatcher{
							pos:        position{line: 260, col: 50, offset: 6046},
							val:        "]",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 54, offset: 6050},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 56, offset: 6052},
							label: "l",
							expr: &zeroOrOneExpr{
								pos: position{line: 260, col: 58, offset: 6054},
								expr: &ruleRefExpr{
									pos:  position{line: 260, col: 58, offset: 6054},
									name: "UnionLength",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 71, offset: 6067},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 260, col: 73, offset: 6069},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 77, offset: 6073},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 79, offset: 6075},
							label: "cs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 260, col: 82, offset: 6078},
								expr: &ruleRefExpr{
									pos:  position{line: 260, col: 82, offset: 6078},
									name: "UnionMember",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 95, offset: 6091},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 260, col: 97, offset: 6093},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "UnionLength",
			pos:  position{line: 272, col: 1, offset: 6269},
			expr: &actionExpr{
				pos: position{line: 272, col: 16, offset: 6284},
				run: (*parser).callonUnionLength1,
				expr: &seqExpr{
					pos: position{line: 272, col: 16, offset: 6284},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 272, col: 16, offset: 6284},
							val:        "with",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 272, col: 23, offset: 6291},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 272, col: 26, offset: 6294},
							val:        "length",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 272, col: 35, offset: 6303},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 272, col: 38, offset: 6306},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 272, col: 40, offset: 6308},
								name: "LengthConstraint",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionMember",
			pos:  position{line: 278, col: 1, offset: 6421},
			expr: &actionExpr{
				pos: position{line: 278, col: 16, offset: 6436},
				run: (*parser).callonUnionMember1,
				expr: &seqExpr{
					pos: position{line: 278, col: 16, offset: 6436},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 278, col: 16, offset: 6436},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 278, col: 18, offset: 6438},
							label: "cse",
							expr: &ruleRefExpr{
								pos:  position{line: 278, col: 22, offset: 6442},
								name: "UnionCase",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 32, offset: 6452},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 278, col: 34, offset: 6454},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 38, offset: 6458},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 278, col: 40, offset: 6460},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 278, col: 43, offset: 6463},
								name: "UnionBody",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 278, col: 53, offset: 6473},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "UnionCase",
			pos:  position{line: 291, col: 1, offset: 6653},
			expr: &choiceExpr{
				pos: position{line: 291, col: 14, offset: 6666},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 291, col: 14, offset: 6666},
						run: (*parser).callonUnionCase2,
						expr: &labeledExpr{
							pos:   position{line: 291, col: 14, offset: 6666},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 291, col: 16, offset: 6668},
								name: "IntList",
							},
						},
					},
					&actionExpr{
						pos: position{line: 293, col: 5, offset: 6698},
						run: (*parser).callonUnionCase5,
						expr: &litMatcher{
							pos:        position{line: 293, col: 5, offset: 6698},
							val:        "default",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "UnionBody",
			pos:  position{line: 304, col: 1, offset: 6946},
			expr: &choiceExpr{
				pos: position{line: 304, col: 14, offset: 6959},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 304, col: 14, offset: 6959},
						run: (*parser).callonUnionBody2,
						expr: &seqExpr{
							pos: position{line: 304, col: 14, offset: 6959},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 304, col: 14, offset: 6959},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 304, col: 16, offset: 6961},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 306, col: 5, offset: 6989},
						run: (*parser).callonUnionBody6,
						expr: &seqExpr{
							pos: position{line: 306, col: 5, offset: 6989},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 306, col: 5, offset: 6989},
									val:        "fail",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 306, col: 12, offset: 6996},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 306, col: 14, offset: 6998},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 308, col: 5, offset: 7034},
						run: (*parser).callonUnionBody11,
						expr: &seqExpr{
							pos: position{line: 308, col: 5, offset: 7034},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 308, col: 5, offset: 7034},
									val:        "ignore",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 308, col: 14, offset: 7043},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 308, col: 16, offset: 7045},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 310, col: 5, offset: 7083},
						run: (*parser).callonUnionBody16,
						expr: &labeledExpr{
							pos:   position{line: 310, col: 5, offset: 7083},
							label: "fs",
							expr: &ruleRefExpr{
								pos:  position{line: 310, col: 8, offset: 7086},
								name: "UnionFields",
							},
						},
					},
				},
			},
		},
		{
			name: "UnionFields",
			pos:  position{line: 314, col: 1, offset: 7120},
			expr: &actionExpr{
				pos: position{line: 314, col: 16, offset: 7135},
				run: (*parser).callonUnionFields1,
				expr: &seqExpr{
					pos: position{line: 314, col: 16, offset: 7135},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 314, col: 16, offset: 7135},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 314, col: 19, offset: 7138},
								expr: &ruleRefExpr{
									pos:  position{line: 314, col: 19, offset: 7138},
									name: "UnionField",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 314, col: 31, offset: 7150},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 314, col: 33, offset: 7152},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 314, col: 35, offset: 7154},
								expr: &ruleRefExpr{
									pos:  position{line: 314, col: 35, offset: 7154},
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
			pos:  position{line: 331, col: 1, offset: 7483},
			expr: &actionExpr{
				pos: position{line: 331, col: 15, offset: 7497},
				run: (*parser).callonUnionField1,
				expr: &seqExpr{
					pos: position{line: 331, col: 15, offset: 7497},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 331, col: 15, offset: 7497},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 17, offset: 7499},
							label: "m",
							expr: &choiceExpr{
								pos: position{line: 331, col: 20, offset: 7502},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 331, col: 20, offset: 7502},
										name: "SMArray",
									},
									&ruleRefExpr{
										pos:  position{line: 331, col: 30, offset: 7512},
										name: "SMInteger",
									},
									&ruleRefExpr{
										pos:  position{line: 331, col: 42, offset: 7524},
										name: "SMString",
									},
									&ruleRefExpr{
										pos:  position{line: 331, col: 53, offset: 7535},
										name: "SMStruct",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 63, offset: 7545},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 331, col: 65, offset: 7547},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExtentSpec",
			pos:  position{line: 339, col: 1, offset: 7659},
			expr: &choiceExpr{
				pos: position{line: 339, col: 15, offset: 7673},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 339, col: 15, offset: 7673},
						run: (*parser).callonExtentSpec2,
						expr: &seqExpr{
							pos: position{line: 339, col: 15, offset: 7673},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 339, col: 15, offset: 7673},
									val:        "...",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 339, col: 21, offset: 7679},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 339, col: 23, offset: 7681},
									val:        ";",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 341, col: 5, offset: 7719},
						run: (*parser).callonExtentSpec7,
						expr: &seqExpr{
							pos: position{line: 341, col: 5, offset: 7719},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 341, col: 5, offset: 7719},
									label: "r",
									expr: &ruleRefExpr{
										pos:  position{line: 341, col: 7, offset: 7721},
										name: "SMRemainder",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 341, col: 19, offset: 7733},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 341, col: 21, offset: 7735},
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
			pos:  position{line: 345, col: 1, offset: 7760},
			expr: &choiceExpr{
				pos: position{line: 345, col: 14, offset: 7773},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 345, col: 14, offset: 7773},
						run: (*parser).callonStructRef2,
						expr: &labeledExpr{
							pos:   position{line: 345, col: 14, offset: 7773},
							label: "s",
							expr: &ruleRefExpr{
								pos:  position{line: 345, col: 16, offset: 7775},
								name: "StructDecl",
							},
						},
					},
					&actionExpr{
						pos: position{line: 348, col: 5, offset: 7880},
						run: (*parser).callonStructRef5,
						expr: &labeledExpr{
							pos:   position{line: 348, col: 5, offset: 7880},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 348, col: 7, offset: 7882},
								name: "StructIdentifier",
							},
						},
					},
				},
			},
		},
		{
			name: "CharType",
			pos:  position{line: 352, col: 1, offset: 7951},
			expr: &actionExpr{
				pos: position{line: 352, col: 13, offset: 7963},
				run: (*parser).callonCharType1,
				expr: &litMatcher{
					pos:        position{line: 352, col: 13, offset: 7963},
					val:        "char",
					ignoreCase: false,
				},
			},
		},
		{
			name: "IntType",
			pos:  position{line: 361, col: 1, offset: 8085},
			expr: &actionExpr{
				pos: position{line: 361, col: 12, offset: 8096},
				run: (*parser).callonIntType1,
				expr: &seqExpr{
					pos: position{line: 361, col: 12, offset: 8096},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 361, col: 12, offset: 8096},
							val:        "u",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 361, col: 16, offset: 8100},
							label: "b",
							expr: &choiceExpr{
								pos: position{line: 361, col: 19, offset: 8103},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 361, col: 19, offset: 8103},
										val:        "8",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 361, col: 25, offset: 8109},
										val:        "16",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 361, col: 32, offset: 8116},
										val:        "32",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 361, col: 39, offset: 8123},
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
			pos:  position{line: 369, col: 1, offset: 8281},
			expr: &actionExpr{
				pos: position{line: 369, col: 18, offset: 8298},
				run: (*parser).callonIntConstraint1,
				expr: &seqExpr{
					pos: position{line: 369, col: 18, offset: 8298},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 369, col: 18, offset: 8298},
							val:        "IN",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 23, offset: 8303},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 369, col: 26, offset: 8306},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 30, offset: 8310},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 369, col: 32, offset: 8312},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 369, col: 34, offset: 8314},
								name: "IntList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 42, offset: 8322},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 369, col: 44, offset: 8324},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "IntList",
			pos:  position{line: 376, col: 1, offset: 8420},
			expr: &actionExpr{
				pos: position{line: 376, col: 12, offset: 8431},
				run: (*parser).callonIntList1,
				expr: &seqExpr{
					pos: position{line: 376, col: 12, offset: 8431},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 376, col: 12, offset: 8431},
							label: "m",
							expr: &ruleRefExpr{
								pos:  position{line: 376, col: 14, offset: 8433},
								name: "IntListMember",
							},
						},
						&labeledExpr{
							pos:   position{line: 376, col: 28, offset: 8447},
							label: "ms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 376, col: 31, offset: 8450},
								expr: &seqExpr{
									pos: position{line: 376, col: 32, offset: 8451},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 376, col: 32, offset: 8451},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 376, col: 34, offset: 8453},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 376, col: 38, offset: 8457},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 376, col: 40, offset: 8459},
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
			pos:  position{line: 387, col: 1, offset: 8749},
			expr: &actionExpr{
				pos: position{line: 387, col: 18, offset: 8766},
				run: (*parser).callonIntListMember1,
				expr: &seqExpr{
					pos: position{line: 387, col: 18, offset: 8766},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 387, col: 18, offset: 8766},
							label: "lo",
							expr: &ruleRefExpr{
								pos:  position{line: 387, col: 21, offset: 8769},
								name: "Integer",
							},
						},
						&labeledExpr{
							pos:   position{line: 387, col: 29, offset: 8777},
							label: "hi",
							expr: &zeroOrOneExpr{
								pos: position{line: 387, col: 32, offset: 8780},
								expr: &seqExpr{
									pos: position{line: 387, col: 34, offset: 8782},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 387, col: 34, offset: 8782},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 387, col: 36, offset: 8784},
											val:        "..",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 387, col: 41, offset: 8789},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 387, col: 43, offset: 8791},
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
			pos:  position{line: 400, col: 1, offset: 8996},
			expr: &actionExpr{
				pos: position{line: 400, col: 12, offset: 9007},
				run: (*parser).callonInteger1,
				expr: &labeledExpr{
					pos:   position{line: 400, col: 12, offset: 9007},
					label: "i",
					expr: &choiceExpr{
						pos: position{line: 400, col: 15, offset: 9010},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 400, col: 15, offset: 9010},
								name: "IntegerConstRef",
							},
							&ruleRefExpr{
								pos:  position{line: 400, col: 33, offset: 9028},
								name: "IntegerLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "IntegerConstRef",
			pos:  position{line: 404, col: 1, offset: 9065},
			expr: &actionExpr{
				pos: position{line: 404, col: 20, offset: 9084},
				run: (*parser).callonIntegerConstRef1,
				expr: &labeledExpr{
					pos:   position{line: 404, col: 20, offset: 9084},
					label: "n",
					expr: &ruleRefExpr{
						pos:  position{line: 404, col: 22, offset: 9086},
						name: "ConstIdentifier",
					},
				},
			},
		},
		{
			name: "IntegerLiteral",
			pos:  position{line: 408, col: 1, offset: 9160},
			expr: &actionExpr{
				pos: position{line: 408, col: 19, offset: 9178},
				run: (*parser).callonIntegerLiteral1,
				expr: &labeledExpr{
					pos:   position{line: 408, col: 19, offset: 9178},
					label: "v",
					expr: &ruleRefExpr{
						pos:  position{line: 408, col: 21, offset: 9180},
						name: "IntLiteral",
					},
				},
			},
		},
		{
			name: "IDRef",
			pos:  position{line: 417, col: 1, offset: 9304},
			expr: &actionExpr{
				pos: position{line: 417, col: 10, offset: 9313},
				run: (*parser).callonIDRef1,
				expr: &seqExpr{
					pos: position{line: 417, col: 10, offset: 9313},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 417, col: 10, offset: 9313},
							label: "s",
							expr: &zeroOrOneExpr{
								pos: position{line: 417, col: 12, offset: 9315},
								expr: &seqExpr{
									pos: position{line: 417, col: 13, offset: 9316},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 417, col: 13, offset: 9316},
											name: "Identifier",
										},
										&litMatcher{
											pos:        position{line: 417, col: 24, offset: 9327},
											val:        ".",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 417, col: 30, offset: 9333},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 417, col: 32, offset: 9335},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "IdentifierList",
			pos:  position{line: 427, col: 1, offset: 9476},
			expr: &actionExpr{
				pos: position{line: 427, col: 19, offset: 9494},
				run: (*parser).callonIdentifierList1,
				expr: &seqExpr{
					pos: position{line: 427, col: 19, offset: 9494},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 427, col: 19, offset: 9494},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 427, col: 21, offset: 9496},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 427, col: 32, offset: 9507},
							label: "ns",
							expr: &zeroOrMoreExpr{
								pos: position{line: 427, col: 35, offset: 9510},
								expr: &seqExpr{
									pos: position{line: 427, col: 36, offset: 9511},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 427, col: 36, offset: 9511},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 427, col: 38, offset: 9513},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 427, col: 42, offset: 9517},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 427, col: 44, offset: 9519},
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
			pos:  position{line: 435, col: 1, offset: 9683},
			expr: &actionExpr{
				pos: position{line: 435, col: 15, offset: 9697},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 435, col: 15, offset: 9697},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 435, col: 15, offset: 9697},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 435, col: 25, offset: 9707},
							expr: &charClassMatcher{
								pos:        position{line: 435, col: 25, offset: 9707},
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
			pos:  position{line: 439, col: 1, offset: 9755},
			expr: &actionExpr{
				pos: position{line: 439, col: 20, offset: 9774},
				run: (*parser).callonConstIdentifier1,
				expr: &seqExpr{
					pos: position{line: 439, col: 20, offset: 9774},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 439, col: 20, offset: 9774},
							val:        "[A-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 439, col: 27, offset: 9781},
							expr: &charClassMatcher{
								pos:        position{line: 439, col: 27, offset: 9781},
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
			pos:  position{line: 445, col: 1, offset: 9841},
			expr: &actionExpr{
				pos: position{line: 445, col: 15, offset: 9855},
				run: (*parser).callonIntLiteral1,
				expr: &choiceExpr{
					pos: position{line: 445, col: 16, offset: 9856},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 445, col: 16, offset: 9856},
							name: "HexLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 445, col: 29, offset: 9869},
							name: "OctalLiteral",
						},
						&ruleRefExpr{
							pos:  position{line: 445, col: 44, offset: 9884},
							name: "DecimalLiteral",
						},
					},
				},
			},
		},
		{
			name: "DecimalLiteral",
			pos:  position{line: 449, col: 1, offset: 9954},
			expr: &oneOrMoreExpr{
				pos: position{line: 449, col: 19, offset: 9972},
				expr: &charClassMatcher{
					pos:        position{line: 449, col: 19, offset: 9972},
					val:        "[0-9]",
					ranges:     []rune{'0', '9'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "HexLiteral",
			pos:  position{line: 451, col: 1, offset: 9980},
			expr: &seqExpr{
				pos: position{line: 451, col: 15, offset: 9994},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 451, col: 15, offset: 9994},
						val:        "0x",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 451, col: 20, offset: 9999},
						expr: &charClassMatcher{
							pos:        position{line: 451, col: 20, offset: 9999},
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
			pos:  position{line: 453, col: 1, offset: 10013},
			expr: &seqExpr{
				pos: position{line: 453, col: 17, offset: 10029},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 453, col: 17, offset: 10029},
						val:        "0",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 453, col: 21, offset: 10033},
						expr: &charClassMatcher{
							pos:        position{line: 453, col: 21, offset: 10033},
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
			pos:  position{line: 457, col: 1, offset: 10054},
			expr: &anyMatcher{
				line: 457, col: 15, offset: 10068,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 458, col: 1, offset: 10070},
			expr: &choiceExpr{
				pos: position{line: 458, col: 12, offset: 10081},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 458, col: 12, offset: 10081},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 458, col: 31, offset: 10100},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 459, col: 1, offset: 10118},
			expr: &seqExpr{
				pos: position{line: 459, col: 21, offset: 10138},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 459, col: 21, offset: 10138},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 459, col: 26, offset: 10143},
						expr: &seqExpr{
							pos: position{line: 459, col: 28, offset: 10145},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 459, col: 28, offset: 10145},
									expr: &litMatcher{
										pos:        position{line: 459, col: 29, offset: 10146},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 459, col: 34, offset: 10151},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 459, col: 48, offset: 10165},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 460, col: 1, offset: 10170},
			expr: &seqExpr{
				pos: position{line: 460, col: 22, offset: 10191},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 460, col: 22, offset: 10191},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 460, col: 27, offset: 10196},
						expr: &seqExpr{
							pos: position{line: 460, col: 29, offset: 10198},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 460, col: 29, offset: 10198},
									expr: &ruleRefExpr{
										pos:  position{line: 460, col: 30, offset: 10199},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 460, col: 34, offset: 10203},
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
			pos:  position{line: 464, col: 1, offset: 10240},
			expr: &oneOrMoreExpr{
				pos: position{line: 464, col: 7, offset: 10246},
				expr: &ruleRefExpr{
					pos:  position{line: 464, col: 7, offset: 10246},
					name: "Skip",
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 465, col: 1, offset: 10252},
			expr: &zeroOrMoreExpr{
				pos: position{line: 465, col: 6, offset: 10257},
				expr: &ruleRefExpr{
					pos:  position{line: 465, col: 6, offset: 10257},
					name: "Skip",
				},
			},
		},
		{
			name: "Skip",
			pos:  position{line: 467, col: 1, offset: 10264},
			expr: &choiceExpr{
				pos: position{line: 467, col: 10, offset: 10273},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 467, col: 10, offset: 10273},
						name: "Whitespace",
					},
					&ruleRefExpr{
						pos:  position{line: 467, col: 23, offset: 10286},
						name: "EOL",
					},
					&ruleRefExpr{
						pos:  position{line: 467, col: 29, offset: 10292},
						name: "Comment",
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 469, col: 1, offset: 10302},
			expr: &charClassMatcher{
				pos:        position{line: 469, col: 15, offset: 10316},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 470, col: 1, offset: 10324},
			expr: &litMatcher{
				pos:        position{line: 470, col: 8, offset: 10331},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 471, col: 1, offset: 10336},
			expr: &notExpr{
				pos: position{line: 471, col: 8, offset: 10343},
				expr: &anyMatcher{
					line: 471, col: 9, offset: 10344,
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

func (c *current) onStructDecl1(n, ctx, ms, e interface{}) (interface{}, error) {
	m := make([]ast.Member, 0)
	for _, i := range ms.([]interface{}) {
		m = append(m, i.(ast.Member))
	}
	if e != nil {
		m = append(m, e.(ast.Member))
	}

	s := &ast.Struct{
		Name:    n.(string),
		Members: m,
	}

	if ctx != nil {
		s.Contexts = ctx.([]interface{})[1].([]string)
	}

	return s, nil
}

func (p *parser) callonStructDecl1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructDecl1(stack["n"], stack["ctx"], stack["ms"], stack["e"])
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

func (c *current) onContextRefs1(ns interface{}) (interface{}, error) {
	return ns, nil
}

func (p *parser) callonContextRefs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextRefs1(stack["ns"])
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

func (c *current) onSMPosition1(n interface{}) (interface{}, error) {
	return &ast.Ptr{
		Name: n.(string),
	}, nil
}

func (p *parser) callonSMPosition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSMPosition1(stack["n"])
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
