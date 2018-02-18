package repository

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	// Truststores que se configuran en build time
	Truststores = ""
)

// Repository es el repositorio nexus
type Repository struct {
	url      string
	user     string
	password string
	client   *http.Client
}

//NewRepository create a Repository with default client HTTP.
func NewRepository(url string, username string, password string, trust []byte) (*Repository, error) {
	var cacerts = ""

	if len(Truststores) > 0 {
		Truststores = strings.Replace(Truststores, " CER", "_CER", -1)
		Truststores = strings.Replace(Truststores, " ", "\n", -1)
		Truststores = strings.Replace(Truststores, "_CER", " CER", -1)
		cacerts = Truststores
	}
	if trust != nil && len(trust) > 0 {
		cacerts += string(trust)
	}

	var client = http.Client{}
	if len(cacerts) > 0 {
		ca := x509.NewCertPool()
		ok := ca.AppendCertsFromPEM([]byte(cacerts))
		if !ok {
			return nil, errors.New("Error leyendo certificado")
		}
		tlsConf := &tls.Config{RootCAs: ca}
		tr := &http.Transport{TLSClientConfig: tlsConf}
		client = http.Client{Transport: tr}
	}

	return &Repository{
		url:      url,
		user:     username,
		password: password,
		client:   &client}, nil

}

// Put Realiza un upload en el repositorio
func (n *Repository) Put(url string, data io.Reader) error {
	req, _ := http.NewRequest("PUT", url, data)
	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(res.Status)
	}
	return nil
}

// Delete borra un archivo en el repositorio
func (n *Repository) Delete(url string) error {
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf(res.Status)
	}
	return nil
}
