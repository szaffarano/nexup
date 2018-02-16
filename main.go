package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Pubfile es la representación del archivo Pubfile
type Pubfile struct {
	System      string
	Application string
	Version     string
	Repository  string
}

// Credentials contiene datos de autenticación
type Credentials struct {
	Username string
	Password string
}

// Repository es el repositorio nexus
type Repository struct {
	url      string
	user     string
	password string
	client   *http.Client
	hash     map[string]*repositoryHash
}

// repositoryHash create hash and has the suffix for the file on repository
type repositoryHash struct {
	suffix string
	//TODO see if we need to a func or variable
	hash func() hash.Hash
}

func main() {
	pubfile := Pubfile{}
	credentials := Credentials{}

	usr, err := user.Current()
	die(err, 1)
	defaultCredentialsPath := fmt.Sprintf("%s%c.Pubfile.credentials", usr.HomeDir, os.PathSeparator)

	pubfilePath := flag.String("pubfile", "Pubfile", "Descriptor del artefacto")
	credentialsPath := flag.String("credentials", defaultCredentialsPath, "Credenciales para autenticarse")

	flag.Parse()

	pubfileData, err := ioutil.ReadFile(*pubfilePath)
	die(err, 2)

	credentialsData, err := ioutil.ReadFile(*credentialsPath)
	die(err, 3)

	err = yaml.Unmarshal([]byte(pubfileData), &pubfile)
	die(err, 4)

	err = yaml.Unmarshal([]byte(credentialsData), &credentials)
	die(err, 5)

	repo := NewRepository(pubfile, credentials)

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "No se indicaron archivos a publicar")
		os.Exit(-1)
	}

	for i := 0; i < len(flag.Args()); i++ {
		path := flag.Arg(i)

		data, err := ioutil.ReadFile(path)
		die(err, 6)

		ext := filepath.Ext(filepath.Base(path))
		name := strings.TrimSuffix(filepath.Base(path), ext)

		url := fmt.Sprintf("%s/%s/%s/%s/%s-%s%s",
			pubfile.Repository,
			pubfile.System,
			pubfile.Application,
			pubfile.Version,
			name,
			pubfile.Version,
			ext)

		fmt.Println(url)

		repo.upload(url, bytes.NewReader(data))
	}
}

//NewRepository create a Repository with default client HTTP.
func NewRepository(pubfile Pubfile, credentials Credentials) *Repository {
	const (
		nameMd5  = "md5"
		nameSha1 = "sha1"
	)

	shaOneAndMdFive := make(map[string]*repositoryHash)

	shaOneAndMdFive[nameMd5] = &repositoryHash{
		suffix: nameMd5,
		hash:   func() hash.Hash { return md5.New() },
	}

	shaOneAndMdFive[nameSha1] = &repositoryHash{
		suffix: nameSha1,
		hash:   func() hash.Hash { return sha1.New() },
	}

	return &Repository{
		url:      pubfile.Repository,
		user:     credentials.Username,
		password: credentials.Password,
		client:   &http.Client{},
		hash:     shaOneAndMdFive}

}

func die(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(code)
	}
}

func (n *Repository) upload(url string, data io.Reader) error {
	const (
		PUT         = "PUT"
		httpSuccess = 201
	)

	req, _ := http.NewRequest(PUT, url, data)
	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != httpSuccess {
		return fmt.Errorf(res.Status)
	}
	return nil
}

func (n *Repository) delete(url string) error {
	const httpSuccess = 204
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != httpSuccess {
		return fmt.Errorf(res.Status)
	}
	return nil
}
