package model

type DevSoft struct {
	Id          int
	Type        int
	SoftSn      string
	DevSn       string
	AgentId     int
	AgentSubId  int
	ProvinceId  int
	CityId      int
	DistrictId  int
	OwnerName   string
	OwnerMobile string
	State       int
	ActiveTime  int
}

func (DevSoft) TableName() string {
	return "dev_soft"
}

type DevSoftPrice struct {
	Id                  int
	DevSoftId           int
	RateCoach           int
	RateAgent           int
	RatePlatform        int
	K2RateCoach         int
	K2RateAgent         int
	K2RatePlatform      int
	UnitPrice           float64
	K2UnitPrice         float64
	PriceAp             float64
	PriceAgent          float64
	K2PriceAp           float64
	K2PriceAgent        float64
	RechargeUnitPrice   float64
	K2RechargeUnitPrice float64
	CreatedAt           int
	UpdatedAt           int
}

func (DevSoftPrice) TableName() string {
	return "dev_soft_price"
}

type DevSoftInfo struct {
	Id                 int
	DevSoftId          int
	Remark             string
	AppVersion         string
	FreeDurationUsed   int
	FreeDurationRemain int
	TotalAmount        float64
	UnpaidAmount       float64
	Day30Amount        float64
	AgentAmount        float64
	Day30FreeDuration  int
	AgentRemark        string
	AgentRemarkCreated int
	SaleUnitPrice      float64
	ExamDurationStatus int
	SchoolId           int
	CreatedAt          int
	UpdatedAt          int
}

func (DevSoftInfo) TableName() string {
	return "dev_soft_info"
}

type DevSoftOpt struct {
	Id         int
	SoftSn     string
	Type       int
	IsPay      int
	Action     string
	AgentId    int
	AgentSubId int
	StartTime  int
	CreatedAt  int
}

func (DevSoftOpt) TableName() string {
	return "dev_soft_opt"
}
