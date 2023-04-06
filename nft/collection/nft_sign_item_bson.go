package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it NFTSignItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         it.Hint().String(),
			"qualification": it.qualification,
			"nft":           it.nft,
			"currency":      it.currency,
		},
	)
}

type NFTSignItemBSONUnmarshaler struct {
	Hint          string   `bson:"_hint"`
	Qualification string   `bson:"qualification"`
	NFT           bson.Raw `bson:"nft"`
	Currency      string   `bson:"currency"`
}

func (it *NFTSignItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of NFTSignItem")

	var u NFTSignItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, ht, u.Qualification, u.NFT, u.Currency)
}
