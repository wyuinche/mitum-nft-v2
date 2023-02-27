package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (form *CollectionRegisterForm) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	tg string,
	sb string,
	nm string,
	ry uint,
	uri string,
	bws []string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionRegisterForm")

	form.BaseHinter = hint.NewBaseHinter(ht)
	form.symbol = extensioncurrency.ContractID(sb)
	form.name = CollectionName(nm)
	form.royalty = nft.PaymentParameter(ry)
	form.uri = nft.URI(uri)

	target, err := base.DecodeAddress(tg, enc)
	if err != nil {
		return e(err, "")
	}
	form.target = target

	whites := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e(err, "")
		}
		whites[i] = white

	}
	form.whites = whites

	return nil
}

func (fact *CollectionRegisterFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bf []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionRegisterFact")

	fact.currency = currency.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	if hinter, err := enc.Decode(bf); err != nil {
		return e(err, "")
	} else if form, ok := hinter.(CollectionRegisterForm); !ok {
		return e(util.ErrWrongType.Errorf("expected CollectionRegisterForm, not %T", hinter), "")
	} else {
		fact.form = form
	}

	return nil
}
