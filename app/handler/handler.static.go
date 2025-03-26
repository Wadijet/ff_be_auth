package handler

import (
	"meta_commerce/app/utility"
	"strconv"

	"github.com/valyala/fasthttp"
)

type StaticHandler struct {
}

func NewStaticHandler() *StaticHandler {
	newHandler := new(StaticHandler)
	return newHandler
}

// ==========================================================================================

func (h *StaticHandler) TestApi(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	utility.JSON(ctx, utility.Payload(true, nil, "Test API ok"))
}

type SystemStaticResponse struct {
	Cpu    interface{} `json:"cpu" bson:"cpu"`
	Memory interface{} `json:"memory" bson:"memory"`
}

func (h *StaticHandler) GetSystemStatic(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	result := new(SystemStaticResponse)
	result.Cpu = utility.GetCpuStatic()
	result.Memory = utility.GetMemoryStatic()

	response = utility.Payload(true, result, "Successful manipulation!")

	ctx.SetStatusCode(fasthttp.StatusOK)
	utility.JSON(ctx, response)
}

func (h *StaticHandler) GetApiStatic(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	buf := string(ctx.FormValue("inseconds"))
	insesonds, err := strconv.ParseInt(buf, 10, 64)
	if err != nil {
		insesonds = 30
	}

	response = utility.Payload(true, utility.GetApiStatic(insesonds), "Successful manipulation!")

	ctx.SetStatusCode(fasthttp.StatusOK)
	utility.JSON(ctx, response)
}
