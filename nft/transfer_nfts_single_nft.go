package nft

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	TransferNFTsItemSingleNFTType   = hint.Type("mitum-currency-create-contract-accounts-single-amount")
	TransferNFTsItemSingleNFTHint   = hint.NewHint(TransferNFTsItemSingleNFTType, "v0.0.1")
	TransferNFTsItemSingleNFTHinter = TransferNFTsItemSingleNFT{
		BaseTransferNFTsItem: BaseTransferNFTsItem{
			BaseHinter: hint.NewBaseHinter(TransferNFTsItemSingleNFTHint),
		},
	}
)

type TransferNFTsItemSingleNFT struct {
	BaseTransferNFTsItem
}

func NewTransferNFTsItemSingleNFT(from base.Address, to base.Address, nft NFTID, cid currency.CurrencyID) TransferNFTsItemSingleNFT {
	return TransferNFTsItemSingleNFT{
		BaseTransferNFTsItem: NewBaseTransferNFTsItem(TransferNFTsItemSingleNFTHint, from, to, []NFTID{nft}, cid),
	}
}

func (it TransferNFTsItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseTransferNFTsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it TransferNFTsItemSingleNFT) Rebuild() TransferNFTsItem {
	it.BaseTransferNFTsItem = it.BaseTransferNFTsItem.Rebuild().(BaseTransferNFTsItem)

	return it
}
