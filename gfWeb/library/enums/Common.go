package enums

type ResultCode int

const (
	CodeSuccess       = 0
	CodeFail          = 1   //前端 弹窗提示
	CodeFail2         = 2   //前端 不主动弹窗提示
	CodeUnauthorized  = 401 //未授权访问
	CodeNoLogin       = 9   //未登录
	CodeLoginOther    = 10  //在其他地方登录
	CodeLoginDisabled = 11  //登录禁止
)

// 账号状态
const (
	Deleted  = iota - 1 //删除
	Disabled            //禁止
	Enabled             //启用
)

// 公告类型
const (
	NoticeTypeMoment = 1 //立即发送
	NoticeTypeClock  = 2 //定时发送
	NoticeTypeLoop   = 3 //循环发送
)

// 盐值
const (
	GmSalt       = "fretj9tnda3gr7t14terg4es5f4ds514f" //gm 盐值
	PasswordSalt = "fdsafafa4cw78c8qwcwce8v7rwc7"
)

// 服务器状态
const (
	ServerStateMaintenance = 1 //维护
	ServerStateNormal      = 2 //正常
	ServerStateHot         = 3 //火爆
)

// 登录类型请求
const (
	LoginTypeSuccess  = 1 // 已登录请求
	LoginTypeNo       = 2 // 未登录请求
	LoginTypeKick     = 3 // 被踢下线
	LoginTypeDisabled = 4 // 请求禁用
)
