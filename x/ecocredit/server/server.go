package server

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/regen-network/regen-ledger/orm"
	"github.com/regen-network/regen-ledger/x/ecocredit"
)

const (
	TradableBalancePrefix  byte = 0x0
	TradableSupplyPrefix   byte = 0x1
	RetiredBalancePrefix   byte = 0x2
	RetiredSupplyPrefix    byte = 0x3
	IDSeqPrefix            byte = 0x4
	ClassInfoTablePrefix   byte = 0x5
	BatchInfoTablePrefix   byte = 0x6
	MaxDecimalPlacesPrefix byte = 0x7
)

type serverImpl struct {
	storeKey sdk.StoreKey

	// we use a single sequence to avoid having the same string/ID identifying a class and batch denom
	idSeq          orm.Sequence
	classInfoTable orm.NaturalKeyTable
	batchInfoTable orm.NaturalKeyTable
}

// Server is the the ecocredits implementation of ADR 031 Msg Service
type Server interface {
	// Msg Services
	CreateClass(goCtx context.Context, req *ecocredit.MsgCreateClassRequest) (*ecocredit.MsgCreateClassResponse, error)
	CreateBatch(goCtx context.Context, req *ecocredit.MsgCreateBatchRequest) (*ecocredit.MsgCreateBatchResponse, error)
	Send(goCtx context.Context, req *ecocredit.MsgSendRequest) (*ecocredit.MsgSendResponse, error)

	Retire(goCtx context.Context, req *ecocredit.MsgRetireRequest) (*ecocredit.MsgRetireResponse, error)
	SetPrecision(goCtx context.Context, request *ecocredit.MsgSetPrecisionRequest) (*ecocredit.MsgSetPrecisionResponse, error)

	// Query methods
	ClassInfo(ctx context.Context, request *ecocredit.QueryClassInfoRequest) (*ecocredit.QueryClassInfoResponse, error)
	BatchInfo(goCtx context.Context, request *ecocredit.QueryBatchInfoRequest) (*ecocredit.QueryBatchInfoResponse, error)
	Balance(goCtx context.Context, request *ecocredit.QueryBalanceRequest) (*ecocredit.QueryBalanceResponse, error)
	Supply(goCtx context.Context, request *ecocredit.QuerySupplyRequest) (*ecocredit.QuerySupplyResponse, error)
	Precision(goCtx context.Context, request *ecocredit.QueryPrecisionRequest) (*ecocredit.QueryPrecisionResponse, error)
}

// NewServer implements the interface for ADR-031
func NewServer(storeKey sdk.StoreKey) Server {
	s := serverImpl{storeKey: storeKey}

	s.idSeq = orm.NewSequence(storeKey, IDSeqPrefix)

	classInfoTableBuilder := orm.NewNaturalKeyTableBuilder(ClassInfoTablePrefix, storeKey, &ecocredit.ClassInfo{}, orm.Max255DynamicLengthIndexKeyCodec{})
	s.classInfoTable = classInfoTableBuilder.Build()

	batchInfoTableBuilder := orm.NewNaturalKeyTableBuilder(BatchInfoTablePrefix, storeKey, &ecocredit.BatchInfo{}, orm.Max255DynamicLengthIndexKeyCodec{})
	s.batchInfoTable = batchInfoTableBuilder.Build()

	return s
}

func RegisterServices(storeKey sdk.StoreKey, cfg module.Configurator) {
	impl := NewServer(storeKey)
	ecocredit.RegisterMsgServer(cfg.MsgServer(), impl)
	ecocredit.RegisterQueryServer(cfg.QueryServer(), impl)
}
