package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	PostNFTsFactType   = hint.Type("mitum-nft-post-nfts-operation-fact")
	PostNFTsFactHint   = hint.NewHint(PostNFTsFactType, "v0.0.1")
	PostNFTsFactHinter = PostNFTsFact{BaseHinter: hint.NewBaseHinter(PostNFTsFactHint)}
	PostNFTsType       = hint.Type("mitum-nft-post-nfts-operation")
	PostNFTsHint       = hint.NewHint(PostNFTsType, "v0.0.1")
	PostNFTsHinter     = PostNFTs{BaseOperation: nft.OperationHinter(PostNFTsHint)}
)

var MaxPostNFTsItem uint = 10

type PostNFTsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []PostNFTsItem
}

func NewPostNFTsFact(token []byte, sender base.Address, items []PostNFTsItem, cid currency.CurrencyID) PostNFTsFact {
	fact := PostNFTsFact{
		BaseHinter: hint.NewBaseHinter(PostNFTsFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact PostNFTsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact PostNFTsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact PostNFTsFact) Bytes() []byte {
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

func (fact PostNFTsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items")
	} else if n > int(MaxPostNFTsItem) {
		return isvalid.InvalidError.Errorf("items, %d over max, %d", n, MaxPostNFTsItem)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nft := fact.items[i].NFT()

		if err := nft.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundNFT[nft.String()]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found, %s", nft.String())
		}

		foundNFT[nft.String()] = true
	}

	return nil
}

func (fact PostNFTsFact) Token() []byte {
	return fact.token
}

func (fact PostNFTsFact) Sender() base.Address {
	return fact.sender
}

func (fact PostNFTsFact) Items() []PostNFTsItem {
	return fact.items
}

func (fact PostNFTsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact PostNFTsFact) Rebuild() PostNFTsFact {
	items := make([]PostNFTsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type PostNFTs struct {
	currency.BaseOperation
}

func NewPostNFTs(fact PostNFTsFact, fs []base.FactSign, memo string) (PostNFTs, error) {
	bo, err := currency.NewBaseOperationFromFact(PostNFTsHint, fact, fs, memo)
	if err != nil {
		return PostNFTs{}, err
	}

	return PostNFTs{BaseOperation: bo}, nil
}
