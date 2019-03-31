package services

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	models "snift-backend/models"
	"snift-backend/utils"
	"strconv"
	"strings"
	"time"
)

// XSSHeader has the XSS Header Name
const XSSHeader = "X-Xss-Protection"

// XFrameHeader has the XFrame Header Name
const XFrameHeader = "X-Frame-Options"

// HSTSHeader has the HSTS Header Name
const HSTSHeader = "Strict-Transport-Security"

// CSPHeader has the CSP Header Name
const CSPHeader = "Content-Security-Policy"

// PKPHeader has the PKP Header Name
const PKPHeader = "Public-Key-Pins"

// RPHeader has the RP Header Name
const RPHeader = "Referrer-Policy"

// XContentTypeHeader has the X-Content-Type Header Name
const XContentTypeHeader = "X-Content-Type-Options"

// Server has the Server Header
const Server = "Server"

// TXTQuery is used to extract all the TXT Records of a Domain
const TXTQuery = "dig @8.8.8.8 +ignore +short +bufsize=1024 domain.com txt"

// DMARCQuery is used to extract all the DMARC Records of a Domain
const DMARCQuery = "dig +short TXT _dmarc.domain.com"

// OpenBugBountyURL is used to query for previous security incidents
const OpenBugBountyURL = "https://www.openbugbounty.org/api/1/search/?domain="

// MaxIncidentResponseTime is the Maximum Incident Response Time taken as 30 days -> 30 * 24 = 720 hours
const MaxIncidentResponseTime = 720

const (
	// HTTPSecure Badge
	HTTPSecure = "HTTP_SECURE"
	// XSSProtect Badge
	XSSProtect = "XSS_PROTECT"
	// HTTPS2 Badge
	HTTPS2 = "HTTP_2.0"
	// TLSSecure Badge
	TLSSecure = "LATEST_TLS"
	// XFRAMEDeny Badge
	XFRAMEDeny = "CLICKJACKING_PROTECT"
	// SeriousSecurity Badge
	SeriousSecurity = "SERIOUS_SECURITY"
)

// XSSValues is used to store the X-Xss-Protection Header values
var XSSValues = [...]string{"0", "1"}

// XFrameValues is used to store the X-Frame-Options Header values
var XFrameValues = [...]string{"deny", "sameorigin", "allow-from"}

// HSTSValues used to store the X-Frame-Options Header values
var HSTSValues = [...]string{"max-age", "includeSubDomains", "preload"}

// ReferrerPolicyValues used to store the Referrer-Policy Header values
var ReferrerPolicyValues = [...]string{"no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-url"}

// XContentTypeHeaderValue is used to store the value for X-Content-Type Options Header
const XContentTypeHeaderValue = "nosniff"

// HTTPVersion is used to store the HTTP Versions
var HTTPVersion = [...]string{"HTTP/2.0", "HTTP/1.1"}

// CalculateProtocolScore returns a score based on whether the protocol is http/https
func CalculateProtocolScore(protocol string) (score int, message string) {
	score = -1
	if strings.Compare(protocol, "http") == 0 {
		score = 0
		message = "Website is unencrypted and hence subjective to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks."
	} else if strings.Compare(protocol, "https") == 0 {
		score = 5
		message = "From the protocol level, Website is secure."
	} else {
		message = "Protocol Not Found"
	}
	return
}

var getDefaultPort = func(protocol string) string {
	// default http port
	port := "80"
	// default https port
	if protocol == "https" {
		port = "443"
	}
	return port
}

// CalculateOverallScore returns the overall score for the incoming request
func CalculateOverallScore(scoresURL string) (*models.ScoreResponse, error) {
	var messages []string
	var score int
	var host string
	var port string
	domain, err := url.Parse(scoresURL)
	if err != nil {
		return nil, err
	}
	protocol := domain.Scheme
	if strings.Contains(domain.Host, ":") {
		host, port, _ = net.SplitHostPort(domain.Host)
	} else {
		host = domain.Host
	}
	if port == "" {
		port = getDefaultPort(protocol)
	}
	protocolScore, protocolMessage := CalculateProtocolScore(protocol)
	messages = append(messages, protocolMessage)
	score = score + protocolScore
	var maximumScore = 5
	headerScore, _, maxScore, ServerDetail, err := GetResponseHeaderScore(scoresURL)
	if err != nil {
		return nil, err
	}
	maximumScore = maximumScore + maxScore
	score = score + headerScore
	mailServerScore, maxScore := GetMailServerConfigurationScore(host)
	score += mailServerScore
	maximumScore += maxScore
	totalScore := math.Ceil((float64(float64(score)/float64(maximumScore)))*100) / 100
	fmt.Println("Protocol Score is " + strconv.Itoa(protocolScore))
	fmt.Println("Message: " + protocolMessage)
	fmt.Println("Final Score for: " + scoresURL + " is " + strconv.Itoa(score) + " out of " + strconv.Itoa(maximumScore))
	certificates, certError := models.GetCertificate(host, port, protocol)
	if certError != nil {
		return nil, certError
	}
	scores := models.GetScores(scoresURL, totalScore, messages)
	response := models.GetScoresResponse(scores, certificates, nil, ServerDetail)
	return response, nil
}

