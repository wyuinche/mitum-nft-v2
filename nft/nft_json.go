package nft

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type NFTIDJSONPacker struct {
	jsonenc.HintedHead
	CL Symbol       `json:"collection"`
	IX currency.Big `json:"id"`
}

func (nid NFTID) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTIDJSONPacker{
		HintedHead: jsonenc.NewHintedHead(nid.Hint()),
		CL:         nid.collection,
		IX:         nid.idx,
	})
}

type NFTIDJSONUnpacker struct {
	CL string       `json:"collection"`
	IX currency.Big `json:"id"`
}

func (nid *NFTID) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uid NFTIDJSONUnpacker
	if err := enc.Unmarshal(b, &uid); err != nil {
		return err
	}

	return nid.unpack(enc, uid.CL, uid.IX)
}

type CopyrighterJSONPacker struct {
	jsonenc.HintedHead
	ST bool         `json:"set"`
	AD base.Address `json:"address"`
}

func (cp Copyrighter) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CopyrighterJSONPacker{
		HintedHead: jsonenc.NewHintedHead(cp.Hint()),
		ST:         cp.set,
		AD:         cp.address,
	})
}

type CopyrighterJSONUnpacker struct {
	ST bool                `json:"set"`
	AD base.AddressDecoder `json:"address"`
}

func (cp *Copyrighter) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ucp CopyrighterJSONUnpacker
	if err := enc.Unmarshal(b, &ucp); err != nil {
		return err
	}

	return cp.unpack(enc, ucp.ST, ucp.AD)
}

type NFTJSONPacker struct {
	jsonenc.HintedHead
	ID NFTID        `json:"id"`
	ON base.Address `json:"owner"`
	HS NFTHash      `json:"hash"`
	UR NFTUri       `json:"uri"`
	AP base.Address `json:"approved"`
	CP Copyrighter  `json:"copyrighter"`
}

func (nft NFT) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(NFTJSONPacker{
		HintedHead: jsonenc.NewHintedHead(nft.Hint()),
		ID:         nft.id,
		ON:         nft.owner,
		HS:         nft.hash,
		UR:         nft.uri,
		AP:         nft.approved,
		CP:         nft.copyrighter,
	})
}

type NFTJSONUnpacker struct {
	ID json.RawMessage     `json:"id"`
	ON base.AddressDecoder `json:"owner"`
	HS string              `json:"hash"`
	UR string              `json:"uri"`
	AP base.AddressDecoder `json:"approved"`
	CP json.RawMessage     `json:"copyrighter"`
}

func (nft *NFT) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var unft NFTJSONUnpacker
	if err := enc.Unmarshal(b, &unft); err != nil {
		return err
	}

	return nft.unpack(enc, unft.ID, unft.ON, unft.HS, unft.UR, unft.AP, unft.CP)
}
