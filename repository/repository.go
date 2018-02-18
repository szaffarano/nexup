package repository

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// Truststores que se configuran en build time
	Truststores = ""
)

// Repository es el repositorio nexus
type Repository struct {
	URL      string
	user     string
	password string
	client   *http.Client
}

// Artifact representa un artefacto que se subirá al repositorio
type Artifact struct {
	System      string
	Application string
	Version     string
	File        string
}

// New create a Repository with default client HTTP.
func New(
	url string,
	username string,
	password string,
	trust string) (*Repository, error) {

	ca := x509.NewCertPool()

	if len(Truststores) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(Truststores)
		if err != nil {
			return nil, err
		}
		ok := ca.AppendCertsFromPEM(decoded)
		if !ok {
			return nil, errors.New("Error leyendo certificados embebidos")
		}
	}
	if len(trust) > 0 {
		ok := ca.AppendCertsFromPEM([]byte(trust))
		if !ok {
			return nil,
				errors.New("Error leyendo certificado desde configuración")
		}
	}

	var client = http.Client{}
	if len(ca.Subjects()) > 0 {
		tlsConf := &tls.Config{RootCAs: ca}
		tr := &http.Transport{TLSClientConfig: tlsConf}
		client = http.Client{Transport: tr}
	}

	return &Repository{
		URL:      url,
		user:     username,
		password: password,
		client:   &client}, nil

}

// Put Realiza un upload en el repositorio
func (n *Repository) Put(url string, data io.Reader) error {
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return err
	}

	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			fmt.Sprintf("Error al realizar el put, estado '%s'", res.Status))
	}
	return nil
}

// Delete borra un archivo en el repositorio
func (n *Repository) Delete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf(
			fmt.Sprintf("Error al realizar el delete, estado '%s'", res.Status))
	}
	return nil
}
