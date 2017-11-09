package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/micro/misc/lib/addr"
)

// Namespace 命名空间
var Namespace = "com.mewe.wolf."

// CurrentIP 当前机器内网IP (100.*.*.*, 10.*.*.*)
var CurrentIP string

// CurrentEnv 当前环境(test:内网测试, pre:预发布环境, tpre:腾讯预发布, pro:正式环境)
var CurrentEnv string

// CurrentZoneArea 当前游戏分区(1微信,2QQ,3游客)
var CurrentZoneArea string

// CurrentZoneAreaId 当前游戏分区(1微信,2QQ,3游客)
var CurrentZoneAreaId int

// CurrentDir 当前的执行目录
var CurrentDir string

// 服务名
var (
	ServiceNameGameRoom        string
	ServiceNameConnector       string
	ServiceNameUser            string
	ServiceNameAPI             string
	ServiceNameAPIAdmin        string
	ServiceNameAPILogin        string
	ServiceNameUserStatus      string
	ServiceNameGameRule        string
	ServiceNameRoomManager     string
	ServiceNamePush            string
	ServiceNameCredit          string
	ServiceNameAdminSocket     string
	ServiceNameAdminStatistics string
	ServiceNameIdgo            string
	ServiceNameSecurity        string
	ServiceNameGameStats       string
	ServiceNameAggregation     string
	ServiceNameCms             string
	ServiceNameBag             string
	ServiceNameMission         string
	ServiceNameMatch           string
	ServiceNameInvitation      string
	ServiceNameRoomStats       string
	ServiceNameSensitiveWord   string
	ServiceNameRelation        string
	ServiceNameAgora           string
	ServiceNameQcloudIm        string
	ServiceNameThirdReport     string
	ServiceNameUserReport      string
	ServiceNameTss             string
	ServiceNameIdip            string
	ServiceNameAPIThirdParty   string
	ServiceNameGroup           string
	ServiceNameMall            string
	ServiceNameAssistant       string
)

// ServiceNameAlias 服务的短别名
var ServiceNameAlias map[string]string

// 消息主题与队列
var (
	// 游戏结束 2017-09-04 游戏结束的消息由战绩分发
	//TopicGameOver string
	// 处理系统消息
	QueueSystemMessage string
	// 处理IM用户注册
	QueueIMRegister string
	// 处理延迟后的举报消息
	QueueCreditDelayReport string
	// 处理游戏结束后的信誉分计算
	QueueCreditGameOver string
	// 处理游戏战绩统计
	QueueGameStatsProcess string
	// 处理更新用户的周统计数据
	QueueGameStatsUpdateWeekStats string
	// 处理更新用户的关系链
	QueueUpdateUserRelation string
	// 战斗结束通知信誉分
	QueueGameOverToCredit string
	// 战斗结束通知活动中心
	QueueGameOverToActivity string
	// 战斗结束通知俱乐部
	QueueGameOverToGroup string
	// 商城购买,支付成功后的通知处理
	QueueMallBuyGoodsSuccess string
	//面杀助手-广播消息
	QueueAssistant string
	// 用户登录成功后的通知处理
	QueueUserLogin string
	// 支付成功后的通知礼物发货处理
	QueueMallDeliverGift string
	// 用户等级改变后的通知处理
	QueueUserLevelUpgrade string
	// 成功添加好友后的通知处理
	QueueAddFriend string
)

