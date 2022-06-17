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
	"context"
	"encoding/json"
	"fmt"
	"rem_cli/dao"
	"rem_cli/model"
	"rem_cli/util"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	//"time"
)

// reportCmd represents the report command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "服务商相关命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "indexSetRedis":
			indexSetRedis()
		case "oneSetRedis":
			agentId, _ := strconv.Atoi(args[1])
			fmt.Println(agentId)
			oneSetRedis(agentId)
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func indexSetRedis() {
	var agents []model.AgentBase

	dao.Db.Where("state = ?", 1).Find(&agents)
	wg := sync.WaitGroup{}

	for _, agent := range agents {
		wg.Add(1)
		go func(a model.AgentBase) {
			setAgentIndex(a)
			wg.Done()
		}(agent)
	}
	wg.Wait()
}

func oneSetRedis(id int) {
	var agent model.AgentBase
	dao.Db.Where("id = ?", id).Find(&agent)
	setAgentIndex(agent)
}

func setAgentIndex(agent model.AgentBase) {
	ret := map[string]interface{}{}

	defaultWhere := ""
	if agent.Level == 1 {
		defaultWhere = "agent_id = " + strconv.Itoa(agent.Id)
	} else {
		defaultWhere = "agent_sub_id = " + strconv.Itoa(agent.Id)
	}
	where := defaultWhere + " AND state in (1, 2)"
	var totalNum int64
	dao.DbSlave.Model(&model.DevSoft{}).Where(where).Or("state = 0 AND ec_agent_id = ?", agent.Id).Count(&totalNum)
	// 全部设备数量 = 已激活+未激活(重置过)+已签收未激活+已禁用
	ret["total_num"] = int(totalNum)

	now := time.Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	monthStartStr := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), 1)
	monthStart, err := time.ParseInLocation("2006-01-02", monthStartStr, loc)
	if err != nil {
		panic(err)
	}
	monthStartTimestamp := int(monthStart.Unix())
	where += " AND active_time >= " + strconv.Itoa(monthStartTimestamp)

	var monthNum int64
	dao.DbSlave.Model(&model.DevSoft{}).Where(where).Count(&monthNum)
	// 本月设备数
	ret["month_num"] = int(monthNum)

	k2MonthDuration, k3MonthDuration, k3HpMonthDuration, k3LpMonthDuration := util.GetPaidDuration(agent, monthStartTimestamp, 0)
	// 本月付费时长（小时）
	ret["total_duration"] = util.Decimal(float64(k2MonthDuration+k3MonthDuration) / 60)
	ret["k2_duration"] = util.Decimal(float64(k2MonthDuration) / 60)
	ret["k3_hp_duration"] = util.Decimal(float64(k3HpMonthDuration) / 60)
	ret["k3_lp_duration"] = util.Decimal(float64(k3LpMonthDuration) / 60)

	// 本月收入
	if agent.Level == 1 {
		ret["month_amount"], ret["k2_month_amount"], ret["k3_month_amount"], ret["k3_month_hp_amount"], ret["k3_month_lp_amount"] = util.GetAgentAmount(agent, monthStartTimestamp, 0)
	} else {
		var devSofts []model.DevSoft
		dao.DbSlave.Model(&model.DevSoft{}).Where(defaultWhere+" AND state in ?", []int{1, 2, 3}).Find(&devSofts)
		ret["month_amount"], ret["k2_month_amount"], ret["k3_month_amount"] = util.GetAgentSubAmount(agent, devSofts, monthStartTimestamp, 0)
	}

	ret["agent_no"] = agent.Id

	ret["unpaid_amount"], ret["k2_unpaid_amount"], ret["k3_unpaid_amount"], ret["unpaid_num"], ret["k2_unpaid_num"], ret["k3_unpaid_num"] = util.GetAgentUnpaid(agent)

	// 激活设备数, 禁用设备数, 未激活设备数(新发货未激活+已重置未激活)
	var activeNum, banNum, inactiveNum int64
	dao.DbSlave.Model(&model.DevSoft{}).Where(defaultWhere + " AND state = 1").Count(&activeNum)
	dao.DbSlave.Model(&model.DevSoft{}).Where(defaultWhere + " AND state = 2").Count(&banNum)
	dao.DbSlave.Model(&model.DevSoft{}).Where("state = 0 AND ec_agent_id = ?", agent.Id).Count(&inactiveNum)
	ret["active_num"] = activeNum
	ret["ban_num"] = banNum
	ret["inactive_num"] = inactiveNum

	// 激活设备单设备付费时长 = 总付费时长 / 激活设备数
	if activeNum > 0 {
		ret["average_duration"] = util.Decimal(float64(k2MonthDuration+k3MonthDuration) / float64(activeNum) / 60)
	} else {
		ret["average_duration"] = 0
	}

	// 单设备付费时长 = 总付费时长 / 总设备数
	if totalNum > 0 {
		ret["total_average_duration"] = util.Decimal(float64(k2MonthDuration+k3MonthDuration) / float64(totalNum) / 60)
	} else {
		ret["total_average_duration"] = 0
	}

	// 激活设备活跃率、付费率
	softSns := util.GetAgentSoftSns(agent)
	ret["dev_active_rate"] = util.GetAgentActiveRate(agent, softSns, monthStartTimestamp, 0)
	ret["dev_paid_rate"] = util.GetAgentPaidRate(agent, softSns, monthStartTimestamp, 0)

	// 全部设备活跃率、付费率
	totalSoftSns := util.GetAgentTotalSoftSns(agent)
	ret["total_dev_active_rate"] = util.GetAgentActiveRate(agent, totalSoftSns, monthStartTimestamp, 0)
	ret["total_dev_paid_rate"] = util.GetAgentPaidRate(agent, totalSoftSns, monthStartTimestamp, 0)

	retStr, _ := json.Marshal(ret)
	fmt.Println(string(retStr))

	SugarLogger.Infof("agent_id: %d ret: %s", agent.Id, retStr)

	ctx := context.Background()
	dao.ARedis.Set(ctx, "REM:AgentBase:Index:"+strconv.Itoa(agent.Id), string(retStr), time.Hour*2)

}
