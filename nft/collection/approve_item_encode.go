package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *ApproveItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ap string,
	bn []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal ApproveItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = currency.CurrencyID(cid)

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return e(err, "")
	}
	it.approved = approved

	if hinter, err := enc.Decode(bn); err != nil {
		return e(err, "")
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
	} else {
		it.nft = n
	}

	return nil
}
