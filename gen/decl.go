package gen

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mmcloughlin/trunnel/ast"
	"golang.org/x/tools/imports"
)

// File generates code for the given AST.
func File(pkg string, f *ast.File) ([]byte, error) {
	buf := &bytes.Buffer{}
	g := &generator{
		pkg: pkg,
		w:   buf,
	}
	g.file(f)
	return imports.Process("", buf.Bytes(), nil)
}

type generator struct {
	pkg string
	w   io.Writer
}

func (g *generator) printf(format string, a ...interface{}) {
	fmt.Fprintf(g.w, format, a...)
}

func (g *generator) file(f *ast.File) {
	g.printf("// Code generated by trunnel. DO NOT EDIT.\n\n")
	g.printf("package %s\n\n", g.pkg)

	for _, c := range f.Constants {
		g.constant(c)
	}
	for _, s := range f.Structs {
		g.structure(s)
	}
}

func (g *generator) constant(c *ast.Constant) {
	g.printf("const %s = %d\n\n", name(c.Name), c.Value)
}

func (g *generator) structure(s *ast.Struct) {
	g.structDecl(s)
	g.parse(s)
}

func (g *generator) structDecl(s *ast.Struct) {
	g.printf("type %s struct {\n", name(s.Name))
	for _, m := range s.Members {
		g.structMemberDecl(m)
	}
	g.printf("}\n\n")
}

func (g *generator) structMemberDecl(m ast.Member) {
	switch m := m.(type) {
	case *ast.NulTermString:
		g.printf("\t%s string\n", name(m.Name))
	case *ast.IntegerMember:
		g.printf("\t%s uint%d\n", name(m.Name), m.Type.Size)
	case *ast.StructMember:
		g.printf("\t%s *%s\n", name(m.Name), name(m.Ref.Name))
	case *ast.FixedArrayMember:
		g.printf("\t%s [%s]%s\n", name(m.Name), integer(m.Size), tipe(m.Base))
	case *ast.VarArrayMember:
		g.printf("\t%s []%s\n", name(m.Name), tipe(m.Base))
	case *ast.UnionMember:
		g.structUnionMemberDecl(m)
	case *ast.EOS:
		// ignore
	default:
		panic(unexpected(m))
	}
}

func (g *generator) structUnionMemberDecl(m *ast.UnionMember) {
	for _, c := range m.Cases {
		for _, f := range c.Members {
			switch f := f.(type) {
			case *ast.Fail, *ast.Ignore:
				// nothing
			default:
				g.structMemberDecl(f)
			}
		}
	}
}

// parse generates a parse function for the type.
func (g *generator) parse(s *ast.Struct) {
	receiver := strings.ToLower(s.Name[:1])
	g.printf("func (%s *%s) Parse(data []byte) ([]byte, error) {\n", receiver, name(s.Name))
	for _, m := range s.Members {
		g.parseMember(receiver, m)
	}
	g.printf("return data, nil\n}\n\n")
}

func (g *generator) parseMember(receiver string, m ast.Member) {
	g.printf("{\n")
	switch m := m.(type) {
	case *ast.NulTermString:
		v := receiver + "." + name(m.Name)
		g.printf("i := bytes.IndexByte(data, 0)\n")
		g.printf("if i < 0 { return nil, errors.New(\"could not parse nul-term string\") }\n")
		g.printf("%s, data = string(data[:i]), data[i+1:]\n", v)

	case *ast.IntegerMember:
		v := receiver + "." + name(m.Name)
		n := m.Type.Size / 8
		g.lengthCheck(n)
		if n == 1 {
			g.printf("%s = data[0]\n", v)
		} else {
			g.printf("%s = binary.BigEndian.Uint%d(data)\n", v, m.Type.Size)
		}
		g.printf("data = data[%d:]", n)

	case *ast.StructMember:
		v := receiver + "." + name(m.Name)
		g.printf("var err error\n")
		g.printf("%s = new(%s)\n", v, name(m.Ref.Name))
		g.printf("data, err = %s.Parse(data)\n", v)
		g.printf("if err != nil { return nil, err }\n")

	case *ast.EOS:
		g.printf("if len(data) > 0 { return nil, errors.New(\"trailing data disallowed\") }\n")

	default:
		// XXX panic(unexpected(m))
		g.printf("// %s\n", unexpected(m))
	}
	g.printf("}\n")
}

func (g *generator) lengthCheck(n int) {
	g.printf("if len(data) < %d { return nil, errors.New(\"data too short\") }\n", n)
}

func tipe(t interface{}) string {
	switch t := t.(type) {
	case *ast.IntType:
		return "uint" + strconv.Itoa(t.Size)
	case *ast.CharType:
		return "byte"
	case *ast.StructRef:
		return "*" + t.Name
	default:
		panic(unexpected(t))
	}
}

func integer(i ast.Integer) string {
	switch i := i.(type) {
	case *ast.IntegerConstRef:
		return name(i.Name)
	case *ast.IntegerLiteral:
		return strconv.FormatInt(i.Value, 10)
	default:
		panic(unexpected(i))
	}
}

func unexpected(t interface{}) string {
	return fmt.Sprintf("unexpected type %T", t)
}