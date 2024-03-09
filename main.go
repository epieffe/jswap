package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

const api = "https://api.adoptium.net/v3"

type AvailableReleases struct {
    Releases []int `json:"available_releases"`
    LTS []int `json:"available_lts_releases"`
    MostRecentLTS int `json:"most_recent_lts"`
    MostRecentRelease int `json:"most_recent_feature_release"`
}

// Fetches JSON data from an url and deserializes response.
func fetchDataFromURL[T any](url string) (*T, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP response status error: %s", resp.Status)
    }

    var result T
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return nil, err
    }

    return &result, nil
}

// Fetches available Java releases from Eclipse Temurin.
func fetchAvailableReleases() (*AvailableReleases, error) {
    url := api + "/info/available_releases"
    return fetchDataFromURL[AvailableReleases](url)
}

func main() {
    releases, err := fetchAvailableReleases()
    if err != nil {
        fmt.Println("Errore:", err)
        return
    }

    fmt.Println("Releases:", releases.Releases)
    fmt.Println("LTS:", releases.LTS)
    fmt.Println("Latest LTS:", releases.MostRecentLTS)
    fmt.Println("Latest release:", releases.MostRecentRelease)
}
