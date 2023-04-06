package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p CollectionPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":   p.Hint().String(),
		"name":    p.name,
		"royalty": p.royalty,
		"uri":     p.uri,
		"whites":  p.whites,
	})
}

type PolicyBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Name    string   `bson:"name"`
	Royalty uint     `bson:"royalty"`
	URI     string   `bson:"uri"`
	Whites  []string `bson:"whites"`
}

func (p *CollectionPolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionPolicy")

	var u PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return p.unmarshal(enc, ht, u.Name, u.Royalty, u.URI, u.Whites)
}
