package nft

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	PostNFTsItemType   = hint.Type("mitum-nft-post-nfts-item")
	PostNFTsItemHint   = hint.NewHint(PostNFTsItemType, "v0.0.1")
	PostNFTsItemHinter = PostNFTsItem{BaseHinter: hint.NewBaseHinter(PostNFTsItemHint)}
)

type PostNFTsItem struct {
	hint.BaseHinter
	id   NFTID
	info PostInfo
	cid  currency.CurrencyID
}

func NewPostNFTsItem(id NFTID, info PostInfo, cid currency.CurrencyID) PostNFTsItem {
	return PostNFTsItem{
		BaseHinter: hint.NewBaseHinter(PostNFTsItemHint),
		id:         id,
		info:       info,
		cid:        cid,
	}
}

func (it PostNFTsItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.id.Bytes(),
		it.info.Bytes(),
		it.cid.Bytes(),
	)
}

func (it PostNFTsItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		it.BaseHinter,
		it.id,
		it.info,
		it.cid); err != nil {
		return err
	}

	return nil
}

func (it PostNFTsItem) NFT() NFTID {
	return it.id
}

func (it PostNFTsItem) Info() PostInfo {
	return it.info
}

func (it PostNFTsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it PostNFTsItem) Rebuild() PostNFTsItem {
	info := it.info.Rebuild()
	it.info = info

	return it
}
