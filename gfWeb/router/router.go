package router

import (
	"gfWeb/app/controllers"
	"gfWeb/app/service/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 你可以将路由注册放到一个文件中管理，
// 也可以按照模块拆分到不同的文件中管理，
// 但统一都放到router目录下。
func init() {
	s := g.Server()

	//s.Use(middleware.CORS, middleware.Auth)

	s.Use(middleware.Auth)                                            // 读静态文件不会到中间件
	s.BindHookHandler("/*", ghttp.HOOK_BEFORE_SERVE, middleware.CORS) // 兼容读取静态文件的操作
	//s.Use(middleware.Auth)

	// 某些浏览器直接请求favicon.ico文件，特别是产生404时
	//s.SetRewrite("/favicon.ico", "/resource/image/favicon.ico")
	s.AddStaticPath(g.Cfg().GetString("database.sqlDownloadPath", "/mysql_back"), g.Cfg().GetString("database.backDir", "mysql_back"))
	s.AddStaticPath("/", g.Cfg().GetString("server.ServerRoot", "views"))
	s.AddStaticPath("/dashboard", g.Cfg().GetString("server.ServerRoot", "views"))

	//登录
	s.Group("/", func(group *ghttp.RouterGroup) {
		loginController := new(controllers.LoginController)
		group.Hook("/login", ghttp.HOOK_BEFORE_SERVE, loginController.ControllerInit)
		group.Hook("/logout", ghttp.HOOK_BEFORE_SERVE, loginController.ControllerInit)
		group.ALL("/login/", loginController.Login)
		group.ALL("/logout", loginController.Logout)
		group.ALL("/wss", wsRequest)
		group.ALL("/ws", wsRequest)
	})

	//角色
	s.Group("/role", func(group *ghttp.RouterGroup) {
		roleController := new(controllers.RoleController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, roleController.ControllerInit)
		group.ALL("list", roleController.List)
		group.ALL("edit/", roleController.Edit)
		group.ALL("delete", roleController.Delete)
		group.ALL("allocateResource", roleController.AllocateResource)
		group.ALL("allocateMenu", roleController.AllocateMenu)
		group.ALL("allocateChannel", roleController.AllocateChannel)
	})

	//资源
	s.Group("/resource", func(group *ghttp.RouterGroup) {
		resourceController := new(controllers.ResourceController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, resourceController.ControllerInit)
		group.ALL("/list", resourceController.List)
		group.ALL("/edit/", resourceController.Edit)
		group.ALL("/getParentResourceList", resourceController.GetParentResourceList)
		group.ALL("/delete", resourceController.Delete)
		group.ALL("/resourceTree", resourceController.ResourceTree)
		group.POST("/checkurlfor", resourceController.CheckUrlFor)
	})

	//菜单
	s.Group("/menu", func(group *ghttp.RouterGroup) {
		menuController := new(controllers.MenuController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, menuController.ControllerInit)
		group.ALL("/list", menuController.List)
		group.ALL("/edit", menuController.Edit)
		group.ALL("/getParentMenuList", menuController.GetParentMenuList)
		group.ALL("/delete", menuController.Delete)
		group.ALL("/menuTree", menuController.MenuTree)
	})

	//用户
	s.Group("/user", func(group *ghttp.RouterGroup) {
		userController := new(controllers.UserController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, userController.ControllerInit)
		group.ALL("/list", userController.List)
		group.ALL("/simple_list", userController.SimpleList)
		group.ALL("/edit", userController.Edit)
		group.ALL("/delete", userController.Delete)
		group.ALL("/remove_state", userController.RemoveState)
		group.ALL("/info", userController.Info)
		group.ALL("/changePassword", userController.ChangePassword)
	})
	//工具
	s.Group("/region", func(group *ghttp.RouterGroup) {
		regionController := new(controllers.RegionController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, regionController.ControllerInit)

		// 设置地区手机号码区号与货币单位
		group.POST("/add", regionController.Add)
		group.GET("/list", regionController.Regions)
		group.PUT("/{id}/edit", regionController.Edit)
		group.DELETE("/{id}/delete", regionController.Delete)
	})

	//工具
	s.Group("/tool", func(group *ghttp.RouterGroup) {

		toolController := new(controllers.ToolController)
		ClientHeartVerifyController := new(controllers.ClientHeartVerifyController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, toolController.ControllerInit)
		group.ALL("/action", toolController.Action)
		group.ALL("/send_prop", toolController.SendProp)
		group.ALL("/set_task", toolController.SetTask)
		group.ALL("/active_function", toolController.ActiveFunction)
		group.ALL("/server_time", toolController.ServerTime)
		group.ALL("/get_ip_origin", toolController.GetIpOrigin)
		group.ALL("/merge", toolController.Merge)
		group.ALL("/get_weixin_args", toolController.GetWeixinArgs)
		group.ALL("/update_weixin_args", toolController.UpdateWeixinArgs)
		group.ALL("/finish_branch_task", toolController.FinishBranchTask)
		group.ALL("/background_gm_activity", toolController.ActivityChange)
		group.ALL("/show_file", toolController.ShowFile)
		group.ALL("/get_sys_robot_info", toolController.GetSysRobotInfo)
		//设置开关
		group.ALL("/get_setting_info", toolController.GetSettingInfo)
		group.ALL("/setting_data_edit", toolController.SettingDataEdit)
		group.ALL("/del_setting_data", toolController.DelSettingData)
		//消息模板(通告用户手机和邮箱)
		group.ALL("/get_msg_temp_info", toolController.GetBackgroundMsgTemplateInfo)
		group.ALL("/msg_temp_edit", toolController.BackgroundMsgTemplateEdit)
		group.ALL("/del_msg_temp", toolController.DelBackgroundMsgTemplate)
		// 邮件中心数据
		group.ALL("/get_mail_data_info", toolController.GetMailDataInfo)
		group.ALL("/mail_data_edit", toolController.MailDataEdit)
		group.ALL("/del_mail_data", toolController.DelMailData)
		// 界面操作权限
		group.ALL("/get_page_change_auth", toolController.GetPageChangeAuth)
		group.ALL("/page_change_auth_edit", toolController.PageChangeAuthEdit)
		group.ALL("/del_page_change_auth", toolController.DelPageChangeAuth)
		group.ALL("/get_change_auth_state", toolController.GetChangeAuthState)
		group.ALL("/get_cron_info", toolController.GetCronInfo)
		// 创建客服

		// 设置前端version
		group.GET("/get_version_list", toolController.GetVersionList)
		group.ALL("/get_version", toolController.GetVersion)
		group.ALL("/set_version", toolController.SetVersion)
		// platform_client_info
		group.ALL("/get_platform_client_info_all_list", toolController.GetPlatformClientInfoAllList)
		group.ALL("/get_platform_client_info_list", toolController.GetPlatformClientInfoList)
		group.ALL("/platform_client_info_edit", toolController.PlatformClientInfoEdit)
		group.ALL("/del_platform_client_info_list", toolController.DelPlatformClientInfo)
		group.ALL("/plus_dailys", toolController.PlusDailyStatistics)
		group.ALL("/test", toolController.Test)
		group.ALL("/manual_game_uselog", toolController.ManaulAddGameUselog)
		group.ALL("/manual_game_monsterlog", toolController.ManaulAddGameMonsterlog)
		group.ALL("/manual_ten_i", toolController.ManaulTenMinuteData)
		group.ALL("/push_list", toolController.GotPushList)
		group.ALL("/genparams_push_list", toolController.GenparamsPushList)
		group.ALL("/modify_push_list", toolController.InsertOrUpPushList)
		group.ALL("/set_nopush_account", toolController.SetNopushAccount)

		group.ALL("/statistic_res_opt", toolController.OptStatisticRes)
		group.ALL("/add_statistic_res", toolController.AddStatisticRes)
		group.ALL("/del_statistic_res", toolController.DelStatisticRes)
		group.ALL("/up_statistic_res", toolController.EditStatisticRes)

		group.POST("/get_client_heartbeat", ClientHeartVerifyController.List)
		group.POST("/set_client_heartbeat", ClientHeartVerifyController.Edit)

		//
		group.POST("/bind_up_customer", toolController.BindUpCustomer)
		group.POST("/del_sbind", toolController.DelCustomerBind)
		group.POST("/show_binds", toolController.ShowBindCustomerList)

		group.POST("/tonicks", toolController.ToNicks)

	})

	//游戏服
	s.Group("/game_server", func(group *ghttp.RouterGroup) {
		gameServerController := new(controllers.GameServerController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, gameServerController.ControllerInit)
		group.ALL("/list", gameServerController.List)
		group.ALL("/edit", gameServerController.Edit)
		group.ALL("/delete", gameServerController.Delete)
		group.ALL("/refresh", gameServerController.Refresh)
		group.ALL("/batch_update_state", gameServerController.BatchUpdateState)
		group.ALL("/open_server", gameServerController.OpenServer)
	})

	//节点
	s.Group("/server_node", func(group *ghttp.RouterGroup) {
		serverNodeController := new(controllers.ServerNodeController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, serverNodeController.ControllerInit)
		group.ALL("/list", serverNodeController.List)
		group.ALL("/edit", serverNodeController.Edit)
		group.ALL("/delete", serverNodeController.Delete)
		//group.ALLde/start", serverNodeController.Start)
		//group.ALLde/stop", serverNodeController.Stop)
		group.ALL("/action", serverNodeController.Action)
		group.ALL("/install", serverNodeController.Install)
	})

	//日志
	s.Group("/log", func(group *ghttp.RouterGroup) {
		logController := new(controllers.LogController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, logController.ControllerInit)
		group.ALL("/player_login_log/", logController.PlayerLoinLogList)
		group.ALL("/online_log/", logController.PlayerOnlineLogList)
		group.ALL("/challenge_mission_log/", logController.PlayerChallengeMissionLogList)
		group.ALL("/prop_log/", logController.PlayerPropLogList)
		group.ALL("/mail_log/", logController.PlayerMailLogList)
		group.ALL("/attr_log/", logController.PlayerAttrLogList)
		group.ALL("/request_log/", logController.RequestLogList)
		group.ALL("/login_log/", logController.LoginLogList)
		group.ALL("/open_server_manage_log/", logController.OpenServerManageLogList)
		group.ALL("/update_platform_version_log/", logController.UpdatePlatformVersionLogList)
		group.ALL("/activity_award_log/", logController.ActivityAwardLogList)
		group.ALL("/impact_rank/", logController.ImpactRankList)
		// 玩家游戏场景日志
		group.ALL("/game_scene_log/", logController.PlayerGameSceneLogList)
		// 客户端打点日志
		group.ALL("/client_log/", logController.ClientLogList)
		// 物品使用次数与宝箱、拉霸、转盘日志
		group.POST("/item_event_log/", logController.ItemEventLogList)
		group.POST("/game_monster_log/", logController.GameMonsterLogList)
		group.POST("/game_rbox_log/", logController.GameRoundBoxLog)

		group.POST("/item_event_log_json/", logController.GotOptLogjson)

		//group.POST("/download_prop_log", logController.DownloadPlayerPropLog)
	})

	//玩家
	s.Group("/player", func(group *ghttp.RouterGroup) {
		playerController := new(controllers.PlayerController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, playerController.ControllerInit)
		group.ALL("/list", playerController.List)
		group.ALL("/one", playerController.One)
		group.ALL("/oneByEncodedId", playerController.GetPlayerByEncodedId)
		group.ALL("/detail/", playerController.Detail)
		group.ALL("/account_detail/", playerController.AccountDetail)
		group.ALL("/set_account_type/", playerController.SetAccountType)
		group.ALL("/add_test_account", playerController.AddTestAccount)
		group.ALL("/test_account_list", playerController.TestAccountList)
	})

	//统计
	s.Group("/statistics", func(group *ghttp.RouterGroup) {
		statisticsController := new(controllers.StatisticsController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, statisticsController.ControllerInit)
		//group.ALL("/online_statistics/", statisticsController.OnlineStatisticsList)
		//group.ALL("/register_statistics/", statisticsController.RegisterStatisticsList)
		group.ALL("/consume_analysis/", statisticsController.ConsumeAnalysis)
		group.ALL("/get_server_generalize/", statisticsController.GetServerGeneralize)
		group.ALL("/real_time_online/", statisticsController.GetRealTimeOnline)
		group.ALL("/daily_statistics/", statisticsController.DailyStatisticsList)
		group.ALL("/charge_statistics/", statisticsController.GetChargeStatistics)
		group.ALL("/income_statistics/", statisticsController.GetIncomeStatistics)
		group.ALL("/get_platform_ding_yue/", statisticsController.GetPlatformDingYue)
		group.ALL("/get_oauth_order/", statisticsController.GetPlatformDingYue)
		group.ALL("/withdrawal_list/", statisticsController.OauthList)
		//group.ALL("/get_active_statistics/", statisticsController.ActiveStatisticsList}
		//ConsumeStaticsController := new(controllers.ConsumeStaticsController)
		//group.POST("/consume_statistics/", ConsumeStaticsController.GetConsumeStatisticsList)
		group.POST("/ruby_log", statisticsController.GetRubyStatistics)

	})
	//封禁
	s.Group("/forbid", func(group *ghttp.RouterGroup) {
		forbidController := new(controllers.ForbidController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, forbidController.ControllerInit)
		group.ALL("/set_forbid/", forbidController.SetForbid)
		group.ALL("/forbid_log/", forbidController.ForbidLogList)
	})

	//公告
	s.Group("/notice", func(group *ghttp.RouterGroup) {
		noticeController := new(controllers.NoticeController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, noticeController.ControllerInit)
		//group.ALL("/send_notice/", noticeController.SendNotice)
		group.ALL("/send_cron_notice/", noticeController.SendCronNotice)
		group.ALL("/notice_log/", noticeController.NoticeLogList)
		group.ALL("/del_notice_log/", noticeController.DelNoticeLog)
		group.ALL("/close_cron_notice/", noticeController.RemoveCronNotice)
		group.ALL("/jiguang_push/", noticeController.JiguangPush)
		group.ALL("/jiguang_push_item/", noticeController.GetJgPushItem)

	})

	//登录公告
	s.Group("/login_notice", func(group *ghttp.RouterGroup) {
		loginNoticeController := new(controllers.LoginNoticeController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, loginNoticeController.ControllerInit)
		group.ALL("/set/", loginNoticeController.SetNotice)
		group.ALL("/batch_set/", loginNoticeController.BatchSetNotice)
		group.ALL("/list/", loginNoticeController.LoginNoticeList)
		group.ALL("/del/", loginNoticeController.DelLoginNotice)
	})

	s.Group("/app_notice", func(group *ghttp.RouterGroup) {
		appNoticeController := new(controllers.AppNoticeController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, appNoticeController.ControllerInit)

		group.GET("/lists", appNoticeController.List)
		group.POST("/create", appNoticeController.Create)
		group.PUT("/{id}/edit", appNoticeController.Update)
		group.DELETE("/{id}/delete", appNoticeController.Delete)
	})

	//邮件
	s.Group("/mail", func(group *ghttp.RouterGroup) {
		mailController := new(controllers.MailController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, mailController.ControllerInit)
		group.ALL("/send_mail/", mailController.SendMail)
		group.ALL("/mail_log/", mailController.MailLogList)
		group.ALL("/del_mail_log/", mailController.DelMailLog)
	})

	//平台
	s.Group("/platform", func(group *ghttp.RouterGroup) {
		platformController := new(controllers.PlatformController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, platformController.ControllerInit)
		group.ALL("/list/", platformController.List)
		group.ALL("/edit/", platformController.Edit)
		group.ALL("/del/", platformController.Del)
		group.ALL("/edit_open_server_manage/", platformController.EditOpenServerManage)
	})

	// 渠道
	s.Group("/channel", func(group *ghttp.RouterGroup) {
		channelController := new(controllers.ChannelController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, channelController.ControllerInit)
		group.ALL("/list/", channelController.List)
		group.ALL("/edit/", channelController.Edit)
		group.ALL("/del/", channelController.Del)
	})

	//开服管理
	//group.ALL("/open_server/list/", &controllers.OpenServerController{}, "*:List)
	//group.ALL("/open_server/edit/", &controllers.OpenServerController{}, "*:Edit)
	//group.ALL("/open_server/del/", &controllers.OpenServerController{}, "*:Del}

	// 充值
	s.Group("/charge", func(group *ghttp.RouterGroup) {
		chargeController := new(controllers.ChargeController)
		group.Group("/data", func(groupData *ghttp.RouterGroup) {
			groupData.Hook("/*", ghttp.HOOK_BEFORE_SERVE, chargeController.ControllerInit)
			groupData.ALL("/charge_list/", chargeController.ChargeList)
			groupData.ALL("/charge_rank/", chargeController.ChargeRankList)
			groupData.ALL("/charge_download/", chargeController.ChargeDownload)
			//groupData.ALL("/charge_statistics/", chargeController.ChargeStatisticsList)
			groupData.ALL("/charge_task_distribution/", chargeController.ChargeTaskDistribution)
			groupData.ALL("/charge_activity_distribution/", chargeController.ChargeActivityDistribution)
			groupData.ALL("/charge_money_distribution/", chargeController.ChargeMoneyDistribution)
			groupData.ALL("/charge_level_distribution/", chargeController.ChargeLevelDistribution)
			groupData.ALL("/get_daily_ltv/", chargeController.GetDailyLTV)
			groupData.ALL("/ltv_money", chargeController.LtvMoney)
			groupData.ALL("/diamond_rank", chargeController.RankDiamond)
			groupData.ALL("/goldcoin_rank", chargeController.RankGoldCoin)

		})
		group.Group("/background", func(groupBackground *ghttp.RouterGroup) {
			//后台充值
			backgroundController := new(controllers.BackgroundController)
			groupBackground.Hook("/*", ghttp.HOOK_BEFORE_SERVE, backgroundController.ControllerInit)
			groupBackground.ALL("/background_charge/", backgroundController.Charge)
			groupBackground.ALL("/background_charge_list/", backgroundController.List)
		})
		group.ALL("/getPaymentByPlatform", chargeController.GetPaymentInfoByPlatform)
		group.ALL("/createOrder", chargeController.CreateOrder)
	})

	// 小功能
	s.Group("/small_function", func(group *ghttp.RouterGroup) {
		smallFunctionController := new(controllers.SmallFunctionController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, smallFunctionController.ControllerInit)
		group.ALL("/mail_lock_state/", smallFunctionController.MailLockState)
		group.ALL("/player_id_platform_list/", smallFunctionController.PlayerIdPlatformList)
		group.ALL("/sms_send/", smallFunctionController.SmeSend)
		group.ALL("/get_client_version/", smallFunctionController.GetClientVersion)
		group.ALL("/update_client_version/", smallFunctionController.UpdateClientVersion)
		group.ALL("/background_mail_send/", smallFunctionController.BackgroundMailSend)
		group.ALL("/updatePlatformDingYueStatistics/", smallFunctionController.UpdatePlatformDingYueStatistics)
		group.ALL("/get_env/", smallFunctionController.GetEnvLoginServer)
		group.ALL("/set_env/", smallFunctionController.SetEnvLoginServer)
		group.ALL("/request_game_rpc/", smallFunctionController.RequestGameRpc)
		group.ALL("/get_background_update_version/", smallFunctionController.GetBackgroundUpdateVersion)
		group.ALL("/background_update_version/", smallFunctionController.BackgroundUpdateVersion)
		group.ALL("/stop_background/", smallFunctionController.StopBackgroundUpdateVersion)
		group.GET("/updateCustomerServiceUrl", smallFunctionController.UpdateLoginCustomerServiceUrl)

	})

	// 数据库数据操作
	s.Group("/sql", func(group *ghttp.RouterGroup) {
		sqlController := new(controllers.SqlController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, sqlController.ControllerInit)
		group.ALL("/del_platform_not_database", sqlController.DelPlatformNotDatabase)
		group.ALL("/sql_query", sqlController.SqlQuery)
		group.ALL("/get_database_name", sqlController.GetDatabaseName)
		group.ALL("/pack_database", sqlController.PackDatabase)
	})

	// 合服工具
	s.Group("/merge", func(group *ghttp.RouterGroup) {
		mergeController := new(controllers.MergeController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, mergeController.ControllerInit)
		group.ALL("/platform_merge", mergeController.PlatformMerge)
		group.ALL("/get_platform_merge", mergeController.GetPlatformMerge)
		group.ALL("/audit_platform_merge", mergeController.AuditPlatformMerge)
		group.ALL("/del_platform_merge", mergeController.DelPlatformMerge)
	})

	// 打包工具
	s.Group("/pack_tool", func(group *ghttp.RouterGroup) {
		packToolController := new(controllers.PackToolController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, packToolController.ControllerInit)
		group.ALL("/pack_tool/", packToolController.PackTool)
		group.ALL("/sync_tool/", packToolController.SyncTool)
		//group.ALL("/update_tool/", packToolController.UpdatePlatformTool)
		group.ALL("/get_platform_simple_list/", packToolController.GetPlatformSimpleList)
		group.ALL("/update_platform_version_cron/", packToolController.UpdatePlatformVersionCron)
		//group.ALL("/change_tool/", packToolController.ChangePlatformTool)
		group.ALL("/get_version_tool_info/", packToolController.GetPlatformVersionInfo)
		group.ALL("/get_branch_path/", packToolController.GetBranchPath)
		group.ALL("/get_platform_version_path/", packToolController.GetPlatformVersionPath)
		group.ALL("/get_version_tool_change/", packToolController.GetVersionToolChange)
		group.ALL("/get_version_tool_change_info/", packToolController.GetVersionToolChangeInfo)
		group.ALL("/branch_path_edit/", packToolController.BranchPathEdit)
		group.ALL("/platform_version_path_edit/", packToolController.PlatformVersionPathEdit)
		group.ALL("/version_tool_change_edit/", packToolController.VersionToolChangeEdit)
		group.ALL("/del_branch_path/", packToolController.DelBranchPath)
		group.ALL("/del_platform_version_path/", packToolController.DelPlatformVersionPath)
		group.ALL("/del_version_tool_change/", packToolController.DeleteVersionToolChange)
		group.ALL("/send_version_tool_change/", packToolController.SendVersionToolChange)
		group.ALL("/get_version_tool_change_cron/", packToolController.GetVersionToolChangeCron)
		group.ALL("/version_tool_change_cron_edit/", packToolController.VersionToolChangeCronEdit)
		group.ALL("/del_version_tool_change_cron/", packToolController.DelVersionToolChangeCron)
		group.ALL("/get_version_tool_change_log/", packToolController.GetVersionToolChangeLog)
	})

	//留存
	s.Group("/remain", func(group *ghttp.RouterGroup) {
		remainController := new(controllers.RemainController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, remainController.ControllerInit)
		group.ALL("/total_remain/", remainController.GetTotalRemain)
		group.ALL("/total_remain/", remainController.GetTotalRemain)
		group.ALL("/active_remain/", remainController.GetActiveRemain)
		group.ALL("/task_remain/", remainController.GetTaskRemain)
		group.ALL("/level_remain/", remainController.GetLevelRemain)
		group.ALL("/time_remain/", remainController.GetTimeRemain)
		group.ALL("/charge_remain/", remainController.GetChargeRemain)
	})

	//资产
	s.Group("/inventory", func(group *ghttp.RouterGroup) {
		group.Group("/server", func(groupServer *ghttp.RouterGroup) {
			inventoryServerController := new(controllers.InventoryServerController)
			groupServer.Hook("/*", ghttp.HOOK_BEFORE_SERVE, inventoryServerController.ControllerInit)
			groupServer.ALL("/all_server_list", inventoryServerController.AllServerList)
			groupServer.ALL("/server_list", inventoryServerController.ServerList)
			groupServer.ALL("/edit_server", inventoryServerController.EditServer)
			groupServer.ALL("/delete_server", inventoryServerController.DeleteServer)
			groupServer.ALL("/create_ansible_inventory", inventoryServerController.CreateAnsibleInventory)
		})
		group.Group("/database", func(groupDatabase *ghttp.RouterGroup) {
			inventoryDatabaseController := new(controllers.InventoryDatabaseController)
			groupDatabase.Hook("/*", ghttp.HOOK_BEFORE_SERVE, inventoryDatabaseController.ControllerInit)
			groupDatabase.ALL("/all_database_list", inventoryDatabaseController.AllDatabaseList)
			groupDatabase.ALL("/database_list", inventoryDatabaseController.DatabaseList)
			groupDatabase.ALL("/edit_database", inventoryDatabaseController.EditDatabase)
			groupDatabase.ALL("/delete_database", inventoryDatabaseController.DeleteDatabase)
		})
	})

	//推广人员管理
	s.Group("/promote", func(group *ghttp.RouterGroup) {
		PromoteController := new(controllers.PromoteController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, PromoteController.ControllerInit)
		// 获取推广员信息列表（可传入参数，进行条件查询） get
		group.POST("/lists", PromoteController.SelectPromote)
		// 创建推广员信息（传入推广员标识、状态、推广员备注） post
		group.POST("/create", PromoteController.PromoteCreate)
		// 修改推广员信息（传入推广员标识、状态、推广员备注） put
		group.PUT("/update", PromoteController.PromoteEdit)
		// 删除推广员信息 （可批量，传入id） delete
		group.DELETE("/delete", PromoteController.DeletePromoteList)
		// 获取推广链接
		group.GET("/link", PromoteController.GetPromoteLink)
	})

	s.Group("/adjust", func(group *ghttp.RouterGroup) {
		AdjustController := new(controllers.AdjustController)
		group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, AdjustController.ControllerInit)
		group.POST("/lists", AdjustController.SelectAdjust)
		group.POST("/create", AdjustController.CreateAdjust)
		group.PUT("/edit", AdjustController.EditAdjust)

	})

	s.Group("/api", func(group *ghttp.RouterGroup) {
		ApiController := new(controllers.ApiController)
		group.ALL("/get_player_info", ApiController.GetPlayerInfo)
		group.ALL("/get_platform_info_list", ApiController.GetPlayerInfoList)
		group.ALL("/get_platform_info", ApiController.GetPlatformInfo)
		group.ALL("/get_adjust_list", ApiController.GetAdjustList)
		group.ALL("/get_client_verify", ApiController.GetClientVerifyList)
		group.ALL("/get_app_info_list", ApiController.GetAppNoticeList)
		group.ALL("/get_all_tracker_info", ApiController.GetTrackerInfoList)
		group.ALL("/get_area_code_list", ApiController.GetAreaCodeList)
	})

	s.Group("/gift", func(group *ghttp.RouterGroup) {
		GiftController := new(controllers.GiftController)
		group.POST("/list", GiftController.List)
		group.POST("/create", GiftController.Create)
		group.PUT("/update", GiftController.Update)
		group.DELETE("/delete", GiftController.Delete)
		group.POST("/download", GiftController.Download)
	})

	s.EnableAdmin("/system/tool")
}

// ws请求处理
func wsRequest(r *ghttp.Request) {
	g.Log().Debug("ws请求处理")
	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error(err)
		r.Exit()
	}
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
