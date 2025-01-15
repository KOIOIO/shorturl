package errmsg

const (
	ERROR   = 500
	SUCCESS = 200
)

const (
	//redis
	ERROR_FAILED_SAVE_TO_REDIS = 1001
	//MYSQL
	ERROR_NOT_FOUND_IN_MYSQL = 2001
	ERROR_OTHER_EMS          = 2002
)

var errmsg = map[int]string{
	SUCCESS:                    "操作成功",
	ERROR:                      "操作失败",
	ERROR_FAILED_SAVE_TO_REDIS: "保存到redis失败",
	ERROR_NOT_FOUND_IN_MYSQL:   "MYSQL中未找到该URL",
	ERROR_OTHER_EMS:            "其他错误",
}
