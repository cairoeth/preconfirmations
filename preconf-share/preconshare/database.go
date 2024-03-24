package preconshare

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrBundleNotCancelled = errors.New("bundle not cancelled")

type DBSbundle struct {
	Hash               []byte         `db:"hash"`
	MatchingHash       []byte         `db:"matching_hash"`
	Signer             []byte         `db:"signer"`
	Cancelled          bool           `db:"cancelled"`
	AllowMatching      bool           `db:"allow_matching"`
	Prematched         bool           `db:"prematched"`
	ReceivedAt         time.Time      `db:"received_at"`
	SimSuccess         bool           `db:"sim_success"`
	SimError           sql.NullString `db:"sim_error"`
	SimulatedAt        sql.NullTime   `db:"simulated_at"`
	SimEffGasPrice     sql.NullString `db:"sim_eff_gas_price"`
	SimProfit          sql.NullString `db:"sim_profit"`
	SimRefundableValue sql.NullString `db:"sim_refundable_value"`
	SimGasUsed         sql.NullInt64  `db:"sim_gas_used"`
	// sum of all simulations gas used
	SimAllSimsGasUsed sql.NullInt64 `db:"sim_all_sims_gas_used"`
	// number of simulations that were run for this bundle
	SimTotalSimCount sql.NullInt64  `db:"sim_total_sim_count"`
	Body             []byte         `db:"body"`
	BodySize         int            `db:"body_size"`
	OriginID         sql.NullString `db:"origin_id"`
	InsertedAt       time.Time      `db:"inserted_at"`
}

var insertBundleQuery = `
INSERT INTO sbundle (hash, matching_hash, signer, cancelled, allow_matching, prematched, received_at, 
                     sim_success, sim_error, simulated_at, sim_eff_gas_price, sim_profit, sim_refundable_value, sim_gas_used,
                     sim_all_sims_gas_used, sim_total_sim_count,
                     body, body_size, origin_id)
VALUES (:hash, :matching_hash, :signer, :cancelled, :allow_matching, :prematched, :received_at, 
        :sim_success, :sim_error, :simulated_at, :sim_eff_gas_price, :sim_profit, :sim_refundable_value, :sim_gas_used,
        :sim_all_sims_gas_used, :sim_total_sim_count,
        :body, :body_size, :origin_id)
ON CONFLICT (hash) DO NOTHING
RETURNING hash`

var selectSimDataBundleQueryForUpdate = `
SELECT hash, sim_success, sim_error, simulated_at, sim_eff_gas_price, sim_profit, sim_refundable_value, sim_gas_used, sim_all_sims_gas_used, sim_total_sim_count
FROM sbundle
WHERE hash = $1
FOR UPDATE`

var updateBundleSimQuery = `
UPDATE sbundle
SET sim_success = :sim_success, sim_error = :sim_error, simulated_at = :simulated_at, 
    sim_eff_gas_price = :sim_eff_gas_price, sim_profit = :sim_profit, sim_refundable_value = :sim_refundable_value, 
    sim_gas_used = :sim_gas_used, sim_all_sims_gas_used = :sim_all_sims_gas_used, sim_total_sim_count = :sim_total_sim_count, body = :body
WHERE hash = :hash`

var getBundleQuery = `
SELECT matching_hash, body
FROM sbundle
WHERE matching_hash = $1 AND allow_matching = true AND cancelled = false limit 1`

type DBSbundleBody struct {
	Hash        []byte `db:"hash"`
	ElementHash []byte `db:"element_hash"`
	Idx         int    `db:"idx"`
	Type        int    `db:"type"`
}

var insertBundleBodyQuery = `
INSERT INTO sbundle_body (hash, element_hash, idx, type)
VALUES (:hash, :element_hash, :idx, :type)
ON CONFLICT (hash, idx) DO NOTHING`

type DBSpreconf struct {
	Hash       []byte    `db:"hash"`
	Block      int64     `db:"block"`
	Signature  []byte    `db:"signature"`
	Time       int64     `db:"time"`
	InsertedAt time.Time `db:"inserted_at"`
}

