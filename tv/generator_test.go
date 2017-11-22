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
				Constraints: map[string]int64{},
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
				Constraints: map[string]int64{},
			},
		},
		"gradient": []Vector{
			{
				Data:        []byte{0x97, 0x1b, 0xbf, 0x64, 0xb1, 0x96},
				Constraints: map[string]int64{},
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
				Constraints: map[string]int64{},
			},
		},
	}
	assert.Equal(t, expect, v)
}
