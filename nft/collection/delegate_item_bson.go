package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it DelegateItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      it.Hint().String(),
			"collection": it.collection,
			"agent":      it.agent,
			"mode":       it.mode,
			"currency":   it.currency,
		},
	)
}

type DelegateItemBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Collection string `bson:"collection"`
	Agent      string `bson:"agent"`
	Mode       string `bson:"mode"`
	Currency   string `bson:"currency"`
}

func (it *DelegateItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of DelegateItem")

	var u DelegateItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unmarshal(enc, ht, u.Collection, u.Agent, u.Mode, u.Currency)
}
