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
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	"rem_cli/dao"
	"rem_cli/model"
	"strconv"
	"time"
	//"time"
)

// reportCmd represents the report command
var insureCmd = &cobra.Command{
	Use:   "insure",
	Short: "insure相关命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "expire":
			insureExpire()
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(insureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func insureExpire() {
	var pays []model.InsureOrderPay
	t := int(time.Now().Unix()) - 3600*24*365
	ctx := context.Background()
	lastId, err := dao.ARedis.Get(ctx, "REM:Insure:Expire:LastId").Result()
	where := "insure_order_pay.pay_time > 0"
	if err != redis.Nil && err != nil && lastId != "" { //todo:条件有误
		where = " and insure_order_pay.id > " + lastId
	}
	where += " and insure_order_pay.pay_time < ? and insure_order_pay.pay_state = 1 and insure_order.insure_state = 1"
	dao.DbSlave.Model(&model.InsureOrderPay{}).Joins("inner join insure_order on insure_order_pay.order_id = insure_order.id").Where(where, t).Order("insure_order_pay.id asc").Find(&pays)
	var orderIds []int
	for _, pay := range pays {
		orderIds = append(orderIds, pay.OrderId)
		lastId = strconv.Itoa(pay.Id)
	}
	fmt.Println(orderIds)
	dao.Db.Model(&model.InsureOrder{}).Where("id in ? and pay_state = 1", orderIds).Updates(map[string]interface{}{"insure_state": 3, "updated_at": int(time.Now().Unix())})

	if lastId != "" {
		dao.ARedis.Set(ctx, "REM:Insure:Expire:LastId", lastId, 0)
	}

	SugarLogger.Infof("insure:expire", lastId)
	fmt.Println("done")
}
