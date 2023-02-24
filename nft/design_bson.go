package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   de.Hint().String(),
			"parent":  de.parent,
			"creator": de.creator,
			"symbol":  de.symbol,
			"active":  de.active,
			"policy":  de.policy,
		})
}

type DesignBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Parent  string   `bson:"parent"`
	Creator string   `bson:"creator"`
	Symbol  string   `bson:"symbol"`
	Active  bool     `bson:"active"`
	Policy  bson.Raw `bson:"policy"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Design")

	var u DesignBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return de.unmarshal(enc, ht, u.Parent, u.Creator, u.Symbol, u.Active, u.Policy)
}
