package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "pocketpilot/internal/models"
    "pocketpilot/internal/services"
    "pocketpilot/internal/utils"
)

type ExpenseHandler struct {
    expenseService *services.ExpenseService
}

func NewExpenseHandler(expenseService *services.ExpenseService) *ExpenseHandler {
    return &ExpenseHandler{expenseService: expenseService}
}

// CreateExpense handles expense creation
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    var req models.CreateExpenseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request data"))
        return
    }

    expense, err := h.expenseService.CreateExpense(userID.(string), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusCreated, utils.SuccessResponse("Expense created successfully", expense))
}

// GetExpenses retrieves user's expenses
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    expenses, err := h.expenseService.GetUserExpenses(userID.(string), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Expenses retrieved successfully", expenses))
}

// GetExpense retrieves a specific expense
func (h *ExpenseHandler) GetExpense(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    expenseID := c.Param("id")
    expense, err := h.expenseService.GetExpense(expenseID, userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Expense retrieved successfully", expense))
}

// UpdateExpense updates an expense
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    expenseID := c.Param("id")
    var req models.UpdateExpenseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request data"))
        return
    }

    expense, err := h.expenseService.UpdateExpense(expenseID, userID.(string), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Expense updated successfully", expense))
}

// DeleteExpense deletes an expense
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    expenseID := c.Param("id")
    err := h.expenseService.DeleteExpense(expenseID, userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Expense deleted successfully", nil))
}

// GetTeamExpenses retrieves team expenses
func (h *ExpenseHandler) GetTeamExpenses(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    teamID := c.Param("teamId")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    expenses, err := h.expenseService.GetTeamExpenses(teamID, userID.(string), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Team expenses retrieved successfully", expenses))
}