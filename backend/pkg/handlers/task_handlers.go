package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/service"
)

type HTTPHandlers struct {
	*entities.User
	AIService      *service.AIPredictService
	HeatMapService *service.HeatMapService
	ProblemService *service.ProblemService
}

func respondError(c *gin.Context, err error, status int) {
	c.Error(err)
	errDTO := newErrDTO(err, time.Now())
	c.JSON(status, errDTO)
	c.Abort()
}

func GuestUser(svc *service.AIPredictService, heatmap *service.HeatMapService, problems *service.ProblemService) gin.HandlerFunc {
	return func(c *gin.Context) {
		guest := &HTTPHandlers{
			&entities.User{
				ID:   0,
				Name: "guest",
				Role: "guest",
			},
			svc,
			heatmap,
			problems,
		}
		c.Set("currentUser", guest)
		c.Next()
	}
}

/*
pattern: /heatmap
method:  GET
info:

succeed:

	status code: 200 created
	response body: json represents heatmap

failed:

	status code: 500 ...
	response body: json with error, time
*/

func (h *HTTPHandlers) GetHeatmap(c *gin.Context) {
	cachemap, err := h.HeatMapService.GetHeatMap(c)
	if err != nil {
		respondError(c, err, http.StatusInternalServerError)
		return
	}

	heatmap := entities.HeatMap{
		Max:        len(cachemap.HeatMap.HeatPoints),
		HeatPoints: cachemap.HeatMap.HeatPoints,
	}

	c.JSON(http.StatusOK, heatmap)
}

/*
pattern: /heatmap
method:  POST
info:

succeed:

	status code: 201 created
	response body: json represents created predicts

failed:

	status code: 500 ...
	response body: json with error, time
*/

func (h *HTTPHandlers) CreateBreefPredicts(c *gin.Context) {
	ctx := c.Request.Context()
	c1 := make(chan string, 1)
	c2 := make(chan string, 2)

	go func() {
		breef, err := h.AIService.PopAnalysis(c, 3072217)
		if err != nil {
			errDTO := newErrDTO(err, time.Now())
			c.JSON(http.StatusInternalServerError, errDTO)
		}

		c1 <- breef.Breef
	}()

	go func() {
		breef, err := h.AIService.PopAnalysis(c, 3390291)
		if err != nil {
			errDTO := newErrDTO(err, time.Now())
			c.JSON(http.StatusInternalServerError, errDTO)
		}

		c2 <- breef.Breef
	}()

	responsesDTO := make([]entities.BreefAIResponse, 0, 2)
	for i := 0; i < 2; i++ {
		select {
		case txt := <-c1:
			responsesDTO = append(responsesDTO, entities.BreefAIResponse{
				DistrictID: 3072217,
				Breef:      txt,
				Status:     "ok",
			})
		case txt := <-c2:
			responsesDTO = append(responsesDTO, entities.BreefAIResponse{
				DistrictID: 3390291,
				Breef:      txt,
				Status:     "ok",
			})
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "request cancelled",
				"time":  time.Now(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, entities.HeatMapResponse{Responses: responsesDTO})
}

/*
pattern: /heatmap/analysis/district/:districtID
method:  GET
info:	 parameters from path

succeed:
	status code: 200 OK
	response body: json represents extended analysis

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/

func (h *HTTPHandlers) GetDistrictPrediction(c *gin.Context) {
	districtID, err := strconv.Atoi(c.Param("districtID"))
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	prediction, err := h.AIService.GetAnalysisByDistrict(c, districtID)
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, prediction)
}

/*
pattern: /heatmap/analysis/type/:typeID
method:  GET
info:	 parameters from path

succeed:

	status code: 200 OK
	response body: json represents extended analysis

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/
func (h *HTTPHandlers) GetTypePrediction(c *gin.Context) {
	typeID, err := strconv.Atoi(c.Param("typeID"))
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	prediction, err := h.AIService.GetAnalysisByType(c, typeID)
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, prediction)
}

/*
pattern: /heatmap/analysis/city/-1 (temporary reserved for city)
method:  GET
info:	 parameters from path

succeed:

	status code: 200 OK
	response body: json represents extended analysis

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/
func (h *HTTPHandlers) GetPredictByCity(c *gin.Context) {
	cityID, err := strconv.Atoi(c.Param("cityID"))
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	if err := validateCityID(cityID); err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	prediction, err := h.AIService.GetAnalysisByCity(c)
	if err != nil {
		respondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, prediction)
}

/*
pattern: /heatmap/problems/:problemID
method:  GET
info:	 parameters from path

succeed:

	status code: 200 OK
	response body: json represents problem

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/
func (h *HTTPHandlers) GetProblem(c *gin.Context) {
	problemID, err := strconv.Atoi(c.Param("problemID"))
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	problem, err := h.ProblemService.GetProblem(c, problemID)
	if err != nil {
		respondError(c, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, problem)
}

/*
pattern: heatmap/districts/:districtID/problems
method:  GET
info:	 parameters in path

succeed:

	status code: 200 OK
	response body: json represents district problems

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/
func (h *HTTPHandlers) ListProblemsByDistrict(c *gin.Context) {
	var distID districtID

	if err := c.ShouldBindUri(&distID); err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}
	problems, err := h.ProblemService.ListProblemsByDistrict(c, distID.DistrictID)
	if err != nil {
		respondError(c, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, problems)
}

/*
pattern: heatmap/districts/:districtID/problems?lat=123&lon=456
method:  POST
info:	 parameters in path + query params

succeed:

	status code: 201 created
	response body: json represents created problem

failed:

	status code: 500, 400 ...
	response body: json with error, time
*/
func (h *HTTPHandlers) CreateProblem(c *gin.Context) {
	var form entities.CreateProblemForm

	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	if err := c.ShouldBind(&form); err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	form.Lat = lat
	form.Lon = lon

	file, err := c.FormFile("file")
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(file.Filename)
	newName := uuid.New().String() + ext
	savePath := filepath.Join("../../uploads", newName)

	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		respondError(c, err, http.StatusInternalServerError)
		return
	}

	form.ImageURL = fmt.Sprintf("http://localhost:8080/uploads/%s", newName)
	err = h.ProblemService.NewProblem(c, form)
	if err != nil {
		respondError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, form)
}
