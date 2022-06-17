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
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"rem_cli/dao"
	"rem_cli/model"
	"rem_cli/util"
	"sync"
	"time"
	//"time"
)

// reportCmd represents the report command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev相关命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "setDevSoftData":
			softSn := ""
			if len(args) > 1 && args[1] != "" {
				softSn = args[1]
			}
			setDevSoftData(softSn)
		case "setDataAsync":
			setDevSoftDataAsync()
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(devCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setDevSoftData(softSn string) {
	var devSofts []struct {
		Id      int
		SoftSn  string
		AgentId int
	}

	if softSn != "" {
		dao.DbSlave.Model(&model.DevSoft{}).Where("soft_sn = ? AND state IN (1, 2, 3)", softSn).Find(&devSofts)
	} else {
		dao.DbSlave.Model(&model.DevSoft{}).Where("state IN (1, 2, 3) AND soft_sn != ?", "").Find(&devSofts)
	}

	ch := make(chan map[string]interface{}, 20)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, devSoft := range devSofts {
			ch <- map[string]interface{}{
				"soft_sn":  devSoft.SoftSn,
				"agent_id": devSoft.AgentId,
			}
			fmt.Println("in chan", devSoft.SoftSn)
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			i, ok := <-ch
			if !ok {
				break
			}
			fmt.Println("out chan", i)
			freeDurationUsed, freeDurationRemain := util.GetDevFreeDuration(i["soft_sn"].(string))
			totalAmount := util.GetDevTotalAmount(i["soft_sn"].(string))
			agentAmount := util.GetDevAgentAmount(i["soft_sn"].(string), i["agent_id"].(int))
			day30Amount := util.GetDev30DayAmount(i["soft_sn"].(string))
			unpaidAmount := util.GetDevUnpaidAmount(i["soft_sn"].(string))
			day30FreeDuration := util.GetDev30DayFreeDuration(i["soft_sn"].(string))
			fmt.Println("consume chan", i["soft_sn"].(string))

			// mysql
			var devSoft model.DevSoft
			dao.DbSlave.Model(&model.DevSoft{}).Where("soft_sn = ?", i["soft_sn"].(string)).First(&devSoft)
			mUpdate := map[string]interface{}{
				"free_duration_used":   freeDurationUsed,
				"free_duration_remain": freeDurationRemain,
				"total_amount":         totalAmount,
				"agent_amount":         agentAmount,
				"day_30_amount":        day30Amount,
				"unpaid_amount":        unpaidAmount,
				"day_30_free_duration": day30FreeDuration,
				"updated_at":           int(time.Now().Unix()),
			}
			dao.Db.Model(&model.DevSoftInfo{}).Where("dev_soft_id = ?", devSoft.Id).Updates(mUpdate)

			// mongodb
			filter := bson.D{{"soft_sn", i["soft_sn"].(string)}}
			update := bson.D{
				{"$set", bson.D{
					{"free_duration_used", freeDurationUsed},
					{"free_duration_remain", freeDurationRemain},
					{"total_amount", totalAmount},
					{"agent_amount", agentAmount},
					{"day_30_amount", day30Amount},
					{"unpaid_amount", unpaidAmount},
					{"day_30_free_duration", day30FreeDuration},
					{"updated_at", int(time.Now().Unix())},
				}},
			}
			updateResult, err := dao.MongoDb.Collection("dev_soft_info").UpdateOne(context.TODO(), filter, update)
			if err != nil {
				fmt.Println(err)
			}
			if updateResult.MatchedCount == 0 {
				insert := bson.D{
					{"soft_sn", i["soft_sn"].(string)},
					{"free_duration_used", freeDurationUsed},
					{"free_duration_remain", freeDurationRemain},
					{"total_amount", totalAmount},
					{"agent_amount", agentAmount},
					{"day_30_amount", day30Amount},
					{"unpaid_amount", unpaidAmount},
					{"day_30_free_duration", day30FreeDuration},
					{"created_at", int(time.Now().Unix())},
					{"updated_at", int(time.Now().Unix())},
				}
				insertResult, err := dao.MongoDb.Collection("dev_soft_info").InsertOne(context.TODO(), insert)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("inserted id", insertResult.InsertedID)
			}
		}
	}()

	wg.Wait()
}

