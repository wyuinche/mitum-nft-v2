package nft

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTJSONMarshaler struct {
	hint.BaseHinter
	ID           NFTID        `json:"id"`
	Active       bool         `json:"active"`
	Owner        base.Address `json:"owner"`
	Hash         NFTHash      `json:"hash"`
	URI          URI          `json:"uri"`
	Approved     base.Address `json:"approved"`
	Creators     Signers      `json:"creators"`
	Copyrighters Signers      `json:"copyrighters"`
}

func (n NFT) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTJSONMarshaler{
		BaseHinter:   n.BaseHinter,
		ID:           n.id,
		Active:       n.active,
		Owner:        n.owner,
		Hash:         n.hash,
		URI:          n.uri,
		Approved:     n.approved,
		Creators:     n.creators,
		Copyrighters: n.copyrighters,
	})
}

type NFTJSONUnmarshaler struct {
	Hint         hint.Hint       `json:"_hint"`
	ID           json.RawMessage `json:"id"`
	Active       bool            `json:"active"`
	Owner        string          `json:"owner"`
	Hash         string          `json:"hash"`
	URI          string          `json:"uri"`
	Approved     string          `json:"approved"`
	Creators     json.RawMessage `json:"creators"`
	Copyrighters json.RawMessage `json:"copyrighters"`
}

func (n *NFT) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFT")

	var u NFTJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return n.unmarshal(enc, u.Hint, u.ID, u.Active, u.Owner, u.Hash, u.URI, u.Approved, u.Creators, u.Copyrighters)
}
