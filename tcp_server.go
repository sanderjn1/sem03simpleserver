package main

import (
        "io"
        "log"
        "net"
        "sync"
        "strings"
	"strconv"
	"fmt"
	"errors"
        "github.com/sanderjn1/is105sem03/mycrypt"
	"github.com/sanderjn1/funtemps/conv"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8100")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Printf("bundet til %s", server.Addr().String())

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			log.Println("før server.Accept() kallet")

			conn, err := server.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				defer conn.Close()

				for {
					buf := make([]byte, 2048)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return
					}

					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))

					msgString := string(dekryptertMelding)

					switch msgString {
					case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, -4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))

					default:
						if strings.HasPrefix(msgString, "Kjevik") {
							newString, err := CelsiusToFahrenheitLine("Kjevik;SN39040;18.03.2022 01:50;6")
							if err != nil {
								log.Fatal(err)
							}

							kryptertMelding := mycrypt.Krypter([]rune(newString), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
							_, err = conn.Write([]byte(string(kryptertMelding)))
						} else {
							kryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
							_, err = c.Write([]byte(string(kryptertMelding)))

						}
					}

					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return
					}
				}
			}(conn)
		}
	}()

	wg.Wait()
}

func CelsiusToFahrenheitLine(line string) (string, error) {

	dividedString := strings.Split(line, ";")
	var err error

	if len(dividedString) == 4 {
		dividedString[3], err = CelsiusToFahrenheitString(dividedString[3])
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("linje har ikke forventet format")
	}
	return strings.Join(dividedString, ";"), nil
}

func CelsiusToFahrenheitString(celsius string) (string, error) {
	var fahrFloat float64
	var err error
	if celsiusFloat, err := strconv.ParseFloat(celsius, 64); err == nil {
		fahrFloat = conv.CelcsiusToFahrenheit(celsiusFloat)
	}
	fahrString := fmt.Sprintf("%.1f", fahrFloat)
	return fahrString, err
}
