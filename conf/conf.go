package conf

import "errors"

// Nexupfile es la representación del archivo Pubfile
type Nexupfile struct {
	System      string
	Application string
	Version     string
	Repository  string
	Truststores string
}

// NexupfileCredentials contiene datos de autenticación
type NexupfileCredentials struct {
	Username string
	Password string
}

// Validate verifica que se hayan recibido todos los parámetros necesarios
func (n *Nexupfile) Validate() error {
	if len(n.Application) == 0 {
		return errors.New("Debe especificar un nombre de aplicación")
	} else if len(n.System) == 0 {
		return errors.New("Debe especificar un nombre de sistema")
	} else if len(n.Version) == 0 {
		return errors.New("Debe especificar una versión")
	} else if len(n.Repository) == 0 {
		return errors.New("Debe especificar un repositorio")
	}
	return nil
}
