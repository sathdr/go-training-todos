package todos

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/pallat/todos/logger"
	"gorm.io/gorm"
)

// NewNewTodoHandler return handler using provided DB
func NewNewTodoHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		db.AutoMigrate(Task{})

		var todo struct {
			Task string `json:"task"`
		}

		logger := logger.Extract(c)
		logger.Info("new task ...")

		if error := c.Bind(&todo); error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(error, "New Task").Error(),
			})
		}

		var task = Task{
			Task: todo.Task,
		}

		if result := db.Create(&task); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "New Task").Error(),
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

// NewListTodoHandler return handler using provided DB
func NewListTodoHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("list task ...")

		var tasks = []Task{}

		if result := db.Find(&tasks); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "List Task").Error(),
			})
		}

		return c.JSON(http.StatusOK, tasks)
	}
}

// NewGetTodoHandler return handler using provided DB
func NewGetTodoHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		idStr := c.Param("id")

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("get task ..." + idStr)

		var id int
		var err error

		if id, err = strconv.Atoi(idStr); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(err, "Get Task").Error(),
			})
		}
		var task = Task{}

		if result := db.Find(&task, id); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "Get Task").Error(),
			})
		} else if result.RowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

// NewUpdateTodoHandler return handler using provided DB
func NewUpdateTodoHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		idStr := c.Param("id")

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("update task ..." + idStr)

		var id int
		var err error

		if id, err = strconv.Atoi(idStr); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(err, "Get Task").Error(),
			})
		}

		var todo struct {
			Task      string `json:"task"`
			Processed bool   `json:"processed"`
		}

		if error := c.Bind(&todo); error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(error, "Update Task").Error(),
			})
		}

		var task = Task{}

		if result := db.Find(&task, id); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "Update Task").Error(),
			})
		} else if result.RowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
		}

		// update fields
		task.Task = todo.Task
		task.Processed = todo.Processed

		if result := db.Save(&task); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "Update Task").Error(),
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

// NewDeleteTodoHandler return handler using provided DB
func NewDeleteTodoHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		idStr := c.Param("id")

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("delete task ..." + idStr)

		var id int
		var err error

		if id, err = strconv.Atoi(idStr); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(err, "Delete Task").Error(),
			})
		}
		var task = Task{}

		if result := db.Find(&task, id); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "Delete Task").Error(),
			})
		} else if result.RowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
		}

		if result := db.Delete(&task); result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(result.Error, "Delete Task").Error(),
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

// Task represent task in the database
type Task struct {
	gorm.Model
	Task      string
	Processed bool
}
