package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type AgentBoxJSONMarshaler struct {
	hint.BaseHinter
	Collection extensioncurrency.ContractID `json:"collection"`
	Agents     []base.Address               `json:"agents"`
}

func (ab AgentBox) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AgentBoxJSONMarshaler{
		BaseHinter: ab.BaseHinter,
		Collection: ab.collection,
		Agents:     ab.agents,
	})
}

type AgentBoxJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Collection string    `json:"collection"`
	Agents     []string  `json:"agents"`
}

func (ab *AgentBox) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of AgentBox")

	var u AgentBoxJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return ab.unmarshal(enc, u.Hint, u.Collection, u.Agents)
}
