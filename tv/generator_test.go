package tv

import (
	"encoding/binary"
	"testing"

	"github.com/mmcloughlin/random"
	"github.com/mmcloughlin/trunnel/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func String(code string) (map[string][]Vector, error) {
	f, err := parse.String(code)
	if err != nil {
		return nil, err
	}
	return Generate(f, WithRandom(random.NewWithSeed(42)))
}

func TestIntType(t *testing.T) {
	v, err := String(`struct color { u8 r; u8 g; u8 b; }`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"color": []Vector{
			{
				Data:        []byte{0x7f, 0x8c, 0x53},
				Constraints: NewConstraints(),
			},
		},
	}
	assert.Equal(t, expect, v)
}

func TestIntConstraint(t *testing.T) {
	f, err := parse.String(`struct date {
		u16 year IN [ 1970..2017 ];
		u8 month IN [ 1, 2..6, 7..12 ];
		u8 day IN [ 1..31 ];
	}`)
	require.NoError(t, err)

	for i := 0; i < 10000; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		b := vs["date"][0].Data
		require.Len(t, b, 4)

		y := binary.BigEndian.Uint16(b)
		m := b[2]
		d := b[3]

		require.True(t, 1970 <= y && y <= 2017)
		require.True(t, 1 <= m && m <= 12)
		require.True(t, 1 <= d && d <= 31)
	}
}

func TestNestedStruct(t *testing.T) {
	v, err := String(`
	struct color { u8 r; u8 g; u8 b; };
	struct gradient {
		struct color from;
		struct color to;
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"color": []Vector{
			{
				Data:        []byte{0x7f, 0x8c, 0x53},
				Constraints: NewConstraints(),
			},
		},
		"gradient": []Vector{
			{
				Data:        []byte{0x97, 0x1b, 0xbf, 0x64, 0xb1, 0x96},
				Constraints: NewConstraints(),
			},
		},
	}
	assert.Equal(t, expect, v)
}

func TestNulTerm(t *testing.T) {
	v, err := String(`
	struct nul_term {
		u16 pre;
		nulterm s;
		u8 post;
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"nul_term": []Vector{
			{
				Data: []byte{
					0x8c, 0x7f, // pre
					'u', 'k', 'p', 't', 't', 0, // s
					0x53, // post
				},
				Constraints: NewConstraints(),
			},
		},
	}
	assert.Equal(t, expect, v)
}

func TestFixedArray(t *testing.T) {
	f, err := parse.String(`
	const NUM_BYTES = 8;
	struct color { u8 r; u8 g; u8 b; }
	struct fixie {
		u8 bytes[NUM_BYTES];
		char letters[NUM_BYTES];
		u16 shortwords[4];
		u32 words[2];
		u64 big_words[2];
		struct color colors[2];
	}`)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		v, err := Generate(f)
		require.NoError(t, err)
		d := v["fixie"][0].Data
		require.Len(t, d, 8+8+2*4+4*2+8*2+3*2)
	}
}

func TestVarArray(t *testing.T) {
	f, err := parse.String(`struct var { u16 n; u32 words[n]; };`)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		v := vs["var"][0]
		n := binary.BigEndian.Uint16(v.Data)
		require.Len(t, v.Data, 2+4*int(n))
	}
}

func TestNestedVar(t *testing.T) {
	f, err := parse.String(`
		struct var  { u16 n; u32 w[n]; };
		struct nest { u16 n; struct var v[n]; };
	`)
	require.NoError(t, err)
	for i := 0; i < 100; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		b := vs["nest"][0].Data

		// should be able to follow the length fields to the end
		n, b := binary.BigEndian.Uint16(b), b[2:]
		for j := 0; j < int(n); j++ {
			l := binary.BigEndian.Uint16(b)
			skip := 2 + 4*int(l)
			require.True(t, len(b) >= skip)
			b = b[skip:]
		}
		require.Len(t, b, 0)
	}
}

func TestLengthDoubleUse(t *testing.T) {
	f, err := parse.String(`struct dbl {
		u8 n;
		u32 a[n];
		u64 b[n];
	};`)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		b := vs["dbl"][0].Data
		n := int(b[0])
		require.Len(t, b, 1+12*n)
	}
}

