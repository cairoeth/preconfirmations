package server

import (
	"net/url"
	"testing"

	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestExtractAuctionPreferenceFromUrl(t *testing.T) {
	tests := map[string]struct {
		url  string
		want URLParameters
		err  error
	}{
		"no auction preference": {
			url: "https://rpc.flashbots.net",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: false,
				originID:   "",
			},
			err: nil,
		},
		"only hash hint": {
			url: "https://rpc.flashbots.net?hint=hash",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"hash"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: true,
				originID:   "",
			},
			err: nil,
		},
		"correct hint preference": {
			url: "https://rpc.flashbots.net?hint=contract_address&hint=function_selector&hint=logs&hint=calldata&hint=hash",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"contract_address", "function_selector", "logs", "calldata", "hash"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: true,
				originID:   "",
			},
			err: nil,
		},
		"incorrect hint preference": {
			url:  "https://rpc.flashbots.net?hint=contract_address&hint=function_selector&hint=logs&hint=incorrect",
			want: URLParameters{},
			err:  ErrIncorrectAuctionHints,
		},
		"rpc endpoint set": {
			url: "https://rpc.flashbots.net?rpc=https://mainnet.infura.io/v3/123",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: false,
				originID:   "",
			},
			err: nil,
		},
		"origin id": {
			url: "https://rpc.flashbots.net?originID=123",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: false,
				originID:   "123",
			},
			err: nil,
		},
		"target builder": {
			url: "https://rpc.flashbots.net?builder=builder1&builder=builder2",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy:  types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}, Builders: []string{"builder1", "builder2"}},
					Validity: types.TxValidityPreferences{},
				},
				prefWasSet: false,
				originID:   "",
			},
			err: nil,
		},
		"set refund": {
			url: "https://rpc.flashbots.net?refund=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:17",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy: types.TxPrivacyPreferences{
						Hints: []string{"hash", "special_logs"},
					},
					Validity: types.TxValidityPreferences{
						Refund: []types.RefundConfig{{Address: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), Percent: 17}},
					},
				},
				prefWasSet: false,
				originID:   "",
			},
			err: nil,
		},
		"set refund, two addresses": {
			url: "https://rpc.flashbots.net?&refund=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:70&refund=0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:10",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy: types.TxPrivacyPreferences{
						Hints: []string{"hash", "special_logs"},
					},
					Validity: types.TxValidityPreferences{
						Refund: []types.RefundConfig{
							{Address: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), Percent: 70},
							{Address: common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), Percent: 10},
						},
					},
				},
				prefWasSet: false,
				originID:   "",
			},
			err: nil,
		},
		"set refund, incorrect query": {
			url: "https://rpc.flashbots.net?refund",
			want: URLParameters{
				pref:       types.PrivateTxPreferences{},
				prefWasSet: false,
				originID:   "",
			},
			err: ErrIncorrectRefundQuery,
		},
		"set refund, incorrect 110": {
			url: "https://rpc.flashbots.net?refund=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:110",
			want: URLParameters{
				pref:       types.PrivateTxPreferences{},
				prefWasSet: false,
				originID:   "",
			},
			err: ErrIncorrectRefundPercentageQuery,
		},
		"set refund, incorrect address": {
			url: "https://rpc.flashbots.net?refund=0xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:80",
			want: URLParameters{
				pref:       types.PrivateTxPreferences{},
				prefWasSet: false,
				originID:   "",
			},
			err: ErrIncorrectRefundAddressQuery,
		},
		"set refund, incorrect 50 + 60": {
			url: "https://rpc.flashbots.net?refund=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:50&refund=0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:60",
			want: URLParameters{
				pref:       types.PrivateTxPreferences{},
				prefWasSet: false,
				originID:   "",
			},
			err: ErrIncorrectRefundTotalPercentageQuery,
		},
		"fast": {
			url: "https://rpc.flashbots.net/fast",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy: types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}, Builders: []string{"builder1", "builder2"}},
				},
				prefWasSet: false,
				fast:       true,
				originID:   "",
			},
			err: nil,
		},
		"fast, ignore builders": {
			url: "https://rpc.flashbots.net/fast?builder=builder3&builder=builder4",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy: types.TxPrivacyPreferences{Hints: []string{"hash", "special_logs"}, Builders: []string{"builder1", "builder2"}},
				},
				prefWasSet: false,
				fast:       true,
				originID:   "",
			},
			err: nil,
		},
		"fast, keep hints": {
			url: "https://rpc.flashbots.net/fast?hint=contract_address&hint=function_selector&hint=logs&hint=calldata&hint=hash",
			want: URLParameters{
				pref: types.PrivateTxPreferences{
					Privacy: types.TxPrivacyPreferences{Hints: []string{"contract_address", "function_selector", "logs", "calldata", "hash"}, Builders: []string{"builder1", "builder2"}},
				},
				prefWasSet: true,
				fast:       true,
				originID:   "",
			},
			err: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			url, err := url.Parse(tt.url)
			if err != nil {
				t.Fatal("failed to parse url: ", err)
			}

			got, err := ExtractParametersFromURL(url, []string{"builder1", "builder2"})
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}

			if tt.err == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
