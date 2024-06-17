package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
)

// Config struct to match the YAML file structure
type Config struct {
	RecipientName          string `yaml:"recipient_name"`
	RoleTitle              string `yaml:"role_title"`
	CompanyName            string `yaml:"company_name"`
	JobBoard               string `yaml:"job_board"`
	CompanySpecificMention string `yaml:"company_specific_mention"`
	YourName               string `yaml:"your_name"`
}

// Function to read the YAML configuration file
func readConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Function to read the template and replace placeholders
func generateCoverLetter(templateFile string, replacements map[string]string) (string, error) {
	// Open the template file
	file, err := os.Open(templateFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	scanner := bufio.NewScanner(file)
	var content strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		content.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Replace placeholders with actual values
	result := content.String()
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result, nil
}

func main() {
	// Read configuration from YAML file
	configFile := "config/input.yaml"
	config, err := readConfig(configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Define the placeholders and their replacements from the config
	replacements := map[string]string{
		"[Recipient's Name]":            config.RecipientName,
		"[role title]":                  config.RoleTitle,
		"[Company's Name]":              config.CompanyName,
		"[Job Board/Company's Website]": config.JobBoard,
		"[company_specific_mention]":    config.CompanySpecificMention,
		"[Your Name]":                   config.YourName,
	}

	// Generate the cover letter
	templateFile := "templates/draft-1.txt"
	coverLetter, err := generateCoverLetter(templateFile, replacements)
	if err != nil {
		fmt.Println("Error generating cover letter:", err)
		return
	}

	// Prepare the output file path
	outputDir := "generated"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			fmt.Println("Error creating output directory:", err)
			return
		}
	}

	// Generate PDF
	outputPDF := fmt.Sprintf("%s/cover_letter_%s.pdf", outputDir, config.CompanyName)
	err = generatePDF(coverLetter, outputPDF)
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}

	fmt.Printf("Cover letter generated successfully: %s\n", outputPDF)
}

// Function to generate a PDF from the cover letter content
// Function to generate a PDF from the cover letter content
func generatePDF(content, outputFile string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Add a basic font
	fontPath := "./assets/Times_New_Roman.ttf"
	err := pdf.AddTTFFont("times_new_roman", fontPath)
	if err != nil {
		return fmt.Errorf("could not add font from %s: %v", fontPath, err)
	}

	// Set the font. The font name here must match the one used in AddTTFFont
	err = pdf.SetFont("times_new_roman", "", 14)
	if err != nil {
		return fmt.Errorf("could not set font: %v", err)
	}

	// Set the margins
	pdf.SetMargins(20, 30, 20, 30)

	// Set the starting position
	pdf.SetX(20)
	pdf.SetY(30)

	// Define a width for the MultiCell. Use the full page width minus margins.
	pageWidth := 595.28
	// pageWidth := pdf.GetX()
	// pageWidth, _ := gopdf.PageSizeA4.Width()
	margin := 20.0
	width := pageWidth - 2*margin

	// Assume a reasonable height for the content cell
	height := pdf.GetY() + 200.0 // Adjust the height as needed

	// Write content to PDF using MultiCell with specified width and height
	rect := &gopdf.Rect{W: width, H: height}
	fmt.Printf("Writing content to PDF with width: %f and height: %f\n", rect.W, rect.H)
	err = pdf.MultiCell(rect, content)
	if err != nil {
		return fmt.Errorf("could not write content to pdf: %v", err)
	}

	// Save the PDF
	err = pdf.WritePdf(outputFile)
	if err != nil {
		return fmt.Errorf("could not write pdf: %v", err)
	}

	return nil
}
