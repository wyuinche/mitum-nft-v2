package collection

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BaseTransferItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"from":     it.from,
				"to":       it.to,
				"nfts":     it.nfts,
				"currency": it.cid,
			}),
	)
}

type BaseTransferItemBSONUnpacker struct {
	FR base.AddressDecoder `bson:"from"`
	TO base.AddressDecoder `bson:"to"`
	NS bson.Raw            `bson:"nfts"`
	CR string              `bson:"currency"`
}

func (it *BaseTransferItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BaseTransferItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.FR, uit.TO, uit.NS, uit.CR)
}
