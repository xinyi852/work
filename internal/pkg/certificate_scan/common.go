package certificate_scan

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"golang.org/x/crypto/ed25519"
	"strconv"
	"strings"
	"time"
)

const (
	OLD_TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256   = uint16(0xcc13)
	OLD_TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256 = uint16(0xcc14)

	DISABLED_TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384 = uint16(0xc024)
	DISABLED_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384   = uint16(0xc028)
	DISABLED_TLS_RSA_WITH_AES_256_CBC_SHA256         = uint16(0x003d)

	FAKE_OLD_TLS_DHE_RSA_WITH_CHACHA20_POLY1305_SHA256 = uint16(0xcc15) // we can try to craft these ciphersuites
	FAKE_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256           = uint16(0x009e) // from existing pieces, if needed

	FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA = uint16(0x0033)
	TLS_RSA_WITH_RC4_128_MD5              = uint16(0x4)
	TLS_DHE_RSA_WITH_AES_128_CBC_SHA      = uint16(0x33)
	TLS_DHE_DSS_WITH_AES_128_CBC_SHA      = uint16(0x32)
	TLS_DHE_RSA_WITH_3DES_EDE_CBC_SHA     = uint16(0x16)
	TLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA     = uint16(0x13)
	TLS_RSA_WITH_DES_CBC_SHA              = uint16(0x9)
	TLS_DHE_RSA_WITH_DES_CBC_SHA          = uint16(0x15)
	TLS_DHE_DSS_WITH_DES_CBC_SHA          = uint16(0x12)

	TLS_RSA_EXPORT_WITH_RC4_40_MD5    = uint16(0x3)
	TLS_RSA_EXPORT_WITH_DES40_CBC_SHA = uint16(0x8)

	TLS_DHE_RSA_EXPORT_WITH_DES40_CBC_SHA = uint16(0x14)
	TLS_DHE_DSS_EXPORT_WITH_DES40_CBC_SHA = uint16(0x11)

	FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA    = uint16(0x0039)
	FAKE_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384 = uint16(0x009f)
	FAKE_TLS_RSA_WITH_RC4_128_MD5            = uint16(0x0004)
	FAKE_TLS_EMPTY_RENEGOTIATION_INFO_SCSV   = uint16(0x00ff)
)

var CipherSuite = map[string]uint16{
	"c02c": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"c02b": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"c00a": tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"c023": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"c009": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"c030": tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"c02f": tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"c027": tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"c013": tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
}

type PCI_DSSAndATSReturn struct {
	PCI_DSS bool `json:"pci_dss"` //pci dss是否可信
	ATS     *ATS `json:"ats"`
}

type ATS struct {
	DomainMatch   bool   `json:"domain_match"`   //域名匹配
	Expire        bool   `json:"expire"`         //是否过期
	Credible      bool   `json:"credible"`       //是否可信
	Cipher        bool   `json:"cipher"`         //是否用了规定的加密套件
	HashAlgorithm bool   `json:"hash_algorithm"` //叶服务器证书哈希算法是否大于SHA-256或更高
	KeyLength     bool   `json:"key_length"`     //RSA，密钥长度最少2048位 	ECC，密钥长度最少256位
	NotBefore     string `json:"not_before"`
	NotAfter      string `json:"not_after"`
}

type Cert struct {
	SerialNumber        string       `json:"serial_number" bson:"serial_number"`
	SubjectCountry      string       `json:"subject_country" bson:"subject_country"`
	SubjectProvince     string       `json:"subject_province" bson:"subject_province"`
	SubjectLocality     string       `json:"subject_locality" bson:"subject_locality"`
	SubjectAddress      string       `json:"subject_address" bson:"subject_address"`
	SubjectPostal       string       `json:"subject_postal" bson:"subject_postal"`
	SubjectCN           string       `json:"subject_cn" bson:"subject_cn"`
	SubjectOrg          string       `json:"subject_org" bson:"subject_org"`
	SubjectOrgUnit      string       `json:"subject_org_unit" bson:"subject_org_unit"`
	SubjectSerialNumber string       `json:"subject_serial_number" bson:"subject_serial_number"`
	IssuerCN            string       `json:"issuer_cn" bson:"issuer_cn"`
	IssuerOrg           string       `json:"issuer_org" bson:"issuer_org"`
	IssuerCountry       string       `json:"issuer_country" bson:"issuer_country"`
	SignatureAlgorithm  string       `json:"signature_algorithm" bson:"signature_algorithm"`
	PublicKeyAlgorithm  string       `json:"public_key_algorithm" bson:"public_key_algorithm"`
	PublicKeyLength     string       `json:"public_key_length" bson:"public_key_length"`
	OCSPServer          string       `json:"ocsp_server" bson:"ocsp_server"`
	CRLServer           string       `json:"crl_server" bson:"crl_server"`
	OCSPStatus          string       `json:"ocsp_status" bson:"ocsp_status"`               //Normal, Expired, Revoked
	OCSPStapling        bool         `json:"ocsp_stapling" bson:"ocsp_stapling"`           //1
	OCSPMustStapling    bool         `json:"ocsp_must_stapling" bson:"ocsp_must_stapling"` //1
	NotBefore           time.Time    `json:"not_before" bson:"not_before"`
	NotAfter            time.Time    `json:"not_after" bson:"not_after"`
	DNSNames            []string     `json:"dns_names" bson:"dns_names"` //1
	DNSCAA              string       `json:"dnscaa" bson:"dnscaa"`       //1
	CertType            string       `json:"cert_type" bson:"cert_type"`
	CertBrand           string       `json:"cert_brand" bson:"cert_brand"`
	SNI                 bool         `json:"sni" bson:"sni"`                 //1
	CTLog               bool         `json:"ct_log" bson:"ct_log"`           //1
	CertChains          []*CertChain `json:"cert_chains" bson:"cert_chains"` //1
	Paid                bool         `json:"paid" bson:"paid"`
}

