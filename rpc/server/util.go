package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func Min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func GetTx(rawTxHex string) (*ethtypes.Transaction, error) {
	if len(rawTxHex) < 2 {
		return nil, errors.New("invalid raw transaction")
	}

	rawTxBytes, err := hex.DecodeString(rawTxHex[2:])
	if err != nil {
		return nil, errors.New("invalid raw transaction")
	}

	tx := new(ethtypes.Transaction)
	if err = tx.UnmarshalBinary(rawTxBytes); err != nil {
		return nil, errors.New("error unmarshalling")
	}

	return tx, nil
}

func GetSenderAddressFromTx(tx *ethtypes.Transaction) (common.Address, error) {
	signer := ethtypes.LatestSignerForChainID(tx.ChainId())
	sender, err := ethtypes.Sender(signer, tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

func GetSenderFromTx(tx *ethtypes.Transaction) (string, error) {
	signer := ethtypes.LatestSignerForChainID(tx.ChainId())
	sender, err := ethtypes.Sender(signer, tx)
	if err != nil {
		return "", err
	}
	return sender.Hex(), nil
}

func GetSenderFromRawTx(tx *ethtypes.Transaction) (string, error) {
	from, err := GetSenderFromTx(tx)
	if err != nil {
		return "", errors.New("error getting from")
	}

	return from, nil
}

func GetTxStatus(txHash string) (*types.PrivateTxAPIResponse, error) {
	privTxAPIURL := fmt.Sprintf("%s/tx/%s", ProtectTxAPIHost, txHash)
	resp, err := http.Get(privTxAPIURL)
	if err != nil {
		return nil, errors.Wrap(err, "privTxApi call failed for "+txHash)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "privTxApi body-read failed for "+txHash)
	}

	respObj := new(types.PrivateTxAPIResponse)
	err = json.Unmarshal(bodyBytes, respObj)
	if err != nil {
		msg := fmt.Sprintf("privTxApi jsonUnmarshal failed for %s - status: %d / body: %s", txHash, resp.StatusCode, string(bodyBytes))
		return nil, errors.Wrap(err, msg)
	}

	return respObj, nil
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		if strings.Contains(forwarded, ",") { // return first entry of list of IPs
			return strings.Split(forwarded, ",")[0]
		}
		return forwarded
	}
	return r.RemoteAddr
}

func GetIPHash(r *http.Request) string {
	hash := md5.Sum([]byte(GetIP(r)))
	return hex.EncodeToString(hash[:])
}

// IsMetamask CHROME_ID: nkbihfbeogaeaoehlefnkodbefgpgknn
func IsMetamask(r *http.Request) bool {
	return r.Header.Get("Origin") == "chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn"
}

// IsMetamaskMoz FIREFOX_ID: webextension@metamask.io
func IsMetamaskMoz(r *http.Request) bool {
	return r.Header.Get("Origin") == "moz-extension://57f9aaf6-270a-154f-9a8a-632d0db4128c"
}

func respBytesToJSONRPCResponse(respBytes []byte) (*types.JSONRPCResponse, error) {
	jsonRPCResp := new(types.JSONRPCResponse)

	// Check if returned an error, if so then convert to standard JSON-RPC error
	errorResp := new(types.RelayErrorResponse)
	if err := json.Unmarshal(respBytes, errorResp); err == nil && errorResp.Error != "" {
		// relay returned an error, convert to standard JSON-RPC error now
		jsonRPCResp.Error = &types.JSONRPCError{Message: errorResp.Error}
		return jsonRPCResp, nil
	}

	// Unmarshall JSON-RPC response and check for error inside
	if err := json.Unmarshal(respBytes, jsonRPCResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return jsonRPCResp, nil
}

func BigIntPtrToStr(i *big.Int) string {
	if i == nil {
		return ""
	}
	return i.String()
}

func AddressPtrToStr(a *common.Address) string {
	if a == nil {
		return ""
	}
	return a.Hex()
}

// GetEnv returns the value of the environment variable named by key, or defaultValue if the environment variable doesn't exist
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
