package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbpkg "github.com/loomnetwork/dashboard/db"
	"github.com/loomnetwork/dashboard/helper"
	"github.com/loomnetwork/dashboard/models"
	"github.com/loomnetwork/dashboard/version"

	"github.com/gin-gonic/gin"
)

func GetApplications(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Application{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	applications := []models.Application{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Application{}, fields)

	if err := db.Select(queryFields).Find(&applications).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(applications) > 0 {
		index = int(applications[len(applications)-1].ID)
	}

	if err := parameter.SetHeaderLink(c, index); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, application := range applications {
			fieldMap, err := helper.FieldToMap(application, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if err := enc.Encode(fieldMap); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		fieldMaps := []map[string]interface{}{}

		for _, application := range applications {
			fieldMap, err := helper.FieldToMap(application, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMaps = append(fieldMaps, fieldMap)
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(200, fieldMaps)
		} else {
			c.JSON(200, fieldMaps)
		}
	}
}

func GetApplication(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Application{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	application := models.Application{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Application{}, fields)

	if err := db.Select(queryFields).First(&application, id).Error; err != nil {
		content := gin.H{"error": "application with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(application, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func CreateApplication(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	application := models.Application{}

	if err := c.Bind(&application); err != nil {
		switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
		case gin.MIMEHTML:
			c.HTML(200, "TODO rebind the form with errors", nil)
		case gin.MIMEJSON:
			c.JSON(400, gin.H{"error": err.Error()})
		}
		return
	}

	if err := db.Create(&application).Error; err != nil {
		switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
		case gin.MIMEHTML:
			c.HTML(200, "TODO rebind the form with errors", nil)
		case gin.MIMEJSON:
			c.JSON(400, gin.H{"error": err.Error()})
		}

		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		c.Redirect(301, fmt.Sprintf("/dashboard/%d", application.ID))
	case gin.MIMEJSON:
		c.JSON(201, application)
	}
}

func UpdateApplication(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	application := models.Application{}

	if db.First(&application, id).Error != nil {
		content := gin.H{"error": "application with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&application); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&application).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, application)
}

func DeleteApplication(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	application := models.Application{}

	if db.First(&application, id).Error != nil {
		content := gin.H{"error": "application with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if err := db.Delete(&application).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
