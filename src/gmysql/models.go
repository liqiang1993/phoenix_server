package gmysql

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserInfoTab struct {
	Id          int64     `json:"id"`
	UserName    string    `json:"user_name"`
	Password    string    `json:"password"`
	NickName    string    `json:"nick_name"`
	Image       string    `json:"image"`
	GmtCreated  time.Time `json:"gmt_created"`
	GmtModified time.Time `json:"gmt_modified"`
}

type OperateHistoryTab struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	OperateType int32     `json:"operate_type"`
	ContentType int32     `json:"content_type"`
	Detail      string    `json:"detail"`
	GmtCreated  time.Time `json:"gmt_created"`
	GmtModified time.Time `json:"gmt_modified"`
}

type AssetClassTab struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	GmtCreated  time.Time `json:"gmt_created"`
	GmtModified time.Time `json:"gmt_modified"`
}

type AssetObjectTab struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	CardNumber   string    `json:"card_number"`
	Owner        string    `json:"owner"`
	Balance      int64     `json:"balance"`
	AssetClassId int64     `json:"asset_class_id"`
	UserId       int64     `json:"user_id"`
	GmtCreated   time.Time `json:"gmt_created"`
	GmtModified  time.Time `json:"gmt_modified"`
}

type IncomeExpenditureClassTab struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        int32     `json:"type"`
	UserId      int64     `json:"user_id"`
	GmtCreated  time.Time `json:"gmt_created"`
	GmtModified time.Time `json:"gmt_modified"`
}

type IncomeExpenditureFlowTab struct {
	Id            int64     `json:"id"`
	FlowName      string    `json:"flow_name"`
	AssetObjectId int64     `json:"asset_object_id"`
	Owner         string    `json:"owner"`
	Amount        int64     `json:"amount"`
	IeClassId     int64     `json:"ie_class_id"`
	UserId        int64     `json:"user_id"`
	GmtCreated    time.Time `json:"gmt_created"`
	GmtModified   time.Time `json:"gmt_modified"`
}