var insertPreconfQuery = `
INSERT INTO spreconf (hash, block, signature)
VALUES (:hash, :block, :signature)
ON CONFLICT (signature) DO NOTHING`

var getPreconfQuery = `
SELECT block, signature, time
FROM spreconf
WHERE hash = $1
ORDER BY block
LIMIT 1`

var updatePreconfQuery = `UPDATE spreconf SET time = $1 WHERE hash = $2`

var filterPreconfQuery = `DELETE FROM spreconf WHERE hash = $1 AND signature != $2`

var ErrBundleNotFound = errors.New("bundle not found")

type DBSbundleHistoricalHint struct {
	ID         int64           `db:"id"`
	Block      int64           `db:"block"`
	Hint       json.RawMessage `db:"hint"`
	InsertedAt time.Time       `db:"inserted_at"`
}

var insertBundleHistoricalHintQuery = `
INSERT INTO sbundle_hint_history (block, hint)
VALUES (:block, :hint)
RETURNING id`

type DBBackend struct {
	db *sqlx.DB

	insertBundle    *sqlx.NamedStmt
	getBundle       *sqlx.Stmt
	insertPreconf   *sqlx.NamedStmt
	getPreconf      *sqlx.Stmt
	filterPreconf   *sqlx.Stmt
	updatePreconf   *sqlx.Stmt
	insertHint      *sqlx.NamedStmt
	updateBundleSim *sqlx.NamedStmt
}

