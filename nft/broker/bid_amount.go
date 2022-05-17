package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	BidNFTFactType   = hint.Type("mitum-nft-bid-nft-operation-fact")
	BidNFTFactHint   = hint.NewHint(BidNFTFactType, "v0.0.1")
	BidNFTFactHinter = BidNFTFact{BaseHinter: hint.NewBaseHinter(BidNFTFactHint)}
	BidNFTType       = hint.Type("mitum-nft-bid-nft-operation")
	BidNFTHint       = hint.NewHint(BidNFTType, "v0.0.1")
	BidNFTHinter     = BidNFT{BaseOperation: nft.OperationHinter(BidNFTHint)}
)

type BidNFTFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	broker nft.Symbol
	nft    nft.NFTID
	amount currency.Amount
}

func NewBidNFTFact(token []byte, sender base.Address, broker nft.Symbol, nft nft.NFTID, amount currency.Amount) BidNFTFact {
	fact := BidNFTFact{
		BaseHinter: hint.NewBaseHinter(BidNFTFactHint),
		token:      token,
		sender:     sender,
		broker:     broker,
		nft:        nft,
		amount:     amount,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact BidNFTFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact BidNFTFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact BidNFTFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.broker.Bytes(),
		fact.nft.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact BidNFTFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for BidNFTFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.broker,
		fact.nft,
		fact.amount); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact BidNFTFact) Token() []byte {
	return fact.token
}

func (fact BidNFTFact) Sender() base.Address {
	return fact.sender
}

func (fact BidNFTFact) Broker() nft.Symbol {
	return fact.broker
}

func (fact BidNFTFact) NFT() nft.NFTID {
	return fact.nft
}

func (fact BidNFTFact) Amount() currency.Amount {
	return fact.amount
}

func (fact BidNFTFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact BidNFTFact) Rebuild() BidNFTFact {
	fact.h = fact.GenerateHash()

	return fact
}

type BidNFT struct {
	currency.BaseOperation
}

func NewBidNFT(fact BidNFTFact, fs []base.FactSign, memo string) (BidNFT, error) {
	bo, err := currency.NewBaseOperationFromFact(BidNFTHint, fact, fs, memo)
	if err != nil {
		return BidNFT{}, err
	}
	return BidNFT{BaseOperation: bo}, nil
}