type d struct {
	Id      int
	SoftSn  string
	AgentId int
}

func setDevSoftDataAsync() {
	var devSofts []d
	offset := 0
	for {
		dao.DbSlave.Model(&model.DevSoft{}).Offset(offset).Limit(500).Order("id asc").Where("state IN (1, 2, 3) AND soft_sn != ?", "").Find(&devSofts)
		if len(devSofts) == 0 {
			break
		}
		wg := sync.WaitGroup{}
		for _, devSoft := range devSofts {
			wg.Add(1)
			go func(dev d) {
				defer wg.Done()
				setData(dev)
			}(devSoft)
		}
		wg.Wait()

		offset += 500
	}

	fmt.Println("done")
}

func setData(devSoft d) {
	freeDurationUsed, freeDurationRemain := util.GetDevFreeDuration(devSoft.SoftSn)
	totalAmount := util.GetDevTotalAmount(devSoft.SoftSn)
	agentAmount := util.GetDevAgentAmount(devSoft.SoftSn, devSoft.AgentId)
	day30Amount := util.GetDev30DayAmount(devSoft.SoftSn)
	unpaidAmount := util.GetDevUnpaidAmount(devSoft.SoftSn)
	day30FreeDuration := util.GetDev30DayFreeDuration(devSoft.SoftSn)
	fmt.Println("consume ", devSoft.SoftSn)

	// mysql
	var dev model.DevSoft
	dao.DbSlave.Model(&model.DevSoft{}).Where("soft_sn = ?", devSoft.SoftSn).Last(&dev)
	mUpdate := map[string]interface{}{
		"free_duration_used":   freeDurationUsed,
		"free_duration_remain": freeDurationRemain,
		"total_amount":         totalAmount,
		"agent_amount":         agentAmount,
		"day_30_amount":        day30Amount,
		"unpaid_amount":        unpaidAmount,
		"day_30_free_duration": day30FreeDuration,
		"updated_at":           int(time.Now().Unix()),
	}
	dao.Db.Model(&model.DevSoftInfo{}).Where("dev_soft_id = ?", dev.Id).Updates(mUpdate)
	fmt.Println(dev.SoftSn, mUpdate)

	// mongodb
	filter := bson.D{{"soft_sn", dev.SoftSn}}
	update := bson.D{
		{"$set", bson.D{
			{"free_duration_used", freeDurationUsed},
			{"free_duration_remain", freeDurationRemain},
			{"total_amount", totalAmount},
			{"agent_amount", agentAmount},
			{"day_30_amount", day30Amount},
			{"unpaid_amount", unpaidAmount},
			{"day_30_free_duration", day30FreeDuration},
			{"updated_at", int(time.Now().Unix())},
		}},
	}
	updateResult, err := dao.MongoDb.Collection("dev_soft_info").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("update count", updateResult.ModifiedCount)
	if updateResult.MatchedCount == 0 {
		insert := bson.D{
			{"soft_sn", dev.SoftSn},
			{"free_duration_used", freeDurationUsed},
			{"free_duration_remain", freeDurationRemain},
			{"total_amount", totalAmount},
			{"agent_amount", agentAmount},
			{"day_30_amount", day30Amount},
			{"unpaid_amount", unpaidAmount},
			{"day_30_free_duration", day30FreeDuration},
			{"created_at", int(time.Now().Unix())},
			{"updated_at", int(time.Now().Unix())},
		}
		insertResult, err := dao.MongoDb.Collection("dev_soft_info").InsertOne(context.TODO(), insert)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("inserted id", insertResult.InsertedID)
	}

	SugarLogger.Infof("dev:setDataSync", dev.SoftSn, mUpdate)
}
