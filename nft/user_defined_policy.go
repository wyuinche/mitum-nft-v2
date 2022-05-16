package nft

import (
	"regexp"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	MinLengthSymbol = 3
	MaxLengthSymbol = 10
	ReValidSymbol   = regexp.MustCompile(`^[A-Z]+$`)
)

type Symbol string

func (s Symbol) Bytes() []byte {
	return []byte(s)
}

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) IsValid([]byte) error {
	if l := len(s); l < MinLengthSymbol || l > MaxLengthSymbol {
		return isvalid.InvalidError.Errorf(
			"invalid length of symbol, %d <= %d <= %d", MinLengthSymbol, l, MaxLengthSymbol)
	} else if !ReValidSymbol.Match([]byte(s)) {
		return isvalid.InvalidError.Errorf("wrong symbol, %q", s)
	}

	return nil
}

type PaymentParameter uint

func (pp PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(pp))
}

func (pp PaymentParameter) IsValid([]byte) error {
	if uint(pp) > 100 {
		return isvalid.InvalidError.Errorf(
			"invalid range of symbol, %d <= %d <= %d", 0, pp, 100)
	}

	return nil
}

var UserDefinedPolicyOptionMap = map[string]UserDefinedPolicyOption{
	"broker":     UserDefinedPolicyOption("nft-broker"),
	"collection": UserDefinedPolicyOption("nft-collection"),
}

type UserDefinedPolicyOption string

func (po UserDefinedPolicyOption) Bytes() []byte {
	return []byte(po)
}

func (po UserDefinedPolicyOption) String() string {
	return string(po)
}

func (po UserDefinedPolicyOption) IsValid([]byte) error {
	if !(po == UserDefinedPolicyOptionMap["broker"] || po == UserDefinedPolicyOptionMap["collection"]) {
		return isvalid.InvalidError.Errorf("invalid user defined policy, %s", string(po))
	}

	return nil
}

var (
	UserDefinedPolicyType             = hint.Type("mitum-nft-user-defined-policy")
	UserDefinedPolicyHint             = hint.NewHint(UserDefinedPolicyType, "v0.0.1")
	BrokerUserDefinedPolicyHinter     = BrokerUserDefinedPolicy{BaseHinter: hint.NewBaseHinter(UserDefinedPolicyHint)}
	CollectionUserDefinedPolicyHinter = CollectionUserDefinedPolicy{BaseHinter: hint.NewBaseHinter(UserDefinedPolicyHint)}
)

type BaseUserDefinedPolicy interface {
	isvalid.IsValider
	Bytes() []byte
	Symbol() Symbol
	Option() BaseUserDefinedPolicy
	Rebuild() BaseUserDefinedPolicy
}

type BrokerUserDefinedPolicy struct {
	hint.BaseHinter
	option    UserDefinedPolicyOption
	symbol    Symbol
	brokerage PaymentParameter
	receiver  base.Address
	royalty   bool
}

func NewBrokerUserDefinedPolicy(symbol Symbol, brokerage PaymentParameter, receiver base.Address, royalty bool) BrokerUserDefinedPolicy {
	return BrokerUserDefinedPolicy{
		BaseHinter: hint.NewBaseHinter(UserDefinedPolicyHint),
		option:     UserDefinedPolicyOptionMap["broker"],
		symbol:     symbol,
		brokerage:  brokerage,
		receiver:   receiver,
		royalty:    royalty,
	}
}

func MustNewBrokerUserDefinedPolicy(symbol Symbol, brokerage PaymentParameter, receiver base.Address, royalty bool) BrokerUserDefinedPolicy {
	broker := NewBrokerUserDefinedPolicy(symbol, brokerage, receiver, royalty)

	if err := broker.IsValid(nil); err != nil {
		panic(err)
	}

	return broker
}

