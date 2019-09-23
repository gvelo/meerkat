package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/schema"
	"net/http"
	"sync"
)

type ApiServer struct {
	schema schema.Schema
	router *gin.Engine
	server *http.Server
	log    zerolog.Logger
	mu     sync.Mutex
}

type ApiError struct {
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	ErrorText string      `json:"error-text"`
	Error     interface{} `json:"error,omitempty"`
}

const (
	indexIDParam = "indexID"
	fieldIDParam = "fieldID"
)

func NewRest(schema schema.Schema) (*ApiServer, error) {

	// TODO(gvelo) set gin to production mode
	server := &ApiServer{
		router: gin.Default(),
		log:    log.With().Str("src", "rest-api").Logger(),
	}

	server.log.Info().Msg("creating rest api")

	server.router.GET("/index", server.getAllIndex)
	server.router.POST("/index", server.createIndex)
	server.router.GET("/index/:indexID", server.getIndex)
	server.router.POST("/index/:indexID", server.updateIndex)
	server.router.DELETE("/index/:indexID", server.deleteIndex)

	server.router.GET("/index/:indexID/fields", server.getFields)
	server.router.POST("/index/:indexID/fields", server.createFields)
	server.router.POST("/index/:indexID/fields/:fieldID", server.updateField)
	server.router.DELETE("/index/:indexID/fields/:fieldID", server.deleteIndex)

	server.router.POST("/index/:indexID/alloc", server.updateAlloc)
	server.router.GET("/index/:indexID/alloc", server.getAlloc)

	server.schema = schema

	return server, nil

}

func (s *ApiServer) Start() {

	s.mu.Lock()
	defer s.mu.Unlock()

	log.Info().Msg("starting rest api server")

	s.server = &http.Server{
		Addr:    ":9090",
		Handler: s.router,
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Error().Err(err).Msg("error listening for connections")
		}
	}()

}

func (s *ApiServer) Shutdown(ctx context.Context) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Info().Msg("stopping rest api server")

	return s.server.Shutdown(ctx)

}

func (s *ApiServer) getAllIndex(c *gin.Context) {
	indexes := s.schema.AllIndex()
	c.JSON(http.StatusOK, indexes)
}

func (s *ApiServer) createIndex(c *gin.Context) {

	indexInfo := &schema.IndexInfo{}

	err := c.ShouldBindJSON(indexInfo)

	if err != nil {
		bindError("cannot process request", c, err)
		return
	}

	index, err := s.schema.CreateIndex(*indexInfo)

	if err != nil {
		appError("error creating index", c, err)
		return
	}

	c.JSON(http.StatusOK, index)
}

func (s *ApiServer) getIndex(c *gin.Context) {

	id := c.Param(indexIDParam)

	index, err := s.schema.Index(id)

	if err != nil {
		appError("cannot retrieve index", c, err)
		return
	}

	c.JSON(http.StatusOK, index)

}

func (s *ApiServer) updateIndex(c *gin.Context) {

	id := c.Param(indexIDParam)
	indexInfo := schema.IndexInfo{}

	err := c.ShouldBindJSON(indexInfo)

	if err != nil {
		bindError("cannot update index", c, err)
		return
	}

	indexInfo.Id = id

	index, err := s.schema.UpdateIndex(indexInfo)

	if err != nil {
		appError("cannot update index", c, err)
		return
	}

	c.JSON(http.StatusOK, index)

}

func (s *ApiServer) deleteIndex(c *gin.Context) {

	id := c.Param(indexIDParam)

	err := s.schema.DeleteIndex(id)

	if err != nil {
		appError("cannot delete index", c, err)
		return
	}

	c.Status(http.StatusOK)

}

func (s *ApiServer) getFields(c *gin.Context) {

	id := c.Param(indexIDParam)

	fields, err := s.schema.AllFields(id)

	if err != nil {
		appError("cannot retrieve fields", c, err)
		return
	}

	c.JSON(http.StatusOK, fields)

}

func (s *ApiServer) createFields(c *gin.Context) {

	indexID := c.Param(indexIDParam)

	_, err := s.schema.Index(indexID)

	if err != nil {
		appError("cannot create fields", c, err)
		return
	}

	field := &schema.Field{}

	err = c.ShouldBindJSON(field)

	if err != nil {
		bindError("cannot process request", c, err)
		return
	}

	field.IndexId = indexID

	r, err := s.schema.CreateFields(indexID, *field)

	if err != nil {
		appError("error creating field", c, err)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (s *ApiServer) updateField(c *gin.Context) {

	indexID := c.Param(indexIDParam)
	fieldID := c.Param(fieldIDParam)

	_, err := s.schema.Index(indexID)

	if err != nil {
		appError("cannot create field", c, err)
		return
	}

	field := schema.Field{}

	err = c.ShouldBindJSON(&field)

	if err != nil {
		bindError("cannot update field", c, err)
		return
	}

	field.Id = fieldID
	field.IndexId = indexID

	err = s.schema.UpdateField(field)

	if err != nil {
		appError("cannot update index", c, err)
		return
	}

	c.JSON(http.StatusOK, field)

}

func (s *ApiServer) updateAlloc(c *gin.Context) {

	indexID := c.Param(indexIDParam)

	_, err := s.schema.Index(indexID)

	if err != nil {
		appError("cannot update alloc", c, err)
		return
	}

	alloc := schema.PartitionAlloc{}

	err = c.ShouldBindJSON(&alloc)

	if err != nil {
		bindError("cannot update alloc", c, err)
		return
	}

	err = s.schema.UpdateAlloc(indexID, alloc)

	if err != nil {
		appError("cannot update alloc", c, err)
		return
	}

	c.JSON(http.StatusOK, alloc)

}

func (s *ApiServer) getAlloc(c *gin.Context) {

	id := c.Param(indexIDParam)

	alloc, err := s.schema.Alloc(id)

	if err != nil {
		appError("cannot retrieve alloc", c, err)
		return
	}

	c.JSON(http.StatusOK, alloc)

}

func appError(status string, c *gin.Context, err error) {
	switch err.(type) {
	case *schema.ValidationError:
		sendError(c, &ApiError{
			Code:      http.StatusBadRequest,
			Status:    status,
			ErrorText: err.Error(),
			Error:     err,
		})
	case *schema.NotFoundError:
		sendError(c, &ApiError{
			Code:      http.StatusNotFound,
			Status:    status,
			ErrorText: err.Error(),
			Error:     err,
		})
	default:
		sendError(c, &ApiError{
			Code:      http.StatusInternalServerError,
			Status:    status,
			ErrorText: err.Error(),
		})
	}
}

func bindError(status string, c *gin.Context, err error) {
	sendError(c, &ApiError{
		Code:      http.StatusBadRequest,
		Status:    status,
		ErrorText: err.Error(),
	})
}

func sendError(c *gin.Context, apiError *ApiError) {
	c.JSON(apiError.Code, apiError)
}