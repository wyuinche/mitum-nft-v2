package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	nftcollection "github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
)

type CollectionRegisterCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender     cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency   cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Target     cmds.AddressFlag    `arg:"" name:"target" help:"target account to register policy" required:"true"`
	Collection string              `arg:"" name:"collection" help:"collection symbol" required:"true"`
	Name       string              `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty    uint                `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	URI        string              `name:"uri" help:"collection uri" optional:""`
	White      cmds.AddressFlag    `name:"white" help:"whitelisted address" optional:""`
	sender     base.Address
	target     base.Address
	form       nftcollection.CollectionRegisterForm
}

func NewCollectionRegisterCommand() CollectionRegisterCommand {
	cmd := NewbaseCommand()
	return CollectionRegisterCommand{baseCommand: *cmd}
}

func (cmd *CollectionRegisterCommand) Run(pctx context.Context) error {
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

func (cmd *CollectionRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Target.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid target format; %q", cmd.Target)
	} else {
		cmd.target = a
	}

	var white base.Address = nil
	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(enc); err != nil {
			return errors.Wrapf(err, "invalid white format; %q", cmd.White)
		} else {
			white = a
		}
	}

	collection := extensioncurrency.ContractID(cmd.Collection)
	if err := collection.IsValid(nil); err != nil {
		return err
	}

	name := nftcollection.CollectionName(cmd.Name)
	if err := name.IsValid(nil); err != nil {
		return err
	}

	royalty := nft.PaymentParameter(cmd.Royalty)
	if err := royalty.IsValid(nil); err != nil {
		return err
	}

	uri := nft.URI(cmd.URI)
	if err := uri.IsValid(nil); err != nil {
		return err
	}

	whites := []base.Address{}
	if white != nil {
		whites = append(whites, white)
	}

	form := nftcollection.NewCollectionRegisterForm(cmd.target, collection, name, royalty, uri, whites)
	if err := form.IsValid(nil); err != nil {
		return err
	}
	cmd.form = form

	return nil
}

func (cmd *CollectionRegisterCommand) createOperation() (base.Operation, error) {
	fact := nftcollection.NewCollectionRegisterFact([]byte(cmd.Token), cmd.sender, cmd.form, cmd.Currency.CID)

	op, err := nftcollection.NewCollectionRegister(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collection-register operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collection-register operation")
	}

	return op, nil
}
