package dmarc

import (
	"encoding/xml"
	"io"
	"fmt"
	"time"
	"strconv"
	"strings"
	"net"
)

const AuthResultType_DKIM = "dkim"
const AuthResultType_SPF  = "spf"


type DateRange struct {
	Begin time.Time
	End   time.Time
}

// Unmarshal the DateRange element into time.Time objects.
// Yahoo tends to put some extra whitespace behind its
// timestamps, so we trim that first.
func (dr *DateRange) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	rangeStr := struct {
		Begin string `xml:"begin"`
		End   string `xml:"end"`
	}{}
	d.DecodeElement(&rangeStr, &start)

	begin, err := strconv.Atoi(strings.TrimSpace(rangeStr.Begin))
	if err != nil {
		return err
	}
	end, err := strconv.Atoi(strings.TrimSpace(rangeStr.End))
	if err != nil {
		return err
	}

	dr.Begin = time.Unix(int64(begin), 0).UTC()
	dr.End = time.Unix(int64(end), 0).UTC()
	return nil
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
	SourceIp        net.IP
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

func (r *Row) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	rowAlias := struct {
		SourceIp        string          `xml:"source_ip"`
		Count           int             `xml:"count"`
		PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
	}{}
	d.DecodeElement(&rowAlias, &start)

	r.Count = rowAlias.Count
	r.PolicyEvaluated = rowAlias.PolicyEvaluated
	r.SourceIp = net.ParseIP(rowAlias.SourceIp)
	if r.SourceIp == nil {
		return fmt.Errorf("Could not parse source_ip")
	}

	return nil
}

type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

type AuthResult struct {
	Type    string // Either SPF or DKIM
	Domain  string `xml:"domain"`
	Result  string `xml:"result"`
}

func (ar *AuthResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	res := struct {
		XMLName xml.Name
		Domain  string `xml:"domain"`
		Result  string `xml:"result"`
	}{}
	d.DecodeElement(&res, &start)

	ar.Domain = res.Domain
	ar.Result = res.Result

	if res.XMLName.Local != AuthResultType_DKIM && res.XMLName.Local != AuthResultType_SPF {
		return fmt.Errorf("Unrecognized AuthResult type: %s", res.XMLName.Local)
	}

	ar.Type = res.XMLName.Local
	return nil
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
	Metadata        ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Record          []Record        `xml:"record"`
}

func ParseReader(xmlFileReader io.Reader) (*FeedbackReport, error) {
	var f *FeedbackReport

	decoder := xml.NewDecoder(xmlFileReader)
	var inElement string

	for {
		t, err := decoder.Token()
		if t == nil && err == io.EOF {
			break;
		}
		if t == nil {
			return nil, err
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "feedback" {
				err := decoder.DecodeElement(&f, &se)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("Unknown root element: %s", inElement)
			}
		default:
		}
	}

	return f, nil
}
