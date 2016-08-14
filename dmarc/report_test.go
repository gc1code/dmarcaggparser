package dmarc

import (
	"testing"
	"strings"
	//"fmt"
	"time"
)

var reportXML = `<?xml version="1.0"?>
<feedback>	
  <report_metadata>	
    <org_name>Yahoo! Inc.</org_name>	
    <email>postmaster@dmarc.yahoo.com</email>	
    <report_id>1405588706.570207</report_id>	
    <date_range>	
      <begin>1405468800</begin>	
      <end>1405555199 </end>	
    </date_range>	
  </report_metadata>	
  <policy_published>	
    <domain>example.com</domain>	
    <adkim>r</adkim>	
    <aspf>r</aspf>	
    <p>none</p>	
    <pct>100</pct>	
  </policy_published>	
  <record>	
    <row>	
      <source_ip>127.0.0.1</source_ip>
      <count>2</count>	
      <policy_evaluated>	
        <disposition>none</disposition>	
        <dkim>fail</dkim>	
        <spf>pass</spf>	
      </policy_evaluated>	
    </row>	
    <identifiers>	
      <header_from>example.com</header_from>	
    </identifiers>	
    <auth_results>	
      <dkim>	
        <domain>example.com</domain>	
        <result>neutral</result>	
      </dkim>	
      <spf>	
        <domain>example.com</domain>	
        <result>pass</result>	
      </spf>	
    </auth_results>	
  </record>	
  <record>	
    <row>	
      <source_ip>127.0.0.2</source_ip>
      <count>988</count>	
      <policy_evaluated>	
        <disposition>none</disposition>	
        <dkim>fail</dkim>	
        <spf>fail</spf>	
      </policy_evaluated>	
    </row>	
    <identifiers>	
      <header_from>idfactor.example.com</header_from>	
    </identifiers>	
    <auth_results>	
      <dkim>	
        <domain>idfactor.example.com</domain>	
        <result>neutral</result>	
      </dkim>	
      <spf>	
        <domain>idfactor.example.com</domain>	
        <result>none</result>	
      </spf>	
    </auth_results>	
  </record>
</feedback>	`

func TestInvalidDocument(t *testing.T) {
	xml := `<?xml version="1.0"?><foobar></foobar>`
	report, err := ParseReader(strings.NewReader(xml))

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	if report != nil {
		t.Errorf("Report should have been <nil> but wasn't")
	}
}

func TestUnmarshalling(t *testing.T) {
	report, err := ParseReader(strings.NewReader(reportXML))

	if err != nil {
		t.Errorf("Got error: %s:", err.Error())
	}

	begin := time.Unix(1405468800, 0).UTC()
	if report.Metadata.DateRange.Begin != begin {
		t.Errorf("report.Metadata.DateRange.Begin did not match expected date")
	}

	end := time.Unix(1405555199, 0).UTC()
	if report.Metadata.DateRange.End != end {
		t.Errorf("report.Metadata.DateRange.end did not match expected date")
	}

	if report.PolicyPublished.Domain != "example.com" {
		t.Errorf("report.PolicyPublished.Domain was '%s', expected 'example.com'", report.PolicyPublished.Domain)
	}

	if len(report.Record) != 2 {
		t.Errorf("report.Record was expected to contain 2 items, only got: %d", len(report.Record))
	}
}