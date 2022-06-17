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
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/cobra"
	"rem_cli/dao"
	"rem_cli/model"
	"strconv"
	"time"
)

// reportCmd represents the report command
var esSyncCmd = &cobra.Command{
	Use:   "esSync",
	Short: "es同步相关命令",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called, args:", args)
		switch args[0] {
		case "initOrderBase":
			initOrderBase()
		case "initBillAgent":
			id := 0
			if len(args) > 1 && args[1] != "" {
				id, _ = strconv.Atoi(args[1])
			}
			initBillAgent(id)
		case "billAgent":
			t := 120
			if len(args) > 1 && args[1] != "" {
				t, _ = strconv.Atoi(args[1])
			}
			billAgent(t)
		case "initBillCoach":
			id := 0
			if len(args) > 1 && args[1] != "" {
				id, _ = strconv.Atoi(args[1])
			}
			initBillCoach(id)
		case "billCoach":
			t := 120
			if len(args) > 1 && args[1] != "" {
				t, _ = strconv.Atoi(args[1])
			}
			billCoach(t)
		case "initBillPlatform":
			id := 0
			if len(args) > 1 && args[1] != "" {
				id, _ = strconv.Atoi(args[1])
			}
			initBillPlatform(id)
		case "billPlatform":
			t := 120
			if len(args) > 1 && args[1] != "" {
				t, _ = strconv.Atoi(args[1])
			}
			billPlatform(t)
		default:
			fmt.Println("未知参数命令")
		}
	},
}

func init() {
	rootCmd.AddCommand(esSyncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initOrderBase() {
	var orderBases []model.OrderBase

	ctx := context.Background()
	id := 0
	for {
		fmt.Println(id)
		dao.Db.Limit(5000).Order("id asc").Where("id > ?", id).Find(&orderBases)
		if len(orderBases) == 0 {
			break
		}
		for _, orderBase := range orderBases {
			put, err := dao.Es().Index().Index("rem_db").Type("order_base").Id(strconv.Itoa(orderBase.Id)).
				BodyJson(orderBase).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
			id = orderBase.Id
		}
	}
}

func initBillAgent(id int) {
	var billAgents []model.BillAgent

	ctx := context.Background()
	for {
		dao.Db.Limit(5000).Order("id asc").Where("id > ?", id).Find(&billAgents)
		if len(billAgents) == 0 {
			break
		}
		for _, billAgent := range billAgents {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_agent").Id(strconv.Itoa(billAgent.Id)).
				BodyJson(billAgent).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
			id = billAgent.Id
		}
	}
}

func billAgent(t int) {
	updatedAt := int(time.Now().Unix()) - t
	var billAgents []model.BillAgent
	dao.Db.Where("updated_at >= ?", updatedAt).Find(&billAgents)
	if len(billAgents) == 0 {
		return
	}

	ctx := context.Background()
	for _, billAgent := range billAgents {
		get, _ := dao.Es().Get().Index("rem_db").Type("bill_agent").Id(strconv.Itoa(billAgent.Id)).Do(ctx)
		//fmt.Println(get, err)
		if get != nil && get.Found {
			getBill := model.BillAgent{}
			getData, _ := get.Source.MarshalJSON()
			_ = json.Unmarshal(getData, &getBill)
			if getBill.UpdatedAt < billAgent.UpdatedAt {
				update, _ := dao.Es().Update().Index("rem_db").Type("bill_agent").Id(strconv.Itoa(billAgent.Id)).Doc(billAgent).Do(ctx)
				fmt.Println(update.Result)
			}
		} else {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_agent").Id(strconv.Itoa(billAgent.Id)).
				BodyJson(billAgent).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
		}
	}
}

func initBillCoach(id int) {
	var billCoaches []model.BillCoach

	ctx := context.Background()
	for {
		dao.Db.Limit(5000).Order("id asc").Where("id > ?", id).Find(&billCoaches)
		if len(billCoaches) == 0 {
			break
		}
		for _, billCoach := range billCoaches {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_coach").Id(strconv.Itoa(billCoach.Id)).
				BodyJson(billCoach).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
			id = billCoach.Id
		}
	}
}

func billCoach(t int) {
	updatedAt := int(time.Now().Unix()) - t
	var billCoaches []model.BillCoach
	dao.Db.Where("updated_at >= ?", updatedAt).Find(&billCoaches)
	if len(billCoaches) == 0 {
		return
	}

	ctx := context.Background()
	for _, billCoach := range billCoaches {
		get, _ := dao.Es().Get().Index("rem_db").Type("bill_coach").Id(strconv.Itoa(billCoach.Id)).Do(ctx)
		if get != nil && get.Found {
			getBill := model.BillCoach{}
			getData, _ := get.Source.MarshalJSON()
			_ = json.Unmarshal(getData, &getBill)
			if getBill.UpdatedAt < billCoach.UpdatedAt {
				update, _ := dao.Es().Update().Index("rem_db").Type("bill_coach").Id(strconv.Itoa(billCoach.Id)).Doc(billCoach).Do(ctx)
				fmt.Println(update.Result)
			}
		} else {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_coach").Id(strconv.Itoa(billCoach.Id)).
				BodyJson(billCoach).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
		}
	}
}

func initBillPlatform(id int) {
	var billPlatforms []model.BillPlatform

	ctx := context.Background()
	for {
		dao.Db.Limit(5000).Order("id asc").Where("id > ?", id).Find(&billPlatforms)
		if len(billPlatforms) == 0 {
			break
		}
		for _, billPlatform := range billPlatforms {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_platform").Id(strconv.Itoa(billPlatform.Id)).
				BodyJson(billPlatform).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
			id = billPlatform.Id
		}
	}
}

func billPlatform(t int) {
	updatedAt := int(time.Now().Unix()) - t
	var billPlatforms []model.BillPlatform
	dao.Db.Where("updated_at >= ?", updatedAt).Find(&billPlatforms)
	if len(billPlatforms) == 0 {
		return
	}

	ctx := context.Background()
	for _, billPlatform := range billPlatforms {
		get, _ := dao.Es().Get().Index("rem_db").Type("bill_platform").Id(strconv.Itoa(billPlatform.Id)).Do(ctx)
		if get != nil && get.Found {
			getBill := model.BillPlatform{}
			getData, _ := get.Source.MarshalJSON()
			_ = json.Unmarshal(getData, &getBill)
			if getBill.UpdatedAt < billPlatform.UpdatedAt {
				update, _ := dao.Es().Update().Index("rem_db").Type("bill_platform").Id(strconv.Itoa(billPlatform.Id)).Doc(billPlatform).Do(ctx)
				fmt.Println(update.Result)
			}
		} else {
			put, err := dao.Es().Index().Index("rem_db").Type("bill_platform").Id(strconv.Itoa(billPlatform.Id)).
				BodyJson(billPlatform).Do(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(put.Id, put.Index, put.Type)
		}
	}
}
