package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTBoxJSONMarshaler struct {
	hint.BaseHinter
	NFTs []nft.NFTID `json:"nfts"`
}

func (nbx NFTBox) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTBoxJSONMarshaler{
		BaseHinter: nbx.BaseHinter,
		NFTs:       nbx.nfts,
	})
}

type NFTBoxJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	NFTs json.RawMessage `json:"nfts"`
}

func (nbx *NFTBox) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTBox")

	var u NFTBoxJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return nbx.unmarshal(enc, u.Hint, u.NFTs)
}
