package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *NFTSignItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	qual string,
	bn []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTSignItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.qualification = Qualification(qual)
	it.currency = currency.CurrencyID(cid)

	if hinter, err := enc.Decode(bn); err != nil {
		return e(err, "")
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
	} else {
		it.nft = n
	}

	return nil
}
