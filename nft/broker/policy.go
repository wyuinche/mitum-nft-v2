package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	BrokerPolicyType   = hint.Type("mitum-nft-broker-policy")
	BrokerPolicyHint   = hint.NewHint(BrokerPolicyType, "v0.0.1")
	BrokerPolicyHinter = BrokerPolicy{BaseHinter: hint.NewBaseHinter(BrokerPolicyHint)}
)

type BrokerPolicy struct {
	hint.BaseHinter
	symbol    nft.Symbol
	brokerage nft.PaymentParameter
	receiver  base.Address
	royalty   bool
}

func NewBrokerPolicy(symbol nft.Symbol, brokerage nft.PaymentParameter, receiver base.Address, royalty bool) BrokerPolicy {
	return BrokerPolicy{
		BaseHinter: hint.NewBaseHinter(BrokerPolicyHint),
		symbol:     symbol,
		brokerage:  brokerage,
		receiver:   receiver,
		royalty:    royalty,
	}
}

func MustNewBrokerPolicy(symbol nft.Symbol, brokerage nft.PaymentParameter, receiver base.Address, royalty bool) BrokerPolicy {
	broker := NewBrokerPolicy(symbol, brokerage, receiver, royalty)

	if err := broker.IsValid(nil); err != nil {
		panic(err)
	}

	return broker
}

func (broker BrokerPolicy) Bytes() []byte {
	if broker.royalty {
		return util.ConcatBytesSlice(
			broker.symbol.Bytes(),
			broker.brokerage.Bytes(),
			broker.receiver.Bytes(),
			[]byte{1},
		)
	}

	return util.ConcatBytesSlice(
		broker.symbol.Bytes(),
		broker.brokerage.Bytes(),
		broker.receiver.Bytes(),
		[]byte{0},
	)
}

func (broker BrokerPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		broker.BaseHinter,
		broker.symbol,
		broker.brokerage,
		broker.receiver); err != nil {
		return err
	}

	return nil
}

func (broker BrokerPolicy) Symbol() nft.Symbol {
	return broker.symbol
}

func (broker BrokerPolicy) Brokerage() nft.PaymentParameter {
	return broker.brokerage
}

func (broker BrokerPolicy) Receiver() base.Address {
	return broker.receiver
}

func (broker BrokerPolicy) Royalty() bool {
	return broker.royalty
}

func (broker BrokerPolicy) Rebuild() BrokerPolicy {
	return broker
}
