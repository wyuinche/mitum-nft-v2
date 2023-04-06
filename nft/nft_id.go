package nft

import (
	"fmt"
	"strconv"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MaxNFTIndex uint64 = 10000

var NFTIDHint = hint.MustNewHint("mitum-nft-nft-id-v0.0.1")

type NFTID struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	index      uint64
}

func NewNFTID(collection extensioncurrency.ContractID, index uint64) NFTID {
	return NFTID{
		BaseHinter: hint.NewBaseHinter(NFTIDHint),
		collection: collection,
		index:      index,
	}
}

func (nid NFTID) IsValid([]byte) error {
	if nid.index > uint64(MaxNFTIndex) {
		return util.ErrInvalid.Errorf("nft-id index over max, %d > %d", nid.index, MaxNFTIndex)
	}

	if nid.index == 0 {
		return util.ErrInvalid.Errorf("zero nft-id index, %q", nid)
	}

	if err := nid.collection.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (nid NFTID) Bytes() []byte {
	return util.ConcatBytesSlice(
		nid.collection.Bytes(),
		util.Uint64ToBytes(nid.index),
	)
}

func (nid NFTID) Collection() extensioncurrency.ContractID {
	return nid.collection
}

func (nid NFTID) Index() uint64 {
	return nid.index
}

func (nid NFTID) Equal(id NFTID) bool {
	return nid.String() == id.String()
}

func (nid NFTID) String() string {
	index := strconv.FormatUint(nid.index, 10)

	l := len(strconv.FormatUint(uint64(MaxNFTIndex), 10)) - len(index)
	for i := 0; i < l; i++ {
		index = "0" + index
	}

	return fmt.Sprintf("%s-%s", nid.collection, index)
}