// InitEnv 初始化当前环境(大区等)
func InitEnv() {
	var err error
	CurrentIP, err = addr.Extract("")
	if err != nil {
		panic(err)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	CurrentDir = filepath.Dir(filepath.Dir(dir))

	data, err := ioutil.ReadFile(filepath.Join(CurrentDir, "env.json"))
	if err != nil {
		panic("无法读取大区配置")
	}

	m := make(map[string]string)
	if err := json.Unmarshal(data, &m); err != nil {
		panic("解析当前大区配置失败:" + err.Error())
	}

	CurrentEnv = m["env"]
	CurrentZoneArea = m["zoneArea"]

	if len(CurrentZoneArea) == 0 {
		panic(errors.New("找不到当前大区配置, 无法启动"))
	}

	if CurrentZoneArea == "wx" {
		CurrentZoneAreaId = 1
	} else if CurrentZoneArea == "qq" {
		CurrentZoneAreaId = 2
	} else if CurrentZoneArea == "guest" {
		CurrentZoneAreaId = 3
	} else {
		panic(errors.New("未知大区配置, 无法启动"))
	}

	initNamespace()
	initQueueNames()
}

// initQueueNames 初始化队列名称
func initQueueNames() {
	//TopicGameOver = CurrentEnv + "-" + CurrentZoneArea + "-game-over"
	QueueSystemMessage = CurrentEnv + "-" + CurrentZoneArea + "-system-message"
	QueueIMRegister = CurrentEnv + "-" + CurrentZoneArea + "-im-register"
	QueueCreditDelayReport = CurrentEnv + "-" + CurrentZoneArea + "-credit-report-delay"
	QueueCreditGameOver = CurrentEnv + "-" + CurrentZoneArea + "-credit-gameover"
	QueueGameStatsProcess = CurrentEnv + "-" + CurrentZoneArea + "-game-stats-process"
	QueueGameStatsUpdateWeekStats = CurrentEnv + "-" + CurrentZoneArea + "-update-weekly-game-stats"
	QueueUpdateUserRelation = CurrentEnv + "-" + CurrentZoneArea + "-update-user-relation"
	QueueGameOverToCredit = CurrentEnv + "-" + CurrentZoneArea + "-gameover-to-credit"
	QueueGameOverToActivity = CurrentEnv + "-" + CurrentZoneArea + "-gameover-to-activity"
	QueueGameOverToGroup = CurrentEnv + "-" + CurrentZoneArea + "-gameover-to-group"
	QueueMallBuyGoodsSuccess = CurrentEnv + "-" + CurrentZoneArea + "-mall-buy-goods-success"
	QueueAssistant = CurrentEnv + "-" + CurrentZoneArea + "-assistant"
	QueueUserLogin = CurrentEnv + "-" + CurrentZoneArea + "-user-login"
	QueueMallDeliverGift = CurrentEnv + "-" + CurrentZoneArea + "-mall-deleiver-gift"
	QueueUserLevelUpgrade = CurrentEnv + "-" + CurrentZoneArea + "-user-level-upgrade"
	QueueAddFriend = CurrentEnv + "-" + CurrentZoneArea + "-add-friend"
}

func initNamespace() {
	// 更新namespace
	Namespace = Namespace + CurrentEnv + "." + CurrentZoneArea + "."
	// 更新服务名
	ServiceNameGameRoom = Namespace + "gameroom"
	ServiceNameConnector = Namespace + "connector"
	ServiceNameUser = Namespace + "user"
	ServiceNameAPI = Namespace + "api"
	ServiceNameAPIAdmin = Namespace + "api.admin"
	ServiceNameAPILogin = Namespace + "api.login"
	ServiceNameUserStatus = Namespace + "userStatus"
	ServiceNameGameRule = Namespace + "gamerule"
	ServiceNameRoomManager = Namespace + "roomManager"
	ServiceNamePush = Namespace + "push"
	ServiceNameCredit = Namespace + "credit"
	ServiceNameAdminSocket = Namespace + "admin.socket"
	ServiceNameAdminStatistics = Namespace + "admin.statistics"
	ServiceNameIdgo = Namespace + "idgo"
	ServiceNameSecurity = Namespace + "security"
	ServiceNameGameStats = Namespace + "gamestats"
	ServiceNameAggregation = Namespace + "aggregation"
	ServiceNameCms = Namespace + "cms"
	ServiceNameBag = Namespace + "bag"
	ServiceNameMission = Namespace + "mission"
	ServiceNameMatch = Namespace + "match"
	ServiceNameInvitation = Namespace + "invitation"
	ServiceNameRoomStats = Namespace + "roomStats"
	ServiceNameSensitiveWord = Namespace + "sensitiveWord"
	ServiceNameRelation = Namespace + "relation"
	ServiceNameAgora = Namespace + "agora"
	ServiceNameQcloudIm = Namespace + "qloudim"
	ServiceNameThirdReport = Namespace + "thirdReport"
	ServiceNameUserReport = Namespace + "userReport"
	ServiceNameTss = Namespace + "tss"
	ServiceNameIdip = Namespace + "idip"
	ServiceNameAPIThirdParty = Namespace + "api.third"
	ServiceNameGroup = Namespace + "group"
	ServiceNameMall = Namespace + "mall"
	ServiceNameAssistant = Namespace + "assistant"
	// 更新服务别名
	ServiceNameAlias = map[string]string{
		ServiceNameGameRoom:        "gameroom",
		ServiceNameConnector:       "connector",
		ServiceNameUser:            "user",
		ServiceNameAPI:             "api",
		ServiceNameAPIAdmin:        "apiadmin",
		ServiceNameAPILogin:        "apilogin",
		ServiceNameUserStatus:      "userstatus",
		ServiceNameRoomManager:     "roommanager",
		ServiceNamePush:            "push",
		ServiceNameCredit:          "credit",
		ServiceNameAdminStatistics: "statistics",
		ServiceNameIdgo:            "idgo",
		ServiceNameGameStats:       "gamestats",
		ServiceNameAggregation:     "aggregation",
		ServiceNameCms:             "cms",
		ServiceNameBag:             "bag",
		ServiceNameMission:         "mission",
		ServiceNameMatch:           "match",
		ServiceNameRelation:        "relation",
		ServiceNameAgora:           "agora",
		ServiceNameQcloudIm:        "qcloudim",
		ServiceNameUserReport:      "report",
		ServiceNameAdminSocket:     "adminsocket",
		ServiceNameGameRule:        "gamerule",
		ServiceNameSecurity:        "security",
		ServiceNameTss:             "tss",
		ServiceNameIdip:            "idip",
		ServiceNameAPIThirdParty:   "apithird",
		ServiceNameGroup:           "group",
		ServiceNameMall:            "mall",
		ServiceNameAssistant:       "assistant",
	}
}
