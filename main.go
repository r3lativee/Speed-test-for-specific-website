package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: speedtest <URL>")
        return
    }

    url := os.Args[1]

    // Clear the terminal and print the header
    clearScreen()
    printHeader()

    fmt.Printf("\nTesting download speed from %s...\n", url)

    response, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error making HTTP request: %v\n", err)
        return
    }
    defer response.Body.Close()

    buffer := make([]byte, 32*1024) // 32KB buffer
    var downloaded int64 = 0
    start := time.Now()

    ticker := time.NewTicker(100 * time.Millisecond) // Update every 100ms

    previousTime := start
    previousDownloaded := downloaded

    go func() {
        for {
            select {
            case <-ticker.C:
                currentTime := time.Now()
                elapsed := currentTime.Sub(previousTime).Seconds()
                if elapsed == 0 {
                    elapsed = 0.1 // Avoid division by zero
                }
                speed := float64(downloaded-previousDownloaded) / elapsed / (1024 * 1024) // Speed in MB/s

                printProgress(downloaded, speed)

                previousTime = currentTime
                previousDownloaded = downloaded
            }
        }
    }()

    // Read data until the download completes
    for {
        n, err := response.Body.Read(buffer)
        if n > 0 {
            downloaded += int64(n)
        }
        if err != nil {
            if err == io.EOF {
                // End of file, download complete
                break
            }
            fmt.Printf("Error reading response body: %v\n", err)
            return
        }
    }

    elapsed := time.Since(start).Seconds()
    speed := float64(downloaded) / elapsed / (1024 * 1024) // Speed in MB/s
    fmt.Printf("\n\nDownloaded %.2f MB in %.2f seconds\n", float64(downloaded)/(1024*1024), elapsed)
    fmt.Printf("Average download speed: %.2f MB/s\n", speed)
}

// just extra ignore
func printHeader() {
    fmt.Println("\t\t\t======================================")
    fmt.Println("\t\t\t               Speed Test             ")
    fmt.Println("\t\t\t          Developed by r3lative       ")
    fmt.Println("\t\t\t======================================")
}

// clearScreen clears the terminal screen
func clearScreen() {
    fmt.Print("\033[H\033[2J")
}

// printProgress shows the current download speed
func printProgress(downloaded int64, speed float64) {
    fmt.Printf("\r%.2f MB downloaded - Speed: %.2f MB/s", float64(downloaded)/(1024*1024), speed)
}

