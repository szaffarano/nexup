package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docCmd = &cobra.Command{
		Use:   "doc",
		Short: "Genera la documentación del sistema",
		Long: `Genera documentación de los comandos y el modo de uso.  Posibles
				valores: MD (MarkDown), RS (ReStructured Text) y MAN (man pages)`,
		Run: func(cmd *cobra.Command, args []string) {
			info, err := os.Stat(output)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.Mkdir(output, os.ModeDir|0775)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				} else {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if !info.IsDir() {
				fmt.Println(fmt.Sprintf("%s: Debe indicar un directorio", output))
				os.Exit(1)
			}

			format = strings.ToUpper(format)
			switch format {
			case "MD":
				err := doc.GenMarkdownTree(rootCmd, output)
				if err != nil {
					fmt.Println(err)
				}
			case "MAN":
				header := &doc.GenManHeader{
					Title:   "Nexup",
					Section: "5",
				}
				err := doc.GenManTree(rootCmd, header, output)
				if err != nil {
					fmt.Println(err)
				}
			case "RS":
				err := doc.GenReSTTree(rootCmd, output)
				if err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Println(fmt.Sprintf("%s: Formato desconocido", format))
			}
			if verbose {
				fmt.Println(fmt.
					Sprintf("Se generó la documentación en %s", output))
			}
		},
	}
	format string
	output string
)

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.
		PersistentFlags().
		StringVarP(&format,
			"format",
			"f",
			"md",
			"Formato de la documentación.  Por default Markdown ")
	docCmd.
		PersistentFlags().
		StringVarP(&output,
			"output",
			"o",
			"./doc",
			`Directorio donde se almacenará la documentación 
			 generada.  Por defecto ./doc`)
}
