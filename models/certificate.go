package models

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"strings"
	"time"
)

// TimeoutSeconds references the total Time Out duration for the Handshake
var TimeoutSeconds = 3

const defaultPort = "443"

// SplitHostPort returns a Host and Port seperately in a host-port URL
func SplitHostPort(hostport string) (string, string, error) {
	if !strings.Contains(hostport, ":") {
		return hostport, defaultPort, nil
	}

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", "", err
	}

	if port == "" {
		port = defaultPort
	}

	return host, port, nil
}

//Cert holds the certificate details
type Cert struct {
	DomainName string   `json:"domain_name"`
	IP         string   `json:"ip_address"`
	Issuer     string   `json:"issuer"`
	CommonName string   `json:"common_name"`
	SANs       []string `json:"sans"`
	NotBefore  string   `json:"not_before"`
	NotAfter   string   `json:"not_after"`
	Error      string   `json:"error"`
	certChain  []*x509.Certificate
}

var serverCert = func(host, port string) ([]*x509.Certificate, string, error) {
	d := &net.Dialer{
		Timeout: time.Duration(TimeoutSeconds) * time.Second,
	}
	conn, err := tls.DialWithDialer(d, "tcp", host+":"+port, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return []*x509.Certificate{&x509.Certificate{}}, "", err
	}
	defer conn.Close()

	addr := conn.RemoteAddr()
	ip, _, _ := net.SplitHostPort(addr.String())
	cert := conn.ConnectionState().PeerCertificates

	return cert, ip, nil
}

// GetCertificate returns the Certificate associated with a host-port
func GetCertificate(hostport string) *Cert {
	host, port, err := SplitHostPort(hostport)
	if err != nil {
		return &Cert{DomainName: host, Error: err.Error()}
	}
	certChain, ip, err := serverCert(host, port)
	if err != nil {
		return &Cert{DomainName: host, Error: err.Error()}
	}
	cert := certChain[0]

	var loc *time.Location
	loc = time.UTC // Setting UTC as Standard Time

	return &Cert{
		DomainName: host,
		IP:         ip,
		Issuer:     cert.Issuer.CommonName,
		CommonName: cert.Subject.CommonName,
		SANs:       cert.DNSNames, // Subject Alternative Name
		NotBefore:  cert.NotBefore.In(loc).String(),
		NotAfter:   cert.NotAfter.In(loc).String(),
		Error:      "",
		certChain:  certChain,
	}
}