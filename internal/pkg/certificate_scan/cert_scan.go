package certificate_scan

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net"
	"time"
)

var timeout = 1 * time.Second

func CheckTcpHandshake(site, port string) error {
	ipConn, err := tcpHandshake(site, port)
	if err != nil {
		return err
	}
	defer ipConn.Close()
	return nil
}

func tcpHandshake(site, port string) (net.Conn, error) {
	certSite := site + ":" + port
	return net.DialTimeout("tcp", certSite, timeout)
}

type GetCertResponse struct {
	OCSPStapling int8     `json:"ocsp_stapling"`
	RawCerts     []string `json:"raw_certs"`
}

// GetCertBySite 获取证书信息
func GetCertBySite(site, port string) (*GetCertResponse, error) {
	ipConn, err := tcpHandshake(site, port)
	if err != nil {
		return nil, err
	}
	defer ipConn.Close()

	var conf tls.Config
	if net.ParseIP(site) == nil {
		conf.InsecureSkipVerify = false
		conf.ServerName = site
	} else {
		conf.InsecureSkipVerify = true
	}
	conn := tls.Client(ipConn, &conf)
	err = conn.Handshake()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connState := conn.ConnectionState()
	var certs []*x509.Certificate

	if len(connState.VerifiedChains) > 0 {
		certs = connState.VerifiedChains[0]
	} else {
		certs = connState.PeerCertificates
	}

	resp := &GetCertResponse{}
	resp.RawCerts = make([]string, len(certs))
	if len(connState.OCSPResponse) > 0 {
		resp.OCSPStapling = 1
	} else {
		resp.OCSPStapling = 0
	}

	for i, cert := range certs {
		resp.RawCerts[i] = base64.StdEncoding.EncodeToString(cert.Raw)
	}
	return resp, nil
}

func FormatCertInfo(certs []*x509.Certificate) *Cert {
	respCert := &Cert{}
	respCert.CertChains = make([]*CertChain, len(certs))
	for i, cert := range certs {
		if i == 0 {
			respCert.SerialNumber = cert.SerialNumber.String()
			respCert.SubjectCN = cert.Subject.CommonName
			respCert.SubjectSerialNumber = cert.Subject.SerialNumber
			respCert.SubjectCountry = firstElement(cert.Subject.Country)
			respCert.SubjectProvince = firstElement(cert.Subject.Province)
			respCert.SubjectLocality = firstElement(cert.Subject.Locality)
			respCert.SubjectAddress = firstElement(cert.Subject.StreetAddress)
			respCert.SubjectPostal = firstElement(cert.Subject.PostalCode)
			respCert.SubjectOrg = firstElement(cert.Subject.Organization)
			respCert.SubjectOrgUnit = firstElement(cert.Subject.OrganizationalUnit)
			respCert.IssuerCN = cert.Issuer.CommonName
			respCert.IssuerOrg = firstElement(cert.Issuer.Organization)
			respCert.IssuerCountry = firstElement(cert.Issuer.Country)
			respCert.setCertType()
			respCert.setCertBrand()
			respCert.SignatureAlgorithm = cert.SignatureAlgorithm.String()
			respCert.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
			respCert.PublicKeyLength = getCertPublicKeyByType(cert.PublicKey)
			respCert.OCSPServer = firstElement(cert.OCSPServer)
			respCert.OCSPStatus = "Normal"
			respCert.OCSPMustStapling = false
			respCert.NotBefore = cert.NotBefore
			respCert.NotAfter = cert.NotAfter
			respCert.DNSNames = cert.DNSNames
			respCert.SNI = true
			respCert.DNSCAA = ""
			respCert.CTLog = true
			respCert.setPaid()
		}
		tmp := &CertChain{
			cert.Subject.CommonName,
			cert.NotBefore,
			cert.NotAfter,
			cert.PublicKeyAlgorithm.String(),
			getCertPublicKeyByType(cert.PublicKey),
			cert.Issuer.CommonName,
			cert.SignatureAlgorithm.String(),
		}
		respCert.CertChains[i] = tmp
	}
	return respCert
}

func firstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
