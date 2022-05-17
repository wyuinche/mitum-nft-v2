package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	AddPolicyFactType   = hint.Type("mitum-nft-add-collection-policy-operation-fact")
	AddPolicyFactHint   = hint.NewHint(AddPolicyFactType, "v0.0.1")
	AddPolicyFactHinter = AddPolicyFact{BaseHinter: hint.NewBaseHinter(AddPolicyFactHint)}
	AddPolicyType       = hint.Type("mitum-nft-add-collection-policy-operation")
	AddPolicyHint       = hint.NewHint(AddPolicyType, "v0.0.1")
	AddPolicyHinter     = AddUserDefinedPolicy{BaseOperation: nft.OperationHinter(AddPolicyHint)}
)

type AddPolicyFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	target base.Address
	policy CollectionPolicy
	cid    currency.CurrencyID
}

func NewAddUserDefinedPolicyFact(token []byte, sender base.Address, target base.Address, policy CollectionPolicy, cid currency.CurrencyID) AddPolicyFact {
	fact := AddPolicyFact{
		BaseHinter: hint.NewBaseHinter(AddPolicyFactHint),
		token:      token,
		sender:     sender,
		target:     target,
		policy:     policy,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact AddPolicyFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact AddPolicyFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddPolicyFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.policy.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact AddPolicyFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for AddPolicyFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.target,
		fact.policy,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact AddPolicyFact) Token() []byte {
	return fact.token
}

func (fact AddPolicyFact) Sender() base.Address {
	return fact.sender
}

func (fact AddPolicyFact) Target() base.Address {
	return fact.target
}

func (fact AddPolicyFact) Policy() CollectionPolicy {
	return fact.policy
}

func (fact AddPolicyFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.target

	return as, nil
}

func (fact AddPolicyFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact AddPolicyFact) Rebuild() AddPolicyFact {
	policy := fact.policy.Rebuild()
	fact.policy = policy

	fact.h = fact.GenerateHash()

	return fact
}

type AddUserDefinedPolicy struct {
	currency.BaseOperation
}

func NewAddUserDefinedPolicy(fact AddPolicyFact, fs []base.FactSign, memo string) (AddUserDefinedPolicy, error) {
	bo, err := currency.NewBaseOperationFromFact(AddPolicyHint, fact, fs, memo)
	if err != nil {
		return AddUserDefinedPolicy{}, err
	}
	return AddUserDefinedPolicy{BaseOperation: bo}, nil
}
