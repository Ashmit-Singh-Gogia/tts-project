package services

import (
	"fmt"
	"os/exec"
)

func ExtractTextFromPDF(pdfFilePath string) (string, error) {
	cmd := exec.Command("pdftotext", pdfFilePath, "-")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		// If it fails (e.g., password-protected PDF, or poppler isn't installed), we return the error
		return "", fmt.Errorf("exit code: %w, details: %s", err, string(outputBytes))
	}
	return string(outputBytes), nil
}
