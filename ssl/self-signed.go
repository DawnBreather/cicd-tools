package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type Ssl struct{
	organization string
	country string
	province string
	locality string
	streetAddress string
	postalCode string

	ca *x509.Certificate
	caPrivKey *rsa.PrivateKey
	caPrivKeyPEM []byte
	caPubKeyPEM []byte
}

func (s *Ssl) SetOrganization(organization string) *Ssl {
	s.organization = organization
	return s
}
func (s *Ssl) SetCountry(country string) *Ssl {
	s.country = country
	return s
}
func (s *Ssl) SetProvince(province string) *Ssl {
	s.province = province
	return s
}
func (s *Ssl) SetLocality(locality string) *Ssl {
	s.locality = locality
	return s
}
func (s *Ssl) SetStreetAddress(streetAddress string) *Ssl {
	s.streetAddress = streetAddress
	return s
}
func (s *Ssl) SetPostalCode(postalCode string) *Ssl {
	s.postalCode = postalCode
	return s
}

func (s *Ssl) InitializeCertificateAuthority() *Ssl {
	s.ca = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			CommonName: "self-signed",
			Country:            []string{s.country},
			Organization:       []string{s.organization},
			Locality:           []string{s.locality},
			Province:           []string{s.province},
			StreetAddress:      []string{s.streetAddress},
			PostalCode:         []string{s.postalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}


	s.caPrivKey, _ = rsa.GenerateKey(rand.Reader, 4096)
	//if err != nil {
	//	return err
	//}

	caBytes, _ := x509.CreateCertificate(rand.Reader, s.ca, s.ca, &s.caPrivKey.PublicKey, s.caPrivKey)
	//if err != nil {
	//	return err
	//}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.caPrivKey),
	})

	s.caPubKeyPEM = caPEM.Bytes()
	s.caPrivKeyPEM = caPrivKeyPEM.Bytes()

	return s

}


func (s *Ssl) GenerateSignedCertificate(DnsNames []string) (pubPEM, privPEM []byte, certKeypair tls.Certificate){
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName: DnsNames[0],
			Organization:  []string{s.organization},
			Country:       []string{s.country},
			Province:      []string{s.province},
			Locality:      []string{s.locality},
			StreetAddress: []string{s.streetAddress},
			PostalCode:    []string{s.postalCode},
		},
		DNSNames: DnsNames,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	certBytes, _ := x509.CreateCertificate(rand.Reader, cert, s.ca, &certPrivKey.PublicKey, s.caPrivKey)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	certKeypair, _ = tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())

	return certPEM.Bytes(), certPrivKeyPEM.Bytes(), certKeypair
}

func (s *Ssl) GetClientTlsConfig() *tls.Config{
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(s.caPubKeyPEM)
	return &tls.Config{
		RootCAs: certpool,
	}
}