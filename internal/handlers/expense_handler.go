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

// @Summary Create expense
// @Description Create a new expense
// @Tags Expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param expense body models.CreateExpenseRequest true "Expense payload"
// @Success 201 {object} models.Expense
// @Failure 400 {object} models.ErrorResponse
// @Router /api/expenses [post]
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

// @Summary Get expenses
// @Description Retrieve user's expenses
// @Tags Expenses
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Expense
// @Router /api/expenses [get]
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

// @Summary Get expense
// @Description Retrieve a specific expense by ID
// @Tags Expenses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expense ID"
// @Success 200 {object} models.Expense
// @Failure 404 {object} models.ErrorResponse
// @Router /api/expenses/{id} [get]
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

// @Summary Update expense
// @Description Update an existing expense
// @Tags Expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expense ID"
// @Param expense body models.UpdateExpenseRequest true "Expense payload"
// @Success 200 {object} models.Expense
// @Failure 400 {object} models.ErrorResponse
// @Router /api/expenses/{id} [put]
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

// @Summary Delete expense
// @Description Delete an expense
// @Tags Expenses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expense ID"
// @Success 204
// @Failure 404 {object} models.ErrorResponse
// @Router /api/expenses/{id} [delete]
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

// @Summary Get team expenses
// @Description Retrieve expenses for a team
// @Tags Expenses
// @Produce json
// @Security BearerAuth
// @Param teamId path string true "Team ID"
// @Success 200 {array} models.Expense
// @Router /api/expenses/team/{teamId} [get]
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