package gpt

import (
	"backend/controllers"
	"backend/models"
	"encoding/json"
	"strconv"
)

type UpdateOutlineRequest struct {
	Outline string `json:"outline"`
}

// @Title UpdateOutline
// @Description Update an outline by ID
// @Param   id      path    int                       true    "Outline ID"
// @Param   body    body    gpt.UpdateOutlineRequest  true    "Request containing new outline"
// @Success 200 {object} controllers.Response "Update successful"
// @Failure 400 {object} controllers.Response "Invalid parameter"
// @Failure 500 {object} controllers.Response "Internal server error"
// @Router /gpt/outline/{id} [post]
func (this *Controller) UpdateOutline() {
	var request UpdateOutlineRequest
	json.NewDecoder(this.Ctx.Request.Body).Decode(&request)
	id_ := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(id_)
	if err != nil {
		this.Data["json"] = controllers.MakeResponse(controllers.Err, "参数错误", nil)
		this.ServeJSON()
		return
	}
	_, err = models.UpdateOutline(id, request.Outline)
	if err != nil {
		this.Data["json"] = controllers.MakeResponse(controllers.Err, err.Error(), nil)
		this.ServeJSON()
		return
	}

	this.Data["json"] = controllers.MakeResponse(controllers.OK, "success", nil)
	this.ServeJSON()

}
