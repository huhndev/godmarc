package storage

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/huhndev/godmarc/model"
	"github.com/huhndev/godmarc/parser"
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

	successCount := 0
	failureCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

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

		filePath := filepath.Join(l.ConfigDir, filename)

		// Handle gzipped files
		if strings.HasSuffix(strings.ToLower(filename), ".gz") {
			xmlPath, err := decompressGzip(filePath)
			if err != nil {
				parseErrors = append(
					parseErrors,
					fmt.Sprintf("Error decompressing %s: %v", filename, err),
				)
				failureCount++
				continue
			}
			defer os.Remove(xmlPath)
			filePath = xmlPath
		} else if !strings.HasSuffix(strings.ToLower(filename), ".xml") {
			continue
		}

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

	if len(reports) == 0 {
		if len(parseErrors) > 0 {
			return nil, fmt.Errorf(
				"failed to parse any reports: %s",
				strings.Join(parseErrors, "; "),
			)
		}
		return nil, ErrNoReports
	}

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

// decompressGzip decompresses a gzipped file and returns the path to a temporary XML file
func decompressGzip(gzPath string) (string, error) {
	f, err := os.Open(gzPath)
	if err != nil {
		return "", fmt.Errorf("could not open gzip file: %w", err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer gr.Close()

	tmpFile, err := os.CreateTemp("", "godmarc-*.xml")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %w", err)
	}

	// Limit decompressed size to 50MB to prevent decompression bombs
	limited := io.LimitReader(gr, 50*1024*1024)
	if _, err := io.Copy(tmpFile, limited); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("could not decompress file: %w", err)
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}

// SortReportsByDate sorts reports by date (newest first)
func SortReportsByDate(reports []model.DMARCReport) {
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ReportMetadata.DateRange.Begin.After(
			reports[j].ReportMetadata.DateRange.Begin,
		)
	})
}
