package model

type BillAgent struct {
	Id                  int     `json:"id"`
	Type                int     `json:"type"`
	AgentId             int     `json:"agent_id"`
	DevSn               string  `json:"dev_sn"`
	SoftSn              string  `json:"soft_sn"`
	Duration            int     `json:"duration"`
	Amount              float64 `json:"amount"`
	Ym                  int     `json:"ym"`
	CurrentMonthBalance float64 `json:"current_month_balance"`
	XxId                int     `json:"xx_id"`
	State               int     `json:"state"`
	AgentSubId          int     `json:"agent_sub_id"`
	SubjectType         int     `json:"subject_type"`
	Category            int     `json:"category"`
	IsHp                int     `json:"is_hp"`
	CreatedAt           int     `json:"created_at"`
	UpdatedAt           int     `json:"updated_at"`
}

func (BillAgent) TableName() string {
	return "bill_agent"
}

type BillCoach struct {
	Id          int     `json:"id"`
	Type        int     `json:"type"`
	CoachId     int     `json:"coach_id"`
	DevSn       string  `json:"dev_sn"`
	SoftSn      string  `json:"soft_sn"`
	Amount      float64 `json:"amount"`
	XxId        int     `json:"xx_id"`
	Duration    int     `json:"duration"`
	Ym          int     `json:"ym"`
	State       int     `json:"state"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   int     `json:"updated_at"`
	AgentSubId  int     `json:"agent_sub_id"`
	SubjectType int     `json:"subject_type"`
	Category    int     `json:"category"`
	IsHp        int     `json:"is_hp"`
}

func (BillCoach) TableName() string {
	return "bill_coach"
}

type BillPlatform struct {
	Id          int     `json:"id"`
	AgentId     int     `json:"agent_id"`
	DevSn       string  `json:"dev_sn"`
	SoftSn      string  `json:"soft_sn"`
	Amount      float64 `json:"amount"`
	XxId        int     `json:"xx_id"`
	State       int     `json:"state"`
	Duration    int     `json:"duration"`
	Ym          int     `json:"ym"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   int     `json:"updated_at"`
	AgentSubId  int     `json:"agent_sub_id"`
	SubjectType int     `json:"subject_type"`
	Category    int     `json:"category"`
	IsHp        int     `json:"is_hp"`
}

func (BillPlatform) TableName() string {
	return "bill_platform"
}
