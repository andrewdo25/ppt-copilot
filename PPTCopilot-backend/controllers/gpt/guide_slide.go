package gpt

import (
	"backend/conf"
	"backend/controllers"
	"encoding/json"
	"strings"
)

type GuideSlideRequest struct {
	Outline string `json:"outline"`
}

// @Title GuideSlide
// @Description Generate a guided slide based on an outline
// @Param   body    body    gpt.GuideSlideRequest     true    "Request containing outline"
// @Success 200 {object} controllers.Response "Generated guide slide"
// @Failure 400 {object} controllers.Response "Invalid parameter"
// @Failure 500 {object} controllers.Response "Internal server error"
// @Router /gpt/guide_slide [post]
func (this *Controller) GuideSlide() {
	var request GuideSlideRequest
	json.NewDecoder(this.Ctx.Request.Body).Decode(&request)

	template := conf.GetGuideSinglePromptTemplate()
	template = strings.ReplaceAll(template, "{{outline}}", request.Outline)

	guide_slide, err := RequestGpt(template, SectionXML{}) // <section></section>
	if err != nil {
		this.Data["json"] = controllers.MakeResponse(controllers.Err, err.Error(), nil)
		this.ServeJSON()
		return
	}

	guide_slide = strings.ReplaceAll(guide_slide, "\n", "")
	this.Data["json"] = controllers.MakeResponse(controllers.OK, "success", guide_slide)
	this.ServeJSON()

}
