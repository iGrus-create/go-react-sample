package controller

import (
	"net/http"
	"practice/model"
	"practice/usecase"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ITaskController interface {
	GetAllTask(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	tu usecase.ITaskUsecase
}

func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	return &taskController{tu}
}

// タスク全件取得
func (tc *taskController) GetAllTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	taskRes, err := tc.tu.GetAllTask(uint(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, taskRes)
}

// タスクを一件取得
func (tc *taskController) GetTaskById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	// リクエストパラメータからタスクIDを取得
	id := c.Param("task_id")
	// 文字列を数値に変換
	taskId, _ := strconv.Atoi(id)

	taskRes, err := tc.tu.GetTaskById(uint(userId), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

// タスクを作成
func (tc *taskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task.UserId = uint(userId)
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, taskRes)
}

// タスクを更新
func (tc *taskController) UpdateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	id := c.Param("task_id")
	taskId, _ := strconv.Atoi(id)

	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	taskRes, err := tc.tu.UpdateTask(task, uint(userId), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

// タスクを削除
func (tc *taskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	id := c.Param("task_id")
	taskId, _ := strconv.Atoi(id)

	if err := tc.tu.DeleteTask(uint(userId), uint(taskId)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}