package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (nid *NFTID) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	idx uint64,
) error {
	nid.BaseHinter = hint.NewBaseHinter(ht)
	nid.collection = extensioncurrency.ContractID(col)
	nid.index = idx

	return nil
}
