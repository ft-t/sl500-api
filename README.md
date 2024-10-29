
# SL500 Card Reader API

The `sl500_api` package provides a Go-based interface for interacting with the SL500 card reader, allowing applications to communicate with the reader via serial connections. This package includes functions for initializing, configuring, and controlling the SL500, as well as reading and writing data on compatible RFID cards.

## Features

- Initialize and manage serial connections with SL500.
- Control the reader's antenna, LED, and buzzer.
- RFID card operations, including reading, writing, and authentication.
- Compatibility with MIFARE and ISO15693 tags.

## Installation

To include `sl500_api` in your project, use:

```bash
go get github.com/your_username/sl500_api
```

## Usage

### Initializing a Connection

To start, create a new connection to the SL500 card reader using the `NewConnection` function. Specify the serial port path, baud rate, logging option, and timeout.

```go
package main

import (
    "fmt"
    "time"
    "github.com/ft-t/sl500-api"
)

func main() {
    // Initialize SL500 with the appropriate serial port path, baud rate, logging, and timeout
    reader, err := sl500_api.NewConnection("/dev/ttyUSB0", sl500_api.Baud.Baud9600, true, 3*time.Second)
    if err != nil {
        fmt.Println("Error initializing SL500:", err)
        return
    }
    defer reader.Close()
}
```

### Basic Commands

Below are some primary commands for controlling the SL500:

- **Open and Close Connection**
  ```go
  err := reader.Open()      // Opens the serial connection
  err = reader.Close()      // Closes the connection
  ```

- **Antenna Control**
  ```go
  reader.RfAntennaSta(sl500_api.AntennaOn)    // Turn antenna on
  reader.RfAntennaSta(sl500_api.AntennaOff)   // Turn antenna off
  ```

- **RFID Card Commands**
  ```go
  reader.RfRequest(sl500_api.RequestAll)     // Scan for all cards
  reader.RfAnticoll()                        // Anti-collision function
  ```

- **MIFARE Operations**
  ```go
  blockNumber := byte(4)
  data, _ := reader.RfM1Read(blockNumber)    // Read data from block
  fmt.Printf("Data: %x\n", data)

  writeData := []byte{0x01, 0x02, 0x03, 0x04}
  reader.RfM1Write(blockNumber, writeData)   // Write data to block
  ```

### Advanced Commands

#### Setting Device Type
```go
reader.RfInitType(sl500_api.Type_A)
```

#### Enabling Beep and LED Light
```go
reader.RfBeep(100)                        // Beep for 100 milliseconds
reader.RfLight(sl500_api.ColorGreen)      // Set LED to green
```

## Error Handling

Each function returns an error object; check this for successful execution. Example:
```go
if err := reader.RfAntennaSta(sl500_api.AntennaOn); err != nil {
    fmt.Println("Antenna error:", err)
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue for feedback or feature requests.
