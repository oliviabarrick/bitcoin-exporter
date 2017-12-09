package usps

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/xml"
    "net/url"
)

type TrackResponse struct {
    TrackInfo TrackInfo
}

/*
<TrackInfo ID="9410810298370122910773">
    <Class>Priority Mail&lt;SUP&gt;&amp;reg;&lt;/SUP&gt;</Class>
    <ClassOfMailCode>PM</ClassOfMailCode>
    <DestinationCity>SAN FRANCISCO</DestinationCity>
    <DestinationState>CA</DestinationState>
    <DestinationZip>94108</DestinationZip>
    <EmailEnabled>false</EmailEnabled>
    <KahalaIndicator>false</KahalaIndicator>
    <MailTypeCode>DM</MailTypeCode>
    <MPDATE>2017-04-03 15:37:41.000000</MPDATE>
    <MPSUFFIX>367568793</MPSUFFIX>
    <OriginCity>SANTA FE SPRINGS</OriginCity>
    <OriginState>CA</OriginState>
    <OriginZip>90670</OriginZip>
    <PodEnabled>true</PodEnabled>
    <RestoreEnabled>false</RestoreEnabled>
    <RramEnabled>false</RramEnabled>
    <RreEnabled>false</RreEnabled>
    <Service>Signature Confirmation&lt;SUP&gt;&amp;#153;&lt;/SUP&gt;</Service>
    <Service>Up to $100 insurance included</Service>
    <ServiceTypeCode>108</ServiceTypeCode>
    <Status>Delivered, Individual Picked Up at Postal Facility</Status>
    <StatusCategory>Delivered</StatusCategory>
    <StatusSummary>Your item was picked up at a postal facility at 11:27 am on April 8, 2017 in SAN FRANCISCO, CA 94133. The item was signed for by J JUSTIN.</StatusSummary>
    <TABLECODE>T</TABLECODE>
    <TrackSummary></TrackSummary>
    <TrackDetail></TrackDetail>
</TrackInfo>
*/
type TrackInfo struct {
    ID string `xml:"ID,attr"`
    Class string
    ClassOfMailCode string
    DestinationCity string
    DestinationState string
    DestinationZip string
    EmailEnabled bool
    KahalaIndicator bool
    MailTypeCode string
    MPDATE string
    MPSUFFIX string
    OriginCity string
    OriginState string
    OriginZip string
    PodEnabled bool
    RestoreEnabled bool
    RramEnabled bool
    RreEnabled bool
    Service []string
    ServiceTypeCode int
    Status []string
    StatusCategory string
    StatusSummary string
    TABLECODE string
    TrackSummary TrackSummary
    TrackDetails []TrackDetail `xml:"TrackDetail"`
}

/*
    <TrackSummary>
      <EventTime>11:27 am</EventTime>
      <EventDate>April 8, 2017</EventDate>
      <Event>Delivered, Individual Picked Up at Postal Facility</Event>
      <EventCity>SAN FRANCISCO</EventCity>
      <EventState>CA</EventState>
      <EventZIPCode>94133</EventZIPCode>
      <EventCountry/>
      <FirmName/>
      <Name>J JUSTIN</Name>
      <AuthorizedAgent>false</AuthorizedAgent>
      <EventCode>01</EventCode>
      <DeliveryAttributeCode>09</DeliveryAttributeCode>
    </TrackSummary>
*/
type TrackSummary struct {
    EventTime string
    EventDate string
    Event string
    EventCity string
    EventState string
    EventZIPCode string
    EventCount string
    FirmName string
    Name string
    AuthorizedAgent bool
    EventCode int
    DeliveryAttributeCode int
}

/*
    <TrackDetail>
      <EventTime>5:02 pm</EventTime>
      <EventDate>April 3, 2017</EventDate>
      <Event>Shipment Received, Package Acceptance Pending</Event>
      <EventCity>SANTA FE SPRINGS</EventCity>
      <EventState>CA</EventState>
      <EventZIPCode>90670</EventZIPCode>
      <EventCountry/>
      <FirmName/>
      <Name/>
      <AuthorizedAgent>false</AuthorizedAgent>
      <EventCode>TM</EventCode>
    </TrackDetail>
*/
type TrackDetail struct {
    EventTime string
    EventDate string
    Event string
    EventCity string
    EventState string
    EventZIPCode string
    EventCountry string
    FirmName string
    Name string
    AuthorizedAgent bool
    EventCode string
}

/*
<?xml version="1.0" encoding="UTF-8" ?>
<TrackFieldRequest USERID="820NA0004722">
  <Revision>1</Revision>
  <ClientIp>111.0.0.1</ClientIp>
  <SourceId>hello</SourceId>
  <TrackID ID="9410810298370122910773"></TrackID>
</TrackFieldRequest>
*/
type TrackFieldRequest struct {
    UserId string `xml:"USERID,attr"`
    Revision int
    ClientIp string
    SourceId string
    TrackID TrackID
}

type TrackID struct {
    Id string `xml:"ID,attr"`
}

func (tr TrackFieldRequest) construct_url() (usps_url string, err error) {
    proc_inst := xml.ProcInst{
        Target: "xml",
        Inst:   []byte("version=\"1.0\" encoding=\"UTF-8\""),
    }

    encoded := &bytes.Buffer{}

    enc := xml.NewEncoder(encoded)
    enc.Indent(" ", "    ")

    if err = enc.EncodeToken(proc_inst); err != nil {
        return
    }

    if err = enc.Encode(tr); err != nil {
        return
    }

    parsed_url, err := url.Parse("https://secure.shippingapis.com/ShippingAPI.dll")
    if err != nil {
        return
    }

    parameters := url.Values{}
    parameters.Add("XML", encoded.String())
    parameters.Add("API", "TrackV2")
    parsed_url.RawQuery = parameters.Encode()

    usps_url = parsed_url.String()
    return
}

type PackageTracker struct {
    client *http.Client
}

func (pt *PackageTracker) Fetch(track_id string) (track_response TrackResponse, err error) {
    var id TrackID
    id.Id = track_id

    var track_request TrackFieldRequest
    track_request.UserId = "820NA0004722"
    track_request.TrackID = id
    track_request.Revision = 1
    track_request.ClientIp = "111.0.0.1"
    track_request.SourceId = "hello"

    usps_url, err := track_request.construct_url()
    if err != nil {
        return
    }

    if pt.client == nil {
        pt.client = &http.Client{}
    }

    resp, err := pt.client.Get(usps_url)
    if err != nil {
        return
    }

    if resp.StatusCode != 200 {
        err = fmt.Errorf("Response status %d != 200\n", resp.StatusCode)
        return
    }

    if err = xml.NewDecoder(resp.Body).Decode(&track_response); err != nil {
        return
    }

    return
}