func NewDBBackend(postgresDSN string) (*DBBackend, error) {
	db, err := sqlx.Connect("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(20)

	insertBundle, err := db.PrepareNamed(insertBundleQuery)
	if err != nil {
		return nil, err
	}
	getBundle, err := db.Preparex(getBundleQuery)
	if err != nil {
		return nil, err
	}
	insertPreconf, err := db.PrepareNamed(insertPreconfQuery)
	if err != nil {
		return nil, err
	}
	getPreconf, err := db.Preparex(getPreconfQuery)
	if err != nil {
		return nil, err
	}
	filterPreconf, err := db.Preparex(filterPreconfQuery)
	if err != nil {
		return nil, err
	}
	updatePreconf, err := db.Preparex(updatePreconfQuery)
	if err != nil {
		return nil, err
	}
	insertHint, err := db.PrepareNamed(insertBundleHistoricalHintQuery)
	if err != nil {
		return nil, err
	}

	updateBundleSim, err := db.PrepareNamed(updateBundleSimQuery)
	if err != nil {
		return nil, err
	}

	return &DBBackend{
		db:              db,
		insertBundle:    insertBundle,
		getBundle:       getBundle,
		insertPreconf:   insertPreconf,
		getPreconf:      getPreconf,
		filterPreconf:   filterPreconf,
		updatePreconf:   updatePreconf,
		insertHint:      insertHint,
		updateBundleSim: updateBundleSim,
	}, nil
}

func (b *DBBackend) GetBundleByMatchingHash(ctx context.Context, hash common.Hash) (*SendRequestArgs, error) {
	var dbSbundle DBSbundle
	err := b.getBundle.GetContext(ctx, &dbSbundle, hash.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBundleNotFound
	} else if err != nil {
		return nil, err
	}

	var bundle SendRequestArgs
	err = json.Unmarshal(dbSbundle.Body, &bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (b *DBBackend) GetPreconfByMatchingHash(ctx context.Context, hash common.Hash) (*int64, *hexutil.Bytes, *int64, error) {
	var dbSpreconf DBSpreconf

	// First, we select the best preconf.
	err := b.getPreconf.GetContext(ctx, &dbSpreconf, hash.Bytes())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, ErrBundleNotFound
	} else if err != nil {
		return nil, nil, nil, err
	}

	var signature hexutil.Bytes
	err = json.Unmarshal(dbSpreconf.Signature, &signature)
	if err != nil {
		return nil, nil, nil, err
	}

	// Then, we remove the other preconfs.
	err = b.filterPreconf.GetContext(ctx, &dbSpreconf, hash.Bytes(), dbSpreconf.Signature)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, err
	}

	return &dbSpreconf.Block, &signature, &dbSpreconf.Time, nil
}

func (b *DBBackend) UpdatePreconfTimeBySignature(ctx context.Context, time int64, hash common.Hash) error {
	var dbSpreconf DBSpreconf

	// First, we select the best preconf.
	err := b.updatePreconf.GetContext(ctx, &dbSpreconf, time, hash.Bytes())
	if err != nil {
		return err
	}

	return err
}

// InsertBundleForStats inserts a bundle into the database.
// When called for the second time for the known bundle, it will return known = true and update bundle simulation
// results with the last inserted simulation results.
func (b *DBBackend) InsertBundleForStats(ctx context.Context, bundle *SendRequestArgs) (known bool, err error) { //nolint:gocognit
	var dbBundle DBSbundle
	if bundle.Metadata == nil {
		return known, ErrNilBundleMetadata
	}
	dbBundle.Hash = bundle.Metadata.RequestHash.Bytes()
	dbBundle.MatchingHash = bundle.Metadata.MatchingHash.Bytes()
	dbBundle.Signer = bundle.Metadata.Signer.Bytes()
	dbBundle.AllowMatching = bundle.Privacy != nil && bundle.Privacy.Hints.HasHint(HintHash)
	dbBundle.Prematched = bundle.Metadata.Prematched
	dbBundle.Cancelled = false
	dbBundle.ReceivedAt = time.UnixMicro(int64(bundle.Metadata.ReceivedAt))
	// dbBundle.SimSuccess = result.Success
	// dbBundle.SimError = sql.NullString{String: result.Error, Valid: result.Error != ""}
	// dbBundle.SimulatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	// dbBundle.SimEffGasPrice = sql.NullString{String: dbIntToEth(&result.MevGasPrice), Valid: result.Success}
	// dbBundle.SimProfit = sql.NullString{String: dbIntToEth(&result.Profit), Valid: result.Success}
	// dbBundle.SimRefundableValue = sql.NullString{String: dbIntToEth(&result.RefundableValue), Valid: result.Success}
	// dbBundle.SimGasUsed = sql.NullInt64{Int64: int64(result.GasUsed), Valid: true}
	// dbBundle.SimAllSimsGasUsed = sql.NullInt64{Int64: int64(result.GasUsed), Valid: true}
	// dbBundle.SimTotalSimCount = sql.NullInt64{Int64: 1, Valid: true}
	dbBundle.Body, err = json.Marshal(bundle)
	if err != nil {
		return known, err
	}
	dbBundle.BodySize = len(bundle.Body)
	dbBundle.OriginID = sql.NullString{String: bundle.Metadata.OriginID, Valid: bundle.Metadata.OriginID != ""}

	dbTx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return known, err
	}
	// get hash from db
	var hash []byte
	err = dbTx.NamedStmtContext(ctx, b.insertBundle).GetContext(ctx, &hash, dbBundle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// bundle is known so we update it with fresh simulation results
			known = true
			// 1. get bundle from db
			var storedBundle DBSbundle
			err = dbTx.GetContext(ctx, &storedBundle, selectSimDataBundleQueryForUpdate, dbBundle.Hash)
			if err != nil {
				_ = dbTx.Rollback()
				return known, err
			}

			// var shouldUpdateSim bool
			// if storedBundle.SimSuccess {
			// 	shouldUpdateSim = result.Success
			// } else {
			// 	shouldUpdateSim = true
			// }

			// if shouldUpdateSim {
			// 	storedBundle.SimSuccess = result.Success
			// 	storedBundle.SimError = sql.NullString{String: result.Error, Valid: result.Error != ""}
			// 	storedBundle.SimulatedAt = sql.NullTime{Time: time.Now(), Valid: true}
			// 	storedBundle.SimEffGasPrice = sql.NullString{String: dbIntToEth(&result.MevGasPrice), Valid: result.Success}
			// 	storedBundle.SimProfit = sql.NullString{String: dbIntToEth(&result.Profit), Valid: result.Success}
			// 	storedBundle.SimRefundableValue = sql.NullString{String: dbIntToEth(&result.RefundableValue), Valid: result.Success}
			// 	storedBundle.SimGasUsed = sql.NullInt64{Int64: int64(result.GasUsed), Valid: true}
			// }

			// if storedBundle.SimTotalSimCount.Valid {
			// 	storedBundle.SimAllSimsGasUsed = sql.NullInt64{Int64: storedBundle.SimAllSimsGasUsed.Int64 + int64(result.GasUsed), Valid: true}
			// } else {
			// 	storedBundle.SimAllSimsGasUsed = sql.NullInt64{Int64: int64(result.GasUsed), Valid: true}
			// }
			// if storedBundle.SimTotalSimCount.Valid {
			// 	storedBundle.SimTotalSimCount = sql.NullInt64{Int64: storedBundle.SimTotalSimCount.Int64 + 1, Valid: true}
			// } else {
			// 	storedBundle.SimTotalSimCount = sql.NullInt64{Int64: 1, Valid: true}
			// }
			// 2. update bundle
			// NOTE: we update bundle body as well to make sure we have the latest bundle body in the db.
			// since we are processing bundle every block (and thus updating in database every block) we'll
			// have bundle body with the biggest maxBlock in database, which is the desired behavior.
			// There are cornercases when system crashes, so we do not record latest bundle inclusion in db
			storedBundle.Body = dbBundle.Body
			_, err := dbTx.NamedStmtContext(ctx, b.updateBundleSim).ExecContext(ctx, storedBundle)
			if err != nil {
				_ = dbTx.Rollback()
				return known, err
			}

			_ = dbTx.Commit()
			return known, nil
		}
		_ = dbTx.Rollback()
		return known, err
	}

	// insert body
	bodyElements := make([]DBSbundleBody, len(bundle.Metadata.BodyHashes))
	for i, hash := range bundle.Metadata.BodyHashes {
		var bodyType int
		if i < len(bundle.Body) {
			if bundle.Body[i].Tx != nil {
				bodyType = 1
			}
		}
		bodyElements[i] = DBSbundleBody{Hash: bundle.Metadata.RequestHash.Bytes(), ElementHash: hash.Bytes(), Idx: i, Type: bodyType}
	}

	_, err = dbTx.NamedExecContext(ctx, insertBundleBodyQuery, bodyElements)
	if err != nil {
		_ = dbTx.Rollback()
		return known, err
	}
	return known, dbTx.Commit()
}

func (b *DBBackend) InsertPreconf(ctx context.Context, preconf *ConfirmRequestArgs) error {
	var dbPreconf DBSpreconf

	dbPreconf.Hash = preconf.Preconf.Hash.Bytes()
	dbPreconf.Block = int64(uint64(preconf.Preconf.Block))
	byteSignature, err := json.Marshal(preconf.Signature)
	if err != nil {
		return err
	}
	dbPreconf.Signature = byteSignature
	dbPreconf.Time = 0

	_, err = b.insertPreconf.ExecContext(ctx, dbPreconf)
	if err != nil {
		fmt.Printf("error writing to preconf db")
		return err
	}
	return err
}

func (b *DBBackend) InsertHistoricalHint(ctx context.Context, currentBlock uint64, hint *Hint) error {
	var dbHint DBSbundleHistoricalHint

	dbHint.Block = int64(currentBlock)

	byteHint, err := json.Marshal(hint)
	if err != nil {
		return err
	}
	dbHint.Hint = byteHint

	_, err = b.insertHint.ExecContext(ctx, dbHint)
	return err
}

func (b *DBBackend) Close() error {
	return b.db.Close()
}
