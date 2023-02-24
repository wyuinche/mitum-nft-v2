package nft

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (sgn *Signer) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ac string,
	sh uint,
	sg bool,
) error {
	e := util.StringErrorFunc("failed to unmarshal Signer")

	sgn.BaseHinter = hint.NewBaseHinter(ht)
	sgn.share = sh
	sgn.signed = sg

	account, err := base.DecodeAddress(ac, enc)
	if err != nil {
		return e(err, "")
	}
	sgn.account = account

	return nil
}
