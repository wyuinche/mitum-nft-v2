package collection

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
	ApproveNFTsFactType   = hint.Type("mitum-nft-approve-nfts-operation-fact")
	ApproveNFTsFactHint   = hint.NewHint(ApproveNFTsFactType, "v0.0.1")
	ApproveNFTsFactHinter = ApproveNFTsFact{BaseHinter: hint.NewBaseHinter(ApproveNFTsFactHint)}
	ApproveNFTsType       = hint.Type("mitum-nft-approve-nfts-operation")
	ApproveNFTsHint       = hint.NewHint(ApproveNFTsType, "v0.0.1")
	ApproveNFTsHinter     = ApproveNFTs{BaseOperation: nft.OperationHinter(ApproveNFTsHint)}
)

type ApproveNFTsFact struct {
	hint.BaseHinter
	h        valuehash.Hash
	token    []byte
	sender   base.Address
	approved base.Address
	nfts     []nft.NFTID
	cid      currency.CurrencyID
}

func NewApproveNFTsFact(token []byte, sender base.Address, approved base.Address, nfts []nft.NFTID, cid currency.CurrencyID) ApproveNFTsFact {
	fact := ApproveNFTsFact{
		BaseHinter: hint.NewBaseHinter(ApproveNFTsFactHint),
		token:      token,
		sender:     sender,
		approved:   approved,
		nfts:       nfts,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact ApproveNFTsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact ApproveNFTsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveNFTsFact) Bytes() []byte {
	ns := make([][]byte, len(fact.nfts))
	for i := range fact.nfts {
		ns[i] = fact.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.approved.Bytes(),
		fact.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (fact ApproveNFTsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for ApproveNFTsFact")
	} else if n := len(fact.nfts); n < 1 {
		return errors.Errorf("empty nfts")
	}

	if err := isvalid.Check(
		nil, false, fact.h,
		fact.sender, fact.approved, fact.cid); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.nfts {
		if err := fact.nfts[i].IsValid(nil); err != nil {
			return err
		}
		nft := fact.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return errors.Errorf("duplicated nft found, %s", nft)
		}
		foundNFT[nft] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact ApproveNFTsFact) Token() []byte {
	return fact.token
}

func (fact ApproveNFTsFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveNFTsFact) Approved() base.Address {
	return fact.approved
}

func (fact ApproveNFTsFact) NFTs() []nft.NFTID {
	return fact.nfts
}

func (fact ApproveNFTsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.Sender()
	as[1] = fact.Approved()

	return as, nil
}

func (fact ApproveNFTsFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact ApproveNFTsFact) Rebuild() ApproveNFTsFact {
	fact.h = fact.GenerateHash()

	return fact
}

type ApproveNFTs struct {
	currency.BaseOperation
}

func NewApproveNFTs(fact ApproveNFTsFact, fs []base.FactSign, memo string) (ApproveNFTs, error) {
	bo, err := currency.NewBaseOperationFromFact(ApproveNFTsHint, fact, fs, memo)
	if err != nil {
		return ApproveNFTs{}, err
	}
	return ApproveNFTs{BaseOperation: bo}, nil
}
