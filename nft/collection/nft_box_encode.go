package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (nbx *NFTBox) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	bns []byte,
) error {
	e := util.StringErrorFunc("failed to unmarshal NFTBox")

	nbx.BaseHinter = hint.NewBaseHinter(ht)

	hns, err := enc.DecodeSlice(bns)
	if err != nil {
		return e(err, "")
	}

	nfts := make([]nft.NFTID, len(hns))
	for i, hinter := range hns {
		n, ok := hinter.(nft.NFTID)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected NFTID, not %T", hinter), "")
		}

		nfts[i] = n
	}
	nbx.nfts = nfts

	return nil
}
