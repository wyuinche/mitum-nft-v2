package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *DelegateItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	ag string,
	md string,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal DelegateItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	it.collection = extensioncurrency.ContractID(col)
	it.mode = DelegateMode(md)
	it.currency = currency.CurrencyID(cid)

	agent, err := base.DecodeAddress(ag, enc)
	if err != nil {
		return e(err, "")
	}
	it.agent = agent

	return nil
}
