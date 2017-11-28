package gen

import "github.com/mmcloughlin/trunnel/tv"

// Fuzzers generates fuzzing functions for the types in the corpus.
func Fuzzers(pkg string, c *tv.Corpus) ([]byte, error) {
	p := &printer{}

	p.markgenerated()
	p.printf("// +build gofuzz\n\n")
	p.pkg(pkg)

	for _, s := range c.Suites {
		if constrained(s.Vectors) {
			continue
		}
		p.printf("func Fuzz%s(data []byte) int {\n", name(s.Type))
		p.printf("_, err := Parse%s(data)\n", name(s.Type))
		p.printf("if err != nil { return 0 }\n")
		p.printf("return 1")
		p.printf("}\n\n")
	}

	return p.imported()
}
