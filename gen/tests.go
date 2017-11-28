package gen

import "github.com/mmcloughlin/trunnel/tv"

// CorpusTests generates a test file based on a corpus of test vectors.
func CorpusTests(pkg string, c *tv.Corpus) ([]byte, error) {
	p := &printer{}
	p.header(pkg)

	for _, s := range c.Suites {
		if constrained(s.Vectors) {
			continue
		}

		p.printf("func Test%sCorpus(t *testing.T) {\n", name(s.Type))

		// cases
		p.printf("cases := []struct{\nData []byte\n}{\n")
		for _, v := range s.Vectors {
			p.printf("{\nData: %#v,\n},\n", v.Data)
		}
		p.printf("}\n")

		// test each one
		p.printf("for _, c := range cases {\n")
		p.printf("_, err := Parse%s(c.Data)\n", name(s.Type))
		p.printf("if err != nil { t.Fatal(err) }\n")
		p.printf("}\n")

		p.printf("}\n\n")
	}

	return p.imported()
}

func constrained(vs []tv.Vector) bool {
	for _, v := range vs {
		if len(v.Constraints) > 0 {
			return true
		}
	}
	return false
}
