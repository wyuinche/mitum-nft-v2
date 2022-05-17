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
	SettleAuctionFactType   = hint.Type("mitum-nft-settle-auction-operation-fact")
	SettleAuctionFactHint   = hint.NewHint(SettleAuctionFactType, "v0.0.1")
	SettleAuctionFactHinter = SettleAuctionFact{BaseHinter: hint.NewBaseHinter(SettleAuctionFactHint)}
	SettleAuctionType       = hint.Type("mitum-nft-settle-auction-operation")
	SettleAuctionHint       = hint.NewHint(SettleAuctionType, "v0.0.1")
	SettleAuctionHinter     = SettleAuction{BaseOperation: nft.OperationHinter(SettleAuctionHint)}
)

type SettleAuctionFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	nft    nft.NFTID
	cid    currency.CurrencyID
}

func NewSettleAuctionFact(token []byte, sender base.Address, nft nft.NFTID, cid currency.CurrencyID) SettleAuctionFact {
	fact := SettleAuctionFact{
		BaseHinter: hint.NewBaseHinter(SettleAuctionFactHint),
		token:      token,
		sender:     sender,
		nft:        nft,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact SettleAuctionFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact SettleAuctionFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SettleAuctionFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.nft.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact SettleAuctionFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for SettleAuctionFact")
	}

	if err := isvalid.Check(
		nil, false, fact.h,
		fact.sender, fact.nft, fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact SettleAuctionFact) Token() []byte {
	return fact.token
}

func (fact SettleAuctionFact) Sender() base.Address {
	return fact.sender
}

func (fact SettleAuctionFact) NFT() nft.NFTID {
	return fact.nft
}

func (fact SettleAuctionFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact SettleAuctionFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact SettleAuctionFact) Rebuild() SettleAuctionFact {
	fact.h = fact.GenerateHash()

	return fact
}

type SettleAuction struct {
	currency.BaseOperation
}

func NewSettleAuction(fact SettleAuctionFact, fs []base.FactSign, memo string) (SettleAuction, error) {
	bo, err := currency.NewBaseOperationFromFact(SettleAuctionHint, fact, fs, memo)
	if err != nil {
		return SettleAuction{}, err
	}
	return SettleAuction{BaseOperation: bo}, nil
}
