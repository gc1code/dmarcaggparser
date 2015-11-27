package main

import (
	"encoding/xml"
	"fmt"
)

func main() {

	type DateRange struct {
		Begin int `xml:"begin"`
		End   int `xml:"end"`
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

	v := FeedbackReport{}

	data := `
	<?xml version="1.0" encoding="UTF-8" ?>
	<feedback>
	  <report_metadata>
	    <org_name>google.com</org_name>
	    <email>noreply-dmarc-support@google.com</email>
	    <extra_contact_info>https://support.google.com/a/answer/2466580</extra_contact_info>
	    <report_id>11295852759969162400</report_id>
	    <date_range>
	      <begin>1448236800</begin>
	      <end>1448323199</end>
	    </date_range>
	  </report_metadata>
	  <policy_published>
	    <domain>colbs.net</domain>
	    <adkim>r</adkim>
	    <aspf>r</aspf>
	    <p>reject</p>
	    <sp>reject</sp>
	    <pct>100</pct>
	  </policy_published>
	  <record>
	    <row>
	      <source_ip>2607:f8b0:400c:c05::249</source_ip>
	      <count>1</count>
	      <policy_evaluated>
	        <disposition>none</disposition>
	        <dkim>pass</dkim>
	        <spf>fail</spf>
	      </policy_evaluated>
	    </row>
	    <identifiers>
	      <header_from>colbs.net</header_from>
	    </identifiers>
	    <auth_results>
	      <dkim>
	        <domain>google.com</domain>
	        <result>pass</result>
	      </dkim>
	      <dkim>
	        <domain>colbs.net</domain>
	        <result>pass</result>
	      </dkim>
	      <spf>
	        <domain>calendar-server.bounces.google.com</domain>
	        <result>pass</result>
	      </spf>
	    </auth_results>
	  </record>
	  <record>
	    <row>
	      <source_ip>2a00:1450:400c:c09::235</source_ip>
	      <count>1</count>
	      <policy_evaluated>
	        <disposition>none</disposition>
	        <dkim>pass</dkim>
	        <spf>pass</spf>
	      </policy_evaluated>
	    </row>
	    <identifiers>
	      <header_from>colbs.net</header_from>
	    </identifiers>
	    <auth_results>
	      <dkim>
	        <domain>colbs.net</domain>
	        <result>pass</result>
	      </dkim>
	      <spf>
	        <domain>colbs.net</domain>
	        <result>pass</result>
	      </spf>
	    </auth_results>
	  </record>
	</feedback>
	`

	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("XMLName: %#v\n", v)
	//  fmt.Printf("org_name: %#q\n", v.ReportMetadata.OrgName)
	//  fmt.Printf("domain: %#q\n", v.PolicyPublished.Domain)
	//	fmt.Printf("Name: %q\n", v.Name)
	//	fmt.Printf("Phone: %q\n", v.Phone)
	//	fmt.Printf("Email: %v\n", v.Email)
	//	fmt.Printf("Groups: %v\n", v.Groups)
	//	fmt.Printf("Address: %v\n", v.Address)
}
