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
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	"rem_cli/dao"
	"rem_cli/model"
	"rem_cli/model_s"
	"rem_cli/util"
	"strconv"
	"sync"
	"time"
	//"time"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "报表统计命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "setAgent":
			setAgent()
		case "temp":
			temp()
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setAgent() {
	fmt.Println("setAgent 参数命令")
	//var agent model.AgentBase
	var agents []model.AgentBase

	dao.Db.Where("state = ?", 1).Find(&agents)
	wg := sync.WaitGroup{}
	for _, agent := range agents {
		wg.Add(1)
		go func(a model.AgentBase) {
			setAgentData(a, 0)
			wg.Done()
		}(agent)

		wg.Add(1)
		go func(a model.AgentBase) {
			setAgentData(a, 1)
			wg.Done()
		}(agent)
	}
	wg.Wait()
}

func setAgentData(agent model.AgentBase, cat int) {
	var ret struct {
		totalDevNum int
		devNum      int
		activeRate  float64
		paidRate    float64
		amount      float64
		duration    int
		amountRate  string
	}

	defaultWhere := ""
	if agent.Level == 1 {
		defaultWhere = "agent_id = " + strconv.Itoa(agent.Id)
	} else {
		defaultWhere = "agent_sub_id = " + strconv.Itoa(agent.Id)
	}
	where := defaultWhere + " AND state = 1"
	var totalDevNum int64
	dao.DbSlave.Model(&model.DevSoft{}).Where(where).Count(&totalDevNum)
	// 总激活设备数
	ret.totalDevNum = int(totalDevNum)

	now := time.Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	nowStr := now.Format("2006-01-02")
	endTimeObj, err := time.ParseInLocation("2006-01-02", nowStr, loc)
	if err != nil {
		panic(err)
	}
	endTime := int(endTimeObj.Unix())
	startTime := 0
	if cat == 1 {
		startTime = endTime - 86400*31
	} else {
		startTime = endTime - 86400*8
	}
	where += " AND active_time > " + strconv.Itoa(startTime)
	where += " AND active_time < " + strconv.Itoa(endTime)

	var devNum int64
	dao.DbSlave.Model(&model.DevSoft{}).Where(where).Count(&devNum)
	// 时间段内激活设备数
	ret.devNum = int(devNum)

	// devNum 激活+重置设备数
	where = defaultWhere + " AND state IN (1, 2)"
	dao.DbSlave.Model(&model.DevSoft{}).Where(where).Count(&devNum)

	softSns := util.GetAgentSoftSns(agent)
	ret.activeRate = util.GetAgentActiveRate(agent, softSns, startTime, endTime) * 100
	ret.paidRate = util.GetAgentPaidRate(agent, softSns, startTime, endTime) * 100

	//缩小取数范围
	tempIdx := model.OrderBase{}
	dao.DbSlave.Model(&model.OrderBase{}).Where("updated_at >= " + strconv.Itoa(startTime)).Take(&tempIdx)
	if tempIdx.Id > 0 {
		defaultWhere += " AND id >" + strconv.Itoa(tempIdx.Id)
	}

	where = defaultWhere + " AND created_at > " + strconv.Itoa(startTime) + " AND created_at < " + strconv.Itoa(endTime) + " AND pay_state = 1 AND is_free = 0"
	type r struct {
		Duration int64
		Amount   float64
	}
	var sumR r
	dao.DbSlave.Model(&model.OrderBase{}).Select("sum(duration) as duration, sum(amount) as amount").Where(where).First(&sumR)
	ret.amount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", sumR.Amount), 64)
	ret.duration = int(sumR.Duration)
	ret.amountRate = ""
	fmt.Println(cat, agent.Id, ret, devNum)

	reportAgent := model_s.ReportAgent{
		TotalDevNum: ret.totalDevNum,
		DevNum:      ret.devNum,
		ActiveRate:  ret.activeRate,
		PaidRate:    ret.paidRate,
		Duration:    ret.duration,
		Amount:      ret.amount,
		AmountRate:  ret.amountRate,
		UpdatedAt:   int(time.Now().Unix()),
	}

	result := dao.DbS.Model(&model_s.ReportAgent{}).Where("agent_id = ? AND agent_level = ? AND type = ? ", agent.Id, agent.Level, cat).Updates(&reportAgent)
	fmt.Println("update", result.RowsAffected, reportAgent)
	if result.RowsAffected < 1 {
		// insert
		reportAgent.CreatedAt = int(time.Now().Unix())
		reportAgent.AgentId = agent.Id
		reportAgent.AgentLevel = agent.Level
		reportAgent.Type = cat
		reportAgent.ProvinceId = agent.ProvinceId
		reportAgent.CityId = agent.CityId
		dao.DbS.Create(&reportAgent)

		fmt.Println("create", reportAgent.Id)
	}
}

func temp() {
	ch1 := make(chan int, 10)

	wg1 := sync.WaitGroup{}
	wg1.Add(2)
	go func() {
		defer wg1.Done()
		//	wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			//wg.Add(1)
			//defer wg.Done()
			ch1 <- i
			fmt.Println("in", i)
		}
		close(ch1)
		//wg.Wait()
	}()
	//wg1.Done()
	//wg1.Wait()

	/*for i := range ch1 {
		fmt.Println("out", i)
	}*/

	//wg1.Add(1)
	go func() {
		defer wg1.Done()
		for {
			e, ok := <-ch1
			fmt.Println(e, ok)
			if !ok {
				break
			}

			/*if len(ch1) == 0 {
				close(ch1)
				fmt.Println("len 0 close")
			}*/
		}
	}()
	//wg1.Done()
	wg1.Wait()

	fmt.Println("done")

}
