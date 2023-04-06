package nft

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (nid NFTID) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      nid.Hint().String(),
			"collection": nid.collection,
			"index":      nid.index,
		},
	)
}

type NFTIDBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Collection string `bson:"collection"`
	Index      uint64 `bson:"idx"`
}

func (nid *NFTID) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of NFTID")

	var u NFTIDBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	return nid.unmarshal(enc, ht, u.Collection, u.Index)
}
