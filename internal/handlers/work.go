package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/helpers"
	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/bjorndonald/lasgcce/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type WorkHandler struct {
	deps       *bootstrap.AppDependencies
	pageRepo   repository.PageRepositoryInterface
	labourRepo repository.LabourRepositoryInterface
	plantRate  repository.PlantRepositoryInterface
}

func NewWorkHandler(deps *bootstrap.AppDependencies,
) *WorkHandler {
	return &WorkHandler{
		deps:       deps,
		pageRepo:   repository.NewPageRepository(deps.DatabaseService),
		labourRepo: repository.NewLabourRepository(deps.DatabaseService),
		plantRate:  repository.NewPlantRepository(deps.DatabaseService),
	}
}

type PageInput struct {
	Name      string `json:"name"`
	PageTitle string `json:"pagetitle" validate:"required"`
	SubTitle  string `json:"subtitle"`
	Section   string `json:"section"`
}

type LabourRateInput struct {
	WorkIdentifier string  `json:"work_identifier"`
	Title          string  `json:"title"`
	SubTitle       string  `json:"subtitle"`
	Labour         string  `json:"labour"`
	Unit           string  `json:"unit"`
	PriceInNaira   float32 `json:"price_in_naira"`
}

type CurrentPlantRateInput struct {
	WorkIdentifier string  `json:"work_identifier"`
	Title          string  `json:"title"`
	SubTitle       string  `json:"subtitle"`
	Equipment      string  `json:"equipment"`
	Rental         float32 `json:"rental"`
	Diesel         float32 `json:"diesel"`
	TotalPrice     float32 `json:"total_price"`
}

type MaterialPriceInput struct {
	WorkIdentifier string  `json:"work_identifier"`
	Title          string  `json:"title"`
	SubTitle       string  `json:"subtitle"`
	Material       string  `json:"material"`
	Unit           string  `json:"unit"`
	Price          float32 `json:"price"`
}

