package cmds

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/ProtoconNet/mitum2/util/ps"
	"github.com/rs/zerolog"
)

type baseCommand struct {
	enc  *jsonenc.Encoder
	encs *encoder.Encoders
	log  *zerolog.Logger
	Out  io.Writer `kong:"-"`
}

func NewbaseCommand() *baseCommand {
	return &baseCommand{
		Out: os.Stdout,
	}
}

func (cmd *baseCommand) prepare(pctx context.Context) (context.Context, error) {
	pps := ps.NewPS("cmd")

	_ = pps.
		AddOK(launch.PNameEncoder, PEncoder, nil)

	_ = pps.POK(launch.PNameEncoder).
		PostAddOK(launch.PNameAddHinters, PAddHinters)

	var log *logging.Logging
	if err := util.LoadFromContextOK(pctx, launch.LoggingContextKey, &log); err != nil {
		return pctx, err
	}

	cmd.log = log.Log()

	pctx, err := pps.Run(pctx) //revive:disable-line:modifies-parameter
	if err != nil {
		return pctx, err
	}

	return pctx, util.LoadFromContextOK(pctx,
		launch.EncodersContextKey, &cmd.encs,
		launch.EncoderContextKey, &cmd.enc,
	)
}

func (cmd *baseCommand) print(f string, a ...interface{}) {
	_, _ = fmt.Fprintf(cmd.Out, f, a...)
	_, _ = fmt.Fprintln(cmd.Out)
}

func PAddHinters(ctx context.Context) (context.Context, error) {
	e := util.StringErrorFunc("failed to add hinters")

	var enc encoder.Encoder
	if err := util.LoadFromContextOK(ctx, launch.EncoderContextKey, &enc); err != nil {
		return ctx, e(err, "")
	}
	var benc encoder.Encoder
	if err := util.LoadFromContextOK(ctx, BEncoderContextKey, &benc); err != nil {
		return ctx, e(err, "")
	}

	if err := LoadHinters(enc); err != nil {
		return ctx, e(err, "")
	}

	if err := LoadHinters(benc); err != nil {
		return ctx, e(err, "")
	}

	return ctx, nil
}
