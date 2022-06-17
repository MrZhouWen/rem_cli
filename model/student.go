package model

type StudentBase struct {
	Id                        int
	OriginId                  string
	Name                      string
	Mobile                    string
	Avatar                    string
	FreeDuration              int
	FreeDurationUsed          int
	AppOriginId               int
	Token                     string
	FreeBroadcastDuration     int
	FreeBroadcastDurationUsed int
	LastLoginAt               int
	CreatedAt                 int
	UpdatedAt                 int
}

func (StudentBase) TableName() string {
	return "student_base"
}

type StudentLogin struct {
	Id        int    `json:"id"`
	StudentId int    `json:"student_id"`
	CoachId   int    `json:"coach_id"`
	DevSn     string `json:"dev_sn"`
	SoftSn    string `json:"soft_sn"`
	State     int    `json:"state"`
	IsLock    int    `json:"is_lock"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

func (StudentLogin) TableName() string {
	return "student_login"
}
