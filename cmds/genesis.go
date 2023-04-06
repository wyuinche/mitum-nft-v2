package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	isaacblock "github.com/ProtoconNet/mitum2/isaac/block"
	isaacoperation "github.com/ProtoconNet/mitum2/isaac/operation"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type GenesisBlockGenerator struct {
	local    base.LocalNode
	enc      encoder.Encoder
	db       isaac.Database
	proposal base.ProposalSignFact
	ivp      base.INITVoteproof
	avp      base.ACCEPTVoteproof
	*logging.Logging
	dataroot  string
	networkID base.NetworkID
	facts     []base.Fact
	ops       []base.Operation
}

func NewGenesisBlockGenerator(
	local base.LocalNode,
	networkID base.NetworkID,
	enc encoder.Encoder,
	db isaac.Database,
	dataroot string,
	facts []base.Fact,
) *GenesisBlockGenerator {
	return &GenesisBlockGenerator{
		Logging: logging.NewLogging(func(zctx zerolog.Context) zerolog.Context {
			return zctx.Str("module", "genesis-block-generator")
		}),
		local:     local,
		networkID: networkID,
		enc:       enc,
		db:        db,
		dataroot:  dataroot,
		facts:     facts,
	}
}

func (g *GenesisBlockGenerator) Generate() (base.BlockMap, error) {
	e := util.StringErrorFunc("failed to generate genesis block")

	if err := g.generateOperations(); err != nil {
		return nil, e(err, "")
	}

	if err := g.newProposal(nil); err != nil {
		return nil, e(err, "")
	}

	if err := g.process(); err != nil {
		return nil, e(err, "")
	}

	fsreader, err := isaacblock.NewLocalFSReaderFromHeight(g.dataroot, base.GenesisHeight, g.enc)
	if err != nil {
		return nil, e(err, "")
	}

	switch blockmap, found, err := fsreader.BlockMap(); {
	case err != nil:
		return nil, e(err, "")
	case !found:
		return nil, errors.Errorf("blockmap not found")
	default:
		if err := blockmap.IsValid(g.networkID); err != nil {
			return nil, e(err, "")
		}

		g.Log().Info().Interface("blockmap", blockmap).Msg("genesis block generated")

		if err := g.closeDatabase(); err != nil {
			return nil, e(err, "")
		}

		return blockmap, nil
	}
}

func (g *GenesisBlockGenerator) generateOperations() error {
	g.ops = make([]base.Operation, len(g.facts))

	types := map[string]struct{}{}

	for i, fact := range g.facts {
		var err error

		hinter, ok := fact.(hint.Hinter)
		if !ok {
			return errors.Errorf("fact does not support Hinter")
		}

		switch ht := hinter.Hint(); {
		case ht.IsCompatible(isaacoperation.SuffrageGenesisJoinFactHint):
			if _, found := types[ht.String()]; found {
				return errors.Errorf("multiple join operation found")
			}

			g.ops[i], err = g.joinOperation(fact)
		case ht.IsCompatible(isaacoperation.GenesisNetworkPolicyFactHint):
			if _, found := types[ht.String()]; found {
				return errors.Errorf("multiple network policy operation found")
			}

			g.ops[i], err = g.networkPolicyOperation(fact)
		case ht.IsCompatible(extensioncurrency.GenesisCurrenciesFactHint):
			if _, found := types[ht.String()]; found {
				return errors.Errorf("multiple GenesisCurrencies operation found")
			}
			g.ops[i], err = g.genesisCurrenciesOperation(fact, g.networkID)
		}

		if err != nil {
			return err
		}

		types[hinter.Hint().String()] = struct{}{}
	}

	return nil
}

func (g *GenesisBlockGenerator) joinOperation(i base.Fact) (base.Operation, error) {
	e := util.StringErrorFunc("failed to make join operation")

	basefact, ok := i.(isaacoperation.SuffrageGenesisJoinFact)
	if !ok {
		return nil, e(nil, "expected SuffrageGenesisJoinFact, not %T", i)
	}

	fact := isaacoperation.NewSuffrageGenesisJoinFact(basefact.Nodes(), g.networkID)

	if err := fact.IsValid(g.networkID); err != nil {
		return nil, e(err, "")
	}

	op := isaacoperation.NewSuffrageGenesisJoin(fact)
	if err := op.Sign(g.local.Privatekey(), g.networkID); err != nil {
		return nil, e(err, "")
	}

	g.Log().Debug().Interface("operation", op).Msg("genesis join operation created")

	return op, nil
}

func (g *GenesisBlockGenerator) networkPolicyOperation(i base.Fact) (base.Operation, error) {
	e := util.StringErrorFunc("failed to make join operation")

	basefact, ok := i.(isaacoperation.GenesisNetworkPolicyFact)
	if !ok {
		return nil, e(nil, "expected GenesisNetworkPolicyFact, not %T", i)
	}

	fact := isaacoperation.NewGenesisNetworkPolicyFact(basefact.Policy())

	if err := fact.IsValid(nil); err != nil {
		return nil, e(err, "")
	}

	op := isaacoperation.NewGenesisNetworkPolicy(fact)
	if err := op.Sign(g.local.Privatekey(), g.networkID); err != nil {
		return nil, e(err, "")
	}

	g.Log().Debug().Interface("operation", op).Msg("genesis network policy operation created")

	return op, nil
}

