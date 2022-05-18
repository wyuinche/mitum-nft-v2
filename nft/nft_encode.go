package nft

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (nid *NFTID) unpack(
	enc encoder.Encoder,
	collection string,
	idx currency.Big,
) error {
	nid.collection = Symbol(collection)
	nid.idx = idx

	return nil
}

func (cp *Copyrighter) unpack(
	enc encoder.Encoder,
	set bool,
	bAddress base.AddressDecoder,
) error {
	cp.set = set

	address, err := bAddress.Encode(enc)
	if err != nil {
		return err
	}

	cp.address = address

	return nil
}

func (nft *NFT) unpack(
	enc encoder.Encoder,
	bId []byte,
	bOwner base.AddressDecoder,
	hash string,
	uri string,
	bApproved base.AddressDecoder,
	bCopyrighter []byte,
) error {
	if hinter, err := enc.Decode(bId); err != nil {
		return err
	} else if id, ok := hinter.(NFTID); !ok {
		return util.WrongTypeError.Errorf("not Copyrighter; %T", hinter)
	} else {
		nft.id = id
	}

	owner, err := bOwner.Encode(enc)
	if err != nil {
		return err
	}
	nft.owner = owner

	nft.hash = NFTHash(hash)
	nft.uri = NFTUri(uri)

	approved, err := bApproved.Encode(enc)
	if err != nil {
		return err
	}
	nft.approved = approved

	if hinter, err := enc.Decode(bCopyrighter); err != nil {
		return err
	} else if copyrighter, ok := hinter.(Copyrighter); !ok {
		return util.WrongTypeError.Errorf("not Copyrighter; %T", hinter)
	} else {
		nft.copyrighter = copyrighter
	}

	return nil
}
