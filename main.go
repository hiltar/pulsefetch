package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "strings"
)

type Token struct {
    Address            string `json:"address"`
    CirculatingMarketCap string `json:"circulating_market_cap"`
    Decimals           string `json:"decimals"`
    ExchangeRate       string `json:"exchange_rate"`
    Holders            string `json:"holders"`
    IconURL            string `json:"icon_url"`
    Name               string `json:"name"`
    Symbol             string `json:"symbol"`
    TotalSupply        string `json:"total_supply"`
    Type               string `json:"type"`
}

type TokenBalance struct {
    Token         Token  `json:"token"`
    TokenID       string `json:"token_id"`
    TokenInstance string `json:"token_instance"`
    Value         string `json:"value"`
}

// formatWithCommas adds thousand separators to a numeric string
func formatWithCommas(s string) string {
    s = strings.TrimLeft(s, "0")
    if s == "" {
        return "0"
    }
    n := len(s)
    // Calculate the length with commas
    commaCount := (n - 1) / 3
    result := make([]byte, n+commaCount)
    // Build the string from right to left
    for i, j, k := n-1, len(result)-1, 0; i >= 0; i, j = i-1, j-1 {
        result[j] = s[i]
        k++
        if k%3 == 0 && i > 0 {
            j--
            result[j] = ','
        }
    }
    return string(result)
}

func main() {
    // Check if address is provided as a command-line argument
    if len(os.Args) < 2 {
        fmt.Println("Please provide an address as a command-line argument")
        fmt.Println("Usage: go run main.go <address>")
        os.Exit(1)
    }
    address := os.Args[1]

    // Construct API URL
    url := fmt.Sprintf("https://api.scan.pulsechain.com/api/v2/addresses/%s/token-balances", address)

    // Make HTTP GET request
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error fetching data: %v\n", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response: %v\n", err)
        os.Exit(1)
    }

    // Parse JSON response
    var balances []TokenBalance
    if err := json.Unmarshal(body, &balances); err != nil {
        fmt.Printf("Error parsing JSON: %v\n", err)
        os.Exit(1)
    }

    type row struct {
        name    string
        balance string
    }

    var rows []row
    maxNameLen := len("Token Name")
    maxBalanceLen := len("Balance")

    for _, balance := range balances {
        // Get token name, use "Unknown" if name is null
        name := balance.Token.Name
        if name == "" {
            name = "Unknown"
        }

        // Skip unknown named tokens
        if name == "Unknown" {
            continue
        }

        // Convert value to integer string, adjusting for decimals
        valueStr, err := adjustValue(balance.Value, balance.Token.Decimals)
        if err != nil {
            fmt.Printf("Error adjusting value for token %s: %v\n", name, err)
            continue
        }

        // Format value with commas
        formattedValue := formatWithCommas(valueStr)

        // Update max lengths
        if len(name) > maxNameLen {
            maxNameLen = len(name)
        }
        if len(formattedValue) > maxBalanceLen {
            maxBalanceLen = len(formattedValue)
        }

        rows = append(rows, row{name: name, balance: formattedValue})
    }

    if len(rows) == 0 {
        fmt.Println("No tokens found.")
        return
    }

    // Print header
    fmt.Printf("| %-*s | %*s |\n", maxNameLen, "Token Name", maxBalanceLen, "Balance")

    // Print separator
    fmt.Printf("|-%s-|-%s-|\n", strings.Repeat("-", maxNameLen), strings.Repeat("-", maxBalanceLen))

    // Print rows
    for _, r := range rows {
        fmt.Printf("| %-*s | %*s |\n", maxNameLen, r.name, maxBalanceLen, r.balance)
    }
}

// adjustValue returns the integer part (as string) by removing decimal places
func adjustValue(valueStr, decimalsStr string) (string, error) {
    valueStr = strings.TrimLeft(valueStr, "0")
    if valueStr == "" {
        return "0", nil
    }

    if decimalsStr == "" {
        decimalsStr = "0"
    }

    decimals, err := strconv.Atoi(decimalsStr)
    if err != nil {
        return "", fmt.Errorf("invalid decimals: %v", err)
    }

    if decimals < 0 {
        return "", fmt.Errorf("negative decimals")
    }

    valueLen := len(valueStr)
    if valueLen <= decimals {
        return "0", nil
    }

    // Take the integer part by trimming the last 'decimals' digits
    intPart := valueStr[:valueLen-decimals]
    return intPart, nil
}
