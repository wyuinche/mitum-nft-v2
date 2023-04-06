package collection

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (ab AgentBox) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":      ab.Hint().String(),
		"collection": ab.collection,
		"agents":     ab.agents,
	})
}

type AgentBoxBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Collection string   `bson:"collection"`
	Agents     []string `bson:"agents"`
}

func (ab *AgentBox) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of AgentBox")

	var u AgentBoxBSONUnmarshaler
	if err := bsonenc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return ab.unmarshal(enc, ht, u.Collection, u.Agents)
}
