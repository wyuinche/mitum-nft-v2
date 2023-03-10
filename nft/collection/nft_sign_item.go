package collection

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	CreatorQualification     = Qualification("creator")
	CopyrighterQualification = Qualification("copyrighter")
)

type Qualification string

func (q Qualification) Bytes() []byte {
	return []byte(q)
}

func (q Qualification) String() string {
	return string(q)
}

func (q Qualification) IsValid([]byte) error {
	if !(q == CreatorQualification || q == CopyrighterQualification) {
		return util.ErrInvalid.Errorf("invalid qualification, %q", q)
	}
	return nil
}

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	qualification Qualification
	nft           nft.NFTID
	currency      currency.CurrencyID
}

func NewNFTSignItem(q Qualification, n nft.NFTID, currency currency.CurrencyID) NFTSignItem {
	return NFTSignItem{
		BaseHinter:    hint.NewBaseHinter(NFTSignItemHint),
		qualification: q,
		nft:           n,
		currency:      currency,
	}
}

func (it NFTSignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.qualification.Bytes(),
		it.nft.Bytes(),
		it.currency.Bytes(),
	)
}

func (it NFTSignItem) IsValid([]byte) error {
	return util.CheckIsValiders(nil, false, it.BaseHinter, it.qualification, it.nft, it.currency)
}

func (it NFTSignItem) Qualification() Qualification {
	return it.qualification
}

func (it NFTSignItem) NFT() nft.NFTID {
	return it.nft
}

func (it NFTSignItem) Currency() currency.CurrencyID {
	return it.currency
}
