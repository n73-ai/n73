package main

import (
    "fmt"
    "os"
    
    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/launcher"
)

func main() {
    // Validar argumentos
    if len(os.Args) < 2 {
        fmt.Println("Uso: ./screenshot <url>")
        fmt.Println("Ejemplo: ./screenshot https://google.com")
        os.Exit(1)
    }

    url := os.Args[1]

    // Rod descarga Chrome automáticamente si no existe
    path, _ := launcher.LookPath()
    u := launcher.New().
        Bin(path).
        Headless(true).
        NoSandbox(true).
        MustLaunch()
    
    browser := rod.New().ControlURL(u).MustConnect()
    defer browser.MustClose()

    // Tomar screenshot
    fmt.Printf("Tomando screenshot de: %s\n", url)
    page := browser.MustPage(url)
    page.MustWaitLoad()
    page.MustWaitIdle()
    page.MustScreenshot("/app/screenshot.png")
    
    fmt.Println("✓ Screenshot guardado en screenshot.png")
}
