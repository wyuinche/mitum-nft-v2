package nft

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (sgns *Signers) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	tt uint,
	bsns []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal Signers")

	sgns.BaseHinter = hint.NewBaseHinter(ht)
	sgns.total = tt

	hinters, err := enc.DecodeSlice(bsns)
	if err != nil {
		return e(err, "")
	}

	signers := make([]Signer, len(hinters))
	for i, hinter := range hinters {
		signer, ok := hinter.(Signer)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected Signer, not %T", hinter), "")
		}

		signers[i] = signer
	}
	sgns.signers = signers

	return nil
}
