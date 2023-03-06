package cmds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

type SignerFlag struct {
	address string
	share   uint
}

func (v *SignerFlag) UnmarshalText(b []byte) error {
	l := strings.SplitN(string(b), ",", 2)
	if len(l) != 2 {
		return fmt.Errorf("invalid signer; %q", string(b))
	}

	v.address = l[0]

	if share, err := strconv.ParseUint(l[1], 10, 8); err != nil {
		return err
	} else if share > uint64(nft.MaxSignerShare) {
		return errors.Errorf("share is over max; %d > %d", share, nft.MaxSignerShare)
	} else {
		v.share = uint(share)
	}

	return nil
}

func (v *SignerFlag) String() string {
	s := fmt.Sprintf("%s,%d", v.address, v.share)
	return s
}

func (v *SignerFlag) Encode(enc encoder.Encoder) (base.Address, error) {
	return base.DecodeAddress(v.address, enc)
}
