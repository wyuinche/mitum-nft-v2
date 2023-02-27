package collection

import (
	"encoding/json"

	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CollectionDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CollectionDesign CollectionDesign `json:"collectiondesign"`
}

func (s CollectionDesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		CollectionDesignStateValueJSONMarshaler(s),
	)
}

type CollectionDesignStateValueJSONUnmarshaler struct {
	Hint             hint.Hint       `json:"_hint"`
	CollectionDesign json.RawMessage `json:"collectiondesign"`
}

func (s *CollectionDesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CollectionDesignStateValue")

	var u CollectionDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var cd CollectionDesign
	if err := cd.DecodeJSON(u.CollectionDesign, enc); err != nil {
		return e(err, "")
	}
	s.CollectionDesign = cd

	return nil
}
