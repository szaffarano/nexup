package conf

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
