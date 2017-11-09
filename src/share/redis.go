package share

import (
	"share/config"

	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Address  string `json:"addr"`
	Password string `json:"password"`
	DBNum    int    `json:"db"`
}

func NewRedisPool(conf config.RedisConfig, maxIdle, maxActive int, timeout time.Duration) *redis.Pool {
	return &redis.Pool{
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		IdleTimeout: timeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			// return redis.DialURL(rawurl)
			// return redis.Dial("tcp", addr, redis.DialPassword(password), redis.DialDatabase(dbNum))
			return redis.Dial("tcp", conf.Address, redis.DialPassword(conf.Password), redis.DialDatabase(conf.DBNum))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func RedisGet(pool *redis.Pool, key string, v interface{}, fn func() (interface{}, error)) error {
	if pool == nil {
		return errors.New("empty pool")
	}
	conn := pool.Get()
	value, _ := redis.Bytes(conn.Do("GET", key))
	conn.Close()
	if value == nil {
		ptr, err := fn()
		if err != nil {
			return err
		}

		value, err = json.Marshal(ptr)
		if err != nil {
			return err
		}

		conn := pool.Get()
		conn.Do("SET", key, value)
		conn.Close()
	}

	return json.Unmarshal(value, v)
}

func RedisGetEx(pool *redis.Pool, key string, expire int, v interface{}, fn func() (interface{}, error)) error {
	if pool == nil {
		return errors.New("empty pool")
	}
	conn := pool.Get()
	value, _ := redis.Bytes(conn.Do("GET", key))
	conn.Close()
	if value == nil {
		ptr, err := fn()
		if err != nil {
			return err
		}
		value, err = json.Marshal(ptr)
		if err != nil {
			return err
		}

		conn := pool.Get()
		conn.Do("SETEX", key, expire, value)
		conn.Close()
	}
	return json.Unmarshal(value, v)
}

func RedisDel(pool *redis.Pool, key string) error {
	if pool == nil {
		return errors.New("empty pool")
	}
	conn := pool.Get()
	_, err := conn.Do("DEL", key)
	conn.Close()
	return err
}

var (
	//"redis.log(redis.LOG_NOTICE,tostring(ARGV[i])..' '..tostring(ARGV[i+1])) " +
	LuaAtomicHmset = `
		local key = KEYS[1]
		local num_of_field = tonumber(ARGV[1])
		local cmps = {}
		local fields = {}
		for i=2,2*num_of_field+1,2 do
			table.insert(fields,ARGV[i])
			table.insert(cmps,ARGV[i+1])
		end
		local values = redis.call('hmget',key,unpack(fields))
		for k,v in pairs(cmps) do
			if v ~= values[k] then
				return 0
			end
		end
		local swaps = {}
		for i=2*num_of_field+2,table.getn(ARGV),1 do
			table.insert(swaps,ARGV[i])
		end
		redis.call('hmset',key,unpack(swaps))
		return 1
	`
	HashAtomicScript = redis.NewScript(1, LuaAtomicHmset)
)

//对hash进行原子性的compare and swap cmps-待比较的map swaps-待交换的数据
func RedisAtomicHashCompareAndSwap(pool *redis.Pool, key string, cmps map[string]interface{}, swaps map[string]interface{}) (bool, error) {
	if pool == nil {
		return false, errors.New("empty pool")
	}

	args := redis.Args{}
	args = args.Add(key, len(cmps))
	for k, v := range cmps {
		args = args.Add(k, v)
	}
	for k, v := range swaps {
		args = args.Add(k, v)
	}

	conn := pool.Get()
	defer conn.Close()

	ret, err := redis.Int(HashAtomicScript.Do(conn, args...))

	//log.Println("args:", fmt.Sprintln(args...))
	return ret == 1, err
}

// 对hash进行原子性的compare and swap cmps-待比较的map swaps-待交换的数据
func RedisHashcompareAndSwap(pool *redis.Pool, key string, cmps map[string]string, swaps map[string]string) (bool, error) {
	if pool == nil {
		return false, errors.New("empty pool")
	}
	cmpValues := make([]string, 0)
	cmpArgs := redis.Args{}.Add(key)
	for k, v := range cmps {
		cmpArgs = cmpArgs.Add(k)
		cmpValues = append(cmpValues, v)
	}
	swapArgs := redis.Args{}.Add(key)
	for k, v := range swaps {
		swapArgs = swapArgs.Add(k).Add(v)
	}

	conn := pool.Get()
	values, err := redis.Strings(conn.Do("HMGET", cmpArgs...))
	conn.Close()
	if err != nil {
		return false, err
	}
	for i, v := range cmpValues {
		if v != values[i] {
			return false, nil
		}
	}
	conn = pool.Get()
	_, err = conn.Do("HMSET", swapArgs...)
	conn.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}