func (g *GenesisBlockGenerator) genesisCurrenciesOperation(i base.Fact, token []byte) (base.Operation, error) {
	e := util.StringErrorFunc("failed to make genesisCurrencies operation")

	basefact, ok := i.(extensioncurrency.GenesisCurrenciesFact)
	if !ok {
		return nil, e(nil, "expected GenesisCurrenciesFact, not %T", i)
	}
	acks, err := currency.NewBaseAccountKeys(basefact.Keys().Keys(), basefact.Keys().Threshold())
	if err != nil {
		return nil, e(err, "")
	}
	fact := extensioncurrency.NewGenesisCurrenciesFact(token, basefact.GenesisNodeKey(), acks, basefact.Currencies())
	if err := fact.IsValid(g.networkID); err != nil {
		return nil, e(err, "")
	}
	op := extensioncurrency.NewGenesisCurrencies(fact)
	if err := op.Sign(g.local.Privatekey(), g.networkID); err != nil {
		return nil, e(err, "")
	}
	g.Log().Debug().Interface("operation", op).Msg("genesis join operation created")

	return op, nil
}

func (g *GenesisBlockGenerator) newProposal(ops []util.Hash) error {
	e := util.StringErrorFunc("failed to make genesis proposal")

	nops := make([]util.Hash, len(ops)+len(g.ops))
	copy(nops[:len(ops)], ops)

	for i := range g.ops {
		nops[i+len(ops)] = g.ops[i].Hash()
	}

	fact := isaac.NewProposalFact(base.GenesisPoint, g.local.Address(), nops)
	sign := isaac.NewProposalSignFact(fact)

	if err := sign.Sign(g.local.Privatekey(), g.networkID); err != nil {
		return e(err, "")
	}

	if err := sign.IsValid(g.networkID); err != nil {
		return e(err, "")
	}

	g.proposal = sign

	g.Log().Debug().Interface("proposal", sign).Msg("proposal created for genesis")

	return nil
}

func (g *GenesisBlockGenerator) initVoetproof() error {
	e := util.StringErrorFunc("failed to make genesis init voteproof")

	fact := isaac.NewINITBallotFact(base.GenesisPoint, nil, g.proposal.Fact().Hash(), nil)
	if err := fact.IsValid(nil); err != nil {
		return e(err, "")
	}

	sf := isaac.NewINITBallotSignFact(fact)
	if err := sf.NodeSign(g.local.Privatekey(), g.networkID, g.local.Address()); err != nil {
		return e(err, "")
	}

	if err := sf.IsValid(g.networkID); err != nil {
		return e(err, "")
	}

	vp := isaac.NewINITVoteproof(fact.Point().Point)
	vp.
		SetMajority(fact).
		SetSignFacts([]base.BallotSignFact{sf}).
		SetThreshold(base.MaxThreshold).
		Finish()

	if err := vp.IsValid(g.networkID); err != nil {
		return e(err, "")
	}

	g.ivp = vp

	g.Log().Debug().Interface("init_voteproof", vp).Msg("init voteproof created for genesis")

	return nil
}

func (g *GenesisBlockGenerator) acceptVoteproof(proposal, newblock util.Hash) error {
	e := util.StringErrorFunc("failed to make genesis accept voteproof")

	fact := isaac.NewACCEPTBallotFact(base.GenesisPoint, proposal, newblock, nil)
	if err := fact.IsValid(nil); err != nil {
		return e(err, "")
	}

	sf := isaac.NewACCEPTBallotSignFact(fact)
	if err := sf.NodeSign(g.local.Privatekey(), g.networkID, g.local.Address()); err != nil {
		return e(err, "")
	}

	if err := sf.IsValid(g.networkID); err != nil {
		return e(err, "")
	}

	vp := isaac.NewACCEPTVoteproof(fact.Point().Point)
	vp.
		SetMajority(fact).
		SetSignFacts([]base.BallotSignFact{sf}).
		SetThreshold(base.MaxThreshold).
		Finish()

	if err := vp.IsValid(g.networkID); err != nil {
		return e(err, "")
	}

	g.avp = vp

	g.Log().Debug().Interface("init_voteproof", vp).Msg("init voteproof created for genesis")

	return nil
}

func (g *GenesisBlockGenerator) process() error {
	e := util.StringErrorFunc("failed to process")

	if err := g.initVoetproof(); err != nil {
		return e(err, "")
	}

	pp, err := g.newProposalProcessor()
	if err != nil {
		return e(err, "")
	}

	_ = pp.SetLogging(g.Logging)

	switch m, err := pp.Process(context.Background(), g.ivp); {
	case err != nil:
		return e(err, "")
	default:
		if err := m.IsValid(g.networkID); err != nil {
			return e(err, "")
		}

		g.Log().Info().Interface("manifest", m).Msg("genesis block generated")

		if err := g.acceptVoteproof(g.proposal.Fact().Hash(), m.Hash()); err != nil {
			return e(err, "")
		}
	}

	if err := pp.Save(context.Background(), g.avp); err != nil {
		return e(err, "")
	}

	return nil
}

func (g *GenesisBlockGenerator) closeDatabase() error {
	e := util.StringErrorFunc("failed to close database")

	if err := g.db.MergeAllPermanent(); err != nil {
		return e(err, "failed to merge temps")
	}

	return nil
}

func (g *GenesisBlockGenerator) newProposalProcessor() (*isaac.DefaultProposalProcessor, error) {
	return isaac.NewDefaultProposalProcessor(
		g.proposal,
		nil,
		launch.NewBlockWriterFunc(g.local, g.networkID, g.dataroot, g.enc, g.db),
		func(key string) (base.State, bool, error) {
			return nil, false, nil
		},
		func(_ context.Context, operationhash util.Hash) (base.Operation, error) {
			for _, op := range g.ops {
				if operationhash.Equal(op.Hash()) {
					return op, nil
				}
			}

			return nil, util.ErrNotFound.Errorf("operation not found")
		},
		func(base.Height, hint.Hint) (base.OperationProcessor, error) {
			return nil, nil
		},
	)
}
