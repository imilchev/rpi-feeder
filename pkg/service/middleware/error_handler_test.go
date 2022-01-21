package middleware

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/tests/utils"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type ErrorHandlerSuite struct {
	suite.Suite
	app *fiber.App
	ctx *fasthttp.RequestCtx
	err error
}

func (suite *ErrorHandlerSuite) SetupTest() {
	fCfg := fiber.Config{
		ErrorHandler: ErrorHandler,
	}
	suite.app = fiber.New(fCfg)
	suite.app.Get("/", func(c *fiber.Ctx) error {
		return suite.err
	})
	suite.ctx = &fasthttp.RequestCtx{}
}

func (suite *ErrorHandlerSuite) TestDoesNotExistErrorHandling() {
	msg := utils.RandString(10)
	suite.err = models.NewDoesNotExistError(msg, msg, msg)

	suite.app.Handler()(suite.ctx)

	suite.Equal(http.StatusNotFound, suite.ctx.Response.Header.StatusCode())
	res := &models.ApiError{}
	suite.NoError(json.Unmarshal(suite.ctx.Response.Body(), res))

	expected := models.NewDoesNotExistError(msg, msg, msg)
	suite.Equal(expected.Message, res.Message)
}

func (suite *ErrorHandlerSuite) TestAlreadyExistsErrorHandling() {
	msg := utils.RandString(10)
	suite.err = models.NewAlreadyExistsError(msg, msg, msg)

	suite.app.Handler()(suite.ctx)

	suite.Equal(http.StatusConflict, suite.ctx.Response.Header.StatusCode())
	res := &models.ApiError{}
	suite.NoError(json.Unmarshal(suite.ctx.Response.Body(), res))

	expected := models.NewAlreadyExistsError(msg, msg, msg)
	suite.Equal(expected.Message, res.Message)
}

func (suite *ErrorHandlerSuite) TestValidationErrorHandling() {
	msg := utils.RandString(10)
	suite.err = models.NewValidationError(msg)

	suite.app.Handler()(suite.ctx)

	suite.Equal(http.StatusBadRequest, suite.ctx.Response.Header.StatusCode())
	res := &models.ApiError{}
	suite.NoError(json.Unmarshal(suite.ctx.Response.Body(), res))

	suite.Equal(msg, res.Message)
}

func (suite *ErrorHandlerSuite) TestRandomErrorHandling() {
	msg := utils.RandString(10)
	suite.err = fmt.Errorf("%s", msg)

	suite.app.Handler()(suite.ctx)

	suite.Equal(http.StatusInternalServerError, suite.ctx.Response.Header.StatusCode())
	res := &models.ApiError{}
	suite.NoError(json.Unmarshal(suite.ctx.Response.Body(), res))

	suite.Equal(msg, res.Message)
}

func TestErrorHandlerSuite(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	suite.Run(t, new(ErrorHandlerSuite))
}
