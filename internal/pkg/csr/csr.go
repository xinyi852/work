package csr

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"net"
	"plesk/internal/pkg/common"
	"strings"

	"racent.com/pkg/helpers"
)

var oidSubjectEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func GenerateCSR(params *common.GenerateCSRRequest) (*common.GenerateCSRResult, error) {

	privateKey, sigAlgo, err := helpers.GeneratePrivateKey(params.KeyAlgorithm, params.KeySize, params.SignatureHashAlgorithm)
	if err != nil {
		return nil, err
	}
	key, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("x509: certificate private key does not implement crypto.Signer")
	}

	subject := pkix.Name{
		CommonName:         params.CommonName,
		Country:            []string{params.Country},
		Province:           []string{params.Province},
		Locality:           []string{params.Locality},
		Organization:       []string{params.Organization},
		OrganizationalUnit: []string{params.OrganizationalUnit},
		StreetAddress:      []string{params.Address},
		PostalCode:         []string{params.Postcode},
	}
	if !helpers.Empty(params.EmailAddress) {
		subject.ExtraNames = []pkix.AttributeTypeAndValue{{Type: oidSubjectEmailAddress, Value: params.EmailAddress}}
	}
	rawSub := subject.ToRDNSequence()
	asn1Sub, err := asn1.Marshal(rawSub)
	if err != nil {
		return nil, err
	}
	dnsNames := []string{""}
	if params.DNSNames != "" {
		dnsNames = strings.Split(params.DNSNames, ",")
	}
	email := []string{""}
	if !helpers.Empty(params.Email) {
		email = strings.Split(params.Email, ",")
	}

	iPAddresses := make([]net.IP, 0)
	if !helpers.Empty(params.IpAddress) {
		for _, address := range strings.Split(params.IpAddress, ",") {
			isIP := net.ParseIP(address)
			if isIP != nil {
				iPAddresses = append(iPAddresses, isIP)
			}
		}
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Sub,
		EmailAddresses:     email,
		DNSNames:           dnsNames,
		SignatureAlgorithm: sigAlgo,
		IPAddresses:        iPAddresses,
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, err
	}
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	privatePem := helpers.GenBytesByPrivateKey(privateKey, params.KeyAlgorithm)
	result := &common.GenerateCSRResult{}
	result.CSRPem = string(csrPem)
	result.PrivatePem = string(privatePem)
	return result, nil
}
