package projects

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	basic_controller "task-management/internal/controller/http/v1/_basic_controller"
	"task-management/internal/service/projects"
	project_usecase "task-management/internal/usecase/projects"
)

type Controller struct {
	useCase *project_usecase.UseCase
}

func NewController(useCase *project_usecase.UseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (cl Controller) ProjectGetList(c *gin.Context) {
	var filter projects.Filter
	query := c.Request.URL.Query()

	ownerIdQ := query["owner_id"]
	if len(ownerIdQ) > 0 {
		queryInt, err := strconv.Atoi(ownerIdQ[0])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "owner_id must be integer!",
				"status":  false,
			})

			return
		}
		filter.OwnerId = &queryInt
	}

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		queryInt, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Limit must be a number",
				"status":  false,
			})
			return
		}

		filter.Limit = &queryInt
	}

	offsetQ := query["offset"]
	if len(offsetQ) > 0 {
		page, err := strconv.Atoi(offsetQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "offset must be number!",
				"status":  false,
			})
			return
		}
		offset := (page - 1) * *filter.Limit
		filter.Offset = &offset
	}

	ctx := context.Background()

	list, count, err := cl.useCase.ProjectGetList(ctx, filter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
	})
}

func (cl Controller) ProjectGetDetail(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "id must be a number!",
			"status":  false,
		})

		return
	}

	ctx := context.Background()

	detail, err := cl.useCase.ProjectGetDetail(ctx, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl Controller) ProjectCreate(c *gin.Context) {
	var data projects.Create

	err := c.ShouldBind(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	ctx := context.Background()

	detail, err := cl.useCase.ProjectCreate(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl Controller) ProjectUpdate(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Id must be a number!",
			"status":  false,
		})

		return
	}

	var data projects.Update

	err = c.ShouldBind(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	if data.Id == nil {
		data.Id = &id
	}
	ctx := context.Background()

	detail, err := cl.useCase.ProjectUpdate(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl Controller) ProjectDelete(c *gin.Context) {
	ctx, data, err := basic_controller.BasicDelete(c)
	if err != nil {
		return
	}

	err = cl.useCase.ProjectDelete(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}
