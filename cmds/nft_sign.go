package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type NFTSignCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender        cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	NFT           NFTIDFlag           `arg:"" name:"nft" help:"target nft; \"<symbol>,<idx>\""`
	Currency      cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Qualification string              `name:"qualification" help:"target qualification; creator | copyrighter" optional:""`
	sender        base.Address
	nft           nft.NFTID
	qualification collection.Qualification
}

func NewNFTSignCommand() NFTSignCommand {
	cmd := NewbaseCommand()
	return NFTSignCommand{baseCommand: *cmd}
}

func (cmd *NFTSignCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *NFTSignCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	cmd.nft = n

	if cmd.Qualification == "" {
		cmd.qualification = collection.CreatorQualification
	} else {
		q := collection.Qualification(cmd.Qualification)
		if err := q.IsValid(nil); err != nil {
			return err
		}
		cmd.qualification = q
	}

	return nil

}

func (cmd *NFTSignCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create nft-sign operation")

	item := collection.NewNFTSignItem(cmd.qualification, cmd.nft, cmd.Currency.CID)

	fact := collection.NewNFTSignFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.NFTSignItem{item},
	)

	op, err := collection.NewNFTSign(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
