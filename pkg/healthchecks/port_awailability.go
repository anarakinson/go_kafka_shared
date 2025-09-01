package healthchecks

import (
	"fmt"
	"net"
)

// isPortAvailable - утилитарная функция проверки доступности одного порта
func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// CheckPortsAvailability - проверка доступности указанных портов
// requiredPorts - порты, которые нужно проверить
func CheckPortsAvailability(requiredPorts []string) error {

	for _, port := range requiredPorts {
		if !isPortAvailable(port) {
			return fmt.Errorf("port %s is already in use", port)
		}
	}

	return nil

}
