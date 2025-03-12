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

package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/huhnsystems/godmarc/internal/model"
)

// ParseDMARCReport parses a DMARC report XML file
func ParseDMARCReport(filepath string) (model.DMARCReport, error) {
	var report model.DMARCReport

	// Validate file exists
	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return report, fmt.Errorf("file not found: %s", filepath)
		}
		return report, fmt.Errorf("error accessing file %s: %w", filepath, err)
	}

	// Check file size
	if info.Size() == 0 {
		return report, fmt.Errorf("file is empty: %s", filepath)
	}

	// Read file with more context if it fails
	data, err := os.ReadFile(filepath)
	if err != nil {
		return report, fmt.Errorf("could not read file %s: %w", filepath, err)
	}

	// Check if data seems to be XML
	if !hasXMLHeader(data) && !hasRootElement(data) {
		return report, fmt.Errorf(
			"file %s does not appear to be valid XML",
			filepath,
		)
	}

	// Handle possible XML declaration
	data = []byte(strings.TrimSpace(string(data)))

	// Create a secure XML decoder with entity expansion disabled
	decoder := xml.NewDecoder(bytes.NewReader(data))

	// Disable entity expansion to prevent XXE attacks
	decoder.Entity = xml.HTMLEntity

	// Impose reasonable limits to prevent billion laughs attack
	decoder.Strict = true

	// Parse XML using the secure decoder
	err = decoder.Decode(&report)
	if err != nil {
		// Add more context to XML parsing errors
		return report, fmt.Errorf("invalid XML in file %s: %w", filepath, err)
	}

	// Validate required fields
	if err := validateReport(report); err != nil {
		return report, fmt.Errorf(
			"invalid DMARC report in file %s: %w",
			filepath,
			err,
		)
	}

	return report, nil
}

// hasXMLHeader checks if the data starts with an XML declaration
func hasXMLHeader(data []byte) bool {
	return bytes.HasPrefix(bytes.TrimSpace(data), []byte("<?xml"))
}

// hasRootElement checks if the data contains what appears to be an XML root element
func hasRootElement(data []byte) bool {
	s := string(bytes.TrimSpace(data))
	return strings.Contains(s, "<feedback") || strings.Contains(s, "<report")
}

// validateReport checks that essential fields are present
func validateReport(report model.DMARCReport) error {
	if report.ReportMetadata.ReportID == "" {
		return fmt.Errorf("missing report ID")
	}

	if report.ReportMetadata.OrgName == "" {
		return fmt.Errorf("missing organization name")
	}

	if report.PolicyPublished.Domain == "" {
		return fmt.Errorf("missing domain in policy")
	}

	// Make sure we have a valid date range
	zeroTime := time.Time{}
	if report.ReportMetadata.DateRange.Begin == zeroTime ||
		report.ReportMetadata.DateRange.End == zeroTime {
		return fmt.Errorf("invalid date range")
	}

	return nil
}
