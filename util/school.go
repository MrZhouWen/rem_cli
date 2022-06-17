package util

import (
	"rem_cli/dao"
	"rem_cli/model"
)

func SchoolDevDurations(softSn string) (int, int, int, int) {
	type r struct {
		Duration    int
		SubjectType int
		IsHp        int
	}
	var sumR []r
	dao.DbSlave.Model(&model.OrderBaseDs{}).Where("soft_sn = ?", softSn).Select("sum(duration) as duration, subject_type, is_hp").Group("subject_type, is_hp").Find(&sumR)
	duration, k2Duration, hpDuration, totalDuration := 0, 0, 0, 0
	for _, s := range sumR {
		if s.SubjectType == 2 {
			k2Duration = s.Duration
		} else if s.IsHp == 1 {
			hpDuration = s.Duration
		} else {
			duration = s.Duration
		}
	}
	totalDuration = k2Duration + duration + hpDuration
	return duration, k2Duration, hpDuration, totalDuration
}

func SchoolCoachDurations(coachId int) (int, int, int, int) {
	var sumR []struct {
		Duration    int
		SubjectType int
		IsHp        int
	}
	dao.DbSlave.Model(&model.OrderBaseDs{}).Where("coach_id = ?", coachId).Select("sum(duration) as duration, subject_type, is_hp").Group("subject_type, is_hp").Find(&sumR)
	duration, k2Duration, hpDuration, totalDuration := 0, 0, 0, 0
	for _, s := range sumR {
		if s.SubjectType == 2 {
			k2Duration = s.Duration
		} else if s.IsHp == 1 {
			hpDuration = s.Duration
		} else {
			duration = s.Duration
		}
	}
	totalDuration = k2Duration + duration + hpDuration
	return duration, k2Duration, hpDuration, totalDuration
}

func SchoolStudentDurations(studentId int) (int, int, int, int) {
	var sumR []struct {
		Duration    int
		SubjectType int
		IsHp        int
	}
	dao.DbSlave.Model(&model.OrderBaseDs{}).Where("student_id = ?", studentId).Select("sum(duration) as duration, subject_type, is_hp").Group("subject_type, is_hp").Find(&sumR)
	duration, k2Duration, hpDuration, totalDuration := 0, 0, 0, 0
	for _, s := range sumR {
		if s.SubjectType == 2 {
			k2Duration = s.Duration
		} else if s.IsHp == 1 {
			hpDuration = s.Duration
		} else {
			duration = s.Duration
		}
	}
	totalDuration = k2Duration + duration + hpDuration
	return duration, k2Duration, hpDuration, totalDuration
}
