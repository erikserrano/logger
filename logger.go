// Paquete logger encargado de respaldar en archivos de texto mensajes
// y mostrar en consola
package logger

import (
	"os"
	"os/user"
	"time"
	"fmt"
)

// Estructura Logger encargada de almacenar la instancia para almacenamiento de mensajes
type Logger struct {
	FileName string
	FilePath string
}

// Función encargada de crear/abrir la ruta del archivo de respaldo
func openFile(filePath, fileName string) (*os.File, error) {
	// Creamos directorio
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		return nil, err
	}

	// Abrimos archivo
	file, err := os.OpenFile(filePath+fileName, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		// Creamos archivo
		file, err = os.Create(filePath + fileName)
		if err != nil {
			return nil, err
		} else {
			initText := "#version: 0.1\n"
			initText += "#creation: " + time.Now().String() + "\n"
			initText += "#config: [datetime][userid:username][alert/error/info/etc] message\n"

			// Cerramos archivo
			defer file.Close()

			// Escribimos encabezado
			if _, err := file.WriteString(initText); err != nil {
				return nil, err
			}
		}
	}

	return file, nil
}

// Método encargado de escribir un mensaje sobre el archivo LOG
func (log *Logger) WriteLine(kind, message string) (int, error) {

	// Abrimos archivo
	file, err := openFile(log.FilePath, log.FileName)
	if err != nil {
		return 0, nil
	}
	// Cerramos archivo
	defer file.Close()

	// Obtenemos la información del usuario del Sistema Operativo
	user, err := user.Current()
	if err != nil {
		return 0, err
	}

	userMessage := user.Uid + ":" + user.Username
	timeMessage := time.Now().Format("02/01/2006 15:04:05.99999")

	// Guardamos mensaje
	n, err := file.WriteString(fmt.Sprintf("[%s][%s][%s] %s\n", timeMessage, userMessage, kind, message))
	if err != nil {
		return 0, err
	}

	// Imprimitmos mensaje
	println(timeMessage + logMessage)

	return n, nil
}

// Función encargada de crear una nueva estructura para el respaldo de los mensajes
func New(filePath, fileName string) (*Logger, error) {
	log := &Logger{fileName, filePath}
	if _, err := openFile(filePath, fileName); err != nil {
		return nil, err
	}
	return log, nil
}