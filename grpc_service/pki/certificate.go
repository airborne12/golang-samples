package pki

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/airborne12/golang-samples/grpc_service/config"
	"github.com/airborne12/golang-samples/grpc_service/log"
)

var (
	gCertKey *CertKey
	glog     = log.Get()
	gconfig  = config.Get()
)

func init() {
	gCertKey = NewCertKey()
}

//CertKey store key pair and cert pool
type CertKey struct {
	KeyPair  *tls.Certificate
	CertPool *x509.CertPool
}

// NewCertKey creates a certicicate pair and pool and package them into CertKey struct
func NewCertKey() *CertKey {
	if gconfig.Insecure {
		return nil
	}
	serverCrt, err := ioutil.ReadFile(gconfig.ServerCert)
	if err != nil {
		glog.Fatal(err)
	}
	serverKey, err := ioutil.ReadFile(gconfig.ServerKey)
	if err != nil {
		glog.Fatal(err)
	}

	pair, err := tls.X509KeyPair(serverCrt, serverKey)
	if err != nil {
		glog.Fatal(err)
	}
	keyPair := &pair
	var clientCertPool *x509.CertPool
	if gconfig.ClientCA != "" {
		caCert, err := ioutil.ReadFile(gconfig.ClientCA)
		if err != nil {
			glog.Fatalf("failed to load client ca: %v", err)
		}
		clientCertPool = x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			glog.Fatal("failed to append client ca certs")
		}
	}

	return &CertKey{
		KeyPair:  keyPair,
		CertPool: clientCertPool,
	}
}

//Get returns a CertKey object
func Get() *CertKey {
	return gCertKey
}
