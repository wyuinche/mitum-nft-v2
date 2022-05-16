package nft

import (
	"fmt"
	"regexp"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	NFTIDType   = hint.Type("mitum-nft-nft-id")
	NFTIDHint   = hint.NewHint(NFTIDType, "v0.0.1")
	NFTIDHinter = NFTID{BaseHinter: hint.NewBaseHinter(NFTIDHint)}
)

type NFTID struct {
	hint.BaseHinter
	collection Symbol
	idx        currency.Big
}

func NewNFTID(collection Symbol, idx currency.Big) NFTID {
	return NFTID{
		BaseHinter: hint.NewBaseHinter(NFTIDHint),
		collection: collection,
		idx:        idx,
	}
}

func MustNewNFTID(collection Symbol, idx currency.Big) NFTID {
	id := NewNFTID(collection, idx)

	if err := id.IsValid(nil); err != nil {
		panic(err)
	}

	return id
}

func (nid NFTID) Bytes() []byte {
	return util.ConcatBytesSlice(
		nid.collection.Bytes(),
		nid.idx.Bytes(),
	)
}

func (nid NFTID) IsValid([]byte) error {
	if !nid.idx.OverZero() {
		return isvalid.InvalidError.Errorf("zero collection idx; %s", nid.idx.String())
	}

	if err := isvalid.Check(nil, false,
		nid.BaseHinter,
		nid.collection,
		nid.idx,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid NFTID: %w", err)
	}

	return nil
}

func (nid NFTID) Symbol() Symbol {
	return nid.collection
}

func (nid NFTID) Idx() currency.Big {
	return nid.idx
}

func (nid NFTID) String() string {
	return fmt.Sprintf("%s-%s)", nid.collection.String(), nid.idx.String())
}

type NFTUri string

func (uri NFTUri) Bytes() []byte {
	return []byte(uri)
}

func (uri NFTUri) String() string {
	return string(uri)
}

func (uri NFTUri) IsValid([]byte) error {
	if len(uri) == 0 {
		return isvalid.InvalidError.Errorf("empty nft uri")
	}

	return nil
}

type NFTHash string

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

func (hs NFTHash) IsValid([]byte) error {
	if len(hs) == 0 {
		return isvalid.InvalidError.Errorf("empty nft hash")
	}

	return nil
}

type PostCloseTime string

var RevalidPostCloseTime = regexp.MustCompile(`^\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1-2]\d|3[0-1])T(?:[0-1]\d|2[0-3]):[0-5]\d:[0-5]\dZ$`)

func (pct PostCloseTime) Bytes() []byte {
	return []byte(pct)
}

func (pct PostCloseTime) String() string {
	return string(pct)
}

func (pct PostCloseTime) IsValid([]byte) error {
	if !RevalidPostCloseTime.Match([]byte(pct)) {
		return isvalid.InvalidError.Errorf("wrong post close time, %q", pct)
	}

	return nil
}

var (
	SellPostOption    = PostOption("sell")
	AuctionPostOption = PostOption("auction")
)

type PostOption string

func (po PostOption) Bytes() []byte {
	return []byte(po)
}

func (po PostOption) String() string {
	return string(po)
}

func (po PostOption) IsValid([]byte) error {
	if !(po == SellPostOption || po == AuctionPostOption) {
		return isvalid.InvalidError.Errorf("invalid post option, %s", po)
	}

	return nil
}

var (
	PostInfoType   = hint.Type("mitum-nft-post-info")
	PostInfoHint   = hint.NewHint(PostInfoType, "v0.0.1")
	PostInfoHinter = PostInfo{BaseHinter: hint.NewBaseHinter(PostInfoHint)}
)

type PostInfo struct {
	hint.BaseHinter
	posting   bool
	option    PostOption
	broker    Symbol
	closeTime PostCloseTime
	price     currency.Amount
}

func NewPostInfo(posting bool, option PostOption, broker Symbol, closeTime PostCloseTime, price currency.Amount) PostInfo {
	if posting {
		return PostInfo{
			posting:   true,
			broker:    broker,
			option:    option,
			closeTime: closeTime,
			price:     price,
		}
	}

	return PostInfo{
		posting:   false,
		option:    option,
		broker:    "",
		closeTime: "",
		price:     currency.NewAmount(currency.NewBig(0), ""),
	}
}

func MustNewPostInfo(posting bool, option PostOption, broker Symbol, closeTime PostCloseTime, price currency.Amount) PostInfo {
	info := NewPostInfo(posting, option, broker, closeTime, price)

	if err := info.IsValid(nil); err != nil {
		panic(err)
	}

	return info
}

