package nft

import (
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseTransferNFTsItem struct {
	hint.BaseHinter
	from base.Address
	to   base.Address
	nfts []NFTID
	cid  currency.CurrencyID
}

func NewBaseTransferNFTsItem(ht hint.Hint, from base.Address, to base.Address, nfts []NFTID, cid currency.CurrencyID) BaseTransferNFTsItem {
	return BaseTransferNFTsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		from:       from,
		to:         to,
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseTransferNFTsItem) Bytes() []byte {
	ns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		ns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.from.Bytes(),
		it.to.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (it BaseTransferNFTsItem) IsValid([]byte) error {
	if n := len(it.nfts); n == 0 {
		return errors.Errorf("empty nfts")
	}

	if err := isvalid.Check(nil, false, it.BaseHinter, it.from, it.to, it.cid); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}
		nft := it.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return errors.Errorf("duplicated nft found, %s", nft)
		}
		foundNFT[nft] = true
	}

	return nil
}

func (it BaseTransferNFTsItem) From() base.Address {
	return it.from
}

func (it BaseTransferNFTsItem) To() base.Address {
	return it.to
}

func (it BaseTransferNFTsItem) Addresses() []base.Address {
	as := make([]base.Address, 2)
	as[0] = it.From()
	as[1] = it.To()
	return as
}

func (it BaseTransferNFTsItem) NFTs() []NFTID {
	return it.nfts
}

func (it BaseTransferNFTsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseTransferNFTsItem) Rebuild() TransferNFTsItem {
	nfts := make([]NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
