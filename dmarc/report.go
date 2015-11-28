package dmarc

import (
	"encoding/xml"
	"io"
)

type DateRange struct {
	// TODO: should be int but Y! trailing spaces
	Begin string `xml:"begin"`
	End   string `xml:"end"`
}

type ReportMetadata struct {
	OrgName      string    `xml:"org_name"`
	Email        string    `xml:"email"`
	ExtraContact string    `xml:"extra_contact_info"`
	ReportId     string    `xml:"report_id"`
	DateRange    DateRange `xml:"date_range"`
}

type PolicyPublished struct {
	Domain          string `xml:"domain"`
	Adkim           string `xml:"adkim"`
	Aspf            string `xml:"aspf"`
	Policy          string `xml:"p"`
	SubdomainPolicy string `xml:"sp"`
	Percentage      int    `xml:"pct"`
}

type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	Dkim        string `xml:"dkim"`
	Spf         string `xml:"spf"`
}

type Row struct {
	// TODO: Figure out how to cast this to an IP
	SourceIp        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

type AuthResult struct {
	// FIXME: this could be either DKIM or SPF
	XMLName xml.Name
	Domain  string `xml:"domain"`
	Result  string `xml:"result"`
}

type AuthResults struct {
	AuthResult []AuthResult `xml:",any"`
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

type FeedbackReport struct {
	XMLName         xml.Name        `xml:"feedback"`
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Record          []Record        `xml:"record"`
}

func ParseReader(xmlFileReader io.Reader) FeedbackReport {
	var f FeedbackReport

	decoder := xml.NewDecoder(xmlFileReader)
	var inElement string

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "feedback" {
				decoder.DecodeElement(&f, &se)
				//				xmlerr := decoder.DecodeElement(&f, &se)
				//				if xmlerr != nil {
				//					fmt.Printf("decode error: %v\n", xmlerr)
				//				}
				//				fmt.Printf("XMLName: %#v\n", f)
			}
		default:
		}
	}
	return f
}
