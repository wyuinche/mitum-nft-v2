package nft

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
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
