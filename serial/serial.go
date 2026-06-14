package serial

import (
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
)

type IRSensor struct {
	One   Coordinate
	Two   Coordinate
	Three Coordinate
	Four  Coordinate
}

type Coordinate struct {
	XCoord string
	YCoord string
}

func GetOutput(hello chan []byte) {

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	// Print the list of detected ports
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 19200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[1], mode)
	if err != nil {
		fmt.Println("Busy port", err)
		return
	}

	defer port.Close()

	// Send the string "10,20,30\n\r" to the serial port
	//n, err := port.Write([]byte("10,20,30\n\r"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Sent %v bytes\n", n)

	// Read and print the response
	buff := make([]byte, 100)
	ir := IRSensor{}
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			fmt.Println("error on port", err)
			return
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		if n > 1 {
			coords := getCoordinates(buff, n)
			if len(coords) == 8 {
				ir.One = Coordinate{
					XCoord: coords[0],
					YCoord: coords[1],
				}
				ir.Two = Coordinate{
					XCoord: coords[2],
					YCoord: coords[3],
				}
				ir.Three = Coordinate{
					XCoord: coords[4],
					YCoord: coords[5],
				}
				ir.Four = Coordinate{
					XCoord: coords[6],
					YCoord: coords[7],
				}
				// fmt.Println("Cords", ir)
				message := fmt.Sprintf("Coordinate 1: x:%s y:%s | (2) x:%s, y:%s | (3) x:%s, y:%s | (4) x:%s, y:%s",
					ir.One.XCoord,
					ir.One.YCoord,
					ir.Two.XCoord,
					ir.Two.YCoord,
					ir.Three.XCoord,
					ir.Three.YCoord,
					ir.Four.XCoord,
					ir.Four.YCoord)
				hello <- []byte(message)
			}
		}

		//fmt.Printf("-read: %d, buff: %s-\n", n, string(buff[:n]))
	}
}

func getCoordinates(buffedStr []byte, buffedLen int) []string {
	//input := string(buffedStr[:buffedLen])
	var output []byte
	for _, b := range buffedStr[:buffedLen] {
		if ('0' <= b && b <= '9') || b == ',' {
			output = append(output, b)
		}
	}
	// fmt.Println("getCoordinates", output)
	return strings.Split(string(output), ",")
}
