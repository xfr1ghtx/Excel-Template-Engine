package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/models"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/services"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
)

// ActHandler handles HTTP requests for acts
type ActHandler struct {
	service services.ActService
	config  *config.Config
}

// NewActHandler creates a new ActHandler
func NewActHandler(service services.ActService, cfg *config.Config) *ActHandler {
	return &ActHandler{
		service: service,
		config:  cfg,
	}
}

// CreateAct handles POST /api/act/create
func (h *ActHandler) CreateAct(c *gin.Context) {
	utils.LogMethodInit("ActHandler.CreateAct")
	utils.LogInfo("Received request to create act from IP: %s", c.ClientIP())

	var act models.Act

	// Parse JSON request body
	if err := c.ShouldBindJSON(&act); err != nil {
		utils.LogError("Error binding JSON: %v", err)
		utils.LogMethodError("ActHandler.CreateAct", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate that BigAct exists
	if act.BigAct == nil {
		utils.LogError("BigAct is required but not provided")
		utils.RespondWithError(c, http.StatusBadRequest, "BigAct is required")
		return
	}

	// Create act
	id, err := h.service.CreateAct(c.Request.Context(), &act)
	if err != nil {
		utils.LogMethodError("ActHandler.CreateAct", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create act")
		return
	}

	// Return success response
	utils.LogInfo("Successfully created act via API, ID: %s", id)
	utils.LogMethodSuccess("ActHandler.CreateAct")
	utils.RespondWithJSON(c, http.StatusCreated, gin.H{
		"id": id,
	})
}

// GenerateAct handles GET /api/act/generate?id=xxx
func (h *ActHandler) GenerateAct(c *gin.Context) {
	utils.LogMethodInit("ActHandler.GenerateAct")

	// Get ID from query parameter
	actID := c.Query("id")
	if actID == "" {
		utils.LogError("ID parameter is missing in request")
		utils.RespondWithError(c, http.StatusBadRequest, "ID parameter is required")
		return
	}

	utils.LogInfo("Received request to generate act with ID: %s from IP: %s", actID, c.ClientIP())

	// Generate act
	downloadLink, err := h.service.GenerateAct(c.Request.Context(), actID)
	if err != nil {
		utils.LogMethodError("ActHandler.GenerateAct", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate act")
		return
	}

	// Return download link
	utils.LogInfo("Successfully generated act via API, download link: %s", downloadLink)
	utils.LogMethodSuccess("ActHandler.GenerateAct")
	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"downloadLink": downloadLink,
	})
}

// DownloadAct handles GET /api/act/download/:filename
func (h *ActHandler) DownloadAct(c *gin.Context) {
	utils.LogMethodInit("ActHandler.DownloadAct")

	// Get filename from URL parameter
	filename := c.Param("filename")
	if filename == "" {
		utils.LogError("Filename is missing in request")
		utils.RespondWithError(c, http.StatusBadRequest, "Filename is required")
		return
	}

	utils.LogInfo("Received request to download file: %s from IP: %s", filename, c.ClientIP())

	// Construct file path
	filePath := filepath.Join(h.config.GeneratedPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.LogError("File not found: %s", filePath)
		utils.LogMethodError("ActHandler.DownloadAct", err)
		utils.RespondWithError(c, http.StatusNotFound, "File not found")
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	utils.LogInfo("Sending file to client: %s", filename)

	// Send file
	c.File(filePath)

	utils.LogMethodSuccess("ActHandler.DownloadAct")
}