func (c *Cert) setCertType() {
	var certType string
	if c.SubjectOrg == "" {
		certType = "DV"
	} else if strings.Index(c.IssuerCN, "Extended Validation") >= 0 {
		certType = "EV"
	} else if strings.Index(c.IssuerCN, "EV") >= 0 {
		certType = "EV"
	} else {
		certType = "OV"
	}
	c.CertType = certType
}

func (c *Cert) setPaid() {
	//paid or free
	validDays := int64(c.NotAfter.Sub(c.NotBefore).Hours())
	if validDays < 93*24 || c.IssuerOrg == "Let's Encrypt" || strings.Contains(c.IssuerCN, "Encryption Everywhere") {
		c.Paid = false
		return
	}
	if strings.Contains(c.IssuerCN, "TrustAsia") && c.SubjectOrg == "" {
		dnsNamesLen := len(c.DNSNames)
		if dnsNamesLen > 2 {
			c.Paid = true
			return
		}

		if dnsNamesLen == 2 {
			firstName := strings.ToLower(c.DNSNames[0])
			lastName := strings.ToLower(c.DNSNames[1])
			if strings.Contains(firstName, "www") {
				c.Paid = (firstName != "www"+lastName)
				return
			}
			if strings.Contains(lastName, "www") {
				c.Paid = (lastName != "www"+firstName)
				return
			}
			c.Paid = true
			return
		}

		c.Paid = strings.Contains(c.SubjectCN, "*")
		return
	}
	c.Paid = true
}

func (c *Cert) setCertBrand() {
	switch c.IssuerOrg {
	case "ACCV":
		c.CertBrand = "ACCV"
		return
	case "Amazon":
		c.CertBrand = "Amazon"
		return
	case "Chunghwa Telecom Co., Ltd.":
		c.CertBrand = "Chunghwa"
		return
	case "Corporation Service Company":
		c.CertBrand = "Trusted Secure"
		return
	case "sslTrus (上海锐成信息科技有限公司)":
		c.CertBrand = "sslTrus"
		return
	}

	if (strings.Index(c.IssuerCN, "Secure Site Pro") > -1) || (strings.Index(c.IssuerCN, "Secure Site Pro") > -1) {
		c.CertBrand = "Digicert"
		return
	}

	nIndex := strings.Index(c.IssuerCN, " ")
	var brand string
	if nIndex != -1 {
		brand = strings.TrimSpace(c.IssuerCN[0:nIndex])
	} else {
		brand = c.IssuerCN
	}

	if brand == "Encryption" {
		c.CertBrand = "Encryption Everywhere"
		return
	}
	if brand[len(brand)-1:] == "," {
		c.CertBrand = brand[0 : len(brand)-2]
	} else {
		c.CertBrand = brand
	}
}

type CertChain struct {
	SubjectCN          string    `json:"subject_cn" bson:"subject_cn"`
	NotBefore          time.Time `json:"not_before" bson:"not_before"`
	NotAfter           time.Time `json:"not_after" bson:"not_after"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm" bson:"public_key_algorithm"`
	PublicKeyLength    string    `json:"public_key_length" bson:"public_key_length"`
	IssuerCN           string    `json:"issuer_cn" bson:"issuer_cn"`
	SignatureAlgorithm string    `json:"signature_algorithm" bson:"signature_algorithm"`
}

func getCertPublicKeyByType(publicKey any) string {
	var keyLength string
	switch publicKey.(type) {
	case *rsa.PublicKey:
		keyLength = strconv.Itoa(publicKey.(*rsa.PublicKey).N.BitLen())
	case *ecdsa.PublicKey:
		keyLength = strconv.Itoa(publicKey.(*ecdsa.PublicKey).Curve.Params().BitSize)
	case ed25519.PublicKey:
		keyLength = strconv.Itoa(ed25519.PublicKeySize)
	default:
		keyLength = ""
	}
	return keyLength
}

func formatTLSVersionName(version uint16) string {
	switch version {
	case tls.VersionSSL30:
		return "SSL30"
	case tls.VersionTLS10:
		return "TLS10"
	case tls.VersionTLS11:
		return "TLS11"
	case tls.VersionTLS12:
		return "TLS12"
	case tls.VersionTLS13:
		return "TLS13"
	default:
		return "Unknown"
	}
}
