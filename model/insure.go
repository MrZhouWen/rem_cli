package model

type InsureOrder struct {
	Id             int     `json:"id"`
	StudentId      int     `json:"student_id"`
	AgentId        int     `json:"agent_id"`
	CoachId        int     `json:"coach_id"`
	SoftSn         string  `json:"soft_sn"`
	OrderNo        string  `json:"order_no"`
	InsureNo       string  `json:"insure_no"`
	Amount         float64 `json:"amount"`
	InsureAmount   float64 `json:"insure_amount"`
	DurationAmount float64 `json:"duration_amount"`
	CostAmount     float64 `json:"cost_amount"`
	InsureState    int     `json:"insure_state"`
	AuditState     int     `json:"audit_state"`
	PayState       int     `json:"pay_state"`
	PurchaseTime   int     `json:"purchase_time"`
	AccountingTime int     `json:"accounting_time"`
	Rate           string  `json:"rate"`
	UnitPrice      float64 `json:"unit_price"`
	AgentSubId     int     `json:"agent_sub_id"`
	GroupId        int     `json:"group_id"`
	GroupTime      int     `json:"group_time"`
	CreatedAt      int     `json:"created_at"`
	UpdatedAt      int     `json:"updated_at"`
}

func (InsureOrder) TableName() string {
	return "insure_order"
}

type InsureOrderBill struct {
	Id             int     `json:"id"`
	AgentId        int     `json:"agent_id"`
	AgentSubId     int     `json:"agent_sub_id"`
	Ym             int     `json:"ym"`
	AgentAmount    float64 `json:"agent_amount"`
	CoachAmount    float64 `json:"coach_amount"`
	PlatformAmount float64 `json:"platform_amount"`
	AgentSubAmount float64 `json:"agent_sub_amount"`
	State          int     `json:"state"`
	CreatedAt      int     `json:"created_at"`
}

func (InsureOrderBill) TableName() string {
	return "insure_order_bill"
}

type InsureOrderPay struct {
	Id           int     `json:"id"`
	OrderId      int     `json:"order_id"`
	PayType      string  `json:"pay_type"`
	PayTime      int     `json:"pay_time"`
	PayState     int     `json:"pay_state"`
	PayPlatform  string  `json:"pay_platform"`
	SerialNumber string  `json:"serial_number"`
	Amount       float64 `json:"amount"`
	RealAmount   float64 `json:"real_amount"`
	CreatedAt    int     `json:"created_at"`
	UpdatedAt    int     `json:"updated_at"`
}

func (InsureOrderPay) TableName() string {
	return "insure_order_pay"
}
