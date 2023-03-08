package collection

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (ab *AgentBox) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	bags []string,
) error {
	e := util.StringErrorFunc("failed to unmarshal AgentBox")

	ab.BaseHinter = hint.NewBaseHinter(ht)
	ab.collection = extensioncurrency.ContractID(col)

	agents := make([]base.Address, len(bags))
	for i, bag := range bags {
		agent, err := base.DecodeAddress(bag, enc)
		if err != nil {
			return e(err, "")
		}
		agents[i] = agent
	}
	ab.agents = agents

	return nil
}
