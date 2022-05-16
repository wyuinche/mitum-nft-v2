package nft

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsTransferNFTsItemMultiNFTs = 10

var (
	TransferNFTsItemMultiNFTsType   = hint.Type("mitum-nft-transfer-nfts-multi-nfts")
	TransferNFTsItemMultiNFTsHint   = hint.NewHint(TransferNFTsItemMultiNFTsType, "v0.0.1")
	TransferNFTsItemMultiNFTsHinter = TransferNFTsItemMultiNFTs{
		BaseTransferNFTsItem: BaseTransferNFTsItem{
			BaseHinter: hint.NewBaseHinter(TransferNFTsItemMultiNFTsHint),
		},
	}
)

type TransferNFTsItemMultiNFTs struct {
	BaseTransferNFTsItem
}

func NewTransferNFTsItemMultiNFTs(from base.Address, to base.Address, nfts []NFTID, cid currency.CurrencyID) TransferNFTsItemMultiNFTs {
	return TransferNFTsItemMultiNFTs{
		BaseTransferNFTsItem: NewBaseTransferNFTsItem(TransferNFTsItemMultiNFTsHint, from, to, nfts, cid),
	}
}

func (it TransferNFTsItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseTransferNFTsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n > MaxNFTsTransferNFTsItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", n, MaxNFTsTransferNFTsItemMultiNFTs)
	}

	return nil
}

func (it TransferNFTsItemMultiNFTs) Rebuild() TransferNFTsItem {
	it.BaseTransferNFTsItem = it.BaseTransferNFTsItem.Rebuild().(BaseTransferNFTsItem)

	return it
}
