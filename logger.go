// Package logger encargado de respaldar en archivos de texto mensajes y mostrar en consola
package logger

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
)

// Logger estructura encargada de almacenar la instancia para almacenamiento de mensajes
type Logger struct {
	FileName string
	FilePath string
	Output   bool
}

// openFile open/create file for loggin
func openFile(filePath, fileName string) (*os.File, error) {

	if strings.Contains(filePath, "\\") && !strings.HasSuffix(filePath, "\\") {
		filePath += "\\"
	} else if strings.Contains(filePath, "/") && !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}

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
		}
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

	return file, nil
}

// WriteLine write message in log file
func (log *Logger) WriteLine(message string, kind ...string) (int, error) {

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

	kinds := "info"
	if len(kind) > 0 {
		kinds = strings.Join(kind, ",")
	}

	// Guardamos mensaje
	newLine := fmt.Sprintf("[%s][%s][%s] %s\n", timeMessage, userMessage, kinds, message)
	n, err := file.WriteString(newLine)
	if err != nil {
		return 0, err
	}

	// Imprimitmos mensaje
	if log.Output {
		fmt.Printf("[%s][%s][%s] %s\n", timeMessage, userMessage, kind, message)
	}
	return n, nil
}

// New create instance for loggin messages
func New(filePath, fileName string, output bool) (*Logger, error) {
	log := &Logger{fileName, filePath, output}
	if _, err := openFile(filePath, fileName); err != nil {
		return nil, err
	}
	return log, nil
}
