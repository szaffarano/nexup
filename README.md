# Cliente Nexus

## Modo de uso

```sh
nexup [-v] archivo1 archivo2 ... archivoN
```

Para más información consultar [la documentación](./doc/nexup.md).

## Configuración

Todos los parámetros de configuración además de leerse de los archivos que se
detallan a continuación, podrán ser especificados mediante variables de
ambiente con el prefijo ```NEXUP```.  Por ejemplo, el parámetro ```version```
(ver archivo Nexupfile), podrá indicarse con la variable ```NEXUP_VERSION```,
eso aplica para toda la configuración.

### Archivo Nexupfile

Este archivo tiene la configuración de los artefactos a subir.  Se puede 
especificar con el parámetro --conf, o bien se lee de la ubicación por
defeto en el directorio actual.  Ejemplo:

```yaml
system: supersistema
application: frontend
version: 4.0.0
repository: http://localhost:8081/repository/test-raw
```

Opcionalmente se le puede incluir una entrada ```truststores``` en donde se le
configura una lista de certificados de CAs en las cuales confiar.

### Archivo .nexup-credentials

Este archivo contiene las credenciales para autenticarse en el servidor.  Se
especifica con el parámetro --cred o bien se intenta leer de la ubicación por
defecto que es ```$HOME/.nexup-credentials```.  Ejemplo

```yaml
username: nombre-de-usuario
password: contraseña
```

Si no se encuentra el archivo ni los valores se especifican mediante variables
de ambiente, se le preguntará al usuario por los datos de autenticación.