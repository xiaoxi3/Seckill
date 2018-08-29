package controller

import "github.com/astaxie/beego"

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	p.Data["json"] = "Sec Kill"
	p.ServeJSON()
}

func (p *SkillController) SecInfo() {
	p.Data["json"] = "Sec Info"
	p.ServeJSON()
}
