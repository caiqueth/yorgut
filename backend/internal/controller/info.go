package controller

import (
	"net/http"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type InfoController struct {
	infoService service.InfoService
}

func NewInfoController(infoService service.InfoService) *InfoController {
	return &InfoController{
		infoService: infoService,
	}
}

// FindInfo			godoc
// @ID					findInstanceInfo
// @Summary				Gets this instance's public details.
// @Description			Public: how to reach the administrator. Fields the administrator didn't set are null.
// @Produce				application/json
// @Tags				Info
// @Success				200 {object} transport.InstanceInfoResponse{}
// @Router				/info [get]
func (i *InfoController) FindInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, i.infoService.FindInstanceInfo())
}
