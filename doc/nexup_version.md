## nexup version

Imprime la versión de nexup

### Synopsis


Imprime la versión de nexup

```
nexup version [flags]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
  -c, --credentials string   Archivo de credenciales para autenticarse en el repositorio, por 
			 defecto se usa $HOME/.nexup-credentials

			 Ejemplo

			 username: NombreDeUsuario
			 password: Contraseña
		 	
  -n, --nexupfile string     Archivo de configuración con la información del repositorio
		 	 y de lo que se quiere subir, por defecto se usa ./Nexupfile

			 Ejemplo

			 system: nombre-del-sistema
			 application: nombre-de-la-aplicacion
			 version: 1.2.3
			 repository: http://host:puerto/repository/etc/etc
			 truststores: |
			  .......
		 	
  -t, --truststores string   Archivo con certificados de las CAs en las cuales confiar
  -v, --verbose              Imprime información extra
```

### SEE ALSO
* [nexup](nexup.md)	 - Nexup es un comando para subir contenido a repositorios

###### Auto generated by spf13/cobra on 21-Feb-2018