func (info PostInfo) Bytes() []byte {
	if info.posting {
		return util.ConcatBytesSlice(
			[]byte{1},
			info.option.Bytes(),
			info.broker.Bytes(),
			info.closeTime.Bytes(),
			info.price.Bytes(),
		)
	}

	return []byte{0}
}

func (info PostInfo) IsValid([]byte) error {
	if err := info.BaseHinter.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid PostInfo: %w", err)
	}

	if !info.posting {
		return nil
	}

	if err := isvalid.Check(
		nil, false,
		info.option,
		info.broker,
		info.closeTime,
		info.price,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid PostInfo: %w", err)
	}
	return nil
}

func (info PostInfo) Option() PostOption {
	return info.option
}

func (info PostInfo) Broker() Symbol {
	return info.broker
}

func (info PostInfo) Price() currency.Amount {
	return info.price
}

func (info PostInfo) CloseTime() PostCloseTime {
	return info.closeTime
}

func (info PostInfo) Rebuild() PostInfo {
	return info
}

var (
	CopyrighterType   = hint.Type("mitum-nft-copyrighter")
	CopyrighterHint   = hint.NewHint(CopyrighterType, "v0.0.1")
	CopyrighterHinter = NFT{BaseHinter: hint.NewBaseHinter(CopyrighterHint)}
)

type Copyrighter struct {
	hint.BaseHinter
	set     bool
	address base.Address
}

func NewCopyrighter(set bool, address base.Address) Copyrighter {
	if set {
		return Copyrighter{
			BaseHinter: hint.NewBaseHinter(CopyrighterHint),
			set:        set,
			address:    nil,
		}
	}

	return Copyrighter{
		BaseHinter: hint.NewBaseHinter(CopyrighterHint),
		set:        set,
		address:    address,
	}
}

func MustNewCopyrighter(set bool, address base.Address) Copyrighter {
	copyrighter := NewCopyrighter(set, address)

	if err := copyrighter.IsValid(nil); err != nil {
		panic(err)
	}

	return copyrighter
}

func (cpr Copyrighter) Bytes() []byte {
	if cpr.set {
		return util.ConcatBytesSlice(
			[]byte{1},
			cpr.address.Bytes(),
		)
	}

	return []byte{0}
}

func (cpr Copyrighter) String() string {
	if cpr.set {
		return cpr.address.String()
	}

	return ""
}

func (cpr Copyrighter) IsValid([]byte) error {
	if err := cpr.BaseHinter.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid Copyrighter: %w", err)
	}

	if !cpr.set {
		return nil
	}

	if err := cpr.address.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid Copyrighter: %w", err)
	}

	return nil
}

var (
	NFTType   = hint.Type("mitum-nft-post-info")
	NFTHint   = hint.NewHint(NFTType, "v0.0.1")
	NFTHinter = NFT{BaseHinter: hint.NewBaseHinter(NFTHint)}
)

type NFT struct {
	hint.BaseHinter
	id          NFTID
	owner       base.Address
	hash        NFTHash
	uri         NFTUri
	approved    base.Address
	copyrighter Copyrighter
	post        PostInfo
}

func NewNFT(id NFTID, owner base.Address, hash NFTHash, uri NFTUri, approved base.Address, copyrighter Copyrighter, post PostInfo) NFT {
	return NFT{
		BaseHinter:  hint.NewBaseHinter(NFTHint),
		id:          id,
		owner:       owner,
		hash:        hash,
		uri:         uri,
		approved:    approved,
		copyrighter: copyrighter,
		post:        post,
	}
}

func MustNewNFT(id NFTID, owner base.Address, hash NFTHash, uri NFTUri, approved base.Address, copyrighter Copyrighter, post PostInfo) NFT {
	nft := NewNFT(id, owner, hash, uri, approved, copyrighter, post)

	if err := nft.IsValid(nil); err != nil {
		panic(err)
	}

	return nft
}

func (nft NFT) Bytes() []byte {
	return util.ConcatBytesSlice(
		nft.id.Bytes(),
		nft.owner.Bytes(),
		nft.hash.Bytes(),
		nft.uri.Bytes(),
		nft.approved.Bytes(),
		nft.copyrighter.Bytes(),
		nft.post.Bytes(),
	)
}

func (nft NFT) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		nft.id,
		nft.owner,
		nft.hash,
		nft.uri,
		nft.approved,
		nft.copyrighter,
		nft.post,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid NFT: %w", err)
	}
	return nil
}

func (nft NFT) ID() NFTID {
	return nft.id
}

func (nft NFT) Owner() base.Address {
	return nft.owner
}

func (nft NFT) Hash() NFTHash {
	return nft.hash
}

func (nft NFT) Uri() NFTUri {
	return nft.uri
}

func (nft NFT) Approved() base.Address {
	return nft.approved
}

func (nft NFT) Copyrighter() Copyrighter {
	return nft.copyrighter
}
