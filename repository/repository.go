package repository

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	var client = http.Client{}
	if trust != nil {
		ca := x509.NewCertPool()
		ok := ca.AppendCertsFromPEM([]byte(trust))
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
