package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var NFTTransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type NFTTransferItem struct {
	hint.BaseHinter
	receiver base.Address
	nft      nft.NFTID
	currency currency.CurrencyID
}

func NewNFTTransferItem(receiver base.Address, n nft.NFTID, currency currency.CurrencyID) NFTTransferItem {
	return NFTTransferItem{
		BaseHinter: hint.NewBaseHinter(NFTTransferItemHint),
		receiver:   receiver,
		nft:        n,
		currency:   currency,
	}
}

func (it NFTTransferItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.receiver, it.nft, it.currency)
}

func (it NFTTransferItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.receiver.Bytes(),
		it.nft.Bytes(),
		it.currency.Bytes(),
	)
}

func (it NFTTransferItem) Receiver() base.Address {
	return it.receiver
}

func (it NFTTransferItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it NFTTransferItem) NFT() nft.NFTID {
	return it.nft
}

func (it NFTTransferItem) Currency() currency.CurrencyID {
	return it.currency
}
