package server

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	DefaultAuctionHint = []string{"hash", "special_logs"}

	ErrEmptyHintQuery                      = errors.New("Hint query must be non-empty if set.")
	ErrEmptyTargetBuilderQuery             = errors.New("Target builder query must be non-empty if set.")
	ErrIncorrectAuctionHints               = errors.New("Incorrect auction hint, must be one of: contract_address, function_selector, logs, calldata, default_logs.")
	ErrIncorrectOriginID                   = errors.New("Incorrect origin id, must be less then 255 char.")
	ErrIncorrectRefundQuery                = errors.New("Incorrect refund query, must be 0xaddress:percentage.")
	ErrIncorrectRefundAddressQuery         = errors.New("Incorrect refund address.")
	ErrIncorrectRefundPercentageQuery      = errors.New("Incorrect refund percentage.")
	ErrIncorrectRefundTotalPercentageQuery = errors.New("Incorrect refund total percentage, must be bellow 100%.")
)

type URLParameters struct {
	pref       types.PrivateTxPreferences
	prefWasSet bool
	originID   string
	fast       bool
}

// ExtractParametersFromURL extracts the auction preference from the url query
// Allowed query params:
//   - hint: mev share hints, can be set multiple times, default: hash, special_logs
//   - originID: origin id, default: ""
//   - builder: target builder, can be set multiple times, default: empty (only send to flashbots builders)
//   - refund: refund in the form of 0xaddress:percentage, default: empty (will be set by default when backrun is produced)
//     example: 0x123:80 - will refund 80% of the backrun profit to 0x123
func ExtractParametersFromURL(url *url.URL, allBuilders []string) (params URLParameters, err error) {
	if strings.HasPrefix(url.Path, "/fast") {
		params.fast = true
	}
	var hint []string
	hintQuery, ok := url.Query()["hint"]
	if ok {
		if len(hintQuery) == 0 {
			return params, ErrEmptyHintQuery
		}
		for _, hint := range hintQuery {
			// valid hints are: "hash", "contract_address", "function_selector", "logs", "calldata"
			if hint != "hash" && hint != "contract_address" && hint != "function_selector" && hint != "logs" && hint != "calldata" && hint != "default_logs" {
				return params, ErrIncorrectAuctionHints
			}
		}
		hint = hintQuery
		params.prefWasSet = true
	} else {
		hint = DefaultAuctionHint
	}
	params.pref.Privacy.Hints = hint

	originIDQuery, ok := url.Query()["originID"]
	if ok {
		if len(originIDQuery) == 0 {
			return params, ErrIncorrectOriginID
		}
		params.originID = originIDQuery[0]
	}

	targetBuildersQuery, ok := url.Query()["builder"]
	if ok {
		if len(targetBuildersQuery) == 0 {
			return params, ErrEmptyTargetBuilderQuery
		}
		params.pref.Privacy.Builders = targetBuildersQuery
	}
	if params.fast {
		// set all builders no matter what's in the url
		params.pref.Privacy.Builders = allBuilders
	}

	refundAddressQuery, ok := url.Query()["refund"]
	if ok {
		if len(refundAddressQuery) == 0 {
			return params, ErrIncorrectRefundQuery
		}

		var (
			addresses = make([]common.Address, len(refundAddressQuery))
			percents  = make([]int, len(refundAddressQuery))
		)

		for i, refundAddress := range refundAddressQuery {
			split := strings.Split(refundAddress, ":")
			if len(split) != 2 {
				return params, ErrIncorrectRefundQuery
			}
			if !common.IsHexAddress(split[0]) {
				return params, ErrIncorrectRefundAddressQuery
			}
			address := common.HexToAddress(split[0])
			percent, err := strconv.Atoi(split[1])
			if err != nil {
				return params, ErrIncorrectRefundPercentageQuery
			}
			if percent <= 0 || percent >= 100 {
				return params, ErrIncorrectRefundPercentageQuery
			}
			addresses[i] = address
			percents[i] = percent
		}

		// should not exceed 100%
		var totalRefund int
		for _, percent := range percents {
			totalRefund += percent
		}
		if totalRefund > 100 {
			return params, ErrIncorrectRefundTotalPercentageQuery
		}

		refundConfig := make([]types.RefundConfig, len(percents))
		for i, percent := range percents {
			refundConfig[i] = types.RefundConfig{
				Address: addresses[i],
				Percent: percent,
			}
		}

		params.pref.Validity.Refund = refundConfig
	}

	return params, nil
}

func AuctionPreferenceErrorToJSONRPCResponse(jsonReq *types.JSONRPCRequest, err error) *types.JSONRPCResponse {
	message := fmt.Sprintf("Invalid auction preference in the rpc endpoint url. %s", err.Error())
	return &types.JSONRPCResponse{
		ID:      jsonReq.ID,
		Error:   &types.JSONRPCError{Code: types.JSONRPCInvalidRequest, Message: message},
		Version: "2.0",
	}
}
