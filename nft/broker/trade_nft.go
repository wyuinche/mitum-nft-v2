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
	TradeNFTFactType   = hint.Type("mitum-nft-trade-nft-operation-fact")
	TradeNFTFactHint   = hint.NewHint(TradeNFTFactType, "v0.0.1")
	TradeNFTFactHinter = TradeNFTFact{BaseHinter: hint.NewBaseHinter(TradeNFTFactHint)}
	TradeNFTType       = hint.Type("mitum-nft-trade-nft-operation")
	TradeNFTHint       = hint.NewHint(TradeNFTType, "v0.0.1")
	TradeNFTHinter     = TradeNFT{BaseOperation: nft.OperationHinter(TradeNFTHint)}
)

type TradeNFTFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	nft    nft.NFTID
	cid    currency.CurrencyID
}

func NewTradeNFTFact(token []byte, sender base.Address, nft nft.NFTID, cid currency.CurrencyID) TradeNFTFact {
	fact := TradeNFTFact{
		BaseHinter: hint.NewBaseHinter(TradeNFTFactHint),
		token:      token,
		sender:     sender,
		nft:        nft,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact TradeNFTFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact TradeNFTFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TradeNFTFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.nft.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact TradeNFTFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for TradeNFTFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.nft,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact TradeNFTFact) Token() []byte {
	return fact.token
}

func (fact TradeNFTFact) Sender() base.Address {
	return fact.sender
}

func (fact TradeNFTFact) NFT() nft.NFTID {
	return fact.nft
}

func (fact TradeNFTFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact TradeNFTFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact TradeNFTFact) Rebuild() TradeNFTFact {
	fact.h = fact.GenerateHash()

	return fact
}

type TradeNFT struct {
	currency.BaseOperation
}

func NewTradeNFT(fact TradeNFTFact, fs []base.FactSign, memo string) (TradeNFT, error) {
	bo, err := currency.NewBaseOperationFromFact(TradeNFTHint, fact, fs, memo)
	if err != nil {
		return TradeNFT{}, err
	}
	return TradeNFT{BaseOperation: bo}, nil
}
