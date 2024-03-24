package preconshare

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestMergeInclusionIntervals(t *testing.T) {
	cases := []struct {
		name        string
		top         RequestInclusion
		bottom      RequestInclusion
		expectedTop RequestInclusion
		err         error
	}{
		{
			name: "same intervals",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			err: nil,
		},
		{
			name: "overlap, top to the right",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(3),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(4),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
			err: nil,
		},
		{
			name: "overlap, top to the left",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(4),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(3),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
			err: nil,
		},
		{
			name: "overlap, bottom inside top",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(4),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
			err: nil,
		},
		{
			name: "overlap, top inside bottom",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(4),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(2),
				MaxBlock:     hexutil.Uint64(3),
			},
		},
		{
			name: "no overlap, top to the right",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(3),
				MaxBlock:     hexutil.Uint64(4),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			err: ErrInvalidInclusion,
		},
		{
			name: "no overlap, top to the left",
			top: RequestInclusion{
				DesiredBlock: hexutil.Uint64(3),
				MaxBlock:     hexutil.Uint64(4),
			},
			bottom: RequestInclusion{
				DesiredBlock: hexutil.Uint64(1),
				MaxBlock:     hexutil.Uint64(2),
			},
			expectedTop: RequestInclusion{
				DesiredBlock: hexutil.Uint64(3),
				MaxBlock:     hexutil.Uint64(4),
			},
			err: ErrInvalidInclusion,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bottomCopy := c.bottom
			err := MergeInclusionIntervals(&c.top, &bottomCopy)
			if c.err != nil {
				require.ErrorIs(t, err, c.err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, c.expectedTop, c.top)
			require.Equal(t, c.bottom, bottomCopy)
		})
	}
}
