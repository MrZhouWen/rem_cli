package util

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
	"rem_cli/dao"
	"rem_cli/model"
	"strconv"
	"time"
)

// 获取免费订单时长
func GetDevFreeDuration(softSn string) (int, int) {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewTermQuery("is_free", 1))
	query.Filter(elastic.NewTermQuery("category", 0))
	aggs := elastic.NewSumAggregation().Field("duration")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("sum_duration", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_duration")
	usedDuration := int(*agg.Value)

	id, _ := strconv.Atoi(softSn[1:])
	id -= 10000
	remainDuration := 0
	if id > 1 {
		var student model.StudentBase
		dao.DbSlave.Model(&model.StudentBase{}).First(&student, id)
		remainDuration = student.FreeDuration - usedDuration
	}
	return usedDuration, remainDuration
}

// 获取设备总付费金额（包括科三保）
func GetDevTotalAmount(softSn string) float64 {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermQuery("category", 0))
	query.Filter(elastic.NewTermQuery("pay_state", 1))
	aggs := elastic.NewSumAggregation().Field("amount")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("sum_amount", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_amount")
	amount := *agg.Value

	type r struct {
		Amount float64
	}
	var sumR r
	dao.DbSlave.Model(&model.InsureOrder{}).Where("soft_sn = ? AND pay_state = ?", softSn, 1).Select("sum(duration_amount) AS amount").First(&sumR)

	return Decimal(amount + sumR.Amount)
}

// 获取设备服务商收入
func GetDevAgentAmount(softSn string, agentId int) float64 {
	/*var sumR struct{Amount float64}
	dao.DbSlave.Model(&model.BillAgent{}).Where("soft_sn = ? AND agent_id = ?", softSn, agentId).Select("sum(amount) AS amount").First(&sumR)
	return Decimal(sumR.Amount)*/

	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewTermQuery("agent_id", agentId))
	aggs := elastic.NewSumAggregation().Field("amount")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("bill_agent").
		Query(query).
		Aggregation("sum_amount", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_amount")
	amount := *agg.Value
	return Decimal(amount)
}

// 获取设备30天支付流水
func GetDev30DayAmount(softSn string) float64 {
	startTime := int(time.Now().Unix()) - 30*24*3600
	category := ToInterfaceSlice([]int{model.CatExam, model.CAT_RECHARGE})

	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermsQuery("category", category...))
	query.Filter(elastic.NewTermQuery("pay_state", 1))
	query.Filter(elastic.NewRangeQuery("created_at").Gte(startTime))
	aggs := elastic.NewSumAggregation().Field("amount")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("sum_amount", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_amount")
	amount := *agg.Value

	var sumR struct {
		Amount float64
	}
	dao.DbSlave.Model(&model.InsureOrder{}).Where("soft_sn = ? AND pay_state = ? AND created_at >= ?", softSn, 1, startTime).Select("sum(duration_amount) AS amount").First(&sumR)

	return Decimal(amount + sumR.Amount)
}

// 近30天免费播报时长
func GetDev30DayFreeDuration(softSn string) int {
	startTime := int(time.Now().Unix()) - 30*24*3600
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewRangeQuery("created_at").Gte(startTime))
	aggs := elastic.NewSumAggregation().Field("duration")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base1").
		Query(query).
		Aggregation("sum_duration", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_duration")
	amount := *agg.Value
	return int(amount)
}

// 获取设备未支付金额
func GetDevUnpaidAmount(softSn string) float64 {
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewMatchQuery("soft_sn", softSn))
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermQuery("category", 0))
	query.Filter(elastic.NewTermQuery("pay_state", 0))
	aggs := elastic.NewSumAggregation().Field("amount")
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("sum_amount", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	agg, _ := result.Aggregations.Sum("sum_amount")
	return *agg.Value
}

// 设备登录教练id
func DevCoachIds(softSn string) []int {
	var r []struct {
		CoachId int
	}
	dao.DbSlave.Model(&model.CoachLogin{}).Where("soft_sn = ?", softSn).Distinct("coach_id").Find(&r)
	var coachIds []int
	for _, coachId := range r {
		coachIds = append(coachIds, coachId.CoachId)
	}
	return coachIds
}

// 设备登录学员id
func DevStudentIds(softSn string) []int {
	var r []struct {
		StudentId int
	}
	dao.DbSlave.Model(&model.StudentLogin{}).Where("soft_sn = ?", softSn).Distinct("student_id").Find(&r)
	var studentIds []int
	for _, studentId := range r {
		studentIds = append(studentIds, studentId.StudentId)
	}
	return studentIds
}
