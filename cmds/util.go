package cmds

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/isaac"
	isaacblock "github.com/spikeekips/mitum/isaac/block"
	isaacdatabase "github.com/spikeekips/mitum/isaac/database"
	isaacnetwork "github.com/spikeekips/mitum/isaac/network"
	isaacoperation "github.com/spikeekips/mitum/isaac/operation"
	isaacstates "github.com/spikeekips/mitum/isaac/states"
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/network/quicmemberlist"
	"github.com/spikeekips/mitum/network/quicstream"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/logging"
	"github.com/spikeekips/mitum/util/ps"
)

var (
	PNameDigestDesign           = ps.Name("digest-design")
	PNameOperationProcessorsMap = ps.Name("mitum-currency-operation-processors-map")
	PNameGenerateGenesis        = ps.Name("mitum-currency-generate-genesis")
	PNameDigestAPIHandlers      = ps.Name("mitum-currency-digest-api-handlers")
	PNameDigesterFollowUp       = ps.Name("mitum-currency-followup_digester")
	BEncoderContextKey          = util.ContextKey("bencoder")
)

func LoadFromStdInput() ([]byte, error) {
	var b []byte
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			b = append(b, sc.Bytes()...)
			b = append(b, []byte("\n")...)
		}

		if err := sc.Err(); err != nil {
			return nil, err
		}
	}

	return bytes.TrimSpace(b), nil
}

func GenerateED25519Privatekey() (ed25519.PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)

	return priv, err
}

func GenerateTLSCertsPair(host string, key ed25519.PrivateKey) (*pem.Block, *pem.Block, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     []string{host},
		NotBefore:    time.Now().Add(time.Minute * -1),
		NotAfter:     time.Now().Add(time.Hour * 24 * 1825),
	}

	if i := net.ParseIP(host); i != nil {
		template.IPAddresses = []net.IP{i}
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		key.Public().(ed25519.PublicKey),
		key,
	)
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	return &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes},
		&pem.Block{Type: "CERTIFICATE", Bytes: certDER},
		nil
}

func GenerateTLSCerts(host string, key ed25519.PrivateKey) ([]tls.Certificate, error) {
	k, c, err := GenerateTLSCertsPair(host, key)
	if err != nil {
		return nil, err
	}

	certificate, err := tls.X509KeyPair(pem.EncodeToMemory(c), pem.EncodeToMemory(k))
	if err != nil {
		return nil, err
	}

	return []tls.Certificate{certificate}, nil
}

type NetworkIDFlag []byte

func (v *NetworkIDFlag) UnmarshalText(b []byte) error {
	*v = b

	return nil
}

func (v NetworkIDFlag) NetworkID() base.NetworkID {
	return base.NetworkID(v)
}

func PrettyPrint(out io.Writer, i interface{}) {
	var b []byte
	b, err := enc.Marshal(i)
	if err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(out, string(b))
}

