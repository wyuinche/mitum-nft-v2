package nft

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	NFTInfoType   = hint.Type("mitum-nft-nft-info")
	NFTInfoHint   = hint.NewHint(NFTInfoType, "v0.0.1")
	NFTInfoHinter = NFTInfo{BaseHinter: hint.NewBaseHinter(NFTInfoHint)}
)

type NFTInfo struct {
	hint.BaseHinter
	hash        NFTHash
	uri         NFTUri
	copyrighter Copyrighter
}

func (info NFTInfo) Bytes() []byte {
	return util.ConcatBytesSlice(
		info.hash.Bytes(),
		info.uri.Bytes(),
		info.copyrighter.Bytes(),
	)
}

func (info NFTInfo) IsValid([]byte) error {
	if err := info.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := isvalid.Check(
		nil, false,
		info.BaseHinter,
		info.hash,
		info.uri,
		info.copyrighter); err != nil {
		return err
	}

	return nil
}

var (
	MintNFTFactType   = hint.Type("mitum-nft-mint-nft-operation-fact")
	MintNFTFactHint   = hint.NewHint(MintNFTFactType, "v0.0.1")
	MintNFTFactHinter = MintNFTFact{BaseHinter: hint.NewBaseHinter(MintNFTFactHint)}
	MintNFTType       = hint.Type("mitum-nft-mint-nft-operation")
	MintNFTHint       = hint.NewHint(MintNFTType, "v0.0.1")
	MintNFTHinter     = MintNFT{BaseOperation: operationHinter(MintNFTHint)}
)

type MintNFTFact struct {
	hint.BaseHinter
	h          valuehash.Hash
	token      []byte
	sender     base.Address
	collection Symbol
	nft        NFTInfo
	cid        currency.CurrencyID
}

func NewMintNFTFact(token []byte, sender base.Address, collection Symbol, nft NFTInfo, cid currency.CurrencyID) MintNFTFact {
	fact := MintNFTFact{
		BaseHinter: hint.NewBaseHinter(MintNFTFactHint),
		token:      token,
		sender:     sender,
		collection: collection,
		nft:        nft,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact MintNFTFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact MintNFTFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.h.Bytes())
}

func (fact MintNFTFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.collection.Bytes(),
		fact.nft.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact MintNFTFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for MintNFTFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.collection,
		fact.nft,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact MintNFTFact) Token() []byte {
	return fact.token
}

func (fact MintNFTFact) Sender() base.Address {
	return fact.sender
}

func (fact MintNFTFact) Collection() Symbol {
	return fact.collection
}

func (fact MintNFTFact) NFT() NFTInfo {
	return fact.nft
}

func (fact MintNFTFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.Sender()

	return as, nil
}

func (fact MintNFTFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact MintNFTFact) Rebuild() MintNFTFact {
	fact.h = fact.GenerateHash()

	return fact
}

type MintNFT struct {
	currency.BaseOperation
}

func NewMintNFT(fact MintNFTFact, fs []base.FactSign, memo string) (MintNFT, error) {
	bo, err := currency.NewBaseOperationFromFact(MintNFTHint, fact, fs, memo)
	if err != nil {
		return MintNFT{}, err
	}
	return MintNFT{BaseOperation: bo}, nil
}