// GetResponseHeaderScore returns the Response Header Score for the HTTP Request
func GetResponseHeaderScore(url string) (totalScore int, XSSReportURL string, maxScore int, serverInfo *models.ServerDetail, err error) {
	err = utils.IsValidURL(url)
	if err != nil {
		return 0, "", 0, nil, err
	}
	var responseHeaderMap map[string]string
	response, err := http.Head(url)
	if err != nil {
		return 0, "", 0, nil, err
	}
	responseHeaderMap = make(map[string]string)
	for k, v := range response.Header {
		value := strings.Join(v, ",")
		responseHeaderMap[k] = value
	}
	totalScore = 0
	var score = 0
	maxScore = 0
	if val, ok := responseHeaderMap[XSSHeader]; ok {
		score, XSSReportURL = GetXSSScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 1
	if val, ok := responseHeaderMap[XFrameHeader]; ok {
		score = GetXFrameScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 2
	if val, ok := responseHeaderMap[HSTSHeader]; ok {
		score = GetHSTSScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 3
	if _, ok := responseHeaderMap[CSPHeader]; ok {
		score = 5
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 3
	if _, ok := responseHeaderMap[PKPHeader]; ok {
		score = 5
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 2
	if val, ok := responseHeaderMap[RPHeader]; ok {
		score = GetReferrerPolicyScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	maxScore = maxScore + 5
	score = 0
	if val, ok := responseHeaderMap[XContentTypeHeader]; ok {
		score = GetXContentTypeScore(val)
	}
	totalScore += score
	maxScore = maxScore + 5
	if val, ok := responseHeaderMap[Server]; ok {
		serverInfo = getServerInformation(val)
	}
	totalScore += GetHTTPVersionScore(response.Proto)
	maxScore = maxScore + 5
	if response.TLS != nil {
		totalScore += GetTLSVersionScore(response.TLS.Version)
	}
	return
}

// GetXSSScore returns the XSS Score of the URL
func GetXSSScore(XSSValue string) (score int, XSSReportURL string) {
	XSSValue = strings.TrimSpace(XSSValue)
	if strings.Compare(XSSValue, XSSValues[0]) == 0 {
		score = 0
	} else if strings.HasPrefix(XSSValue, XSSValues[1]) {
		score = 5
	}
	XSSValueReport := strings.Split(XSSValue, "report=")
	if len(XSSValueReport) == 2 {
		XSSReportURL = XSSValueReport[1]
	}
	return
}

// GetXFrameScore returns the HTTP X-Frame-Options Response Header Score of the URL
func GetXFrameScore(XFrameValue string) (score int) {
	XFrameValue = strings.TrimSpace(strings.ToLower(XFrameValue))
	if strings.Compare(XFrameValue, XFrameValues[0]) == 0 || strings.Compare(XFrameValue, XFrameValues[1]) == 0 {
		score = 5
	} else if strings.HasPrefix(XFrameValue, XFrameValues[2]) {
		score = 4
	}
	return
}

// GetHSTSScore returns the HTTP Strict-Transport-Security Response Header Score of the URL
func GetHSTSScore(HSTS string) (score int) {
	if strings.HasPrefix(HSTS, HSTSValues[0]) {
		score = 4
		if strings.Contains(HSTS, HSTSValues[1]) || strings.Contains(HSTS, HSTSValues[2]) {
			score = 5
		}
	}
	return
}

// GetReferrerPolicyScore returns the HTTP Referrer-Policy Response Header Score of the URL
func GetReferrerPolicyScore(ReferrerPolicy string) (score int) {
	ReferrerPolicy = strings.TrimSpace(strings.ToLower(ReferrerPolicy))
	if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[0]) == 0 {
		score = 5
	} else if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[1]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[2]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[3]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[4]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[5]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[6]) == 0 {
		score = 4
	} else if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[7]) == 0 {
		score = 2
	}
	return
}

// GetXContentTypeScore returns the score for X-Content-Type-Options Header
func GetXContentTypeScore(XContentType string) (score int) {
	score = 0
	if strings.EqualFold(XContentType, XContentTypeHeaderValue) {
		score = 5
	}
	return
}

// GetHTTPVersionScore returns the score for HTTP Version
func GetHTTPVersionScore(Proto string) (score int) {
	score = 0
	if strings.EqualFold(Proto, HTTPVersion[0]) {
		score = 5
	} else if strings.EqualFold(Proto, HTTPVersion[1]) {
		score = 2
	}
	return
}

// GetTLSVersionScore returns the score for TLS Version
func GetTLSVersionScore(Version uint16) (score int) {
	score = 0
	if Version == tls.VersionTLS12 {
		score = 5
	} else if Version == tls.VersionTLS11 {
		score = 3
	} else if Version == tls.VersionTLS10 {
		score = 1
	}
	return
}

// GetMailServerConfigurationScore returns the Mail Server Configuration Score of a Domain
func GetMailServerConfigurationScore(host string) (totalScore int, maximumScore int) {
	maximumScore = 0
	totalScore = 0
	if strings.HasPrefix(host, "www.") {
		host = strings.Replace(host, "www.", "", -1)
	}
	spfScore, maxScore := GetSPFScore(host)
	maximumScore = maximumScore + maxScore
	totalScore = totalScore + spfScore
	totalScore += GetDMARCScore(host)
	maximumScore += 5
	return
}

// GetSPFScore returns the Sender Policy Framework Score of the Domain
func GetSPFScore(domain string) (totalScore int, maxScore int) {
	command := strings.Replace(TXTQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	txtRecords := string(out[:])

	if err != nil {
		fmt.Println("Unexpected Error Occured while extracting TXT Records", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(txtRecords))
	scanner.Split(bufio.ScanLines)

	spfRecordCount := 0
	totalScore = 0

	for scanner.Scan() {
		txtRecord := scanner.Text()
		// Removing Surrounding Quotes and trimming spaces
		txtRecord = strings.TrimSpace(txtRecord[1 : len(txtRecord)-1])
		if strings.HasSuffix(txtRecord, "-all") {
			totalScore = totalScore + 5
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "~all") {
			totalScore = totalScore + 3
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "?all") {
			totalScore = totalScore + 2
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "+all") {
			spfRecordCount++
		}
	}
	maxScore = spfRecordCount * 5
	return
}

// GetDMARCScore returns the DMARC Score of the Domain
func GetDMARCScore(domain string) (score int) {
	command := strings.Replace(DMARCQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	dmarcRecord := string(out[:])

	score = 0

	if err != nil {
		log.Fatal("Unexpected Error Occured while extracting DMARC Records", err)
	}
	if len(dmarcRecord) > 2 {
		dmarcRecord = strings.TrimSpace(dmarcRecord[1 : len(dmarcRecord)-1])
		if strings.HasPrefix(dmarcRecord, "v=DMARC") {
			score = 5
		}
	}
	return
}

// GetPreviousVulnerabilitiesScore gets the score for Previous Vulnerabilities taken from openbugbounty.org
func GetPreviousVulnerabilitiesScore(host string) (totalScore int, maxScore int, IncidentList []models.Incident) {
	if strings.HasPrefix(host, "www.") {
		host = strings.Replace(host, "www.", "", -1)
	}
	resp, err := http.Get(OpenBugBountyURL + host)
	if err != nil {
		log.Fatalln("Error Occured while sending Request ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error Occured while reading HTTP Response ", err)
	}

	var incidents models.Incidents
	err = xml.Unmarshal(body, &incidents)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling XML Response", err)
	}
	maxScore = len(incidents.IncidentList) * 10
	totalScore = 0
	for _, incident := range incidents.IncidentList {
		if incident.Fixed {
			ReportedDate, _ := time.Parse(time.RFC1123Z, incident.ReportedDate)
			FixedDate, _ := time.Parse(time.RFC1123Z, incident.FixedDate)
			diff := FixedDate.Sub(ReportedDate)
			if diff.Hours() > MaxIncidentResponseTime {
				totalScore += 5
			} else {
				totalScore += 10
			}
		}
	}
	return totalScore, maxScore, incidents.IncidentList
}

func getServerInformation(server string) (serverInfo *models.ServerDetail) {
	jsonFile, err := os.Open("resources/web_servers.json")
	if err != nil {
		log.Fatal("Error Occured while opening JSON File ", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal("Error Occured while reading JSON File ", err)
	}

	values := make([]models.WebServer, 0)
	err = json.Unmarshal(jsonValue, &values)
	if err != nil {
		log.Fatal("Error Occured while parsing JSON ", err)
	}

	for _, serverValue := range values {
		if strings.HasPrefix(server, serverValue.Prefix) {
			serverInfo = serverValue.ServerDetail
			break
		}
	}
	return
}
