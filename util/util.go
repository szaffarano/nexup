package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

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
