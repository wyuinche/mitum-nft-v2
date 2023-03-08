package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	approved base.Address
	nft      nft.NFTID
	currency currency.CurrencyID
}

func NewApproveItem(approved base.Address, n nft.NFTID, currency currency.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		approved:   approved,
		nft:        n,
		currency:   currency,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.approved,
		it.nft,
		it.currency,
	)
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.approved.Bytes(),
		it.nft.Bytes(),
		it.currency.Bytes(),
	)
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFT() nft.NFTID {
	return it.nft
}

func (it ApproveItem) Currency() currency.CurrencyID {
	return it.currency
}