func (broker BrokerUserDefinedPolicy) Bytes() []byte {
	if broker.royalty {
		return util.ConcatBytesSlice(
			broker.option.Bytes(),
			broker.symbol.Bytes(),
			broker.brokerage.Bytes(),
			broker.receiver.Bytes(),
			[]byte{1},
		)
	}

	return util.ConcatBytesSlice(
		broker.option.Bytes(),
		broker.symbol.Bytes(),
		broker.brokerage.Bytes(),
		broker.receiver.Bytes(),
		[]byte{0},
	)
}

func (broker BrokerUserDefinedPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		broker.BaseHinter,
		broker.option,
		broker.symbol,
		broker.brokerage,
		broker.receiver); err != nil {
		return err
	}

	return nil
}

func (broker BrokerUserDefinedPolicy) Option() UserDefinedPolicyOption {
	return broker.option
}

func (broker BrokerUserDefinedPolicy) Symbol() Symbol {
	return broker.symbol
}

func (broker BrokerUserDefinedPolicy) Brokerage() PaymentParameter {
	return broker.brokerage
}

func (broker BrokerUserDefinedPolicy) Receiver() base.Address {
	return broker.receiver
}

func (broker BrokerUserDefinedPolicy) Royalty() bool {
	return broker.royalty
}

func (broker BrokerUserDefinedPolicy) Rebuild() BrokerUserDefinedPolicy {
	return broker
}

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
			"invalid length of collection name, %d <= %d <= %d", MinLengthCollectionName, l, MaxLengthCollectionName)
	} else if !ReValidSymbol.Match([]byte(cn)) {
		return isvalid.InvalidError.Errorf("wrong collection name, %q", cn)
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

type CollectionUserDefinedPolicy struct {
	hint.BaseHinter
	option  UserDefinedPolicyOption
	symbol  Symbol
	name    CollectionName
	creator base.Address
	royalty PaymentParameter
	uri     CollectionUri
}

func NewCollectionUserDefinedPolicy(symbol Symbol, name CollectionName, creator base.Address, royalty PaymentParameter, uri CollectionUri) CollectionUserDefinedPolicy {
	return CollectionUserDefinedPolicy{
		BaseHinter: hint.NewBaseHinter(UserDefinedPolicyHint),
		option:     UserDefinedPolicyOptionMap["collection"],
		symbol:     symbol,
		name:       name,
		creator:    creator,
		royalty:    royalty,
		uri:        uri,
	}
}

func MustNewCollectionUserDefinedPolicy(symbol Symbol, name CollectionName, creator base.Address, royalty PaymentParameter, uri CollectionUri) CollectionUserDefinedPolicy {
	collection := NewCollectionUserDefinedPolicy(symbol, name, creator, royalty, uri)

	if err := collection.IsValid(nil); err != nil {
		panic(err)
	}

	return collection
}

func (collection CollectionUserDefinedPolicy) Bytes() []byte {
	return util.ConcatBytesSlice(
		collection.option.Bytes(),
		collection.symbol.Bytes(),
		collection.name.Bytes(),
		collection.creator.Bytes(),
		collection.royalty.Bytes(),
		collection.uri.Bytes(),
	)
}

func (collection CollectionUserDefinedPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		collection.option,
		collection.symbol,
		collection.name,
		collection.creator,
		collection.royalty,
		collection.uri); err != nil {
		return err
	}

	return nil
}

func (collection CollectionUserDefinedPolicy) Option() UserDefinedPolicyOption {
	return collection.option
}

func (collection CollectionUserDefinedPolicy) Symbol() Symbol {
	return collection.symbol
}

func (collection CollectionUserDefinedPolicy) Name() CollectionName {
	return collection.name
}

func (collection CollectionUserDefinedPolicy) Creator() base.Address {
	return collection.creator
}

func (collection CollectionUserDefinedPolicy) Royalty() PaymentParameter {
	return collection.royalty
}

func (collection CollectionUserDefinedPolicy) Uri() CollectionUri {
	return collection.uri
}

func (collection CollectionUserDefinedPolicy) Rebuild() CollectionUserDefinedPolicy {
	return collection
}
