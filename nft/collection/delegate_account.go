package collection

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
	DelegateAccountFactType   = hint.Type("mitum-nft-delegate-account-operation-fact")
	DelegateAccountFactHint   = hint.NewHint(DelegateAccountFactType, "v0.0.1")
	DelegateAccountFactHinter = DelegateAccountFact{BaseHinter: hint.NewBaseHinter(DelegateAccountFactHint)}
	DelegateAccountType       = hint.Type("mitum-nft-delegate-account-operation")
	DelegateAccountHint       = hint.NewHint(DelegateAccountType, "v0.0.1")
	DelegateAccountHinter     = DelegateAccount{BaseOperation: operationHinter(DelegateAccountHint)}
)

type DelegateAccountFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	agent  base.Address
	cid    currency.CurrencyID
}

func NewDelegateAccountFact(token []byte, sender base.Address, agent base.Address, cid currency.CurrencyID) DelegateAccountFact {
	fact := DelegateAccountFact{
		BaseHinter: hint.NewBaseHinter(DelegateAccountFactHint),
		token:      token,
		sender:     sender,
		agent:      agent,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact DelegateAccountFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact DelegateAccountFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DelegateAccountFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.agent.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact DelegateAccountFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for DelegateAccountFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.agent,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact DelegateAccountFact) Token() []byte {
	return fact.token
}

func (fact DelegateAccountFact) Sender() base.Address {
	return fact.sender
}

func (fact DelegateAccountFact) Agent() base.Address {
	return fact.agent
}

func (fact DelegateAccountFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact DelegateAccountFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact DelegateAccountFact) Rebuild() DelegateAccountFact {
	fact.h = fact.GenerateHash()

	return fact
}

type DelegateAccount struct {
	currency.BaseOperation
}

func NewDelegateAccount(fact DelegateAccountFact, fs []base.FactSign, memo string) (DelegateAccount, error) {
	bo, err := currency.NewBaseOperationFromFact(DelegateAccountHint, fact, fs, memo)
	if err != nil {
		return DelegateAccount{}, err
	}
	return DelegateAccount{BaseOperation: bo}, nil
}
