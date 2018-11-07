package ztsclient

import (
	"crypto/tls"
	"github.com/aws/aws-sdk-go/aws/credentials"
	zts "github.com/yahoo/athenz/clients/go/zts"
	"io/ioutil"
	"log"
	"net/http"
)

type ZtsConfig struct {
	AthensDomain string
	AwsRole      string
	ZtsUrl       string
	KeyFile      string
	CertFile     string
}

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ZtsClient(ztsUrl, keyFile, certFile string) (*zts.ZTSClient, error) {
	keypem, err := ioutil.ReadFile(keyFile)
	errorCheck(err)

	certpem, err := ioutil.ReadFile(certFile)
	errorCheck(err)

	config, err := tlsConfiguration(keypem, certpem)
	errorCheck(err)
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	client := zts.NewClient(ztsUrl, tr)
	return &client, nil
}

func tlsConfiguration(keypem, certpem []byte) (*tls.Config, error) {
	config := &tls.Config{}
	if certpem != nil && keypem != nil {
		mycert, err := tls.X509KeyPair(certpem, keypem)
		errorCheck(err)

		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0] = mycert
	}
	return config, nil
}

func GetTemporaryCredentials(ztsCfg ZtsConfig) (*credentials.Credentials, error) {
	domainName := ztsCfg.AthensDomain
	roleName := ztsCfg.AwsRole
	ztsUrl := ztsCfg.ZtsUrl
	keyFile := ztsCfg.KeyFile
	certFile := ztsCfg.CertFile
	ztsClient, err := ZtsClient(ztsUrl, keyFile, certFile)
	errorCheck(err)

	awsCreds, err := ztsClient.GetAWSTemporaryCredentials(zts.DomainName(domainName), zts.CompoundName(roleName))
	errorCheck(err)

	creds := credentials.NewStaticCredentials(awsCreds.AccessKeyId, awsCreds.SecretAccessKey, awsCreds.SessionToken)
	return creds, nil
}
