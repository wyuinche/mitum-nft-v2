package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CollectionPolicyUpdaterFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	col string,
	bpo []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionPolicyUpdaterFact")

	fact.collection = extensioncurrency.ContractID(col)
	fact.currency = currency.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e(err, "")
	}
	fact.sender = sender

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if policy, ok := hinter.(CollectionPolicy); !ok {
		return e(util.ErrWrongType.Errorf("expected CollectionPolicy, not %T", hinter), "")
	} else {
		fact.policy = policy
	}

	return nil
}