// CreatePage is a route handler that handles creating new for posting data page
//
// # This endpoint is used to create the page
//
// @Summary Create page
// @Description Creates some details about the page
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body PageInput true "create page"
// @Security BearerAuth
// @Success 201 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/page/ [post]
func (u *WorkHandler) CreatePage(c *gin.Context) {
	var input PageInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(PageInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	_, found, err := u.pageRepo.FindByCondition("name = ?", input.Name)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("page already exists"), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	page := &models.Page{
		ID:        id,
		Name:      input.Name,
		Section:   input.Section,
		SubTitle:  input.SubTitle,
		PageTitle: input.PageTitle,
		UpdatedAt: time.Now(),
	}

	err = u.pageRepo.Create(page)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Page created successfully", page, http.StatusCreated)
}

// GetPages is a route handler that handles gettings list of pages
//
// # This endpoint is used to get the pages
//
// @Summary Get page
// @Description Get pages
// @Tags Work
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/page [get]
func (u *WorkHandler) GetPages(c *gin.Context) {

	pages, err := u.pageRepo.Find()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Page created successfully", pages, http.StatusCreated)
}

// UpdatePage is a route handler that handles updating the page
//
// # This endpoint is used to update the page
//
// @Summary Update page
// @Description Updates some details about the page
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body UpdatePageInput true "update page"
// @Param id path string true "Page ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/page/{id} [put]
func (u *WorkHandler) UpdatePage(c *gin.Context) {
	var input PageInput
	id := c.Param("id")

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(PageInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	page, found, err := u.pageRepo.FindByCondition("id = ?", id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("page not found"), http.StatusNotFound)
		return
	}

	page.Name = input.Name
	page.PageTitle = input.PageTitle
	page.Section = input.Section
	page.SubTitle = input.SubTitle
	page.UpdatedAt = time.Now()

	_, err = u.pageRepo.Save(page)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Page updated successfully", page, http.StatusOK)
}

// DeletePage is a route handler that handles deleting the page
//
// # This endpoint is used to update the page
//
// @Summary Delete page
// @Description Deletes some details about the page
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "Page ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/page/{id} [delete]
func (u *WorkHandler) DeletePage(c *gin.Context) {
	id := c.Param("id")

	page, err := u.pageRepo.Delete(id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	helpers.ReturnJSON(c, "Page deleted successfully", page, http.StatusOK)
}

// GetLabourRates is a route handler that handles getting list of labour rates
//
// # This endpoint is used to get the labour rates
//
// @Summary Get labour rates
// @Description Get labour rates
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "Labour Rate ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/labour-rate [get]
func (u *WorkHandler) GetLabourRates(c *gin.Context) {
	labours, err := u.labourRepo.Find()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Labour rates created successfully", labours, http.StatusCreated)
}

// CreateLabourRateRow is a route handler that handles creating new row for labour rate for a particulkar work identifier. The work identifier will come from a list
//
// # This endpoint is used to create the labour rate row
//
// @Summary Create labour rate
// @Description Creates labour rate row
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body LabourRateInput true "create labour rate row"
// @Security BearerAuth
// @Success 201 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/labour-rate [post]
func (u *WorkHandler) CreateLabourRate(c *gin.Context) {
	var input LabourRateInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(LabourRateInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	_, found, err := u.pageRepo.FindByCondition("labour = ?", input.Labour)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("labour already exists"), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	labour := &models.LabourRate{
		ID:             id,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Labour:         input.Labour,
		Unit:           input.Unit,
		PriceInNaira:   input.PriceInNaira,
	}

	err = u.labourRepo.Create(labour)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Page created successfully", labour, http.StatusCreated)
}

// UpdateLabourRate is a route handler that handles updating the page
//
// # This endpoint is used to update the labour rate
//
// @Summary Update labour rate
// @Description Updates labour rate row for particular work identifier
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body LabourRateInput true "update labour rate"
// @Param id path string true "Labour Rate ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/labour-rate/{id} [put]
func (u *WorkHandler) UpdateLabourRate(c *gin.Context) {
	var input LabourRateInput
	id := c.Param("id")

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(LabourRateInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	labour, found, err := u.labourRepo.FindByCondition("id = ?", id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("page not found"), http.StatusNotFound)
		return
	}

	labour = &models.LabourRate{
		ID:             labour.ID,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Labour:         input.Labour,
		Unit:           input.Unit,
		PriceInNaira:   input.PriceInNaira,
	}

	_, err = u.labourRepo.Save(labour)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Labour Rate updated successfully", labour, http.StatusOK)
}

// DeletePlantRate is a route handler that handles deleting the labour rate
//
// # This endpoint is used to delete the labour rate
//
// @Summary Delete labour rate
// @Description Delete labour rate
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "Plant Rate ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/plant-rate/{id} [delete]
func (u *WorkHandler) DeletePlantRate(c *gin.Context) {
	id := c.Param("id")

	page, err := u.labourRepo.Delete(id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	helpers.ReturnJSON(c, "Plant rate deleted successfully", page, http.StatusOK)
}

// GetPlantRates is a route handler that handles gettings list of plant rates
//
// # This endpoint is used to get the plant rates
//
// @Summary Get labour rates
// @Description Get labour rates
// @Tags Work
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/plant-rate [get]
func (u *WorkHandler) GetPlantRates(c *gin.Context) {
	rates, err := u.plantRate.Find()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Plant rates retrieved successfully", rates, http.StatusCreated)
}

// CreatePlantRate is a route handler that handles creating new row for current plant rate for a particulkar work identifier. The work identifier will come from a list
//
// # This endpoint is used to create the plant rate row
//
// @Summary Create plant rate
// @Description Creates plant rate row
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body LabourRateInput true "create plant rate row"
// @Security BearerAuth
// @Success 201 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/plant-rate [post]
func (u *WorkHandler) CreatePlantRate(c *gin.Context) {
	var input CurrentPlantRateInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(CurrentPlantRateInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	_, found, err := u.pageRepo.FindByCondition("equipment = ?", input.Equipment)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("labour already exists"), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	plantRate := &models.CurrentPlantRate{
		ID:             id,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Equipment:      input.Equipment,
		Rental:         input.Rental,
		Diesel:         input.Diesel,
		TotalPrice:     input.TotalPrice,
	}

	err = u.plantRate.Create(plantRate)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Plant rate created successfully", plantRate, http.StatusCreated)
}

// UpdatePlantRate is a route handler that handles updating the plant rate
//
// # This endpoint is used to update the plant rate
//
// @Summary Update plant rate
// @Description Updates plant rate row for particular work identifier
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body PlantRateInput true "update plant rate"
// @Param id path string true "Plant Rate ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/plant-rate/{id} [put]
func (u *WorkHandler) UpdatePlantRate(c *gin.Context) {
	var input CurrentPlantRateInput
	id := c.Param("id")

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(CurrentPlantRateInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	plantRate, found, err := u.plantRate.FindByCondition("id = ?", id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("page not found"), http.StatusNotFound)
		return
	}

	plantRate = &models.CurrentPlantRate{
		ID:             plantRate.ID,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Equipment:      input.Equipment,
		Rental:         input.Rental,
		Diesel:         input.Diesel,
		TotalPrice:     input.TotalPrice,
	}

	_, err = u.plantRate.Save(plantRate)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Plant Rate updated successfully", plantRate, http.StatusOK)
}

// DeleteLabourRate is a route handler that handles deleting the material price
//
// # This endpoint is used to delete the material price
//
// @Summary Delete material price
// @Description Delete material price
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "Material Price ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/material-price/{id} [delete]
func (u *WorkHandler) DeleteMaterialPrice(c *gin.Context) {
	id := c.Param("id")

	page, err := repository.NewMaterialRepository(u.deps.DatabaseService).Delete(id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	helpers.ReturnJSON(c, "Material price deleted successfully", page, http.StatusOK)
}

// GetMaterialPrices is a route handler that handles gettings list of material prices
//
// # This endpoint is used to get the material prices
//
// @Summary Get material prices
// @Description Get material prices
// @Tags Work
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/material-price [post]
func (u *WorkHandler) GetMaterialPrices(c *gin.Context) {
	rates, err := repository.NewMaterialRepository(u.deps.DatabaseService).Find()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Plant rates retrieved successfully", rates, http.StatusCreated)
}

// CreateMaterialPrice is a route handler that handles creating new row for market price for a particular work identifier. The work identifier will come from a list
//
// # This endpoint is used to create the market price row
//
// @Summary Create market price
// @Description Creates market price row
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body LabourRateInput true "create market price row"
// @Security BearerAuth
// @Success 201 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/material-price [post]
func (u *WorkHandler) CreateMaterialPrice(c *gin.Context) {
	var input MaterialPriceInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(MaterialPriceInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	_, found, err := repository.NewMaterialRepository(u.deps.DatabaseService).FindByCondition("material = ?", input.Material)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("material already exists"), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	materialPrice := &models.MaterialPrice{
		ID:             id,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Material:       input.Material,
		Unit:           input.Unit,
		Price:          input.Price,
	}

	err = repository.NewMaterialRepository(u.deps.DatabaseService).Create(materialPrice)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Material price created successfully", materialPrice, http.StatusCreated)
}

// UpdateMaterialPrice is a route handler that handles updating the material price
//
// # This endpoint is used to update the material price
//
// @Summary Update material price
// @Description Updates material price row for particular work identifier
// @Tags Work
// @Accept json
// @Produce json
// @Param credentials body MaterialPriceInput true "update material price"
// @Param id path string true "Material Price ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /work/material-price/{id} [put]
func (u *WorkHandler) UpdateMaterialPrice(c *gin.Context) {
	var input MaterialPriceInput
	id := c.Param("id")

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(MaterialPriceInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	materialPrice, found, err := repository.NewMaterialRepository(u.deps.DatabaseService).FindByCondition("id = ?", id)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusForbidden)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("page not found"), http.StatusNotFound)
		return
	}

	materialPrice = &models.MaterialPrice{
		ID:             materialPrice.ID,
		WorkIdentifier: input.WorkIdentifier,
		Title:          input.Title,
		SubTitle:       input.SubTitle,
		Material:       input.Material,
		Unit:           input.Unit,
		Price:          input.Price,
	}

	_, err = repository.NewMaterialRepository(u.deps.DatabaseService).Save(materialPrice)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Material price updated successfully", materialPrice, http.StatusOK)
}
