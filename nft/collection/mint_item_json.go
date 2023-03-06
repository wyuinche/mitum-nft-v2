package collection

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type MintFormJSONMarshaler struct {
	hint.BaseHinter
	Hash         nft.NFTHash `json:"hash"`
	URI          nft.URI     `json:"uri"`
	Creators     nft.Signers `json:"creators"`
	Copyrighters nft.Signers `json:"copyrighters"`
}

func (form MintForm) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintFormJSONMarshaler{
		BaseHinter:   form.BaseHinter,
		Hash:         form.hash,
		URI:          form.uri,
		Creators:     form.creators,
		Copyrighters: form.copyrighters,
	})
}

type MintFormJSONUnmarshaler struct {
	Hint         hint.Hint       `json:"_hint"`
	Hash         string          `json:"hash"`
	URI          string          `json:"uri"`
	Creators     json.RawMessage `json:"creators"`
	Copyrighters json.RawMessage `json:"copyrighters"`
}

func (form *MintForm) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of MintForm")

	var u MintFormJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return form.unmarshal(enc, u.Hint, u.Hash, u.URI, u.Creators, u.Copyrighters)
}

type MintItemJSONMarshaler struct {
	hint.BaseHinter
	Collection extensioncurrency.ContractID `json:"collection"`
	Form       MintForm                     `json:"form"`
	Currency   currency.CurrencyID          `json:"currency"`
}

func (it MintItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Collection: it.collection,
		Form:       it.form,
		Currency:   it.currency,
	})
}

type MintItemJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Collection string          `json:"collection"`
	Form       json.RawMessage `json:"form"`
	Currency   string          `json:"currency"`
}

func (it *MintItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of MintItem")

	var u MintItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, u.Hint, u.Collection, u.Form, u.Currency)
}