func POperationProcessorsMap(ctx context.Context) (context.Context, error) {

	var params *isaac.LocalParams
	var db isaac.Database

	if err := util.LoadFromContextOK(ctx,
		launch.LocalParamsContextKey, &params,
		launch.CenterDatabaseContextKey, &db,
	); err != nil {
		return ctx, err
	}

	limiterf, err := launch.NewSuffrageCandidateLimiterFunc(ctx)
	if err != nil {
		return ctx, err
	}

	set := hint.NewCompatibleSet()

	opr := extensioncurrency.NewOperationProcessor()
	opr.SetProcessor(currency.CreateAccountsHint, extensioncurrency.NewCreateAccountsProcessor())
	opr.SetProcessor(currency.KeyUpdaterHint, extensioncurrency.NewKeyUpdaterProcessor())
	opr.SetProcessor(currency.TransfersHint, extensioncurrency.NewTransfersProcessor())
	opr.SetProcessor(currency.CurrencyRegisterHint, extensioncurrency.NewCurrencyRegisterProcessor(params.Threshold()))
	opr.SetProcessor(currency.CurrencyPolicyUpdaterHint, extensioncurrency.NewCurrencyPolicyUpdaterProcessor(params.Threshold()))
	opr.SetProcessor(currency.SuffrageInflationHint, extensioncurrency.NewSuffrageInflationProcessor(params.Threshold()))
	opr.SetProcessor(extensioncurrency.CreateContractAccountsHint, extensioncurrency.NewCreateContractAccountsProcessor())
	opr.SetProcessor(extensioncurrency.WithdrawsHint, extensioncurrency.NewWithdrawsProcessor())
	opr.SetProcessor(collection.CollectionRegisterHint, collection.NewCollectionRegisterProcessor())
	opr.SetProcessor(collection.CollectionPolicyUpdaterHint, collection.NewCollectionPolicyUpdaterProcessor())
	opr.SetProcessor(collection.MintHint, collection.NewMintProcessor())
	opr.SetProcessor(collection.NFTTransferHint, collection.NewNFTTransferProcessor())
	opr.SetProcessor(collection.DelegateHint, collection.NewDelegateProcessor())
	opr.SetProcessor(collection.ApproveHint, collection.NewApproveProcessor())

	_ = set.Add(currency.CreateAccountsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.KeyUpdaterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.TransfersHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.CurrencyRegisterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.CurrencyPolicyUpdaterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.SuffrageInflationHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(extensioncurrency.CreateContractAccountsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(extensioncurrency.WithdrawsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.CollectionRegisterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.CollectionPolicyUpdaterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.MintHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.NFTTransferHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.DelegateHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(collection.ApproveHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaacoperation.SuffrageCandidateHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageCandidateProcessor(
			height,
			db.State,
			limiterf,
			nil,
			policy.SuffrageCandidateLifespan(),
		)
	})

	_ = set.Add(isaacoperation.SuffrageJoinHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageJoinProcessor(
			height,
			params.Threshold(),
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaac.SuffrageWithdrawOperationHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageWithdrawProcessor(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaacoperation.SuffrageDisjoinHint, func(height base.Height) (base.OperationProcessor, error) {
		return isaacoperation.NewSuffrageDisjoinProcessor(
			height,
			db.State,
			nil,
			nil,
		)
	})

	ctx = context.WithValue(ctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter

	return ctx, nil
}

func PGenerateGenesis(ctx context.Context) (context.Context, error) {
	e := util.StringErrorFunc("failed to generate genesis block")

	var log *logging.Logging
	var design launch.NodeDesign
	var genesisDesign launch.GenesisDesign
	var enc encoder.Encoder
	var local base.LocalNode
	var params *isaac.LocalParams
	var db isaac.Database

	if err := util.LoadFromContextOK(ctx,
		launch.LoggingContextKey, &log,
		launch.DesignContextKey, &design,
		launch.GenesisDesignContextKey, &genesisDesign,
		launch.EncoderContextKey, &enc,
		launch.LocalContextKey, &local,
		launch.LocalParamsContextKey, &params,
		launch.CenterDatabaseContextKey, &db,
	); err != nil {
		return ctx, e(err, "")
	}

	g := NewGenesisBlockGenerator(
		local,
		params.NetworkID(),
		enc,
		db,
		launch.LocalFSDataDirectory(design.Storage.Base),
		genesisDesign.Facts,
	)
	_ = g.SetLogging(log)

	if _, err := g.Generate(); err != nil {
		return ctx, e(err, "")
	}

	return ctx, nil
}

func PEncoder(ctx context.Context) (context.Context, error) {
	e := util.StringErrorFunc("failed to prepare encoders")

	encs := encoder.NewEncoders()
	jenc := jsonenc.NewEncoder()
	benc := bsonenc.NewEncoder()

	if err := encs.AddHinter(jenc); err != nil {
		return ctx, e(err, "")
	}
	if err := encs.AddHinter(benc); err != nil {
		return ctx, e(err, "")
	}

	ctx = context.WithValue(ctx, launch.EncodersContextKey, encs) //revive:disable-line:modifies-parameter
	ctx = context.WithValue(ctx, launch.EncoderContextKey, jenc)  //revive:disable-line:modifies-parameter
	ctx = context.WithValue(ctx, BEncoderContextKey, benc)        //revive:disable-line:modifies-parameter

	return ctx, nil
}

func PNetworkHandlers(pctx context.Context) (context.Context, error) {
	e := util.StringErrorFunc("failed to prepare network handlers")

	var log *logging.Logging
	var encs *encoder.Encoders
	var enc encoder.Encoder
	var design launch.NodeDesign
	var local base.LocalNode
	var params *isaac.LocalParams
	var db isaac.Database
	var pool *isaacdatabase.TempPool
	var proposalMaker *isaac.ProposalMaker
	var memberlist *quicmemberlist.Memberlist
	var syncSourcePool *isaac.SyncSourcePool
	var handlers *quicstream.PrefixHandler
	var nodeinfo *isaacnetwork.NodeInfoUpdater
	var svvotef isaac.SuffrageVoteFunc
	var cb *isaacnetwork.CallbackBroadcaster
	var ballotbox *isaacstates.Ballotbox
	var filternotifymsg launch.FilterMemberlistNotifyMsgFunc

	if err := util.LoadFromContextOK(pctx,
		launch.LoggingContextKey, &log,
		launch.EncodersContextKey, &encs,
		launch.EncoderContextKey, &enc,
		launch.DesignContextKey, &design,
		launch.LocalContextKey, &local,
		launch.LocalParamsContextKey, &params,
		launch.CenterDatabaseContextKey, &db,
		launch.PoolDatabaseContextKey, &pool,
		launch.ProposalMakerContextKey, &proposalMaker,
		launch.MemberlistContextKey, &memberlist,
		launch.SyncSourcePoolContextKey, &syncSourcePool,
		launch.QuicstreamHandlersContextKey, &handlers,
		launch.NodeInfoContextKey, &nodeinfo,
		launch.SuffrageVotingVoteFuncContextKey, &svvotef,
		launch.CallbackBroadcasterContextKey, &cb,
		launch.BallotboxContextKey, &ballotbox,
		launch.FilterMemberlistNotifyMsgFuncContextKey, &filternotifymsg,
	); err != nil {
		return pctx, e(err, "")
	}

	sendOperationFilterf, err := SendOperationFilterFunc(pctx)
	if err != nil {
		return pctx, e(err, "")
	}

	idletimeout := time.Second * 2 //nolint:gomnd //...
	lastBlockMapf := launch.QuicstreamHandlerLastBlockMapFunc(db)
	suffrageNodeConnInfof := launch.QuicstreamHandlerSuffrageNodeConnInfoFunc(db, memberlist)

	handlers.
		Add(isaacnetwork.HandlerPrefixOperation, isaacnetwork.QuicstreamHandlerOperation(encs, idletimeout, pool)).
		Add(isaacnetwork.HandlerPrefixSendOperation,
			isaacnetwork.QuicstreamHandlerSendOperation(
				encs, idletimeout, params, pool,
				db.ExistsInStateOperation,
				sendOperationFilterf,
				svvotef,
				func(id string, b []byte) error {
					return cb.Broadcast(id, b, nil)
				},
			),
		).
		Add(isaacnetwork.HandlerPrefixRequestProposal,
			isaacnetwork.QuicstreamHandlerRequestProposal(encs, idletimeout,
				local, pool, proposalMaker, db.LastBlockMap,
			),
		).
		Add(isaacnetwork.HandlerPrefixProposal,
			isaacnetwork.QuicstreamHandlerProposal(encs, idletimeout, pool),
		).
		Add(isaacnetwork.HandlerPrefixLastSuffrageProof,
			isaacnetwork.QuicstreamHandlerLastSuffrageProof(encs, idletimeout,
				func(last util.Hash) (hint.Hint, []byte, []byte, bool, error) {
					enchint, metabytes, body, found, lastheight, err := db.LastSuffrageProofBytes()

					switch {
					case err != nil:
						return enchint, nil, nil, false, err
					case !found:
						return enchint, nil, nil, false, storage.ErrNotFound.Errorf("last SuffrageProof not found")
					}

					switch h, err := isaacdatabase.ReadHashRecordMeta(metabytes); {
					case err != nil:
						return enchint, nil, nil, true, err
					case last != nil && last.Equal(h):
						nbody, _ := util.NewLengthedBytesSlice(0x01, [][]byte{lastheight.Bytes(), nil})

						return enchint, nil, nbody, false, nil
					default:
						nbody, _ := util.NewLengthedBytesSlice(0x01, [][]byte{lastheight.Bytes(), body})

						return enchint, metabytes, nbody, true, nil
					}
				},
			),
		).
		Add(isaacnetwork.HandlerPrefixSuffrageProof,
			isaacnetwork.QuicstreamHandlerSuffrageProof(encs, idletimeout, db.SuffrageProofBytes),
		).
		Add(isaacnetwork.HandlerPrefixLastBlockMap,
			isaacnetwork.QuicstreamHandlerLastBlockMap(encs, idletimeout, lastBlockMapf),
		).
		Add(isaacnetwork.HandlerPrefixBlockMap,
			isaacnetwork.QuicstreamHandlerBlockMap(encs, idletimeout, db.BlockMapBytes),
		).
		Add(isaacnetwork.HandlerPrefixBlockMapItem,
			isaacnetwork.QuicstreamHandlerBlockMapItem(encs, idletimeout, idletimeout*2, //nolint:gomnd //...
				func(height base.Height, item base.BlockMapItemType) (io.ReadCloser, bool, error) {
					e := util.StringErrorFunc("failed to get BlockMapItem")

					var menc encoder.Encoder

					switch m, found, err := db.BlockMap(height); {
					case err != nil:
						return nil, false, e(err, "")
					case !found:
						return nil, false, e(storage.ErrNotFound.Errorf("BlockMap not found"), "")
					default:
						menc = encs.Find(m.Encoder())
						if menc == nil {
							return nil, false, e(storage.ErrNotFound.Errorf("encoder of BlockMap not found"), "")
						}
					}

					reader, err := isaacblock.NewLocalFSReaderFromHeight(
						launch.LocalFSDataDirectory(design.Storage.Base), height, menc,
					)
					if err != nil {
						return nil, false, e(err, "")
					}
					defer func() {
						_ = reader.Close()
					}()

					return reader.Reader(item)
				},
			),
		).
		Add(isaacnetwork.HandlerPrefixNodeChallenge,
			isaacnetwork.QuicstreamHandlerNodeChallenge(encs, idletimeout, local, params),
		).
		Add(isaacnetwork.HandlerPrefixSuffrageNodeConnInfo,
			isaacnetwork.QuicstreamHandlerSuffrageNodeConnInfo(encs, idletimeout, suffrageNodeConnInfof),
		).
		Add(isaacnetwork.HandlerPrefixSyncSourceConnInfo,
			isaacnetwork.QuicstreamHandlerSyncSourceConnInfo(encs, idletimeout,
				func() ([]isaac.NodeConnInfo, error) {
					members := make([]isaac.NodeConnInfo, syncSourcePool.Len()*2)

					var i int
					syncSourcePool.Actives(func(nci isaac.NodeConnInfo) bool {
						members[i] = nci
						i++

						return true
					})

					return members[:i], nil
				},
			),
		).
		Add(isaacnetwork.HandlerPrefixState,
			isaacnetwork.QuicstreamHandlerState(encs, idletimeout, db.StateBytes),
		).
		Add(isaacnetwork.HandlerPrefixExistsInStateOperation,
			isaacnetwork.QuicstreamHandlerExistsInStateOperation(encs, idletimeout, db.ExistsInStateOperation),
		).
		Add(isaacnetwork.HandlerPrefixNodeInfo,
			isaacnetwork.QuicstreamHandlerNodeInfo(encs, idletimeout, launch.QuicstreamHandlerGetNodeInfoFunc(enc, nodeinfo)),
		).
		Add(isaacnetwork.HandlerPrefixCallbackMessage,
			isaacnetwork.QuicstreamHandlerCallbackMessage(encs, idletimeout, cb),
		).
		Add(isaacnetwork.HandlerPrefixSendBallots,
			isaacnetwork.QuicstreamHandlerSendBallots(
				encs, idletimeout, params,
				func(bl base.BallotSignFact) error {
					switch passed, err := filternotifymsg(bl); {
					case err != nil:
						log.Log().Trace().
							Str("module", "filter-notify-msg-send-ballots").
							Err(err).
							Interface("message", bl).
							Msg("filter error")

						fallthrough
					case !passed:
						log.Log().Trace().
							Str("module", "filter-notify-msg-send-ballots").
							Interface("message", bl).
							Msg("filtered")

						return nil
					}

					_, err := ballotbox.VoteSignFact(bl, params.Threshold())

					return err
				},
			),
		).
		Add(launch.HandlerPrefixPprof, launch.NetworkHandlerPprofFunc(encs))

	return pctx, nil
}
func SendOperationFilterFunc(ctx context.Context) (
	func(base.Operation) (bool, error),
	error,
) {
	var db isaac.Database
	var oprs *hint.CompatibleSet

	if err := util.LoadFromContextOK(ctx,
		launch.CenterDatabaseContextKey, &db,
		launch.OperationProcessorsMapContextKey, &oprs,
	); err != nil {
		return nil, err
	}

	operationfilterf := IsSupportedProposalOperationFactHintFunc()

	return func(op base.Operation) (bool, error) {
		switch hinter, ok := op.Fact().(hint.Hinter); {
		case !ok:
			return false, nil
		case !operationfilterf(hinter.Hint()):
			return false, errors.Errorf("Not supported operation")
		}

		var height base.Height

		switch m, found, err := db.LastBlockMap(); {
		case err != nil:
			return false, err
		case !found:
			return true, nil
		default:
			height = m.Manifest().Height()
		}

		f, closef, err := launch.OperationPreProcess(oprs, op, height)
		if err != nil {
			return false, err
		}

		defer func() {
			_ = closef()
		}()

		_, reason, err := f(context.Background(), db.State)
		if err != nil {
			return false, err
		}

		return reason == nil, reason
	}, nil
}

func IsSupportedProposalOperationFactHintFunc() func(hint.Hint) bool {
	return func(ht hint.Hint) bool {
		for _, hinter := range SupportedProposalOperationFactHinters {
			s := hinter.Hint
			if ht.Type() != s.Type() {
				continue
			}

			return ht.IsCompatible(s)
		}

		return false
	}
}
