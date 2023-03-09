package collection

import (
	"github.com/spikeekips/mitum-currency/currency"
	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact CollectionPolicyUpdaterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      fact.Hint().String(),
			"hash":       fact.BaseFact.Hash().String(),
			"token":      fact.BaseFact.Token(),
			"sender":     fact.sender,
			"collection": fact.collection,
			"policy":     fact.policy,
			"currency":   fact.currency,
		})
}

type CollectionPolicyUpdaterFactBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Sender     string   `bson:"sender"`
	Collection string   `bson:"collection"`
	Policy     bson.Raw `bson:"policy"`
	Currency   string   `bson:"currency"`
}

func (fact *CollectionPolicyUpdaterFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionPolicyUpdaterFact")

	var u currency.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf CollectionPolicyUpdaterFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unmarshal(enc, uf.Sender, uf.Collection, uf.Policy, uf.Currency)
}

func (op CollectionPolicyUpdater) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CollectionPolicyUpdater) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CollectionPolicyUpdater")

	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
