
package entity

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"encoding/json"
	"errors"
)

type Player struct {
	BaseInfo struct {
		Base struct {
			Hp int `json:"hp"`
			Atk int `json:"atk"`
			Df int `json:"df"`
		} `json:"base"`
		BaseAt int `json:"baseAt"`
	} `json:"baseInfo"`
	EquipInfo struct {
		Equip struct {
			Index1 struct {
				ID int `json:"id"`
				Value int `json:"value"`
				Attr []struct {
					ID int `json:"id"`
					Value1 int `json:"value1"`
					Value2 int `json:"value2"`
				} `json:"attr"`
			} `json:"index1"`
		} `json:"equip"`
		EquipAt int `json:"equipAt"`
	} `json:"equipInfo"`
	SkillInfo struct {
		Skill []struct {
			ID int `json:"id"`
			Value int `json:"value"`
			Attr []struct {
				ID int `json:"id"`
				Value1 int `json:"value1"`
				Value2 int `json:"value2"`
			} `json:"attr"`
		} `json:"skill"`
		SkillAt int `json:"skillAt"`
	} `json:"skillInfo"`
	BagInfo struct {
		Bag []struct {
			ItemID int `json:"itemId"`
			Value int `json:"value"`
		} `json:"bag"`
		BagAt int `json:"bagAt"`
	} `json:"bagInfo"`
	Seq int `json:"seq"`
	PlayerAt int `json:"playerAt"`
	TabName string `json:"tabName"`
}

func NewPlayerFromConfig() *Player {
	o, err := InitObjectFromFile(NewPlayer(0),"player.json")
	if err != nil {
		return nil
	}
	p, ok := o.(*Player)
	if ok {
		return p
	}
	return nil
}

func NewPlayer(id int64) *Player {
	p := &Player{}
	p.Seq = int(id)
	p.TabName = "player"
	return p
}

func (p *Player) String() zapcore.Field {
	return zap.Any(p.TabName, p)
}

func (p *Player) Marshal() ([]byte, error) {
	buff, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func (p *Player) UnMarshal(buff []byte) (GameObject, error) {
	err := json.Unmarshal(buff, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Player) GetTime() int64 {
	return int64(p.PlayerAt)
}

func (p *Player) SetTime(t int64) {
	p.PlayerAt = int(t)
}

func (p *Player) ID() int64 {
	return int64(p.Seq)
}

func (p *Player) Name() string {
	return p.TabName
}

func (p *Player) SetID(id int64) {
	p.Seq = int(id)
}

func (p *Player) Instance(o GameObject, err error) (*Player, error) {
	if err != nil {
		return p, err
	}
	p, ok := o.(*Player)
	if ok {
		return p, nil
	}
	return p, errors.New("Player instance fail")
}
