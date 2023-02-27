package collection

import (
	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s CollectionDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            s.Hint().String(),
			"collectiondesign": s.CollectionDesign,
		},
	)
}

type CollectionDesignStateValueBSONUnmarshaler struct {
	Hint             string   `bson:"_hint"`
	CollectionDesign bson.Raw `bson:"collectiondesign"`
}

func (s *CollectionDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionDesignStateValue")

	var u CollectionDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var cd CollectionDesign
	if err := cd.DecodeBSON(u.CollectionDesign, enc); err != nil {
		return e(err, "")
	}

	s.CollectionDesign = cd

	return nil
}
