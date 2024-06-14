package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	outputFile := fmt.Sprintf("%s/cover_letter_%s.txt", outputDir, config.CompanyName)
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	// Write the generated cover letter to the file
	_, err = file.WriteString(coverLetter)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Printf("Cover letter generated successfully: %s\n", outputFile)
}
