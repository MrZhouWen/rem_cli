package model

type CoachBase struct {
	Id          int
	OriginId    string
	Name        string
	Mobile      string
	Avatar      string
	AppOriginId int
	Token       string
	AcToken     string
	CreatedAt   int
	UpdatedAt   int
}

func (CoachBase) TableName() string {
	return "coach_base"
}

type CoachLogin struct {
	Id        int
	CoachId   int
	DevSn     string
	SoftSn    string
	State     int
	CreatedAt int
	UpdatedAt int
}

func (CoachLogin) TableName() string {
	return "coach_login"
}
