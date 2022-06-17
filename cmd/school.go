/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"rem_cli/dao"
	"rem_cli/model"
	"rem_cli/util"
	"sync"
	"time"
	//"time"
)

// reportCmd represents the report command
var schoolCmd = &cobra.Command{
	Use:   "school",
	Short: "school相关命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "dev":
			schoolDevs()
		case "coach":
			schoolCoaches()
		case "student":
			schoolStudents()
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(schoolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type schoolDev struct {
	Id     int
	SoftSn string
	DevSn  string
	State  int
}

func schoolDevs() {
	var devSofts []schoolDev
	offset := 0
	for {
		dao.DbSlave.Model(&model.DevSoft{}).Offset(offset).Limit(500).Order("id asc").Where("type = 2 AND state IN (1, 2, 3) AND soft_sn != ?", "").Find(&devSofts)
		if len(devSofts) == 0 {
			break
		}
		wg := sync.WaitGroup{}
		for _, devSoft := range devSofts {
			wg.Add(1)
			go func(dev schoolDev) {
				defer wg.Done()
				syncSchoolDev(dev)
			}(devSoft)
		}
		wg.Wait()
		offset += 500
	}

	fmt.Println("done")
}

func syncSchoolDev(dev schoolDev) {
	duration, k2Duration, hpDuration, totalDuration := util.SchoolDevDurations(dev.SoftSn)
	var dsDev model.SchoolDsDev
	result := dao.DbSlave.Model(&model.SchoolDsDev{}).Where("dev_soft_id = ?", dev.Id).First(&dsDev)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var r struct {
			SchoolId int
		}
		dao.DbSlave.Model(&model.DevSoftInfo{}).Where("dev_soft_id = ?", dev.Id).Find(&r)
		insert := model.SchoolDsDev{
			DevSoftId:     dev.Id,
			SchoolId:      r.SchoolId,
			SoftSn:        dev.SoftSn,
			DevSn:         dev.DevSn,
			State:         dev.State,
			Duration:      duration,
			K2Duration:    k2Duration,
			HpDuration:    hpDuration,
			TotalDuration: totalDuration,
			CreatedAt:     int(time.Now().Unix()),
			UpdatedAt:     int(time.Now().Unix()),
		}
		dao.Db.Create(&insert)
		SugarLogger.Infof("school:syncSchoolDev", "insert", dev.SoftSn, insert)
	} else {
		update := map[string]interface{}{
			"duration":       duration,
			"k2_duration":    k2Duration,
			"hp_duration":    hpDuration,
			"total_duration": totalDuration,
			"updated_at":     int(time.Now().Unix()),
		}
		dao.Db.Model(&model.SchoolDsDev{}).Where("id = ?", dsDev.Id).Updates(update)
		SugarLogger.Infof("school:syncSchoolDev", "update", dev.SoftSn, update)
	}
}

func schoolCoaches() {
	var dsDevs []model.SchoolDsDev
	offset := 0
	for {
		dao.DbSlave.Model(&model.SchoolDsDev{}).Offset(offset).Limit(500).Order("id asc").Find(&dsDevs)
		if len(dsDevs) == 0 {
			break
		}
		wg := sync.WaitGroup{}
		for _, dsDev := range dsDevs {
			wg.Add(1)
			go func(dev model.SchoolDsDev) {
				defer wg.Done()
				syncSchoolCoach(dev)
			}(dsDev)
		}
		wg.Wait()
		offset += 500
	}

	fmt.Println("done")
}

func syncSchoolCoach(dev model.SchoolDsDev) {
	coachIds := util.DevCoachIds(dev.SoftSn)
	wg := sync.WaitGroup{}
	for _, coachId := range coachIds {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var coach model.CoachBase
			result := dao.DbSlave.First(&coach, id)
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				syncSchoolCoachData(coach, dev.SchoolId)
			}
		}(coachId)
	}
	wg.Wait()
}

