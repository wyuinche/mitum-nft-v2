package collection

import (
	"context"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   ApproveItem
}

func (ipp *ApproveItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	if err := checkExistsState(currency.StateKeyAccount(ipp.item.Approved()), getStateFunc); err != nil {
		return errors.Errorf("approved not found, %q: %w", ipp.item.Approved(), err)
	}

	nid := ipp.item.NFT()

	st, err := existsState(StateKeyCollectionDesign(nid.Collection()), "key of design", getStateFunc)
	if err != nil {
		return errors.Errorf("collection design not found, %q: %w", nid.Collection(), err)
	}

	design, err := StateCollectionDesignValue(st)
	if err != nil {
		return errors.Errorf("collection design value not found, %q: %w", nid.Collection(), err)
	}

	if !design.Active() {
		return errors.Errorf("deactivated collection, %q", nid.Collection())
	}

	st, err = existsState(extensioncurrency.StateKeyContractAccount(design.Parent()), "contract account", getStateFunc)
	if err != nil {
		return errors.Errorf("parent not found, %q: %w", design.Parent(), err)
	}

	ca, err := extensioncurrency.StateContractAccountValue(st)
	if err != nil {
		return errors.Errorf("contract account value not found, %q: %w", design.Parent(), err)
	}

	if !ca.IsActive() {
		return errors.Errorf("deactivated contract account, %q", design.Parent())
	}

	st, err = existsState(StateKeyNFT(nid), "key of nft", getStateFunc)
	if err != nil {
		return errors.Errorf("nft not found, %q: %w", nid, err)
	}

	nv, err := StateNFTValue(st)
	if err != nil {
		return errors.Errorf("nft value not found, %q: %w", nid, err)
	}

	if !nv.Active() {
		return errors.Errorf("burned nft, %q", nid)
	}

	if ipp.item.Approved().Equal(nv.Approved()) {
		return errors.Errorf("already approved, %q", ipp.item.Approved())
	}

	if !nv.Owner().Equal(ipp.sender) {
		if err := checkExistsState(currency.StateKeyAccount(nv.Owner()), getStateFunc); err != nil {
			return errors.Errorf("nft owner not found, %q: %w", nv.Owner(), err)
		}

		st, err = existsState(StateKeyAgentBox(nv.Owner(), nv.ID().Collection()), "key of agents", getStateFunc)
		if err != nil {
			return errors.Errorf("unauthorized sender, %q: %w", ipp.sender, err)
		}

		box, err := StateAgentBoxValue(st)
		if err != nil {
			return errors.Errorf("agent box value not found, %q: %w", StateKeyAgentBox(nv.Owner(), nv.ID().Collection()), err)
		}

		if !box.Exists(ipp.sender) {
			return errors.Errorf("unauthorized sender, %q", ipp.sender)
		}
	}

	return nil
}

func (ipp *ApproveItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	nid := ipp.item.NFT()

	st, err := existsState(StateKeyNFT(nid), "key of nft", getStateFunc)
	if err != nil {
		return nil, errors.Errorf("nft not found, %q: %w", nid, err)
	}

	nv, err := StateNFTValue(st)
	if err != nil {
		return nil, errors.Errorf("nft value not found, %q: %w", nid, err)
	}

	n := nft.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), ipp.item.Approved(), nv.Creators(), nv.Copyrighters())
	if err := n.IsValid(nil); err != nil {
		return nil, err
	}

	sts := []base.StateMergeValue{NewNFTStateMergeValue(st.Key(), NewNFTStateValue(n))}

	return sts, nil
}

func (ipp *ApproveItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveItem{}

	approveItemProcessorPool.Put(ipp)

	return nil
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() extensioncurrency.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new ApproveProcessor")

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
		if !ok {
			return nil, e(nil, "expected ApproveProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Approve")

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, nil, e(nil, "expected ApproveFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected ApproveItemProcessor, not %T", ipc)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess ApproveItem: %w", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Approve")

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return nil, nil, e(nil, "expected ApproveFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected ApproveItemProcessor, not %T", ip)
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process ApproveItem: %w", err), nil
		}
		sts = append(sts, s...)

		ipc.Close()
	}

	fitems := fact.Items()
	items := make([]CollectionItem, len(fitems))
	for i := range fact.Items() {
		items[i] = fitems[i]
	}

	required, err := CalculateCollectionItemsFee(getStateFunc, items)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)

	return nil
}

func CalculateCollectionItemsFee(getStateFunc base.GetStateFunc, items []CollectionItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for _, item := range items {
		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[item.Currency()]; found {
			rq = k
		}

		policy, err := existsCurrencyPolicy(item.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[item.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[item.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
