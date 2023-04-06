package collection

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *NFTTransferItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	rc string,
	bn []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTTransferItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	receiver, err := base.DecodeAddress(rc, enc)
	if err != nil {
		return e(err, "")
	}
	it.receiver = receiver

	if hinter, err := enc.Decode(bn); err != nil {
		return e(err, "")
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
	} else {
		it.nft = n
	}

	it.currency = currency.CurrencyID(cid)

	return nil
}
