package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (cp CollectionPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(cp.Hint()),
		bson.M{
			"symbol":  cp.symbol,
			"name":    cp.name,
			"creator": cp.creator,
			"royalty": cp.royalty,
			"uri":     cp.uri,
		},
	))
}

type CollectionPolicyBSONUnpacker struct {
	SB string              `bson:"symbol"`
	NM string              `bson:"name"`
	CE base.AddressDecoder `bson:"creator"`
	RY uint                `bson:"royalty"`
	UR string              `bson:"uri"`
}

func (cp CollectionPolicy) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ucp CollectionPolicyBSONUnpacker
	if err := enc.Unmarshal(b, &ucp); err != nil {
		return err
	}

	return cp.unpack(enc, ucp.SB, ucp.NM, ucp.CE, ucp.RY, ucp.UR)
}
