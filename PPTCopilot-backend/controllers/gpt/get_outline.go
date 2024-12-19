package gpt

import (
	"backend/controllers"
	"backend/models"
	"strconv"
)

// @Title GetOutline
// @Description Retrieve the outline by its ID
// @Param   id      path    int     true    "Outline ID"
// @Success 200 {object} controllers.Response "Success response with outline"
// @Failure 400 {object} controllers.Response "Invalid parameter"
// @Failure 500 {object} controllers.Response "Internal server error"
// @Router /gpt/outline/{id} [get]
func (this *Controller) GetOutline() {
	id_ := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(id_)
	if err != nil {
		this.Data["json"] = controllers.MakeResponse(controllers.Err, "参数错误", nil)
		this.ServeJSON()
		return
	}

	outline, err := models.GetOutline(id)
	if err != nil {
		this.Data["json"] = controllers.MakeResponse(controllers.Err, err.Error(), nil)
		this.ServeJSON()
		return
	}
	outline = models.RefactOutline(outline)

	this.Data["json"] = controllers.MakeResponse(controllers.OK, "success", outline)
	this.ServeJSON()

}
