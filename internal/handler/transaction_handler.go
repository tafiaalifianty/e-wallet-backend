package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/dto"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

const (
	MIN_TOPUP_AMOUNT    = 50000
	MAX_TOPUP_AMOUNT    = 10000000
	MIN_TRANSFER_AMOUNT = 1000
	MAX_TRANSFER_AMOUNT = 50000000
)

func (h *Handler) initTransactionRoutes(api *gin.RouterGroup) {
	transaction := api.Group("/transactions")
	{
		transaction.GET("/", h.GetTransactionsByWalletNumber)
		transaction.POST("/topup", h.Topup)
		transaction.POST("/transfer", h.Transfer)
	}
}

func (h *Handler) GetTransactionsByWalletNumber(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			custom_error.FailedToGetInfoFromToken{}.Error(),
			nil,
		)
		return
	}
	tokenizedUser := user.(*entity.TokenizedUser)

	search := ctx.DefaultQuery("s", "")
	sortBy := ctx.DefaultQuery("sortBy", "datetime")
	sortMethod := ctx.DefaultQuery("sort", "desc")
	limit, err1 := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, err2 := strconv.Atoi(ctx.DefaultQuery("page", "1"))

	if err1 != nil || err2 != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			nil,
		)
		return
	}

	if limit < 1 || page < 1 {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			nil,
		)
		return
	}

	pagination := &entity.Pagination{
		Limit:  limit,
		Page:   page,
		Search: search,
		SortBy: sortBy,
		Sort:   sortMethod,
	}

	transactions, pagination, err := h.services.Transaction.FindByWalletNumber(
		tokenizedUser.WalletNumber,
		pagination,
	)

	if _, ok := err.(*custom_error.NoDataFound); ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusNotFound,
			err.Error(),
			nil,
		)

		return
	}

	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil,
		)
		return
	}

	resFormatted := dto.FormatGetTransactionsByWalletNumberResponseBody(
		transactions,
		pagination,
		tokenizedUser.WalletNumber,
	)

	helper.WriteSuccessResponse(
		ctx,
		http.StatusOK,
		http.StatusText(http.StatusOK),
		resFormatted,
	)
}

func (h *Handler) Transfer(ctx *gin.Context) {
	var input dto.TransferRequestBody
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.InvalidRequestBody{}.Error(),
			nil,
		)
		return
	}

	if !helper.IsBetweenRange(
		input.Amount,
		MIN_TRANSFER_AMOUNT,
		MAX_TRANSFER_AMOUNT,
	) {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.AmountNotInRange{
				Minimum: MIN_TRANSFER_AMOUNT,
				Maximum: MAX_TRANSFER_AMOUNT,
			}.Error(),
			nil,
		)
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			custom_error.FailedToGetInfoFromToken{}.Error(),
			nil,
		)
		return
	}
	tokenizedUser := user.(*entity.TokenizedUser)

	if tokenizedUser.WalletNumber == input.To {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.CannotTransferToOwnWallet{}.Error(),
			nil,
		)
		return
	}

	transfer := &entity.Transaction{
		Amount:      input.Amount,
		Description: input.Description,
		Type:        entity.Transfer,
		Datetime:    time.Now(),
		From:        tokenizedUser.WalletNumber,
		To:          input.To,
	}

	res, err := h.services.Transaction.CreateTransaction(transfer)

	if _, ok := err.(*custom_error.NoDataFound); ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusNotFound,
			err.Error(),
			nil,
		)

		return
	}

	if _, ok := err.(*custom_error.InsufficientBalance); ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			err.Error(),
			nil,
		)

		return
	}

	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil,
		)
		return
	}

	resFormatted := dto.FormatGetTransaction(res, tokenizedUser.WalletNumber)

	helper.WriteSuccessResponse(
		ctx,
		http.StatusOK,
		http.StatusText(http.StatusOK),
		resFormatted,
	)
}

func (h *Handler) Topup(ctx *gin.Context) {
	var input dto.TopupRequestBody
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.InvalidRequestBody{}.Error(),
			nil,
		)
		return
	}

	if !helper.IsBetweenRange(
		int(input.SourceID),
		int(entity.BankTransfer),
		int(entity.Cash),
	) {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.InvalidRequestBody{}.Error(),
			nil,
		)
		return
	}

	if !helper.IsBetweenRange(
		input.Amount,
		MIN_TOPUP_AMOUNT,
		MAX_TOPUP_AMOUNT,
	) {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.AmountNotInRange{
				Minimum: MIN_TOPUP_AMOUNT,
				Maximum: MAX_TOPUP_AMOUNT,
			}.Error(),
			nil,
		)
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			custom_error.FailedToGetInfoFromToken{}.Error(),
			nil,
		)
		return
	}
	tokenizedUser := user.(*entity.TokenizedUser)

	topup := &entity.Transaction{
		Amount: input.Amount,
		Description: fmt.Sprintf(
			"Top Up from %s",
			entity.SourceOfFundsID(input.SourceID).String(),
		),
		Type:     entity.TopUp,
		Datetime: time.Now(),
		SourceID: &input.SourceID,
		From:     tokenizedUser.WalletNumber,
		To:       tokenizedUser.WalletNumber,
	}

	res, err := h.services.Transaction.CreateTopup(topup)

	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil,
		)
		return
	}

	resFormatted := dto.FormatGetTransaction(res, tokenizedUser.WalletNumber)

	helper.WriteSuccessResponse(
		ctx,
		http.StatusOK,
		http.StatusText(http.StatusOK),
		resFormatted,
	)
}
