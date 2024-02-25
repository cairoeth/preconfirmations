package preconshare

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

var (
	ErrUnsupportedBundleVersion = errors.New("unsupported bundle version")
	ErrBundleTooDeep            = errors.New("bundle too deep")
	ErrInvalidBundleConstraints = errors.New("invalid bundle constraints")
	ErrInvalidBundlePrivacy     = errors.New("invalid bundle privacy")
)

// MergeInclusionIntervals writes to the topLevel inclusion value of overlap between inner and topLevel
// or return error if there is no overlap
func MergeInclusionIntervals(topLevel, inner *RequestInclusion) error {
	if topLevel.MaxBlock < inner.BlockNumber || inner.MaxBlock < topLevel.BlockNumber {
		return ErrInvalidInclusion
	}

	if topLevel.BlockNumber < inner.BlockNumber {
		topLevel.BlockNumber = inner.BlockNumber
	}
	if topLevel.MaxBlock > inner.MaxBlock {
		topLevel.MaxBlock = inner.MaxBlock
	}
	return nil
}

func validateBundleInner(level int, bundle *SendRequestArgs, currentBlock uint64, signer types.Signer) (hash common.Hash, txs int, unmatched bool, err error) { //nolint:gocognit,gocyclo
	if level > MaxNestingLevel {
		return hash, txs, unmatched, ErrBundleTooDeep
	}
	if bundle.Version != "beta-1" && bundle.Version != "v0.1" {
		return hash, txs, unmatched, ErrUnsupportedBundleVersion
	}

	// validate inclusion
	if bundle.Inclusion.MaxBlock == 0 {
		bundle.Inclusion.MaxBlock = bundle.Inclusion.BlockNumber
	}
	minBlock := uint64(bundle.Inclusion.BlockNumber)
	maxBlock := uint64(bundle.Inclusion.MaxBlock)
	if maxBlock < minBlock {
		return hash, txs, unmatched, ErrInvalidInclusion
	}
	if (maxBlock - minBlock) > MaxBlockRange {
		return hash, txs, unmatched, ErrInvalidInclusion
	}
	if currentBlock >= maxBlock {
		return hash, txs, unmatched, ErrInvalidInclusion
	}
	if minBlock > currentBlock+MaxBlockOffset {
		return hash, txs, unmatched, ErrInvalidInclusion
	}

	// validate body
	if len(bundle.Body) == 0 {
		return hash, txs, unmatched, ErrInvalidBundleBodySize
	}

	bodyHashes := make([]common.Hash, 0, len(bundle.Body))
	for _, el := range bundle.Body {
		if el.Tx != nil {
			var tx types.Transaction
			err := tx.UnmarshalBinary(*el.Tx)
			if err != nil {
				return hash, txs, unmatched, err
			}
			bodyHashes = append(bodyHashes, tx.Hash())
			txs++
		}
	}
	if txs > MaxBodySize {
		return hash, txs, unmatched, ErrInvalidBundleBodySize
	}

	if len(bodyHashes) == 1 {
		// special case of bundle with a single tx
		hash = bodyHashes[0]
	} else {
		hasher := sha3.NewLegacyKeccak256()
		for _, h := range bodyHashes {
			hasher.Write(h[:])
		}
		hash = common.BytesToHash(hasher.Sum(nil))
	}

	// validate validity
	if unmatched {
		// refunds should be empty for unmatched bundles
		return hash, txs, unmatched, ErrInvalidBundleConstraints
	}

	// validate privacy
	if unmatched && bundle.Privacy != nil && bundle.Privacy.Hints != HintNone {
		return hash, txs, unmatched, ErrInvalidBundlePrivacy
	}

	if bundle.Privacy != nil {
		if bundle.Privacy.Hints != HintNone {
			bundle.Privacy.Hints.SetHint(HintHash)
		}
	}

	// clean metadata
	// clean fields owned by the node
	bundle.Metadata = &RequestMetadata{}
	bundle.Metadata.BundleHash = hash
	bundle.Metadata.BodyHashes = bodyHashes
	matchingHasher := sha3.NewLegacyKeccak256()
	matchingHasher.Write(hash[:])
	matchingHash := common.BytesToHash(matchingHasher.Sum(nil))
	bundle.Metadata.MatchingHash = matchingHash

	return hash, txs, unmatched, nil
}

func ValidateBundle(bundle *SendRequestArgs, currentBlock uint64, signer types.Signer) (hash common.Hash, unmatched bool, err error) {
	hash, _, unmatched, err = validateBundleInner(0, bundle, currentBlock, signer)
	return hash, unmatched, err
}
