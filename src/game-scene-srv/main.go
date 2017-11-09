package main

import (
	"game-scene-srv/dao/entity"
	"game-scene-srv/dao"
	"share/config"
	"share"
	"time"
	"go.uber.org/zap"
	llog "log"
)

func main() {
	/*
	p := entity.NewPlayer()
	//fmt.Println( "player", p.String())
	cfg := config.Init(share.ConfigNameEnv)
	l := log.Init("test")
	l.Info(cfg.String())
	l.Error("test", p.String())
	fmt.Println(p.String())
	*/
	// "root:505gogogo@tcp(192.168.1.206:3306)/wolf_test_wx?charset=utf8mb4"

	redisConf := config.RedisConfig{}
	redisConf.Address = "127.0.0.1:6379"
	redisConf.Password = ""
	redisConf.DBNum = 0
	redisPool := share.NewRedisPool(redisConf, 2, 5, 300*time.Second)
	dm := dao.NewDataManager(redisPool)
	dm.Init("root:@tcp(127.0.0.1:3306)/wolf_test_wx?charset=utf8mb4")
	//o := entity.NewPlayer(1)
	//o.BaseInfo.BaseAt = 123
	o := entity.NewPlayerFromConfig()
	o.SetID(116)

	llog.Println("name:", o.Name(), "id", o.ID())
	p, err := dm.GetData(o)
	if err != nil {
		llog.Println("111", zap.Error(err))
		err = dm.SetData(o)
		if err != nil {
			llog.Println("222", zap.Error(err))
		}
		llog.Println("...", o.Name(), o.ID())
		o, err = o.Instance(dm.GetData(o))
		if err != nil {
			llog.Println("xxx", zap.Error(err))
			return
		}
		llog.Println("rrr", o)
	}
	o = entity.NewPlayer(2)
	p, err = dm.GetData(o)
	if err != nil {
		llog.Println("333", zap.Error(err))
		return
	}
	player, _ := p.(*entity.Player)

	llog.Println(player)
}


//
func move(x int, y int, z int) {

	// leaveRoom

	// 计算出 该房间在哪台服务器 根据房间坐标划分


	// joinRoom



}

