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

package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/huhnsystems/godmarc/internal/model"
	"github.com/huhnsystems/godmarc/internal/parser"
)

// ReportLoader handles loading DMARC reports from the filesystem
type ReportLoader struct {
	ConfigDir string
}

// ErrNoReports is returned when no reports are found
var ErrNoReports = errors.New("no DMARC reports found")

// NewReportLoader creates a new ReportLoader instance with the default config directory
func NewReportLoader() (*ReportLoader, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	configDir := filepath.Join(homedir, ".godmarc")

	// Ensure config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			return nil, fmt.Errorf(
				"failed to create config directory %s: %w",
				configDir,
				err,
			)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error accessing config directory %s: %w", configDir, err)
	}

	return &ReportLoader{
		ConfigDir: configDir,
	}, nil
}

// ParseResult represents the result of parsing a single file
type ParseResult struct {
	Filename string
	Report   model.DMARCReport
	Error    error
}

// LoadReports loads all DMARC reports from the config directory
func (l *ReportLoader) LoadReports() ([]model.DMARCReport, error) {
	files, err := os.ReadDir(l.ConfigDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read directory %s: %w",
			l.ConfigDir,
			err,
		)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("%w in %s", ErrNoReports, l.ConfigDir)
	}

	var reports []model.DMARCReport
	var parseErrors []string

	// Track success and failure counts
	successCount := 0
	failureCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Get the filename and validate it
		filename := file.Name()

		// Check for suspicious filenames (path traversal attempts)
		cleanPath := filepath.Clean(filename)
		if cleanPath != filename || filepath.IsAbs(cleanPath) ||
			strings.Contains(cleanPath, "..") ||
			strings.Contains(cleanPath, "/") ||
			strings.Contains(cleanPath, "\\") {
			parseErrors = append(
				parseErrors,
				fmt.Sprintf(
					"Security warning: skipping suspicious filename: %s",
					filename,
				),
			)
			failureCount++
			continue
		}

		// Join the path with validated filename
		filePath := filepath.Join(l.ConfigDir, filename)
		report, err := parser.ParseDMARCReport(filePath)
		if err != nil {
			parseErrors = append(
				parseErrors,
				fmt.Sprintf("Error parsing %s: %v", filename, err),
			)
			failureCount++
			continue
		}

		reports = append(reports, report)
		successCount++
	}

	// If we didn't successfully parse any reports, return an error
	if len(reports) == 0 {
		if len(parseErrors) > 0 {
			// Combine all parsing errors into one message
			return nil, fmt.Errorf(
				"failed to parse any reports: %s",
				strings.Join(parseErrors, "; "),
			)
		}
		return nil, ErrNoReports
	}

	// If we had some parsing errors but also some successes, log the errors but return the reports
	if len(parseErrors) > 0 {
		fmt.Printf(
			"Warning: %d of %d files failed to parse\n",
			failureCount,
			successCount+failureCount,
		)
		for _, errMsg := range parseErrors {
			fmt.Println(errMsg)
		}
	}

	return reports, nil
}

// SortReportsByDate sorts reports by date (newest first)
func SortReportsByDate(reports []model.DMARCReport) {
	sort.Slice(reports, func(i, j int) bool {
		// Sort in descending order (newest first)
		return reports[i].ReportMetadata.DateRange.Begin.After(
			reports[j].ReportMetadata.DateRange.Begin,
		)
	})
}
