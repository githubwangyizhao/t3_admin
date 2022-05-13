package service

import (
	"encoding/json"
	"errors"
	"gfWeb/app/models"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
)

func DeleteRegion(params *models.Region) error {
	// 检查数据
	data := models.GetRegionById(params.Id)
	if data == nil {
		g.Log().Errorf("数据不存在error")
		return errors.New("数据不存在")
	}

	// 从数据库中删除数据
	if err := models.DeleteRegion(data); err != nil {
		g.Log().Errorf("DeleteRegion error: %+v", err)
		return err
	}

	// 数据发往中心服
	if err := NoticeCenterRegion(data, 0); err != nil {
		g.Log().Errorf("EditRegion error: %+v", err)
		return err
	}

	return nil
}

func EditRegion(params *models.Region) error {
	// 检查数据
	data := models.GetRegionById(params.Id)
	if data == nil {
		g.Log().Errorf("数据不存在error")
		return errors.New("数据不存在")
	}
	// 将数据库中created_at和created_by的值赋值给params
	params.CreatedBy = data.CreatedBy
	params.CreatedAt = data.CreatedAt

	// 数据入库
	if err := models.EditRegion(params); err != nil {
		g.Log().Errorf("EditRegion error: %+v", err)
		return err
	}

	// 数据发往中心服
	if err := NoticeCenterRegion(params, 1); err != nil {
		g.Log().Errorf("EditRegion error: %+v", err)
		return err
	}

	return nil
}

func AddRegion(params *models.Region) error {
	// 数据入库
	if err := models.AddRegion(params); err != nil {
		g.Log().Errorf("AddRegion error: %+v", err)
		return err
	}

	// 数据发往中心服
	if err := NoticeCenterRegion(params, 1); err != nil {
		g.Log().Errorf("AddRegion error: %+v", err)
		return err
	}

	return nil
}

func NoticeCenterRegion(params *models.Region, stats int) error {
	pool := utils.GetAsyncPool()
	err := pool.Add(func() {
		var request struct {
			Currency string `json:"currency"`
			Region   string `json:"region"`
			AreaCode string `json:"area_code"`
			Stats    int    `json:"stats"`
		}
		request.AreaCode = params.AreaCode
		request.Region = params.Region
		request.Currency = params.Currency
		request.Stats = stats
		g.Log().Infof("ddd: %d %+v", stats, request)

		data, err := json.Marshal(request)
		utils.CheckError(err)
		if err != nil {
			g.Log().Errorf("NoticeCenterAreaCode request failure: %+v", err)
		} else {
			url := utils.GetCenterURL() + "/set_area_code"
			resp, _ := utils.HttpRequest(url, string(data))
			g.Log().Infof("NoticeCenterAreaCode call url: %s response: %+v", url, resp)
		}
	})
	if err != nil {
		g.Log().Errorf("NoticeCenterAreaCode error: %+v", err)
		return err
	}
	return nil
}
