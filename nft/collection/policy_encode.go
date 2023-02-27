package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (p *CollectionPolicy) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	nm string,
	ry uint,
	uri string,
	bws []string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CollectionPoicy")

	p.BaseHinter = hint.NewBaseHinter(ht)
	p.name = CollectionName(nm)
	p.royalty = nft.PaymentParameter(ry)
	p.uri = nft.URI(uri)

	whites := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e(err, "")
		}
		whites[i] = white
	}
	p.whites = whites

	return nil
}
