package model

type SchoolAssignLog struct {
	Id             int     `json:"id"`
	AgentId        int     `json:"agent_id"`
	StudentId      int     `json:"student_id"`
	Duration       int     `json:"duration"`
	AgentAmount    float64 `json:"agent_amount"`
	PlatformAmount float64 `json:"platform_amount"`
	CreatedAt      int     `json:"created_at"`
}

func (SchoolAssignLog) TableName() string {
	return "school_assign_log"
}

type SchoolDsDev struct {
	Id            int    `json:"id"`
	SchoolId      int    `json:"school_id"`
	DevSoftId     int    `json:"dev_soft_id"`
	SoftSn        string `json:"soft_sn"`
	DevSn         string `json:"dev_sn"`
	State         int    `json:"state"`
	K2Duration    int    `json:"k2_duration"`
	Duration      int    `json:"duration"`
	HpDuration    int    `json:"hp_duration"`
	TotalDuration int    `json:"total_duration"`
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"updated_at"`
}

func (SchoolDsDev) TableName() string {
	return "school_ds_dev"
}

type SchoolDsCoach struct {
	Id            int    `json:"id"`
	SchoolId      int    `json:"school_id"`
	CoachId       int    `json:"coach_id"`
	Name          string `json:"name"`
	Mobile        string `json:"mobile"`
	StudentNum    int    `json:"student_num"`
	K2Duration    int    `json:"k2_duration"`
	Duration      int    `json:"duration"`
	HpDuration    int    `json:"hp_duration"`
	TotalDuration int    `json:"total_duration"`
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"updated_at"`
}

func (SchoolDsCoach) TableName() string {
	return "school_ds_coach"
}

type SchoolDsStudent struct {
	Id            int     `json:"id"`
	SchoolId      int     `json:"school_id"`
	StudentId     int     `json:"student_id"`
	Name          string  `json:"name"`
	RealName      string  `json:"real_name"`
	ExamNum       int     `json:"exam_num"`
	K2ExamNum     int     `json:"k2_exam_num"`
	ExamRate      float64 `json:"exam_rate"`
	K2ExamRate    float64 `json:"k2_exam_rate"`
	K2Duration    int     `json:"k2_duration"`
	Duration      int     `json:"duration"`
	HpDuration    int     `json:"hp_duration"`
	TotalDuration int     `json:"total_duration"`
	CreatedAt     int     `json:"created_at"`
	UpdatedAt     int     `json:"updated_at"`
}

func (SchoolDsStudent) TableName() string {
	return "school_ds_student"
}
