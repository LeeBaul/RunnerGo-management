package dal

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/go-omnibus/proof"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	services "kp-management/api"
)

var (
	//grpcClient services.KpControllerClient
	conn *grpc.ClientConn
)

func MustInitGRPC() {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		panic(errors.Wrap(err, "cannot load root CA certs"))
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})

	conn, err = grpc.Dial("kpcontroller.apipost.cn:443", grpc.WithTransportCredentials(creds))

	//var err error
	//conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		proof.Errorf("grpc dial err", err)
	}

	//
	//grpcClient = services.NewKpControllerClient(conn)
}

func ClientGRPC() services.KpControllerClient {
	return services.NewKpControllerClient(conn)
}
