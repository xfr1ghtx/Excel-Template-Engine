package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/models"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/repository"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActService defines the interface for act business logic
type ActService interface {
	CreateAct(ctx context.Context, act *models.Act) (string, error)
	GenerateAct(ctx context.Context, actID string) (string, error)
}

// actService implements ActService
type actService struct {
	repo         repository.ActRepository
	excelService ExcelService
	config       *config.Config
}

// NewActService creates a new ActService
func NewActService(repo repository.ActRepository, excelService ExcelService, cfg *config.Config) ActService {
	return &actService{
		repo:         repo,
		excelService: excelService,
		config:       cfg,
	}
}

// CreateAct creates a new act in the database
func (s *actService) CreateAct(ctx context.Context, act *models.Act) (string, error) {
	utils.LogMethodInit("ActService.CreateAct")
	
	// Set timestamps
	now := time.Now()
	act.CreatedAt = now
	act.UpdatedAt = now

	// Generate IDs for positions if not set
	for i := range act.Positions {
		if act.Positions[i].ID.IsZero() {
			act.Positions[i].ID = primitive.NewObjectID()
		}
	}

	// Save to database
	id, err := s.repo.Create(ctx, act)
	if err != nil {
		utils.LogMethodError("ActService.CreateAct", err)
		return "", fmt.Errorf("failed to create act: %w", err)
	}

	utils.LogInfo("Successfully created act with ID: %s", id)
	utils.LogMethodSuccess("ActService.CreateAct")
	return id, nil
}

// GenerateAct generates an Excel file for an act
func (s *actService) GenerateAct(ctx context.Context, actID string) (string, error) {
	utils.LogMethodInit("ActService.GenerateAct")
	utils.LogInfo("Generating act for ID: %s", actID)
	
	// Fetch act from database
	act, err := s.repo.FindByID(ctx, actID)
	if err != nil {
		utils.LogMethodError("ActService.GenerateAct", err)
		return "", fmt.Errorf("act not found: %w", err)
	}

	// Check if BigAct exists
	if act.BigAct == nil {
		err := fmt.Errorf("act does not have BigAct data")
		utils.LogMethodError("ActService.GenerateAct", err)
		return "", err
	}

	// Check if bigActChanged is true
	if act.BigAct.Changed {
		// Process the act and generate new file
		return s.processAndGenerateAct(ctx, act)
	}

	// Return existing link if not changed
	if act.BigAct.BigActLink != "" {
		utils.LogInfo("Returning existing BigActLink: %s", act.BigAct.BigActLink)
		utils.LogMethodSuccess("ActService.GenerateAct")
		return act.BigAct.BigActLink, nil
	}

	// If no link exists but changed is false, generate anyway
	return s.processAndGenerateAct(ctx, act)
}

// processAndGenerateAct processes the act and generates the Excel file
func (s *actService) processAndGenerateAct(ctx context.Context, act *models.Act) (string, error) {
	utils.LogInfo("Processing and generating act: %s", act.ID.Hex())
	
	// Find positions with current period costs
	positionsWithCurrent := s.findPositionsWithCurrentPeriod(act.Positions)

	var selectedPositions []models.Position
	
	if len(positionsWithCurrent) > 0 {
		// Use positions with current period costs
		selectedPositions = positionsWithCurrent
		utils.LogDebug("Using %d positions with current period costs", len(positionsWithCurrent))
	} else {
		// Fallback to positions with accumulated cost
		positionsWithAccumulated := s.findPositionsWithAccumulated(act.Positions)
		selectedPositions = positionsWithAccumulated
		utils.LogDebug("Using %d positions with accumulated costs", len(positionsWithAccumulated))
	}

	// Calculate totals
	totalCost, totalInspection, totalConsiderations := s.calculateTotals(selectedPositions)
	
	// Update BigAct with totals
	act.BigAct.TotalCost = totalCost
	act.BigAct.TotalCostInspection = totalInspection
	act.BigAct.TotalCostConsiderations = totalConsiderations
	
	// Concatenate position IDs
	act.BigAct.PositionIDs = s.concatenatePositionIDs(selectedPositions)

	// Update act in database
	act.UpdatedAt = time.Now()
	err := s.repo.Update(ctx, act.ID.Hex(), act)
	if err != nil {
		utils.LogMethodError("ActService.processAndGenerateAct", err)
		return "", fmt.Errorf("failed to update act: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("act_%s_%d.xlsx", act.ID.Hex(), timestamp)
	outputPath := fmt.Sprintf("%s/%s", s.config.GeneratedPath, filename)

	// Generate Excel file
	err = s.excelService.GenerateAct(act, outputPath)
	if err != nil {
		utils.LogMethodError("ActService.processAndGenerateAct", err)
		return "", fmt.Errorf("failed to generate Excel: %w", err)
	}

	// Update BigActLink
	downloadLink := fmt.Sprintf("/api/act/download/%s", filename)
	act.BigAct.BigActLink = downloadLink
	act.BigAct.Changed = false // Reset changed flag

	// Update act again with the link
	err = s.repo.Update(ctx, act.ID.Hex(), act)
	if err != nil {
		utils.LogError("Error updating act with BigActLink: %v", err)
		// Don't return error here, file is already generated
	}

	utils.LogInfo("Successfully generated act with download link: %s", downloadLink)
	utils.LogMethodSuccess("ActService.GenerateAct")
	return downloadLink, nil
}

// findPositionsWithCurrentPeriod finds positions with current period costs
func (s *actService) findPositionsWithCurrentPeriod(positions []models.Position) []models.Position {
	var result []models.Position
	for _, pos := range positions {
		if pos.HasCurrentPeriodCost() {
			result = append(result, pos)
		}
	}
	return result
}

// findPositionsWithAccumulated finds positions with accumulated cost
func (s *actService) findPositionsWithAccumulated(positions []models.Position) []models.Position {
	var result []models.Position
	for _, pos := range positions {
		if pos.HasAccumulatedCost() {
			result = append(result, pos)
		}
	}
	return result
}

// concatenatePositionIDs concatenates position IDs into a comma-separated string
func (s *actService) concatenatePositionIDs(positions []models.Position) string {
	var ids []string
	for _, pos := range positions {
		ids = append(ids, pos.ID.Hex())
	}
	return strings.Join(ids, ", ")
}

// calculateTotals calculates total costs from positions
func (s *actService) calculateTotals(positions []models.Position) (float64, float64, float64) {
	var totalCost, totalInspection, totalConsiderations float64

	for _, pos := range positions {
		if pos.CurrentPeriodCost != nil {
			totalCost += *pos.CurrentPeriodCost
		}
		if pos.CurrentPeriodCostInspection != nil {
			totalInspection += *pos.CurrentPeriodCostInspection
		}
		if pos.CurrentPeriodCostConsiderations != nil {
			totalConsiderations += *pos.CurrentPeriodCostConsiderations
		}
	}

	return totalCost, totalInspection, totalConsiderations
}

