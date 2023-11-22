package controllers

import (
	"net/http"
	"notes_application/initializers"
	"notes_application/models"

	"github.com/gin-gonic/gin"
)

// Note structure to represent a note
type Note struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

// GetAllNotes retrieves all notes for the authenticated user
func GetAllNotes(c *gin.Context) {
	// Fetch the authenticated user's ID from the middleware
	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
		return
	}

	var notes []models.Note
	result := initializers.DB.Where("owner_id = ?", userID).Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	var responseNotes []Note
	for _, note := range notes {
		responseNotes = append(responseNotes, Note{
			ID:   note.NoteID,
			Note: note.Note,
		})
	}

	c.JSON(http.StatusOK, gin.H{"notes": responseNotes})
}

// CreateNote creates a new note for the authenticated user
func CreateNote(c *gin.Context) {
	// Fetch the authenticated user's ID from the middleware
	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
		return
	}

	var requestBody struct {
		Note string `json:"note" binding:"required"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	newNote := models.Note{
		Note:    requestBody.Note,
		OwnerID: userID.(uint),
	}

	result := initializers.DB.Create(&newNote)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newNote.NoteID})
}

// DeleteNote deletes a note for the authenticated user
func DeleteNote(c *gin.Context) {
	// Fetch the authenticated user's ID from the middleware
	userID, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
		return
	}

	var requestBody struct {
		ID uint32 `json:"id" binding:"required"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if the note with the given ID belongs to the authenticated user
	result := initializers.DB.Where("note_id = ? AND owner_id = ?", requestBody.ID, userID).Delete(&models.Note{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
