package entity

import (
	"log"
	"io/ioutil"
	"game-scene-srv/utils"
	"go.uber.org/zap/zapcore"
)

type GameObject interface {
	Marshal() ([]byte, error)
	UnMarshal([]byte) (GameObject, error)
	String() zapcore.Field
	GetTime() int64
	SetTime(int64)
	ID() int64
	SetID(int64)
	Name() string
}

func getBuffFromJson(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(utils.GetConfigPath() + filename)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func InitObjectFromFile(o GameObject, fileName string) (GameObject, error) {
	buf, err := getBuffFromJson(fileName)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	retObj, err := o.UnMarshal(buf)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return retObj, nil
}