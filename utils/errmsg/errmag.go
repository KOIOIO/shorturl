package errmsg

const (
	ERROR   = 500
	SUCCESS = 200
)

const (
	//shorturl相关
	ERROR_GET_ERROR_FROM_MYSQL_OFSHORTURL = 1001
	ERROR_GET_ERROR_FROM_REDISOFSHORTURL  = 1002

	//longurl相关
	ERROR_GET_ERROR_FROM_MYSQL_OFLONGURL = 2001
	ERROR_GET_ERROR_FROM_REDISOFLONGURl  = 2002
)

var errmsg = map[int]string{
	SUCCESS:                               "操作成功",
	ERROR:                                 "操作失败",
	ERROR_GET_ERROR_FROM_MYSQL_OFSHORTURL: "mysql获取短链接错误",
	ERROR_GET_ERROR_FROM_REDISOFSHORTURL:  "redis获取短链接错误",
	ERROR_GET_ERROR_FROM_MYSQL_OFLONGURL:  "mysql获取长链接错误",
	ERROR_GET_ERROR_FROM_REDISOFLONGURl:   "redis获取长链接错误",
}
