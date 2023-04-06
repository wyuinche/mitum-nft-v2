package collection

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type ApproveFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address  `json:"sender"`
	Items  []ApproveItem `json:"items"`
}

func (fact ApproveFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type ApproveFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *ApproveFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ApproveFact")

	var uf ApproveFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, uf.Sender, uf.Items)
}

type approveMarshaler struct {
	currency.BaseOperationJSONMarshaler
}

func (op Approve) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(approveMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Approve) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Approve")

	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
