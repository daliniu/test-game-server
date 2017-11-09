package dao

import (
	"github.com/garyburd/redigo/redis"
	"game-scene-srv/dao/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gogs.mcyun.com/lcg635/stk.git"
	syserror "errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	SnapshotKey = "Game:%s:%d" // %d roomid
	MaxCacheTime = 60
	MaxPersistenceTime = 300
)

func (d *DataManager) Init(mysqlDSN string) {
	d.db = sqlx.MustConnect("mysql", mysqlDSN)
	d.db.SetMaxIdleConns(1)
	d.db.SetMaxOpenConns(3)
	d.stmts = make(map[int]*sqlx.Stmt)
	d.Transaction = stk.TransactionFunc(d.db)
}

func GetSnapshotKey(name string, id int64) string {
	return fmt.Sprintf(SnapshotKey, name, id)
}

type DataManager struct {
	pool *redis.Pool
	db          *sqlx.DB // Transaction 事务辅助函数
	Transaction func(fn func(*sqlx.Tx) error) error
	stmts       map[int]*sqlx.Stmt
}

func NewDataManager(pool *redis.Pool) *DataManager {
	return &DataManager{pool: pool}
}

func (d *DataManager) GetData(o entity.GameObject) (entity.GameObject, error) {

	retObj, err:= d.getCacheData(o)

	if err != nil {
		retPersistObj, err1 := d.getPersistenceData(o)
		if err1 != nil {
			return nil, err
		}
		return retPersistObj, nil
	} else {
		interval := time.Now().Unix() - retObj.GetTime()
		if  interval > int64(MaxPersistenceTime) {
			// 取不到数据库数据 缓存有数据仍然认为可用
			retPersistObj, err := d.getPersistenceData(o)
			if err != nil {
				return retObj, nil
			}
			// 重新刷新缓存时间
			retPersistObj.SetTime(time.Now().Unix())
			d.setCacheData(retPersistObj)
			return retPersistObj, nil
		}
	}

	return retObj, nil
}

func (d *DataManager) SetData(o entity.GameObject) error {
	// 设置缓存时间
	o.SetTime(time.Now().Unix())

	_, err := d.setPersistenceData(o)
	if err != nil {
		log.Println("setPersistenceData", o.Name(), o.ID(), err)
		return err
	}

	// 缓存
	return d.setCacheData(o)
}

func (d *DataManager) SetCacheData(o entity.GameObject) error {
	return d.setCacheData(o)
}

func (d *DataManager) GetCacheData(o entity.GameObject) (entity.GameObject, error) {
	return d.getCacheData(o)
}

func (d *DataManager) SetPersistenceData(o entity.GameObject) (int64, error) {
	return d.setPersistenceData(o)
}

func (d *DataManager) GetPersistenceData(o entity.GameObject) (entity.GameObject, error) {
	return d.getPersistenceData(o)
}

func (d *DataManager) getCacheData(o entity.GameObject) (entity.GameObject, error) {
	// 缓存
	key := GetSnapshotKey(o.Name(), o.ID())
	conn := d.pool.Get()
	defer  conn.Close()
	buff, err := redis.Bytes(conn.Do("HGET",  key, "data"))
	if err != nil {
		log.Println("getCacheData", o.Name(), o.ID(), err)
		return nil, err
	}
	return o.UnMarshal(buff)
}

func (d *DataManager) setCacheData(o entity.GameObject) error {
	buff, err := o.Marshal()
	if err != nil {
		return err
	}
	key := GetSnapshotKey(o.Name(), o.ID())

	// 缓存
	conn := d.pool.Get()
	defer conn.Close()
	_, err = conn.Do("HSET", key, "data", buff)

	return err
}

func (d *DataManager) setPersistenceData(o entity.GameObject) (int64, error) {

	buff, err := o.Marshal()
	if err != nil {
		return 0, err
	}

	stmt, err := d.getSQLPrepare(o.Name(), 1)
	if err != nil {
		return 0, err
	}
	ret, err := stmt.Exec(buff, buff)
	if err != nil {
		return 0, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	o.SetID(id)
	return id, err
}

func (d *DataManager) getPersistenceData(o entity.GameObject) (entity.GameObject, error) {
	stmt, err := d.getSQLPrepare(o.Name(), 2)
	if err != nil {
		return nil, err
	}
	buff := []byte{}
	err = stmt.Get(&buff, o.ID())
	if err != nil {
		log.Println("getPersistenceData", o.Name(), o.ID(), err)
		return nil, err
	}
	return o.UnMarshal(buff)
}

func (d *DataManager) getSQLPrepare(tableName string, op int) (*sqlx.Stmt, error) {

	stmt, ok := d.stmts[op]
	if !ok {
		setSqlStr := fmt.Sprintf("INSERT INTO %s(data) VALUES(?) ON DUPLICATE KEY UPDATE data=?", tableName)
		getSqlStr := fmt.Sprintf("SELECT data from %s where id=? LIMIT 1", tableName)
		sqlStr := ""
		if op == 1 {
			sqlStr = setSqlStr
		} else if op == 2 {
			sqlStr = getSqlStr
		} else {
			return nil, syserror.New("error op for sql stmt")
		}
		value, err := d.db.Preparex(sqlStr)
		if err != nil {
			return nil, err
		}
		d.stmts[op] = value
		return d.stmts[op], nil
	}
	return stmt, nil
}
