package factory_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/ElrondNetwork/elrond-go/factory/mock"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkComponentsFactory_NilStatusHandlerShouldErr(t *testing.T) {
	t.Parallel()

	args := getNetworkArgs()
	args.StatusHandler = nil
	ncf, err := factory.NewNetworkComponentsFactory(args)
	require.Nil(t, ncf)
	require.Equal(t, factory.ErrNilStatusHandler, err)
}

func TestNewNetworkComponentsFactory_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	args := getNetworkArgs()
	args.Marshalizer = nil
	ncf, err := factory.NewNetworkComponentsFactory(args)
	require.Nil(t, ncf)
	require.True(t, errors.Is(err, factory.ErrNilMarshalizer))
}

func TestNewNetworkComponentsFactory_OkValsShouldWork(t *testing.T) {
	t.Parallel()
	args := getNetworkArgs()
	ncf, err := factory.NewNetworkComponentsFactory(args)
	require.NoError(t, err)
	require.NotNil(t, ncf)
}

func TestNetworkComponentsFactory_Create_ShouldErrDueToBadConfig(t *testing.T) {
	//TODO remove skip when external library is concurrent safe
	if testing.Short() {
		t.Skip("this test fails with race detector on because of the github.com/koron/go-ssdp lib")
	}

	args := getNetworkArgs()
	args.MainConfig = config.Config{}
	args.P2pConfig = config.P2PConfig{}

	ncf, _ := factory.NewNetworkComponentsFactory(args)

	nc, err := ncf.Create()
	require.Error(t, err)
	require.Nil(t, nc)
}

func TestNetworkComponentsFactory_Create_ShouldWork(t *testing.T) {
	//TODO remove skip when external library is concurrent safe
	if testing.Short() {
		t.Skip("this test fails with race detector on because of the github.com/koron/go-ssdp lib")
	}

	args := getNetworkArgs()
	ncf, _ := factory.NewNetworkComponentsFactory(args)
	ncf.SetListenAddress(libp2p.ListenLocalhostAddrWithIp4AndTcp)

	nc, err := ncf.Create()
	require.NoError(t, err)
	require.NotNil(t, nc)
}

func getNetworkArgs() factory.NetworkComponentsFactoryArgs {
	p2pConfig := config.P2PConfig{
		Node: config.NodeConfig{
			Port: "0",
			Seed: "seed",
		},
		KadDhtPeerDiscovery: config.KadDhtPeerDiscoveryConfig{
			Enabled:                          false,
			RefreshIntervalInSec:             10,
			RandezVous:                       "erd/kad/1.0.0",
			InitialPeerList:                  []string{"peer0", "peer1"},
			BucketSize:                       10,
			RoutingTableRefreshIntervalInSec: 5,
		},
		Sharding: config.ShardingConfig{
			TargetPeerCount:         10,
			MaxIntraShardValidators: 10,
			MaxCrossShardValidators: 10,
			MaxIntraShardObservers:  10,
			MaxCrossShardObservers:  10,
			Type:                    "NilListSharder",
		},
	}

	mainConfig := config.Config{
		Debug: config.DebugConfig{
			Antiflood: config.AntifloodDebugConfig{
				Enabled:                    true,
				CacheSize:                  100,
				IntervalAutoPrintInSeconds: 1,
			},
		},
	}

	appStatusHandler := &mock.AppStatusHandlerMock{}

	return factory.NetworkComponentsFactoryArgs{
		P2pConfig:     p2pConfig,
		MainConfig:    mainConfig,
		StatusHandler: appStatusHandler,
		Marshalizer:   &mock.MarshalizerMock{},
	}
}
