package nft

import (
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MaxNFTHashLength = 1024

type NFTHash string

func (hs NFTHash) IsValid([]byte) error {
	if l := len(hs); l > MaxNFTHashLength {
		return util.ErrInvalid.Errorf("nft hash length over max, %d > %d", l, MaxNFTHashLength)
	}

	if hs != "" && strings.TrimSpace(string(hs)) == "" {
		return util.ErrInvalid.Errorf("empty nft hash")
	}

	return nil
}

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

var NFTHint = hint.MustNewHint("mitum-nft-nft-v0.0.1")

var MaxCreators = 10
var MaxCopyrighters = 10

type NFT struct {
	hint.BaseHinter
	id           NFTID
	active       bool
	owner        base.Address
	hash         NFTHash
	uri          URI
	approved     base.Address
	creators     Signers
	copyrighters Signers
}

func NewNFT(
	id NFTID,
	active bool,
	owner base.Address,
	hash NFTHash,
	uri URI,
	approved base.Address,
	creators Signers,
	copyrighters Signers,
) NFT {
	return NFT{
		BaseHinter:   hint.NewBaseHinter(NFTHint),
		id:           id,
		active:       active,
		owner:        owner,
		hash:         hash,
		uri:          uri,
		approved:     approved,
		creators:     creators,
		copyrighters: copyrighters,
	}
}

func (n NFT) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		n.id,
		n.owner,
		n.hash,
		n.uri,
		n.approved,
		n.creators,
		n.copyrighters,
	); err != nil {
		return err
	}

	if n.uri == "" {
		return util.ErrInvalid.Errorf("empty uri")
	}

	return nil
}

func (n NFT) Bytes() []byte {
	ba := make([]byte, 1)

	if n.active {
		ba[0] = 1
	} else {
		ba[0] = 0
	}

	return util.ConcatBytesSlice(
		n.id.Bytes(),
		ba,
		n.owner.Bytes(),
		n.hash.Bytes(),
		[]byte(n.uri.String()),
		n.approved.Bytes(),
		n.creators.Bytes(),
		n.copyrighters.Bytes(),
	)
}

func (n NFT) ID() NFTID {
	return n.id
}

func (n NFT) Active() bool {
	return n.active
}

func (n NFT) Owner() base.Address {
	return n.owner
}

func (n NFT) NFTHash() NFTHash {
	return n.hash
}

func (n NFT) URI() URI {
	return n.uri
}

func (n NFT) Approved() base.Address {
	return n.approved
}

func (n NFT) Creators() Signers {
	return n.creators
}

func (n NFT) Copyrighters() Signers {
	return n.copyrighters
}

func (n NFT) Equal(cn NFT) bool {
	if !n.ID().Equal(cn.ID()) {
		return false
	}

	if n.Active() != cn.Active() {
		return false
	}

	if !n.Owner().Equal(cn.Owner()) {
		return false
	}

	if n.NFTHash() != cn.NFTHash() {
		return false
	}

	if n.URI() != cn.URI() {
		return false
	}

	if !n.Approved().Equal(cn.Approved()) {
		return false
	}

	if !n.Creators().Equal(cn.Creators()) {
		return false
	}

	if !n.Copyrighters().Equal(cn.Copyrighters()) {
		return false
	}

	return n.ID().Equal(cn.ID())
}

func (n NFT) ExistsApproved() bool {
	return !n.approved.Equal(n.owner)
}
