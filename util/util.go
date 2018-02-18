package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/szaffarano/nexup/repository"

	"golang.org/x/crypto/ssh/terminal"
)

// GetCredentials solicita al usuario las credenciales de autenticación
func GetCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingresar nombre de usuario: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", errors.New("Error obteniendo el nombre de usuario")
	}

	fmt.Print("Ingresar contraseña: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", errors.New("Error obteniendo la contraseña")
	}
	password := string(bytePassword)

	fmt.Println("")

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// GetMavenURL arma la url de un archivo compatible con repositorios maven
func GetMavenURL(system, application, version string) string {
	return "hola"
}

// GetRawURL arma la url de un archivo para un repositorio raw
func GetRawURL(repo repository.Repository, artifact repository.Artifact) string {
	ext := filepath.Ext(artifact.File)
	name := strings.TrimSuffix(artifact.File, ext)

	return fmt.Sprintf("%s/%s/%s/%s/%s-%s%s",
		repo.URL,
		artifact.System,
		artifact.Application,
		artifact.Version,
		name,
		artifact.Version,
		ext)
}
