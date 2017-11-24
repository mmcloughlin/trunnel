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
