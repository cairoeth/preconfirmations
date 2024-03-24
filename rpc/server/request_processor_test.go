package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/cairoeth/preconfirmations/rpc/testutils"
	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/stretchr/testify/require"
)

func setupRedis() {
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	RState, err = NewRedisState(redisServer.Addr())
	if err != nil {
		panic(err)
	}
}

func setupMockTxApi() {
	txApiServer := httptest.NewServer(http.HandlerFunc(testutils.MockTxApiHandler))
	ProtectTxApiHost = txApiServer.URL
	testutils.MockTxApiReset()
}

func setServerTimeNowOffset(td time.Duration) {
	Now = func() time.Time {
		return time.Now().Add(td)
	}
}

func TestRequestshouldSendTxToRelay(t *testing.T) {
	setupRedis()
	setupMockTxApi()

	request := RpcRequest{}
	txHash := "0x0Foo"

	// SEND when not seen before
	shouldSend := !request.blockResendingTxToRelay(txHash)
	require.True(t, shouldSend)

	// Fake a previous send
	err := RState.SetTxSentToRelay(txHash)
	require.Nil(t, err, err)

	// Ensure tx status is UNKNOWN
	txStatusApiResponse, err := GetTxStatus(txHash)
	require.Nil(t, err, err)
	require.Equal(t, types.TxStatusUnknown, txStatusApiResponse.Status)

	// NOT SEND when unknown and time since sent < 5 min
	shouldSend = !request.blockResendingTxToRelay(txHash)
	require.False(t, shouldSend)

	// Set tx status to Failed
	testutils.MockTxApiStatusForHash[txHash] = types.TxStatusFailed
	txStatusApiResponse, err = GetTxStatus(txHash)
	require.Nil(t, err, err)
	require.Equal(t, types.TxStatusFailed, txStatusApiResponse.Status)

	// SEND if failed
	shouldSend = !request.blockResendingTxToRelay(txHash)
	require.True(t, shouldSend)

	// Set tx status to pending
	testutils.MockTxApiStatusForHash[txHash] = types.TxStatusPending
	txStatusApiResponse, err = GetTxStatus(txHash)
	require.Nil(t, err, err)
	require.Equal(t, types.TxStatusPending, txStatusApiResponse.Status)

	// NOT SEND if pending
	shouldSend = !request.blockResendingTxToRelay(txHash)
	require.False(t, shouldSend)

	//
	// SEND if UNKNOWN and 5 minutes have passed
	//
	txHash = "0x0DeadBeef"
	setServerTimeNowOffset(time.Minute * -6)
	defer setServerTimeNowOffset(0)

	err = RState.SetTxSentToRelay(txHash)
	require.Nil(t, err, err)

	timeSent, found, err := RState.GetTxSentToRelay(txHash)
	require.Nil(t, err, err)
	require.True(t, found)
	require.True(t, time.Since(timeSent) > time.Minute*4)

	// Ensure tx status is UNKNOWN
	txStatusApiResponse, err = GetTxStatus(txHash)
	require.Nil(t, err, err)
	require.Equal(t, types.TxStatusUnknown, txStatusApiResponse.Status)

	shouldSend = !request.blockResendingTxToRelay(txHash)
	require.True(t, shouldSend)
}
