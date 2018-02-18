package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/szaffarano/nexup/conf"
)

var (
	verbose   bool
	cfgFile   string
	credFile  string
	trustFile string

	nexupfile            conf.Nexupfile
	nexupfileCredentials conf.NexupfileCredentials

	rootCmd = &cobra.Command{
		Use:   "nexup",
		Short: "Nexup es un comando para subir contenido a repositorios",
		Long: `Con nexup vas a poder interactuar con repositorios para subir 
		       contenido de una forma simple y amigable`,
	}
)

// Execute es el punto de entrada de la cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initconfig)

	// inicializa flags globales
	rootCmd.
		PersistentFlags().
		BoolVarP(
			&verbose,
			"verbose",
			"v",
			false,
			"Imprime información extra")

	rootCmd.
		PersistentFlags().
		StringVarP(
			&cfgFile,
			"nexupfile",
			"n",
			"",
			`Archivo de configuración con la información del repositorio
		 	 y de lo que se quiere subir, por defecto se usa ./Nexupfile

			 Ejemplo

			 system: nombre-del-sistema
			 application: nombre-de-la-aplicacion
			 version: 1.2.3
			 repository: http://host:puerto/repository/etc/etc
			 truststores: |
			  .......
		 	`)

	rootCmd.
		PersistentFlags().
		StringVarP(
			&credFile,
			"credentials",
			"c",
			"",
			`Archivo de credenciales para autenticarse en el repositorio, por 
			 defecto se usa $HOME/.nexup-credentials

			 Ejemplo

			 username: NombreDeUsuario
			 password: Contraseña
		 	`)

	rootCmd.
		PersistentFlags().
		StringVarP(
			&trustFile,
			"truststores",
			"t",
			"",
			"Archivo con certificados de las CAs en las cuales confiar")
}

// initconfig lee la configuración
func initconfig() {
	config := viper.New()
	credentials := viper.New()

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initViper(config, map[string]string{
		"env-prefix":  "NEXUP",
		"file-name":   cfgFile,
		"config-path": ".",
		"config-name": "Nexupfile",
		"config-type": "yaml"})

	initViper(credentials, map[string]string{
		"env-prefix":  "NEXUP",
		"file-name":   credFile,
		"config-path": usr.HomeDir,
		"config-name": ".nexup-credentials",
		"config-type": "yaml"})

	credentials.Unmarshal(&nexupfileCredentials)
	config.Unmarshal(&nexupfile)
}

func initViper(instance *viper.Viper, defaults map[string]string) {
	var f = defaults["file-name"]

	if f != "" {
		instance.SetConfigFile(f)
	} else {
		// Por defecto busca el archivo  en el directorio actual
		// No uso AddconfigPath y SetconfigName porque
		// no funcionan bien sin extensión para autodescubrir
		// el tipo de archivo.
		// @TODO ver si es un bug de viper
		instance.SetConfigFile(path.Join(
			defaults["config-path"],
			defaults["config-name"]))
	}

	instance.AutomaticEnv()
	instance.SetEnvPrefix(defaults["env-prefix"])
	instance.SetConfigType(defaults["config-type"])

	if err := instance.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("Usando configuración:", instance.ConfigFileUsed())
		}
	}
}
