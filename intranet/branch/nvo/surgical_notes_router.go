package nvo

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"gorm.io/gorm"
)

const (
	defaultSurgicalNotesLimit = 50
	maxSurgicalNotesLimit     = 100
)

func SurgicalNotesRegister(router *gin.RouterGroup) {
	router.GET("", SurgicalNoteList)
	router.GET("/", SurgicalNoteList)
	router.POST("", SurgicalNoteCreate)
	router.POST("/", SurgicalNoteCreate)
	router.GET("/:id", SurgicalNoteRetrieve)
	router.PUT("/:id", SurgicalNoteUpdate)
	router.DELETE("/:id", SurgicalNoteDelete)
}

func SurgicalNoteList(c *gin.Context) {
	limit := defaultSurgicalNotesLimit
	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, common.NewError("limit", errors.New("must be a positive integer")))
			return
		}
		limit = parsed
		if limit > maxSurgicalNotesLimit {
			limit = maxSurgicalNotesLimit
		}
	}

	offset := 0
	if value := c.Query("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			c.JSON(http.StatusBadRequest, common.NewError("offset", errors.New("must be a non-negative integer")))
			return
		}
		offset = parsed
	}

	db := common.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", errors.New("database is not initialized")))
		return
	}

	notes := make([]SurgicalNoteModel, 0)
	var count int64
	if err := db.Model(&SurgicalNoteModel{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	if err := db.Order("id DESC").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"surgicalNotes": notes, "surgicalNotesCount": count})
}

func SurgicalNoteCreate(c *gin.Context) {
	var note SurgicalNoteModel
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("body", err))
		return
	}

	db := common.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", errors.New("database is not initialized")))
		return
	}
	note.ID = 0
	note.DeletedAt = gorm.DeletedAt{}
	note.CreatedAt = nil
	note.UpdatedAt = nil
	if err := db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"surgicalNote": note})
}

func SurgicalNoteRetrieve(c *gin.Context) {
	id, ok := surgicalNoteID(c)
	if !ok {
		return
	}
	db := common.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", errors.New("database is not initialized")))
		return
	}

	var note SurgicalNoteModel
	if err := db.First(&note, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewError("surgicalNote", errors.New("not found")))
		} else {
			c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"surgicalNote": note})
}

func SurgicalNoteUpdate(c *gin.Context) {
	id, ok := surgicalNoteID(c)
	if !ok {
		return
	}
	var note SurgicalNoteModel
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("body", err))
		return
	}

	db := common.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", errors.New("database is not initialized")))
		return
	}

	var existing SurgicalNoteModel
	if err := db.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.NewError("surgicalNote", errors.New("not found")))
		} else {
			c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		}
		return
	}
	updates := surgicalNoteUpdateMap(note)
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"surgicalNote": existing})
		return
	}
	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"surgicalNote": existing})
}

func SurgicalNoteDelete(c *gin.Context) {
	id, ok := surgicalNoteID(c)
	if !ok {
		return
	}
	db := common.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", errors.New("database is not initialized")))
		return
	}
	result := db.Delete(&SurgicalNoteModel{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", result.Error))
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, common.NewError("surgicalNote", errors.New("not found")))
		return
	}
	c.Status(http.StatusNoContent)
}

func surgicalNoteID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, common.NewError("id", errors.New("must be a positive integer")))
		return 0, false
	}
	return id, true
}
