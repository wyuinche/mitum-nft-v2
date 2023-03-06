package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (form *MintForm) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	hs string,
	uri string,
	bcrs []byte,
	bcps []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal MintForm")

	form.BaseHinter = hint.NewBaseHinter(ht)
	form.hash = nft.NFTHash(hs)
	form.uri = nft.URI(uri)

	if hinter, err := enc.Decode(bcrs); err != nil {
		return e(err, "")
	} else if creators, ok := hinter.(nft.Signers); !ok {
		return e(util.ErrWrongType.Errorf("expected Signers, not %T", hinter), "")
	} else {
		form.creators = creators
	}

	if hinter, err := enc.Decode(bcps); err != nil {
		return e(err, "")
	} else if copyrighters, ok := hinter.(nft.Signers); !ok {
		return e(util.ErrWrongType.Errorf("expected Signer, not %T", hinter), "")
	} else {
		form.copyrighters = copyrighters
	}

	return nil
}

func (it *MintItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	bf []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal MintItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.collection = extensioncurrency.ContractID(col)

	if hinter, err := enc.Decode(bf); err != nil {
		return e(err, "")
	} else if form, ok := hinter.(MintForm); !ok {
		return e(util.ErrWrongType.Errorf("not MintForm; %T", hinter), "")
	} else {
		it.form = form
	}

	it.currency = currency.CurrencyID(cid)

	return nil
}
