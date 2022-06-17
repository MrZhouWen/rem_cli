package model

type OrderBase struct {
	Id          int     `json:"id"`
	Category    int     `json:"category"`
	PayId       int     `json:"pay_id"`
	ProvinceId  int     `json:"province_id"`
	CityId      int     `json:"city_id"`
	DistrictId  int     `json:"district_id"`
	OrderNo     string  `json:"order_no"`
	DevSn       string  `json:"dev_sn"`
	SoftSn      string  `json:"soft_sn"`
	CoachId     int     `json:"coach_id"`
	StudentId   int     `json:"student_id"`
	AgentId     int     `json:"agent_id"`
	AgentSubId  int     `json:"agent_sub_id"`
	StartTime   int     `json:"start_time"`
	EndTime     int     `json:"end_time"`
	Duration    int     `json:"duration"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
	AmountRate  string  `json:"amount_rate"`
	IsFree      int     `json:"is_free"`
	PayState    int     `json:"pay_state"`
	Skin        int     `json:"skin"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   int     `json:"updated_at"`
	SchoolId    int     `json:"school_id"`
	SubjectType int     `json:"subject_type"`
	LineId      int     `json:"line_id"`
	SiteId      string  `json:"site_id"`
	PaidAt      int     `json:"paid_at"`
}

func (OrderBase) TableName() string {
	return "order_base"
}

const (
	CatExam            = 0
	CAT_FREE_BROADCAST = 1
	CAT_SCHOOL         = 2
	CAT_INSURE         = 3
	CAT_FREE_EXAM      = 4
	CAT_RECHARGE       = 5
	CAT_XUECHE         = 10
)

type OrderBaseDs struct {
	Id          int     `json:"id"`
	BaseId      int     `json:"base_id"`
	DevSn       string  `json:"dev_sn"`
	SoftSn      string  `json:"soft_sn"`
	CoachId     int     `json:"coach_id"`
	StudentId   int     `json:"student_id"`
	AgentId     int     `json:"agent_id"`
	AgentSubId  int     `json:"agent_sub_id"`
	Duration    int     `json:"duration"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
	CoachAmount float64 `json:"coach_amount"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   int     `json:"updated_at"`
	SchoolId    int     `json:"school_id"`
	SubjectType int     `json:"subject_type"`
	IsHp        int     `json:"is_hp"`
	PayState    int     `json:"pay_state"`
}

func (OrderBaseDs) TableName() string {
	return "order_base_ds"
}

type OrderBaseData struct {
	Id          int    `json:"id"`
	StudentId   int    `json:"student_id"`
	CategoryId  int    `json:"category_id"`
	BaseId      int    `json:"base_id"`
	Score       int    `json:"score"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	RuleNo      string `json:"rule_no"`
	RuleContent string `json:"rule_content"`
	CreatedAt   int    `json:"created_at"`
}

func (OrderBaseData) TableName() string {
	return "order_base_data"
}

type OrderBaseDataTag struct {
	Id     int `json:"id"`
	DataId int `json:"data_id"`
	TagId  int `json:"tag_id"`
}

func (OrderBaseDataTag) TableName() string {
	return "order_base_data_tag"
}
