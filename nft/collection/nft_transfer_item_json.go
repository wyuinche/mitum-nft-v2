package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type NFTTransferItemJSONMarshaler struct {
	hint.BaseHinter
	Receiver base.Address        `json:"receiver"`
	NFT      nft.NFTID           `json:"nft"`
	Currency currency.CurrencyID `json:"currency"`
}

func (it NFTTransferItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTTransferItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Receiver:   it.receiver,
		NFT:        it.nft,
		Currency:   it.currency,
	})
}

type NFTTransferItemJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Receiver string          `json:"receiver"`
	NFT      json.RawMessage `json:"nft"`
	Currency string          `json:"currency"`
}

func (it *NFTTransferItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFTTransferItem")

	var u NFTTransferItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, u.Hint, u.Receiver, u.NFT, u.Currency)
}
