package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTIDJSONMarshaler struct {
	hint.BaseHinter
	Collection extensioncurrency.ContractID `json:"collection"`
	Index      uint64                       `json:"index"`
}

func (nid NFTID) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTIDJSONMarshaler{
		BaseHinter: nid.BaseHinter,
		Collection: nid.collection,
		Index:      nid.index,
	})
}

type NFTIDJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Collection string    `json:"collection"`
	Index      uint64    `json:"index"`
}

func (nid *NFTID) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTID")

	var u NFTIDJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return nid.unmarshal(enc, u.Hint, u.Collection, u.Index)
}
