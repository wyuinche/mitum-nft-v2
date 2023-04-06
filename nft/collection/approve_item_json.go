package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveItemJSONMarshaler struct {
	hint.BaseHinter
	Approved base.Address        `json:"approved"`
	NFT      nft.NFTID           `json:"nft"`
	Currency currency.CurrencyID `json:"currency"`
}

func (it ApproveItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Approved:   it.approved,
		NFT:        it.nft,
		Currency:   it.currency,
	})
}

type ApproveItemJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Approved string          `json:"approved"`
	NFT      json.RawMessage `json:"nft"`
	Currency string          `json:"currency"`
}

func (it *ApproveItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed decode json of ApproveItem")

	var u ApproveItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, u.Hint, u.Approved, u.NFT, u.Currency)
}
