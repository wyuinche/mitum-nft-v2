package collection

import (
	"bytes"
	"sort"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

var AgentBoxHint = hint.MustNewHint("mitum-nft-agent-box-v0.0.1")

type AgentBox struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	agents     []base.Address
}

func NewAgentBox(collection extensioncurrency.ContractID, agents []base.Address) AgentBox {
	if agents == nil {
		return AgentBox{BaseHinter: hint.NewBaseHinter(AgentBoxHint), collection: collection, agents: []base.Address{}}
	}
	return AgentBox{BaseHinter: hint.NewBaseHinter(AgentBoxHint), collection: collection, agents: agents}
}

func (ab AgentBox) IsValid([]byte) error {
	for i := range ab.agents {
		if err := ab.agents[i].IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (ab AgentBox) Bytes() []byte {
	bas := make([][]byte, len(ab.agents))

	for i, agent := range ab.agents {
		bas[i] = agent.Bytes()
	}

	return util.ConcatBytesSlice(bas...)
}

func (ab AgentBox) Hash() util.Hash {
	return ab.GenerateHash()
}

func (ab AgentBox) GenerateHash() util.Hash {
	return valuehash.NewSHA256(ab.Bytes())
}

func (ab AgentBox) IsEmpty() bool {
	return len(ab.agents) < 1
}

func (ab AgentBox) Collection() extensioncurrency.ContractID {
	return ab.collection
}

func (ab AgentBox) Equal(b AgentBox) bool {
	ab.Sort(true)
	b.Sort(true)

	for i := range ab.agents {
		if !ab.agents[i].Equal(b.agents[i]) {
			return false
		}
	}

	return true
}

func (ab *AgentBox) Sort(ascending bool) {
	sort.Slice(ab.agents, func(i, j int) bool {
		if ascending {
			return bytes.Compare(ab.agents[j].Bytes(), ab.agents[i].Bytes()) > 0
		}

		return bytes.Compare(ab.agents[j].Bytes(), ab.agents[i].Bytes()) < 0
	})
}

func (ab AgentBox) Exists(ag base.Address) bool {
	if ab.IsEmpty() {
		return false
	}

	for _, agent := range ab.agents {
		if ag.Equal(agent) {
			return true
		}
	}

	return false
}

func (ab AgentBox) Get(ag base.Address) (base.Address, error) {
	for _, agent := range ab.agents {
		if ag.Equal(agent) {
			return agent, nil
		}
	}

	return currency.Address{}, errors.Errorf("account not in agent box, %q", ag)
}

func (ab *AgentBox) Append(ag base.Address) error {
	if err := ag.IsValid(nil); err != nil {
		return err
	}

	if ab.Exists(ag) {
		return errors.Errorf("account already in agent box, %q", ag)
	}

	if len(ab.agents) >= MaxAgents {
		return errors.Errorf("max agents, %v", ag)
	}

	ab.agents = append(ab.agents, ag)

	return nil
}

func (ab *AgentBox) Remove(ag base.Address) error {
	if !ab.Exists(ag) {
		return errors.Errorf("account not in agent box, %q", ag)
	}

	for i := range ab.agents {
		if ag.String() == ab.agents[i].String() {
			ab.agents[i] = ab.agents[len(ab.agents)-1]
			ab.agents[len(ab.agents)-1] = currency.Address{}
			ab.agents = ab.agents[:len(ab.agents)-1]

			return nil
		}
	}
	return nil
}

func (ab AgentBox) Agents() []base.Address {
	return ab.agents
}
