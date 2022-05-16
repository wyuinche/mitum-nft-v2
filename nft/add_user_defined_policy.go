package nft

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	AddUserDefinedPolicyFactType   = hint.Type("mitum-nft-add-user-defined-policy-operation-fact")
	AddUserDefinedPolicyFactHint   = hint.NewHint(AddUserDefinedPolicyFactType, "v0.0.1")
	AddUserDefinedPolicyFactHinter = AddUserDefinedPolicyFact{BaseHinter: hint.NewBaseHinter(AddUserDefinedPolicyFactHint)}
	AddUserDefinedPolicyType       = hint.Type("mitum-nft-add-user-defined-policy-operation")
	AddUserDefinedPolicyHint       = hint.NewHint(AddUserDefinedPolicyType, "v0.0.1")
	AddUserDefinedPolicyHinter     = AddUserDefinedPolicy{BaseOperation: operationHinter(AddUserDefinedPolicyHint)}
)

type AddUserDefinedPolicyFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	target base.Address
	policy BaseUserDefinedPolicy
	cid    currency.CurrencyID
}

func NewAddUserDefinedPolicyFact(token []byte, sender base.Address, target base.Address, policy BaseUserDefinedPolicy, cid currency.CurrencyID) AddUserDefinedPolicyFact {
	fact := AddUserDefinedPolicyFact{
		BaseHinter: hint.NewBaseHinter(AddUserDefinedPolicyFactHint),
		token:      token,
		sender:     sender,
		target:     target,
		policy:     policy,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact AddUserDefinedPolicyFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact AddUserDefinedPolicyFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AddUserDefinedPolicyFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.policy.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact AddUserDefinedPolicyFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for AddUserDefinedPolicyFact")
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

func (fact AddUserDefinedPolicyFact) Token() []byte {
	return fact.token
}

func (fact AddUserDefinedPolicyFact) Sender() base.Address {
	return fact.sender
}

func (fact AddUserDefinedPolicyFact) Target() base.Address {
	return fact.target
}

func (fact AddUserDefinedPolicyFact) Policy() BaseUserDefinedPolicy {
	return fact.policy
}

func (fact AddUserDefinedPolicyFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.target

	return as, nil
}

func (fact AddUserDefinedPolicyFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact AddUserDefinedPolicyFact) Rebuild() AddUserDefinedPolicyFact {
	policy := fact.policy.Rebuild()
	fact.policy = policy

	fact.h = fact.GenerateHash()

	return fact
}

type AddUserDefinedPolicy struct {
	currency.BaseOperation
}

func NewAddUserDefinedPolicy(fact AddUserDefinedPolicyFact, fs []base.FactSign, memo string) (AddUserDefinedPolicy, error) {
	bo, err := currency.NewBaseOperationFromFact(AddUserDefinedPolicyHint, fact, fs, memo)
	if err != nil {
		return AddUserDefinedPolicy{}, err
	}
	return AddUserDefinedPolicy{BaseOperation: bo}, nil
}
