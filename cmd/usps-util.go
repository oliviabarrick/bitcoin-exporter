package main

import (
    "fmt"
    "usps"
)

func main() {
    var pt usps.PackageTracker
    track_response, err := pt.Fetch("9410810298370122910773")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Order status: ", track_response.TrackInfo.Status)
    fmt.Printf("%s, %s -> %s, %s\n", track_response.TrackInfo.OriginCity, track_response.TrackInfo.OriginState, track_response.TrackInfo.DestinationCity, track_response.TrackInfo.DestinationState)
    for _, detail := range track_response.TrackInfo.TrackDetails {
        fmt.Println("\nEvent:", detail.Event)
        fmt.Println("  Date:", detail.EventDate)
        fmt.Println("  City:", detail.EventCity)
    }
}
