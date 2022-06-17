package model

type AgentBase struct {
	Id         int
	Contact    string
	Avatar     string
	Level      int
	ProvinceId int
	CityId     int
}

func (AgentBase) TableName() string {
	return "agent_base"
}

type AgentCity struct {
	Id         int
	AgentId    int
	ProvinceId int
	CityId     int
	Pcd        int
	Price      float64
	K2Price    float64
	CreatedAt  int
	UpdatedAt  int
}

func (AgentCity) TableName() string {
	return "agent_city"
}
