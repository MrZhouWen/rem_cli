package model_s

type ReportAgent struct {
	Id          int
	AgentId     int
	AgentLevel  int
	TotalDevNum int
	DevNum      int
	ActiveRate  float64
	PaidRate    float64
	Duration    int
	Amount      float64
	AmountRate  string
	Type        int
	ProvinceId  int
	CityId      int
	CreatedAt   int
	UpdatedAt   int
}

func (ReportAgent) TableName() string {
	return "report_agent"
}
