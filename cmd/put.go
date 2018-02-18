package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/szaffarano/nexup/repository"
	"github.com/szaffarano/nexup/util"
)

var (
	noSha1    bool
	mavenRepo bool
)

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.
		PersistentFlags().
		BoolVarP(&noSha1,
			"no-sha1",
			"S",
			false,
			"No subir suma de comprobación sha1 junto con el archivo ")
	putCmd.
		PersistentFlags().
		BoolVarP(&mavenRepo,
			"maven-repo",
			"m",
			false,
			`Se subirán los archivos a un repositorio de tipo maven.  Por 
		     defecto se asume que el repositorio es raw`)

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

		repo, err := repository.New(
			nexupfile.Repository, username, password, nexupfile.Truststores)
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
			artifact := repository.Artifact{
				System:      nexupfile.System,
				Application: nexupfile.Application,
				Version:     nexupfile.Version,
				File:        filepath.Base(f),
			}

			url := util.GetRawURL(*repo, artifact)
			urlSha := fmt.Sprintf("%s.sha1", url)

			fmt.Println(url)

			// 1. hacer put del archivo
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

			if !noSha1 {
				// hacer put del sha1
				fileForSha, err := os.Open(f)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer fileForSha.Close()

				hash := sha1.New()
				if _, err := io.Copy(hash, fileForSha); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				sha, err := strings.NewReader(hex.EncodeToString(hash.Sum(nil))), nil
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(urlSha)

				err = repo.Put(urlSha, sha)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	},
}
