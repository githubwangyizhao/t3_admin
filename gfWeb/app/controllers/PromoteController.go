package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
)

type PromoteController struct {
	BaseController
}

const shuffleStr = "wB6Zyk3TlbFDYrJEgH98NV4GCuvKcpnimWAS5MqohdjeXLaRPUfzt2xs7Q"

// 查询推广员
func (c *PromoteController) SelectPromote(r *ghttp.Request) {
	params := &models.PromoteRequest{}
	g.Log().Infof("查询推广员:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	data, count := models.GetPromoteDataList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "查询推广员成功", result)
}

// 创建推广员数据
func (c *PromoteController) PromoteCreate(r *ghttp.Request) {
	params := &models.PromoteData{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	if len(params.Promote) >= 5 {
		c.HttpResult(r, enums.CodeFail, "标识长度不得超过5位", params.Promote)
	}

	g.Log().Infof("创建推广员:%+v", params)
	utils.CheckError(err, "创建推广员")
	//promote := GetPromoteStr()
	//g.Log().Infof("生产的字符串:%+v", promote)
	//params.Promote = promote
	_, err = models.GetPromoteDataOneByPromote(params.Promote)
	if err == nil {
		c.HttpResult(r, enums.CodeFail, "标识不能重复", params.Promote)
	}
	params.CreatedBy = c.curUser.Id
	params.Id = 0
	err = models.AddPromoteDate(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

func (c *PromoteController) GetPromoteLink(r *ghttp.Request) {
	link := "https://www.props-trader.com/share/index.html?invitation="
	c.HttpResult(r, enums.CodeSuccess, "成功!", link)
}

// 编辑推广员数据
func (c *PromoteController) PromoteEdit(r *ghttp.Request) {
	params := &models.PromoteData{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	if len(params.Promote) >= 5 {
		c.HttpResult(r, enums.CodeFail, "标识长度不得超过5位", params.Promote)
	}

	g.Log().Infof("编辑推广员:%+v", params)
	utils.CheckError(err, "编辑推广员")
	params.UpdatedBy = c.curUser.Id
	//var data models.PromoteData

	if params.Id <= 0 {
		c.HttpResult(r, enums.CodeFail, "id输入错误", params.Id)
	}

	data, thisErr := models.GetPromoteDataOne(params.Id)

	g.Log().Infof("数据：%+v,看看更新的错误:%+v", data, thisErr)

	if thisErr != nil {
		c.HttpResult(r, enums.CodeFail, "推广员数据不存在", params.Id)
	}

	//params.Promote = params.promote
	params.CreatedBy = data.CreatedBy
	params.UpdatedBy = c.curUser.Id

	err = models.UpdatePromoteDate(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除推广员列表
func (c *PromoteController) DeletePromoteList(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	utils.CheckError(err)
	PromoteIdList := params
	g.Log().Infof("删除推广员:%+v", PromoteIdList)
	err = models.DeletePromoteDataList(PromoteIdList)
	c.CheckError(err, "删除推广员失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除推广员", PromoteIdList)
}

// 获得Promote字符串
func GetPromoteStr() string {
	value := (gtime.TimestampMilli() - 1615100005987) ^ 5847
	bytes := []byte(shuffleStr)
	promoteLen := int64(len(bytes) - 1)
	var resultBytes []byte
	for {
		value = value / promoteLen
		resultBytes = append(resultBytes, bytes[value%promoteLen])
		if value == 0 {
			return string(resultBytes)
		}
	}
}
