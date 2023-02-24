package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-currency-extension/digest"
	isaacoperation "github.com/ProtoconNet/mitum-currency-extension/isaac"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	digestisaac "github.com/spikeekips/mitum-currency/digest/isaac"
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/encoder"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: currency.BaseStateHint, Instance: currency.BaseState{}},
	{Hint: currency.NodeHint, Instance: currency.BaseNode{}},
	{Hint: currency.AccountHint, Instance: currency.Account{}},
	{Hint: currency.AddressHint, Instance: currency.Address{}},
	{Hint: currency.AmountHint, Instance: currency.Amount{}},
	{Hint: currency.CreateAccountsItemMultiAmountsHint, Instance: currency.CreateAccountsItemMultiAmounts{}},
	{Hint: currency.CreateAccountsItemSingleAmountHint, Instance: currency.CreateAccountsItemSingleAmount{}},
	{Hint: currency.CreateAccountsHint, Instance: currency.CreateAccounts{}},
	{Hint: currency.KeyUpdaterHint, Instance: currency.KeyUpdater{}},
	{Hint: currency.TransfersItemMultiAmountsHint, Instance: currency.TransfersItemMultiAmounts{}},
	{Hint: currency.TransfersItemSingleAmountHint, Instance: currency.TransfersItemSingleAmount{}},
	{Hint: currency.TransfersHint, Instance: currency.Transfers{}},
	{Hint: extensioncurrency.CurrencyDesignHint, Instance: extensioncurrency.CurrencyDesign{}},
	{Hint: extensioncurrency.CurrencyPolicyHint, Instance: extensioncurrency.CurrencyPolicy{}},
	{Hint: extensioncurrency.CurrencyRegisterHint, Instance: extensioncurrency.CurrencyRegister{}},
	{Hint: extensioncurrency.CurrencyPolicyUpdaterHint, Instance: extensioncurrency.CurrencyPolicyUpdater{}},
	{Hint: currency.SuffrageInflationHint, Instance: currency.SuffrageInflation{}},
	{Hint: extensioncurrency.ContractAccountKeysHint, Instance: extensioncurrency.ContractAccountKeys{}},
	{Hint: extensioncurrency.CreateContractAccountsItemMultiAmountsHint, Instance: extensioncurrency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extensioncurrency.CreateContractAccountsItemSingleAmountHint, Instance: extensioncurrency.CreateContractAccountsItemSingleAmount{}},
	{Hint: extensioncurrency.CreateContractAccountsHint, Instance: extensioncurrency.CreateContractAccounts{}},
	{Hint: extensioncurrency.WithdrawsItemMultiAmountsHint, Instance: extensioncurrency.WithdrawsItemMultiAmounts{}},
	{Hint: extensioncurrency.WithdrawsItemSingleAmountHint, Instance: extensioncurrency.WithdrawsItemSingleAmount{}},
	{Hint: extensioncurrency.WithdrawsHint, Instance: extensioncurrency.Withdraws{}},
	// {Hint: currency.FeeOperationFactHint, Instance: currency.FeeOperationFact{}},
	// {Hint: currency.FeeOperationHint, Instance: currency.FeeOperation{}},
	{Hint: extensioncurrency.GenesisCurrenciesFactHint, Instance: extensioncurrency.GenesisCurrenciesFact{}},
	{Hint: extensioncurrency.GenesisCurrenciesHint, Instance: extensioncurrency.GenesisCurrencies{}},
	{Hint: currency.AccountKeysHint, Instance: currency.BaseAccountKeys{}},
	{Hint: currency.AccountKeyHint, Instance: currency.BaseAccountKey{}},
	{Hint: extensioncurrency.NilFeeerHint, Instance: extensioncurrency.NilFeeer{}},
	{Hint: extensioncurrency.FixedFeeerHint, Instance: extensioncurrency.FixedFeeer{}},
	{Hint: extensioncurrency.RatioFeeerHint, Instance: extensioncurrency.RatioFeeer{}},
	{Hint: currency.AccountStateValueHint, Instance: currency.AccountStateValue{}},
	{Hint: currency.BalanceStateValueHint, Instance: currency.BalanceStateValue{}},
	{Hint: extensioncurrency.ContractAccountStateValueHint, Instance: extensioncurrency.ContractAccountStateValue{}},
	{Hint: extensioncurrency.CurrencyDesignStateValueHint, Instance: extensioncurrency.CurrencyDesignStateValue{}},
	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},
	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},
	{Hint: currency.CreateAccountsFactHint, Instance: currency.CreateAccountsFact{}},
	{Hint: currency.KeyUpdaterFactHint, Instance: currency.KeyUpdaterFact{}},
	{Hint: currency.TransfersFactHint, Instance: currency.TransfersFact{}},
	{Hint: extensioncurrency.CurrencyRegisterFactHint, Instance: extensioncurrency.CurrencyRegisterFact{}},
	{Hint: extensioncurrency.CurrencyPolicyUpdaterFactHint, Instance: extensioncurrency.CurrencyPolicyUpdaterFact{}},
	{Hint: currency.SuffrageInflationFactHint, Instance: currency.SuffrageInflationFact{}},
	{Hint: extensioncurrency.CreateContractAccountsFactHint, Instance: extensioncurrency.CreateContractAccountsFact{}},
	{Hint: extensioncurrency.WithdrawsFactHint, Instance: extensioncurrency.WithdrawsFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for i := range Hinters {
		if err := enc.Add(Hinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := enc.Add(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
