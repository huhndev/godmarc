//Copyright (c) 2025, Julian Huhn
//
//Permission to use, copy, modify, and/or distribute this software for any
//purpose with or without fee is hereby granted, provided that the above
//copyright notice and this permission notice appear in all copies.
//
//THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
//WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
//MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
//ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
//WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
//ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
//OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

// DateRange represents a time range with begin and end times
type DateRange struct {
	Begin time.Time
	End   time.Time
}

// DMARCReport represents a parsed DMARC report
type DMARCReport struct {
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"`
}

// ReportMetadata contains metadata about the DMARC report
type ReportMetadata struct {
	OrgName      string    `xml:"org_name"`
	Email        string    `xml:"email"`
	ExtraContact string    `xml:"extra_contact_info"`
	ReportID     string    `xml:"report_id"`
	DateRange    DateRange `xml:"date_range"`
}

// PolicyPublished contains the published DMARC policy
type PolicyPublished struct {
	Domain string `xml:"domain"`
	ADKIM  string `xml:"adkim"`
	ASPF   string `xml:"aspf"`
	P      string `xml:"p"`
	SP     string `xml:"sp"`
	PCT    int    `xml:"pct"`
}

// Record represents a single DMARC record
type Record struct {
	Row struct {
		SourceIP        string `xml:"source_ip"`
		Count           int    `xml:"count"`
		PolicyEvaluated struct {
			Disposition string `xml:"disposition"`
			DKIM        string `xml:"dkim"`
			SPF         string `xml:"spf"`
		} `xml:"policy_evaluated"`
	} `xml:"row"`
	Identifiers struct {
		HeaderFrom string `xml:"header_from"`
	} `xml:"identifiers"`
	AuthResults struct {
		DKIM []struct {
			Domain   string `xml:"domain"`
			Result   string `xml:"result"`
			Selector string `xml:"selector"`
		} `xml:"dkim"`
		SPF []struct {
			Domain string `xml:"domain"`
			Result string `xml:"result"`
			Scope  string `xml:"scope"`
		} `xml:"spf"`
	} `xml:"auth_results"`
}

// UnmarshalXML custom unmarshaler for date range
func (dr *DateRange) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	type dateRange struct {
		Begin string `xml:"begin"`
		End   string `xml:"end"`
	}

	var raw dateRange
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}

	// Convert begin timestamp
	beginInt, err := strconv.ParseInt(raw.Begin, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid begin timestamp: %w", err)
	}
	dr.Begin = time.Unix(beginInt, 0)

	// Convert end timestamp
	endInt, err := strconv.ParseInt(raw.End, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid end timestamp: %w", err)
	}
	dr.End = time.Unix(endInt, 0)

	return nil
}
