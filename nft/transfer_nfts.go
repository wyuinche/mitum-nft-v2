package nft

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	TransferNFTsFactType   = hint.Type("mitum-nft-tranfer-nfts-operation-fact")
	TransferNFTsFactHint   = hint.NewHint(TransferNFTsFactType, "v0.0.1")
	TransferNFTsFactHinter = TransferNFTsFact{BaseHinter: hint.NewBaseHinter(TransferNFTsFactHint)}
	TransferNFTsType       = hint.Type("mitum-nft-transfer-nfts-operation")
	TransferNFTsHint       = hint.NewHint(TransferNFTsType, "v0.0.1")
	TransferNFTsHinter     = TransferNFTs{BaseOperation: operationHinter(TransferNFTsHint)}
)

var MaxTransferNFTsItems uint = 10

type NFTsItem interface {
	NFTs() []NFTID
}

type TransferNFTsItem interface {
	hint.Hinter
	isvalid.IsValider
	NFTsItem
	Bytes() []byte
	From() base.Address
	To() base.Address
	Addresses() []base.Address
	Rebuild() TransferNFTsItem
}

type TransferNFTsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []TransferNFTsItem
}

func NewTransferNFTsFact(token []byte, sender base.Address, items []TransferNFTsItem) TransferNFTsFact {
	fact := TransferNFTsFact{
		BaseHinter: hint.NewBaseHinter(TransferNFTsFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact TransferNFTsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact TransferNFTsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferNFTsFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact TransferNFTsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items")
	} else if n > int(MaxTransferNFTsItems) {
		return isvalid.InvalidError.Errorf("items, %d over max, %d", n, MaxTransferNFTsItems)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nfts := fact.items[i].NFTs()

		for j := range nfts {
			if err := nfts[j].IsValid(nil); err != nil {
				return err
			}

			nft := nfts[j].String()
			if _, found := foundNFT[nft]; found {
				return isvalid.InvalidError.Errorf("duplicated nft found, %s", nft)
			}

			foundNFT[nft] = true
		}
	}

	return nil
}

func (fact TransferNFTsFact) Token() []byte {
	return fact.token
}

func (fact TransferNFTsFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferNFTsFact) Items() []TransferNFTsItem {
	return fact.items
}

func (fact TransferNFTsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)*2+1)

	for i := range fact.items {
		as = append(as, fact.items[i].Addresses()...)
	}

	as[len(fact.items)*2] = fact.Sender()

	return as, nil
}

func (fact TransferNFTsFact) Rebuild() TransferNFTsFact {
	items := make([]TransferNFTsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type TransferNFTs struct {
	currency.BaseOperation
}

func NewTransferNFTs(fact TransferNFTsFact, fs []base.FactSign, memo string) (TransferNFTs, error) {
	bo, err := currency.NewBaseOperationFromFact(TransferNFTsHint, fact, fs, memo)
	if err != nil {
		return TransferNFTs{}, err
	}

	return TransferNFTs{BaseOperation: bo}, nil
}
