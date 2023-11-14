package certificate_scan

import (
	"crypto/tls"
	"net"
	"sync"
	"time"
)

type SupportedSSLCipher struct {
	SSLVersion         string   `json:"ssl_version" bson:"ssl_version"`
	Prefernce          bool     `json:"prefernce" bson:"prefernce"`
	ChaCha20Preference bool     `json:"cha_cha_20_preference" bson:"cha_cha_20_preference"`
	Ciphers            []string `json:"ciphers" bson:"ciphers"`
}

type AsyncSSLCipherParams struct {
	SSLVersion string
	TLSCiphers []uint16
	MaxVersion uint16
	MinVersion uint16
}

var TLS12Cipher = []uint16{
	tls.TLS_RSA_WITH_RC4_128_SHA,
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}

var TLS13Cipher = []uint16{
	tls.TLS_AES_256_GCM_SHA384,
	tls.TLS_AES_128_GCM_SHA256,
	tls.TLS_CHACHA20_POLY1305_SHA256,
}

func GetCipherBySite(site string) []*SupportedSSLCipher {
	// 异步逻辑
	var wg sync.WaitGroup

	params := []AsyncSSLCipherParams{
		{SSLVersion: "SSL30", TLSCiphers: TLS12Cipher, MaxVersion: tls.VersionSSL30, MinVersion: tls.VersionSSL30},
		{SSLVersion: "TLS10", TLSCiphers: TLS12Cipher, MaxVersion: tls.VersionTLS10, MinVersion: tls.VersionTLS10},
		{SSLVersion: "TLS11", TLSCiphers: TLS12Cipher, MaxVersion: tls.VersionTLS11, MinVersion: tls.VersionTLS11},
		{SSLVersion: "TLS12", TLSCiphers: TLS12Cipher, MaxVersion: tls.VersionTLS12, MinVersion: tls.VersionTLS12},
		{SSLVersion: "TLS13", TLSCiphers: TLS13Cipher, MaxVersion: tls.VersionTLS13, MinVersion: tls.VersionTLS13},
	}

	retChan := make(chan *SupportedSSLCipher, len(params))

	for _, param := range params {
		wg.Add(1)
		go asyncGetSupportedSSLCipher(site, param, retChan, &wg)
	}
	wg.Wait()
	close(retChan)
	result := make([]*SupportedSSLCipher, 0)
	for val := range retChan {
		result = append(result, val)
	}
	return result
}

func asyncGetSupportedSSLCipher(site string, SSLCipherParams AsyncSSLCipherParams, retChan chan *SupportedSSLCipher, wg *sync.WaitGroup) {
	sslCipher := asyncSupportedSSLCiphers(site, SSLCipherParams)
	retChan <- &sslCipher
	defer wg.Done()
}

func asyncSupportedSSLCiphers(site string, SSLCipherParams AsyncSSLCipherParams) (result SupportedSSLCipher) {
	result.SSLVersion = SSLCipherParams.SSLVersion
	result.Prefernce = false
	result.ChaCha20Preference = false

	var wg sync.WaitGroup

	retChan := make(chan uint16, len(SSLCipherParams.TLSCiphers))

	for _, cipher := range SSLCipherParams.TLSCiphers {
		conf := &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites:       []uint16{cipher},
			MinVersion:         SSLCipherParams.MinVersion,
			MaxVersion:         SSLCipherParams.MaxVersion,
		}
		wg.Add(1)
		go asyncTLSDial("tcp", site, conf, cipher, retChan, &wg)
	}
	wg.Wait()
	close(retChan)
	for val := range retChan {
		result.Ciphers = append(result.Ciphers, tls.CipherSuiteName(val))
	}
	return
}

func asyncTLSDial(network, site string, config *tls.Config, cipher uint16, retChan chan uint16, wg *sync.WaitGroup) {
	//conn, err := tls.Dial(network, site, config)
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, network, site, config)
	if err == nil {
		if conn.ConnectionState().CipherSuite == cipher {
			retChan <- cipher
		}
		conn.Close()
	}
	defer wg.Done()
}