func TestRemaining(t *testing.T) {
	vs, err := String(`struct rem {
		u32 head;
		u8 tail[];
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"rem": []Vector{
			{
				Data: []byte{
					0x72, 0xe8, 0x9f, 0x5b, 0xb4, 0x4b, 0x9f, 0xbb, 0x97, 0x1b,
				},
				Constraints: NewConstraints(),
			},
		},
	}
	assert.Equal(t, expect, vs)
}

func TestLeftover(t *testing.T) {
	_, err := String(`struct leftover {
		u32 head[2];
		u32 mid[..-8];
		u32 tail[2];
	};`)
	assert.EqualError(t, err, "not implemented")
}

func TestUnionBasic(t *testing.T) {
	vs, err := String(`struct basic {
		u8 tag;
		union u[tag] {
			1: u8 r; u8 g; u8 b;
			2: u16 y; u8 m; u8 d;
		};
		u16 right_after_the_union;
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"basic": []Vector{
			NewVector([]byte{0x01, 0xb1, 0x96, 0x7f, 0x53, 0x8c}),
			NewVector([]byte{0x02, 0x1b, 0x97, 0xbf, 0x64, 0x53, 0x8c}),
		},
	}
	assert.Equal(t, expect, vs)
}

func TestTagDoubleUse(t *testing.T) {
	f, err := parse.String(`struct dbltag {
		u8 tag;
		union u[tag] {
			1: u8 a;
			2: u16 b;
		};
		union w[tag] {
			1: u8 c;
			2: u16 d;
		};
	};`)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		require.Len(t, vs["dbltag"], 2)
		for j, v := range vs["dbltag"] {
			b := v.Data
			tag := j + 1
			require.Equal(t, byte(tag), b[0])
			require.Len(t, b, 1+2*tag)
		}
	}
}

func TestUnionDefault(t *testing.T) {
	vs, err := String(`struct basic {
		u8 tag;
		union u[tag] {
			1:
				u32 a;
			default:
				u8 ukn[];
		};
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"basic": []Vector{
			NewVector([]byte{0x01, 0x53, 0x8c, 0x7f, 0x96}),
			NewVector([]byte{0x64, 0x84, 0x14, 0x5b, 0x9f, 0xe8, 0xbf, 0x64, 0xb1}),
		},
	}
	assert.Equal(t, expect, vs)
}

func TestUnionDefaultRange(t *testing.T) {
	f, err := parse.String(`struct basic {
		u16 tag;
		union u[tag] {
			0..0x2ff: u8 a;
			default: u8 ukn[];
		};
	};`)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		vs, err := Generate(f)
		require.NoError(t, err)
		require.True(t, binary.BigEndian.Uint16(vs["basic"][0].Data) <= uint16(0x2ff))
		require.True(t, binary.BigEndian.Uint16(vs["basic"][1].Data) > uint16(0x2ff))
	}
}

func TestUnionCommands(t *testing.T) {
	vs, err := String(`struct basic {
		u8 tag;
		union u[tag] {
			1: u32 a;
			2..4: ; // empty
			5: ignore;
			default: fail;
		};
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"basic": []Vector{
			NewVector([]byte{0x01, 0x53, 0x8c, 0x7f, 0x96}),
			NewVector([]byte{0x03}),
			NewVector([]byte{0x05}),
			NewVector([]byte{0x05, 0xb1, 0x64, 0xbf}),
		},
	}
	assert.Equal(t, expect, vs)
}

func TestPtr(t *testing.T) {
	vs, err := String(`struct haspos {
		nulterm s;
		@ptr pos1;
		u32 after;
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"haspos": []Vector{
			NewVector([]byte{
				'u', 'k', 'p', 't', 't', 0, // s
				// pos1 occupies no space
				0x53, 0x8c, 0x7f, 0x96, // after
			}),
		},
	}
	assert.Equal(t, expect, vs)
}

func TestEOS(t *testing.T) {
	vs, err := String(`struct haseos {
		u8 r;
		u8 g;
		u8 b;
		eos;
	};`)
	require.NoError(t, err)
	expect := map[string][]Vector{
		"haseos": []Vector{
			NewVector([]byte{0x7f, 0x8c, 0x53}),
		},
	}
	assert.Equal(t, expect, vs)
}
