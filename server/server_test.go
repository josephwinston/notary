package server

import (
	"net"
	"strings"
	"testing"
	"time"

	_ "github.com/docker/distribution/registry/auth/silly"
	"github.com/endophage/gotuf/signed"
	"golang.org/x/net/context"

	"github.com/docker/notary/config"
)

func TestRunBadCerts(t *testing.T) {
	err := Run(
		context.Background(),
		config.ServerConf{},
		signed.NewEd25519(),
	)
	if err == nil {
		t.Fatal("Passed empty certs, Run should have failed")
	}
}

func TestRunBadAddr(t *testing.T) {
	config := config.ServerConf{
		Addr:        "testAddr",
		TLSCertFile: "../fixtures/notary.pem",
		TLSKeyFile:  "../fixtures/notary.key",
		Auth: config.AuthConf{
			Method: "silly",
			Opts: map[string]interface{}{
				"realm":   "testrealm",
				"service": "testservice",
			},
		},
	}
	err := Run(context.Background(), config, signed.NewEd25519())
	if err == nil {
		t.Fatal("Passed bad addr, Run should have failed")
	}
}

func TestRunReservedPort(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	config := config.ServerConf{
		Addr:        "localhost:80",
		TLSCertFile: "../fixtures/notary.pem",
		TLSKeyFile:  "../fixtures/notary.key",
		Auth: config.AuthConf{
			Method: "silly",
			Opts: map[string]interface{}{
				"realm":   "testrealm",
				"service": "testservice",
			},
		},
	}

	err := Run(ctx, config, signed.NewEd25519())

	if _, ok := err.(*net.OpError); !ok {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "bind: permission denied") {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
}

func TestRunGoodCancel(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	config := config.ServerConf{
		Addr:        "localhost:8002",
		TLSCertFile: "../fixtures/notary.pem",
		TLSKeyFile:  "../fixtures/notary.key",
		Auth: config.AuthConf{
			Method: "silly",
			Opts: map[string]interface{}{
				"realm":   "testrealm",
				"service": "testservice",
			},
		},
	}

	go func() {
		time.Sleep(time.Second * 3)
		cancelFunc()
	}()

	err := Run(ctx, config, signed.NewEd25519())

	if _, ok := err.(*net.OpError); !ok {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "use of closed network connection") {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
}