func syncSchoolCoachData(coach model.CoachBase, schoolId int) {
	duration, k2Duration, hpDuration, totalDuration := util.SchoolCoachDurations(coach.Id)
	var dsCoach model.SchoolDsCoach
	result := dao.DbSlave.Model(&model.SchoolDsCoach{}).Where("coach_id = ?", coach.Id).First(&dsCoach)
	studentNum := len(util.CoachStudentIds(coach.Id))
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		insert := model.SchoolDsCoach{
			CoachId:       coach.Id,
			SchoolId:      schoolId,
			Name:          coach.Name,
			Mobile:        coach.Mobile,
			StudentNum:    studentNum,
			Duration:      duration,
			K2Duration:    k2Duration,
			HpDuration:    hpDuration,
			TotalDuration: totalDuration,
			CreatedAt:     int(time.Now().Unix()),
			UpdatedAt:     int(time.Now().Unix()),
		}
		dao.Db.Create(&insert)
		SugarLogger.Infof("school:syncSchoolCoachData", "insert", coach.Id, insert)
	} else {
		update := map[string]interface{}{
			"student_num":    studentNum,
			"name":           coach.Name,
			"mobile":         coach.Mobile,
			"duration":       duration,
			"k2_duration":    k2Duration,
			"hp_duration":    hpDuration,
			"total_duration": totalDuration,
			"updated_at":     int(time.Now().Unix()),
		}
		dao.Db.Model(&model.SchoolDsCoach{}).Where("id = ?", dsCoach.Id).Updates(update)
		SugarLogger.Infof("school:syncSchoolCoachData", "update", coach.Id, update)
	}
}

func schoolStudents() {
	var dsDevs []model.SchoolDsDev
	offset := 0
	for {
		dao.DbSlave.Model(&model.SchoolDsDev{}).Offset(offset).Limit(500).Order("id asc").Find(&dsDevs)
		if len(dsDevs) == 0 {
			break
		}
		wg := sync.WaitGroup{}
		for _, dsDev := range dsDevs {
			wg.Add(1)
			go func(dev model.SchoolDsDev) {
				defer wg.Done()
				syncSchoolStudent(dev)
			}(dsDev)
		}
		wg.Wait()
		offset += 500
	}

	fmt.Println("done")
}

func syncSchoolStudent(dev model.SchoolDsDev) {
	studentIds := util.DevStudentIds(dev.SoftSn)
	wg := sync.WaitGroup{}
	for _, studentId := range studentIds {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var student model.StudentBase
			result := dao.DbSlave.First(&student, id)
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				syncSchoolStudentData(student, dev.SchoolId)
			}
		}(studentId)
	}
	wg.Wait()
}

func syncSchoolStudentData(student model.StudentBase, schoolId int) string {

	if student.Name == "调试账号" {
		return ""
	}

	duration, k2Duration, hpDuration, totalDuration := util.SchoolStudentDurations(student.Id)
	k3OrderNum, k2OrderNum, k3Rate, k2Rate := util.StudentOrderRecord(student.Id)
	var dsStudent model.SchoolDsStudent
	result := dao.DbSlave.Model(&model.SchoolDsStudent{}).Where("student_id = ?", student.Id).First(&dsStudent)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		insert := model.SchoolDsStudent{
			StudentId:     student.Id,
			SchoolId:      schoolId,
			Name:          student.Name,
			ExamNum:       k3OrderNum,
			K2ExamNum:     k2OrderNum,
			ExamRate:      k3Rate,
			K2ExamRate:    k2Rate,
			Duration:      duration,
			K2Duration:    k2Duration,
			HpDuration:    hpDuration,
			TotalDuration: totalDuration,
			CreatedAt:     int(time.Now().Unix()),
			UpdatedAt:     int(time.Now().Unix()),
		}
		dao.Db.Create(&insert)
		SugarLogger.Infof("school:syncSchoolStudentData", "insert", student.Id, insert)
	} else {
		update := map[string]interface{}{
			"name":           student.Name,
			"exam_num":       k3OrderNum,
			"k2_exam_num":    k2OrderNum,
			"exam_rate":      k3Rate,
			"k2_exam_rate":   k2Rate,
			"duration":       duration,
			"k2_duration":    k2Duration,
			"hp_duration":    hpDuration,
			"total_duration": totalDuration,
			"updated_at":     int(time.Now().Unix()),
		}
		dao.Db.Model(&model.SchoolDsStudent{}).Where("id = ?", dsStudent.Id).Updates(update)
		SugarLogger.Infof("school:syncSchoolStudentData", "update", student.Id, update)
	}
	return ""
}
