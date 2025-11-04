package services

import (
	"fmt"
	"regexp"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/models"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
	"github.com/xuri/excelize/v2"
)

// ExcelService defines the interface for Excel operations
type ExcelService interface {
	GenerateAct(act *models.Act, outputPath string) error
}

// excelService implements ExcelService
type excelService struct {
	config *config.Config
}

// NewExcelService creates a new ExcelService
func NewExcelService(cfg *config.Config) ExcelService {
	return &excelService{
		config: cfg,
	}
}

// GenerateAct generates an Excel file from an act using the template
func (s *excelService) GenerateAct(act *models.Act, outputPath string) error {
	utils.LogMethodInit("ExcelService.GenerateAct")
	utils.LogExcelInit(outputPath)

	// Open the template file
	utils.LogInfo("Opening Excel template: %s", s.config.TemplatePath)
	f, err := excelize.OpenFile(s.config.TemplatePath)
	if err != nil {
		utils.LogMethodError("ExcelService.GenerateAct", err)
		return fmt.Errorf("failed to open template: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			utils.LogError("Error closing Excel file: %v", closeErr)
		}
	}()

	// Build template data
	utils.LogDebug("Building template data for act: %s", act.ID.Hex())
	templateData := s.buildTemplateData(act)

	// Process all sheets
	sheets := f.GetSheetList()
	utils.LogInfo("Processing %d sheets in Excel template", len(sheets))
	for _, sheetName := range sheets {
		utils.LogDebug("Processing sheet: %s", sheetName)
		err = s.processSheet(f, sheetName, templateData)
		if err != nil {
			utils.LogError("Error processing sheet %s: %v", sheetName, err)
			utils.LogMethodError("ExcelService.GenerateAct", err)
			return fmt.Errorf("failed to process sheet %s: %w", sheetName, err)
		}
	}

	// Save the file
	utils.LogInfo("Saving Excel file to: %s", outputPath)
	err = f.SaveAs(outputPath)
	if err != nil {
		utils.LogMethodError("ExcelService.GenerateAct", err)
		return fmt.Errorf("failed to save file: %w", err)
	}

	utils.LogExcelComplete(outputPath)
	utils.LogMethodSuccess("ExcelService.GenerateAct")
	return nil
}

// processSheet processes a single sheet, replacing all placeholders
func (s *excelService) processSheet(f *excelize.File, sheetName string, data map[string]interface{}) error {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	// Pattern to match {{key}}
	pattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	// Iterate through all rows and columns
	for rowIdx, row := range rows {
		for colIdx := range row {
			var cellName string
			cellName, err = excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				continue
			}

			var cellValue string
			cellValue, err = f.GetCellValue(sheetName, cellName)
			if err != nil {
				continue
			}

			// Find all matches in the cell
			if pattern.MatchString(cellValue) {
				newValue := pattern.ReplaceAllStringFunc(cellValue, func(match string) string {
					// Extract key from {{key}}
					key := pattern.FindStringSubmatch(match)[1]

					// Get value from data
					if value, ok := data[key]; ok {
						return s.formatValue(value)
					}
					return match // Keep original if not found
				})

				// Set the new value
				err = f.SetCellValue(sheetName, cellName, newValue)
				if err != nil {
					utils.LogError("Error setting cell value at %s: %v", cellName, err)
				}
			}
		}
	}

	return nil
}

// buildTemplateData builds a map of all data that can be used in the template
func (s *excelService) buildTemplateData(act *models.Act) map[string]interface{} {
	data := make(map[string]interface{})

	// Add BigAct data if present
	if act.BigAct != nil {
		// Add numeric values
		data["totalCost"] = act.BigAct.TotalCost
		data["totalCostInspection"] = act.BigAct.TotalCostInspection
		data["totalCostConsiderations"] = act.BigAct.TotalCostConsiderations
		data["positionIds"] = act.BigAct.PositionIDs

		// Add text fields
		if act.BigAct.TextFields != nil {
			for key, value := range act.BigAct.TextFields {
				data[key] = value
			}
		}
	}

	// Add timestamps
	data["createdAt"] = act.CreatedAt.Format("02.01.2006")
	data["updatedAt"] = act.UpdatedAt.Format("02.01.2006")

	// Add act ID
	data["actId"] = act.ID.Hex()

	return data
}

// formatValue formats a value based on its type
func (s *excelService) formatValue(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return utils.FormatNumber(v)
	case float32:
		return utils.FormatNumber(float64(v))
	case int:
		return utils.FormatNumber(float64(v))
	case int64:
		return utils.FormatNumber(float64(v))
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
