package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

func (it ApproveItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"approved": it.approved,
			"nft":      it.nft,
			"currency": it.currency,
		})
}

type ApproveItemBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Approved string   `bson:"approved"`
	NFT      bson.Raw `bson:"nft"`
	Currency string   `bson:"currency"`
}

func (it *ApproveItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ApproveItem")

	var u ApproveItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, ht, u.Approved, u.NFT, u.Currency)
}
