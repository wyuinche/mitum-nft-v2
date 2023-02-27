package collection

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	CollectionDesignStateValueHint = hint.MustNewHint("collection-design-state-value-v0.0.1")
	StateKeyCollectionDesignPrefix = "collection:"
)

type CollectionDesignStateValue struct {
	hint.BaseHinter
	CollectionDesign CollectionDesign
}

func NewCollectionDesignStateValue(design CollectionDesign) CollectionDesignStateValue {
	return CollectionDesignStateValue{
		BaseHinter:       hint.NewBaseHinter(CollectionDesignStateValueHint),
		CollectionDesign: design,
	}
}

func (cs CollectionDesignStateValue) Hint() hint.Hint {
	return cs.BaseHinter.Hint()
}

func (cs CollectionDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionDesignStateValue")

	if err := cs.BaseHinter.IsValid(CollectionDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := cs.CollectionDesign.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (cs CollectionDesignStateValue) HashBytes() []byte {
	return cs.CollectionDesign.Bytes()
}

func StateCollectionDesignValue(st base.State) (CollectionDesign, error) {
	v := st.Value()
	if v == nil {
		return CollectionDesign{}, util.ErrNotFound.Errorf("collection design not found in State")
	}

	d, ok := v.(CollectionDesignStateValue)
	if !ok {
		return CollectionDesign{}, errors.Errorf("invalid collection design value found, %T", v)
	}

	return d.CollectionDesign, nil
}

func IsStateCollectionDesignKey(key string) bool {
	return strings.HasPrefix(key, StateKeyCollectionDesignPrefix)
}

func StateKeyCollectionDesign(id extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyCollectionDesignPrefix, id)
}

type CollectionDesignStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewCollectionDesignStateValueMerger(height base.Height, key string, st base.State) *CollectionDesignStateValueMerger {
	s := &CollectionDesignStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewCollectionDesignStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewCollectionDesignStateValueMerger(height, key, st)
		},
	)
}

var (
	CollectionLastNFTIndexStateValueHint = hint.MustNewHint("collection-last-nft-index-state-value-v0.0.1")
	StateKeyCollectionLastNFTIndexSuffix = ":collectionidx"
)

type CollectionLastNFTIndexStateValue struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	index      uint64
}

func NewCollectionLastNFTIndexStateValue(collection extensioncurrency.ContractID, index uint64) CollectionLastNFTIndexStateValue {
	return CollectionLastNFTIndexStateValue{
		BaseHinter: hint.NewBaseHinter(CollectionLastNFTIndexStateValueHint),
		collection: collection,
		index:      index,
	}
}

func (is CollectionLastNFTIndexStateValue) Hint() hint.Hint {
	return is.BaseHinter.Hint()
}

func (is CollectionLastNFTIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CollectionLastNFTIndexStateValue")

	if err := is.BaseHinter.IsValid(CollectionLastNFTIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := is.collection.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (is CollectionLastNFTIndexStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice(is.collection.Bytes(), util.Uint64ToBytes(is.index))
}

func StateCollectionLastNFTIndexValue(st base.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(CollectionLastNFTIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.index, nil
}

func IsStateCollectionLastNFTIndexKey(key string) bool {
	return strings.HasSuffix(key, StateKeyCollectionLastNFTIndexSuffix)
}

func StateKeyCollectionLastNFTIndex(id extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", id, StateKeyCollectionLastNFTIndexSuffix)
}

type CollectionLastNFTIndexStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewCollectionLastNFTIndexStateValueMerger(height base.Height, key string, st base.State) *CollectionLastNFTIndexStateValueMerger {
	s := &CollectionLastNFTIndexStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewCollectionLastNFTIndexStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewCollectionLastNFTIndexStateValueMerger(height, key, st)
		},
	)
}

func checkExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return base.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return base.NewBaseOperationProcessReasonError("state, %q already exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	var st base.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, base.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = base.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currency.CurrencyID, getStateFunc base.GetStateFunc) (extensioncurrency.CurrencyPolicy, error) {
	var policy extensioncurrency.CurrencyPolicy

	switch st, found, err := getStateFunc(extensioncurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return extensioncurrency.CurrencyPolicy{}, err
	case !found:
		return extensioncurrency.CurrencyPolicy{}, errors.Errorf("currency not found, %v", cid)
	default:
		design, ok := st.Value().(extensioncurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return extensioncurrency.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", st.Value())
		}
		policy = design.CurrencyDesign.Policy()
	}

	return policy, nil
}

func existsCollectionPolicy(id extensioncurrency.ContractID, getStateFunc base.GetStateFunc) (CollectionPolicy, error) {
	var policy CollectionPolicy

	switch st, found, err := getStateFunc(StateKeyCollectionDesign(id)); {
	case err != nil:
		return CollectionPolicy{}, err
	case !found:
		return CollectionPolicy{}, errors.Errorf("collection not found, %v", id)
	default:
		design, ok := st.Value().(CollectionDesignStateValue)
		if !ok {
			return CollectionPolicy{}, errors.Errorf("expected CollectionDesignStateValue, not %T", st.Value())
		}
		p := design.CollectionDesign.Policy()
		policy, ok = p.(CollectionPolicy)
		if !ok {
			return CollectionPolicy{}, errors.Errorf("expected CollectionPolicy, not %T", p)
		}
	}

	return policy, nil
}
