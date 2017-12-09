package usps

import (
    "encoding/xml"
    "github.com/stretchr/testify/assert"
    "fmt"
    "strings"
    "testing"
    "net/url"
)

func TestConstructURL(t *testing.T) {
    var id TrackID
    id.Id = "abcd"

    var track_request TrackFieldRequest
    track_request.UserId = "1234"
    track_request.TrackID = id
    track_request.Revision = 1
    track_request.ClientIp = "111.0.0.1"
    track_request.SourceId = "hello"

    usps_url, err := track_request.construct_url()
    if err != nil {
        t.Error(err)
    }

    parsed_url, err := url.Parse(usps_url)
    if err != nil {
        t.Error(err)
    }

    query, err := url.ParseQuery(parsed_url.RawQuery)
    if err != nil {
        t.Error(err)
    }

    var xml_request TrackFieldRequest
    decoder := xml.NewDecoder(strings.NewReader(query["XML"][0]))
    err = decoder.Decode(&xml_request)
    if err != nil {
        t.Error(err)
    }

    assert.Equal(t, parsed_url.Host, "secure.shippingapis.com", "hostname mismatch")
    assert.Equal(t, parsed_url.Path, "/ShippingAPI.dll", "path mismatch")
    assert.Equal(t, query["API"][0], "TrackV2", "api mismatch")
    assert.Equal(t, xml_request, track_request, "request mismatch")
}

func TestUsps(t *testing.T) {
    var pt PackageTracker
    track_response, err := pt.Fetch("9410810298370122910773")
    if err != nil {
        t.Error(err)
    }

    fmt.Println("Order status: ", track_response.TrackInfo.Status)
    fmt.Printf("%s, %s -> %s, %s\n", track_response.TrackInfo.OriginCity, track_response.TrackInfo.OriginState, track_response.TrackInfo.DestinationCity, track_response.TrackInfo.DestinationState)
    //for _, detail := range track_response.TrackInfo.TrackDetails {
        // fmt.Println("\nEvent:", detail.Event)
        // fmt.Println("  Date:", detail.EventDate)
        // fmt.Println("  City:", detail.EventCity)
    //}
}
