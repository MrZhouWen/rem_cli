package util

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"rem_cli/dao"
)

func StudentOrderRecord(studentId int) (int, int, float64, float64) {
	var r struct {
		K3OrderNum int   `json:"k3_order_num"`
		K2OrderNum int   `json:"k2_order_num"`
		K3Pass     []int `json:"k3_pass"`
		K2Pass     []int `json:"k2_pass"`
	}
	filter := bson.D{{"student_id", studentId}}
	err := dao.MongoDb.Collection("order_record").FindOne(context.TODO(), filter).Decode(&r)
	if err != nil {
		fmt.Println(err)
	}
	k3OrderNum, k2OrderNum := r.K3OrderNum, r.K2OrderNum
	var k3Rate float64
	var k2Rate float64
	k3Num, k2Num := 0, 0
	k3PassNum, k2PassNum := 0, 0
	for i := 0; i < 3; i++ {
		if len(r.K3Pass)-1 >= i {
			k3PassNum += r.K3Pass[i]
			k3Num++
		}
		if len(r.K2Pass)-1 >= i {
			k2PassNum += r.K2Pass[i]
			k2Num++
		}
	}
	if k3Num > 0 {
		k3Rate = Decimal(float64(k3PassNum) / float64(k3Num))
	} else {
		k3Rate = 0
	}
	if k2Num > 0 {
		k2Rate = Decimal(float64(k2PassNum) / float64(k2Num))
	} else {
		k2Rate = 0
	}

	return k3OrderNum, k2OrderNum, k3Rate, k2Rate
}
