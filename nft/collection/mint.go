package collection

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxMintItems = 10

var (
	MintFactHint = hint.MustNewHint("mitum-nft-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-nft-mint-operation-v0.0.1")
)

type MintFact struct {
	base.BaseFact
	sender base.Address
	items  []MintItem
}

func NewMintFact(token []byte, sender base.Address, items []MintItem) MintFact {
	bf := base.NewBaseFact(MintFactHint, token)
	fact := MintFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid(b []byte) error {
	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
	); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return util.ErrInvalid.Errorf("empty items for MintFact")
	} else if l > int(MaxMintItems) {
		return util.ErrInvalid.Errorf("items over allowed, %d > %d", l, MaxMintItems)
	}

	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (fact MintFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MintFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))

	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact MintFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact MintFact) Sender() base.Address {
	return fact.sender
}

func (fact MintFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	for _, item := range fact.items {
		if ads, err := item.Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.sender)

	return as, nil
}

func (fact MintFact) Items() []MintItem {
	return fact.items
}

type Mint struct {
	currency.BaseOperation
}

func NewMint(fact MintFact) (Mint, error) {
	return Mint{BaseOperation: currency.NewBaseOperation(MintHint, fact)}, nil
}

func (op *Mint) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
