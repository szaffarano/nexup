package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/szaffarano/nexup/repository"
	"github.com/szaffarano/nexup/util"
)

func init() {
	rootCmd.AddCommand(putCmd)
}

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Sube un archivo al repositorio",
	Long:  `Sube un archivo al repositorio`,
	Run: func(cmd *cobra.Command, args []string) {
		var username, password string
		var err error

		if len(nexupfileCredentials.Password) != 0 &&
			len(nexupfileCredentials.Username) != 0 {
			if verbose {
				fmt.Println("Usando credenciales provistas en la configuración")
			}
			username = nexupfileCredentials.Username
			password = nexupfileCredentials.Password
		} else {
			username, password, err = util.GetCredentials()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		repo, err := repository.NewRepository(
			nexupfile.Repository, username, password, []byte(nexupfile.Truststores))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) == 0 {
			if verbose {
				fmt.Fprintf(os.Stderr, "No se indicaron archivos a publicar")
			}
		}

		for _, f := range args {
			ext := filepath.Ext(filepath.Base(f))
			name := strings.TrimSuffix(filepath.Base(f), ext)

			// TODO parametrizar la generación de la url
			url := fmt.Sprintf("%s/%s/%s/%s/%s-%s%s",
				nexupfile.Repository,
				nexupfile.System,
				nexupfile.Application,
				nexupfile.Version,
				name,
				nexupfile.Version,
				ext)

			fmt.Println(url)

			file, err := os.Open(f)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()

			err = repo.Put(url, file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
