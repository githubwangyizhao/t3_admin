package service

import (
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"strconv"
	"strings"
)

func PlayerPropLogFromGameDb(params models.PlayerPropLogQueryParam) ([]*models.PlayerPropLog, int, error) {
	data := make([]*models.PlayerPropLog, 0)
	count := 0
	g.Log().Info("getList params: %+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	if err != nil {
		utils.CheckError(err)
		return data, count, err
	}
	node := gameServer.Node
	gameDb, gameDbNodeErr := models.GetGameDbByNode(node)
	if gameDbNodeErr != nil {
		utils.CheckError(gameDbNodeErr)
		return data, count, gameDbNodeErr
	}

	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("prop_id = '52'"))
	if params.StartTime != 0 {
		whereArray = append(whereArray, fmt.Sprintf(" op_time >= '%d'", params.StartTime))
	}
	if params.EndTime != 0 {
		whereArray = append(whereArray, fmt.Sprintf(" op_time < '%d'", params.EndTime))
	}
	PlayerInfo := models.Player{}
	PlayerIdPlayerNameArr := make(map[int]string)
	if params.PlayerName != "" {
		PlayerInfo.Nickname = params.PlayerName
		playerErr := gameDb.Debug().First(&PlayerInfo).Error
		if playerErr != nil {
			utils.CheckError(playerErr)
			return data, count, playerErr
		}
		g.Log().Infof("PlayerInfo: %+v", PlayerInfo)
		PlayerIdPlayerNameArr[PlayerInfo.Id] = PlayerInfo.Nickname
		whereArray = append(whereArray, fmt.Sprintf("player_id = '%d' ", PlayerInfo.Id))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(`SELECT * FROM player_prop_log %s ORDER BY op_time ASC`,
		whereParam)
	err = gameDb.Debug().Raw(sql).Scan(&data).Error
	if err != nil {
		utils.CheckError(err)
		return data, count, err
	}

	if len(PlayerIdPlayerNameArr) == 0 {
		PlayerIds := make([]string, 0)
		exists := make(map[int]bool)
		for _, y := range data {
			if _, ok := exists[y.PlayerId]; !ok {
				exists[y.PlayerId] = true
				PlayerIds = append(PlayerIds, strconv.Itoa(y.PlayerId))
			}
		}

		PlayerList := make([]*models.Player, 0)
		PlayerListErr := gameDb.Debug().Where(fmt.Sprintf(`id in (%s)`, models.GetSQLWhereParam(PlayerIds))).Find(&PlayerList).Error
		if PlayerListErr != nil {
			utils.CheckError(PlayerListErr)
			return data, count, PlayerListErr
		}

		g.Log().Infof("PlayerListErr: %+v", PlayerList)
		for _, i := range PlayerList {
			PlayerIdPlayerNameArr[i.Id] = i.Nickname
		}
	}

	for x, y := range data {
		data[x].PlayerName = PlayerIdPlayerNameArr[y.PlayerId]
	}

	return data, count, nil
}
