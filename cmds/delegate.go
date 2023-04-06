package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	nftcollection "github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type DelegateCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender     cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Collection string              `arg:"" name:"collection" help:"collection symbol" required:"true"`
	Agent      cmds.AddressFlag    `arg:"" name:"agent" help:"agent account address"`
	Currency   cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Mode       string              `name:"mode" help:"delegate mode" optional:""`
	sender     base.Address
	symbol     extensioncurrency.ContractID
	agent      base.Address
	mode       nftcollection.DelegateMode
}

func NewDelegateCommand() DelegateCommand {
	cmd := NewbaseCommand()
	return DelegateCommand{baseCommand: *cmd}
}

func (cmd *DelegateCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *DelegateCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	symbol := extensioncurrency.ContractID(cmd.Collection)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}
	cmd.symbol = symbol

	if a, err := cmd.Agent.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid agent format; %q", cmd.Agent)
	} else {
		cmd.agent = a
	}

	if len(cmd.Mode) < 1 {
		cmd.mode = nftcollection.DelegateAllow
	} else {
		mode := nftcollection.DelegateMode(cmd.Mode)
		if err := mode.IsValid(nil); err != nil {
			return err
		}
		cmd.mode = mode
	}

	return nil

}

func (cmd *DelegateCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create delegate operation")

	items := []nftcollection.DelegateItem{nftcollection.NewDelegateItem(cmd.symbol, cmd.agent, cmd.mode, cmd.Currency.CID)}

	fact := nftcollection.NewDelegateFact([]byte(cmd.Token), cmd.sender, items)

	op, err := nftcollection.NewDelegate(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
