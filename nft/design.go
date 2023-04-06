package nft

import (
	"net/url"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var MaxPaymentParameter uint = 99

type PaymentParameter uint

func (pp PaymentParameter) IsValid([]byte) error {
	if uint(pp) > MaxPaymentParameter {
		return util.ErrInvalid.Errorf("payment parameter over max, %d > %d", pp, MaxPaymentParameter)
	}

	return nil
}

func (pp PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(pp))
}

func (pp PaymentParameter) Uint() uint {
	return uint(pp)
}

var MaxURILength = 1000

type URI string

func (uri URI) IsValid([]byte) error {
	if _, err := url.Parse(string(uri)); err != nil {
		return err
	}

	if l := len(uri); l > MaxURILength {
		return util.ErrInvalid.Errorf("uri length over max, %d > %d", l, MaxURILength)
	}

	if uri != "" && strings.TrimSpace(string(uri)) == "" {
		return util.ErrInvalid.Errorf("empty uri")
	}

	return nil
}

func (uri URI) Bytes() []byte {
	return []byte(uri)
}

func (uri URI) String() string {
	return string(uri)
}

var DesignHint = hint.MustNewHint("mitum-nft-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	parent  base.Address
	creator base.Address
	symbol  extensioncurrency.ContractID
	active  bool
	policy  BasePolicy
}

func NewDesign(parent base.Address, creator base.Address, symbol extensioncurrency.ContractID, active bool, policy BasePolicy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		parent:     parent,
		creator:    creator,
		symbol:     symbol,
		active:     active,
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.parent,
		de.creator,
		de.symbol,
		de.policy,
	); err != nil {
		return err
	}

	if de.parent.Equal(de.creator) {
		return util.ErrInvalid.Errorf("parent and creator are the same, %q == %q", de.parent, de.creator)
	}

	return nil
}

func (de Design) Bytes() []byte {
	ab := make([]byte, 1)
	if de.active {
		ab[0] = 1
	} else {
		ab[0] = 0
	}

	return util.ConcatBytesSlice(
		de.parent.Bytes(),
		de.creator.Bytes(),
		de.symbol.Bytes(),
		ab,
		de.policy.Bytes(),
	)
}

func (de Design) Hash() util.Hash {
	return de.GenerateHash()
}

func (de Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(de.Bytes())
}

func (de Design) Parent() base.Address {
	return de.parent
}

func (de Design) Creator() base.Address {
	return de.creator
}

func (de Design) Symbol() extensioncurrency.ContractID {
	return de.symbol
}

func (de Design) Active() bool {
	return de.active
}

func (de Design) Policy() BasePolicy {
	return de.policy
}

func (de Design) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = de.parent
	as[1] = de.creator

	if ads, err := de.Policy().Addresses(); err != nil {
		return as, err
	} else {
		as = append(as, ads...)
	}

	return as, nil
}

func (de Design) Equal(cd Design) bool {
	if !de.parent.Equal(cd.parent) {
		return false
	}

	if !de.creator.Equal(cd.creator) {
		return false
	}

	if de.symbol != cd.symbol {
		return false
	}

	if de.active != cd.active {
		return false
	}

	if !de.policy.Equal(cd.policy) {
		return false
	}

	if de.Hash() != cd.Hash() {
		return false
	}

	return true
}

type BasePolicy interface {
	util.IsValider
	Bytes() []byte
	Addresses() ([]base.Address, error)
	Equal(c BasePolicy) bool
}
