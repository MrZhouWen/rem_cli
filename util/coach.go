package util

import (
	"rem_cli/dao"
	"rem_cli/model"
)

// 教练下学员id
func CoachStudentIds(coachId int) []int {
	var r []struct {
		StudentId int
	}
	dao.DbSlave.Model(&model.StudentLogin{}).Where("coach_id = ?", coachId).Distinct("student_id").Find(&r)
	var studentIds []int
	for _, studentId := range r {
		studentIds = append(studentIds, studentId.StudentId)
	}
	return studentIds
}
