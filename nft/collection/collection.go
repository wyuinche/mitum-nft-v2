package collection

import (
	"regexp"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MinLengthCollectionName = 3
	MaxLengthCollectionName = 30
	ReValidCollectionName   = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

type CollectionName string

func (cn CollectionName) Bytes() []byte {
	return []byte(cn)
}

func (cn CollectionName) String() string {
	return string(cn)
}

func (cn CollectionName) IsValid([]byte) error {
	if l := len(cn); l < MinLengthCollectionName || l > MaxLengthCollectionName {
		return isvalid.InvalidError.Errorf(
			"invalid length of collection name; %d <= %d <= %d", MinLengthCollectionName, l, MaxLengthCollectionName)
	} else if !nft.ReValidSymbol.Match([]byte(cn)) {
		return isvalid.InvalidError.Errorf("wrong collection name; %q", cn)
	}

	return nil
}

type CollectionUri string

func (cu CollectionUri) Bytes() []byte {
	return []byte(cu)
}

func (cu CollectionUri) String() string {
	return string(cu)
}

func (cu CollectionUri) IsValid([]byte) error {
	if len(cu) == 0 {
		return isvalid.InvalidError.Errorf("empty collection uri")
	}

	return nil
}

var (
	CollectionPolicyType   = hint.Type("mitum-nft-collection-policy")
	CollectionPolicyHint   = hint.NewHint(CollectionPolicyType, "v0.0.1")
	CollectionPolicyHinter = CollectionPolicy{BaseHinter: hint.NewBaseHinter(CollectionPolicyHint)}
)

type CollectionPolicy struct {
	hint.BaseHinter
	symbol  nft.Symbol
	name    CollectionName
	creator base.Address
	royalty nft.PaymentParameter
	uri     CollectionUri
}

func NewCollectionPolicy(symbol nft.Symbol, name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri CollectionUri) CollectionPolicy {
	return CollectionPolicy{
		BaseHinter: hint.NewBaseHinter(CollectionPolicyHint),
		symbol:     symbol,
		name:       name,
		creator:    creator,
		royalty:    royalty,
		uri:        uri,
	}
}

func MustNewCollectionPolicy(symbol nft.Symbol, name CollectionName, creator base.Address, royalty nft.PaymentParameter, uri CollectionUri) CollectionPolicy {
	collection := NewCollectionPolicy(symbol, name, creator, royalty, uri)

	if err := collection.IsValid(nil); err != nil {
		panic(err)
	}

	return collection
}

func (collection CollectionPolicy) Bytes() []byte {
	return util.ConcatBytesSlice(
		collection.symbol.Bytes(),
		collection.name.Bytes(),
		collection.creator.Bytes(),
		collection.royalty.Bytes(),
		collection.uri.Bytes(),
	)
}

func (collection CollectionPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		collection.symbol,
		collection.name,
		collection.creator,
		collection.royalty,
		collection.uri); err != nil {
		return err
	}

	return nil
}

func (collection CollectionPolicy) Symbol() nft.Symbol {
	return collection.symbol
}

func (collection CollectionPolicy) Name() CollectionName {
	return collection.name
}

func (collection CollectionPolicy) Creator() base.Address {
	return collection.creator
}

func (collection CollectionPolicy) Royalty() nft.PaymentParameter {
	return collection.royalty
}

func (collection CollectionPolicy) Uri() CollectionUri {
	return collection.uri
}

func (collection CollectionPolicy) Rebuild() CollectionPolicy {
	return collection
}
