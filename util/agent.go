package util

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
	"rem_cli/dao"
	"rem_cli/model"
	"strconv"
	"time"
)

// 获取服务商付费时长
func GetPaidDuration(agent model.AgentBase, startTime int, endTime int) (int, int, int, int) {
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewRangeQuery("amount").Gt(0))
	if agent.Level == 1 {
		query.Filter(elastic.NewTermQuery("agent_id", agent.Id))
	} else {
		query.Filter(elastic.NewTermQuery("agent_sub_id", agent.Id))
	}

	if startTime > 0 {
		timeObj := time.Unix(int64(startTime), 0)
		ym := fmt.Sprintf("%d%02d", timeObj.Year(), timeObj.Month())
		query.Filter(elastic.NewTermQuery("ym", ym))
	}

	/*aggs := elastic.NewTermsAggregation().Field("subject_type")
	hpAggs := elastic.NewTermsAggregation().Field("is_hp")
	hpAggs.SubAggregation("sum_duration", elastic.NewSumAggregation().Field("duration"))
	aggs.SubAggregation("hp_aggs", hpAggs)*/

	aggs := elastic.NewTermsAggregation().Field("subject_type")
	aggs.SubAggregation("sum_duration", elastic.NewSumAggregation().Field("duration"))

	ctx := context.Background()
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("bill_agent").
		Query(query).
		Aggregation("subject_type_aggs", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println(agent.Id, dao.Es())
		panic(err)
	}

	k2Duration, k3Duration, k3HpDuration, k3LpDuration := 0, 0, 0, 0
	agg, _ := result.Aggregations.Terms("subject_type_aggs")
	for _, bucket := range agg.Buckets {
		a, _ := bucket.Aggregations.Sum("sum_duration")
		if int(bucket.Key.(float64)) == 2 {
			k2Duration = int(*a.Value)
		} else {
			k3Duration = int(*a.Value)
		}
	}

	query.Filter(elastic.NewTermQuery("subject_type", 3))
	query.Filter(elastic.NewTermQuery("is_hp", 1))
	resultHp, err := dao.Es().Search().
		Index("rem_db").
		Type("bill_agent").
		Query(query).
		Aggregation("sum_duration", elastic.NewSumAggregation().Field("duration")).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println(agent.Id, dao.Es())
		panic(err)
	}
	aggHp, _ := resultHp.Aggregations.Sum("sum_duration")
	k3HpDuration = int(*aggHp.Value)
	k3LpDuration = k3Duration - k3HpDuration
	return k2Duration, k3Duration, k3HpDuration, k3LpDuration
}

// 获取一级服务商总收入
func GetAgentAmount(agent model.AgentBase, startTime int, endTime int) (float64, float64, float64, float64, float64) {
	amount, k2Amount, k3Amount, k3HpAmount, k3LpAmount := getAgentOrderAmount(agent, startTime, endTime)
	amount += getAgentSchoolAmount(agent, startTime, endTime)
	amount += getAgentInsureAmount(agent, startTime)
	return Decimal(amount), Decimal(k2Amount), Decimal(k3Amount), Decimal(k3HpAmount), Decimal(k3LpAmount)
}

// 获取一级服务商订单收入
func getAgentOrderAmount(agent model.AgentBase, startTime int, endTime int) (float64, float64, float64, float64, float64) {
	var amount, k2Amount, k3Amount, k3HpAmount, k3LpAmount float64
	where := "agent_id = " + strconv.Itoa(agent.Id)
	if startTime > 0 {
		timeObj := time.Unix(int64(startTime), 0)
		ym := fmt.Sprintf("%d%02d", timeObj.Year(), timeObj.Month())
		where += " AND ym = " + ym
	}

	type r struct {
		Amount      float64
		SubjectType int
	}
	var sumR []r
	dao.DbSlave.Model(&model.BillAgent{}).Where(where).Select("sum(amount) as amount, subject_type").Group("subject_type").Find(&sumR)

	k3Amount, k3HpAmount, k3LpAmount = 0, 0, 0
	for _, s := range sumR {
		if s.SubjectType == 3 {
			k3Amount = s.Amount
		} else if s.SubjectType == 2 {
			k2Amount = s.Amount
		}
	}
	amount = k3Amount + k2Amount

	dao.DbSlave.Model(&model.BillAgent{}).Where(where + " AND is_hp = 1 AND subject_type = 3").Select("sum(amount) as amount").First(&k3HpAmount)
	k3LpAmount = k3Amount - k3HpAmount

	return amount, k2Amount, k3Amount, k3HpAmount, k3LpAmount
}

// 获取一级服务商订单总额
func getOrderAmount(agent model.AgentBase, startTime int, endTime int) (float64, float64, float64) {
	var amount, k2Amount, k3Amount float64
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewTermQuery("agent_id", agent.Id))
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermQuery("pay_state", 1))
	query.Filter(elastic.NewTermQuery("category", model.CatExam))

	createdAtQuery := elastic.NewRangeQuery("paid_at")
	if startTime > 0 {
		createdAtQuery.Gte(startTime)
	}
	if endTime > 0 {
		createdAtQuery.Lt(endTime)
	}
	if startTime > 0 || endTime > 0 {
		query.Filter(createdAtQuery)
	}

	aggs := elastic.NewTermsAggregation().Field("subject_type")
	aggs.SubAggregation("sum_amount", elastic.NewSumAggregation().Field("amount"))

	ctx := context.Background()
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("sum_amount", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println(agent.Id, dao.Es())
		panic(err)
	}
	agg, _ := result.Aggregations.Terms("sum_amount")
	for _, bucket := range agg.Buckets {
		a, _ := bucket.Aggregations.Sum("sum_amount")
		if bucket.Key.(float64) == 3 || bucket.Key.(float64) == 4 {
			k3Amount += *a.Value
		} else {
			k2Amount = *a.Value
		}
	}
	amount = k3Amount + k2Amount
	return amount, k2Amount, k3Amount
}

// 获取一级服务商驾校收入
func getAgentSchoolAmount(agent model.AgentBase, startTime int, endTime int) float64 {
	where := "agent_id = " + strconv.Itoa(agent.Id)
	if startTime > 0 {
		where += " AND created_at >= " + strconv.Itoa(startTime)
	}
	if endTime > 0 {
		where += " AND created_at < " + strconv.Itoa(endTime)
	}

	type r struct {
		Amount float64
	}
	var sumR r
	dao.DbSlave.Model(&model.SchoolAssignLog{}).Where(where).Select("sum(agent_amount) as amount").First(&sumR)
	return sumR.Amount
}

// 获取一级服务商保险收入
func getAgentInsureAmount(agent model.AgentBase, startTime int) float64 {
	where := "agent_id = " + strconv.Itoa(agent.Id)
	if startTime > 0 {
		timeObj := time.Unix(int64(startTime), 0)
		ym := fmt.Sprintf("%d%02d", timeObj.Year(), timeObj.Month())
		where += " AND ym = " + ym
	}

	type r struct {
		Amount float64
	}
	var sumR r
	dao.DbSlave.Model(&model.InsureOrderBill{}).Where(where).Select("sum(agent_amount) as amount").First(&sumR)
	return sumR.Amount
}

// 获取二级服务商保险收入
func getAgentSubInsureAmount(agent model.AgentBase, startTime int) float64 {
	where := "agent_sub_id = " + strconv.Itoa(agent.Id)
	if startTime > 0 {
		timeObj := time.Unix(int64(startTime), 0)
		ym := fmt.Sprintf("%d%02d", timeObj.Year(), timeObj.Month())
		where += " AND ym = " + ym
	}

	type r struct {
		Amount float64
	}
	var sumR r
	dao.DbSlave.Model(&model.InsureOrderBill{}).Where(where).Select("sum(agent_sub_amount) as amount").First(&sumR)
	return sumR.Amount
}

// 获取二级服务商收入总和
func GetAgentSubAmount(agent model.AgentBase, devSofts []model.DevSoft, startTime int, endTime int) (float64, float64, float64) {
	amount, k2Amount, k3Amount := getAgentSubOrderAmount(agent, devSofts, startTime, endTime)
	amount += getAgentSubInsureAmount(agent, startTime)
	return Decimal(amount), Decimal(k2Amount), Decimal(k3Amount)
}

// 获取二级服务商练车订单收入
func getAgentSubOrderAmount(agent model.AgentBase, devSofts []model.DevSoft, startTime int, endTime int) (float64, float64, float64) {
	var amount, k2Amount, k3Amount float64

	ctx := context.Background()
	query := elastic.NewBoolQuery()

	category := ToInterfaceSlice([]int{model.CatExam, model.CAT_RECHARGE})

	query.Filter(elastic.NewTermQuery("agent_sub_id", agent.Id))
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermQuery("pay_state", 1))
	query.Filter(elastic.NewTermsQuery("category", category...))
	createdAtQuery := elastic.NewRangeQuery("paid_at")
	if startTime > 0 {
		createdAtQuery.Gte(startTime)
	}
	if endTime > 0 {
		createdAtQuery.Lt(endTime)
	}
	if startTime > 0 || endTime > 0 {
		query.Filter(createdAtQuery)
	}

	aggsSize := 10
	if len(devSofts) > 10 {
		aggsSize = len(devSofts)
	}
	aggs := elastic.NewTermsAggregation().Field("soft_sn.keyword").Size(aggsSize)
	subjectAggs := elastic.NewTermsAggregation().Field("subject_type")
	subjectAggs.SubAggregation("sum_duration", elastic.NewSumAggregation().Field("duration"))
	aggs.SubAggregation("subject_type_aggs", subjectAggs)

	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("soft_sn_aggs", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println(agent.Id, dao.Es())
		panic(err)
	}
	agg, _ := result.Aggregations.Terms("soft_sn_aggs")
	for _, bucket := range agg.Buckets {
		var devSoft model.DevSoft
		dao.DbSlave.Model(&model.DevSoft{}).Where("soft_sn = ?", bucket.Key.(string)).First(&devSoft)
		k3Price := GetAgentPrice(agent, devSoft.DistrictId, devSoft.CityId, 3)
		k2Price := GetAgentPrice(agent, devSoft.DistrictId, devSoft.CityId, 2)
		k3Duration, k2Duration := 0, 0
		subjectTypeAgg, _ := bucket.Aggregations.Terms("subject_type_aggs")
		for _, subjectBucket := range subjectTypeAgg.Buckets {
			a, _ := subjectBucket.Aggregations.Sum("sum_duration")
			if subjectBucket.Key.(float64) == 3 || subjectBucket.Key.(float64) == 4 {
				k3Duration += int(*a.Value)
			} else {
				k2Duration = int(*a.Value)
			}
		}
		k2Amount += k2Price * float64(k2Duration) / 60
		k3Amount += k3Price * float64(k3Duration) / 60
	}
	amount = k2Amount + k3Amount
	return amount, k2Amount, k3Amount
}

func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// 获取服务商价格
func GetAgentPrice(agent model.AgentBase, districtId int, cityId int, subjectType int) float64 {
	var agentCity model.AgentCity
	dao.DbSlave.Model(&model.AgentCity{}).Where("agent_id = ? AND pcd = ?", agent.Id, districtId).First(&agentCity)
	if agentCity.Id == 0 {
		dao.DbSlave.Model(&model.AgentCity{}).Where("agent_id = ? AND pcd = ?", agent.Id, cityId).First(&agentCity)
	}
	if agentCity.Id == 0 {
		return float64(0)
	}
	var price float64
	if subjectType == 2 {
		price = agentCity.K2Price
	} else {
		price = agentCity.Price
	}
	return price
}

// 获取服务商待支付
func GetAgentUnpaid(agent model.AgentBase) (float64, float64, float64, int64, int64, int64) {
	query := elastic.NewBoolQuery()
	if agent.Level == 1 {
		query.Filter(elastic.NewTermQuery("agent_id", agent.Id))
	} else {
		query.Filter(elastic.NewTermQuery("agent_sub_id", agent.Id))
	}
	query.Filter(elastic.NewTermQuery("is_free", 0))
	query.Filter(elastic.NewTermQuery("pay_state", 0))

	aggs := elastic.NewTermsAggregation().Field("subject_type")
	aggs.SubAggregation("sum_amount", elastic.NewSumAggregation().Field("amount"))

	ctx := context.Background()
	result, err := dao.Es().Search().
		Index("rem_db").
		Type("order_base").
		Query(query).
		Aggregation("subject_aggs", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		fmt.Println(agent.Id, dao.Es())
		panic(err)
	}

	var k2UnpaidAmount, k3UnpaidAmount, unpaidAmount float64
	var k2UnpaidNum, k3UnpaidNum, unpaidNum int64
	agg, _ := result.Aggregations.Terms("subject_aggs")
	for _, bucket := range agg.Buckets {
		a, _ := bucket.Aggregations.Sum("sum_amount")
		if bucket.Key.(float64) == 3 || bucket.Key.(float64) == 4 {
			k3UnpaidAmount += *a.Value
			k3UnpaidNum += bucket.DocCount
		} else {
			k2UnpaidAmount = *a.Value
			k2UnpaidNum = bucket.DocCount
		}
	}
	unpaidAmount = k2UnpaidAmount + k3UnpaidAmount
	unpaidNum = k2UnpaidNum + k3UnpaidNum

	return Decimal(unpaidAmount), Decimal(k2UnpaidAmount), Decimal(k3UnpaidAmount), unpaidNum, k2UnpaidNum, k3UnpaidNum
}

// 获取服务商设备号
func GetAgentSoftSns(agent model.AgentBase) []string {
	where := ""
	if agent.Level == 1 {
		where = "agent_id = " + strconv.Itoa(agent.Id)
	} else {
		where = "agent_sub_id = " + strconv.Itoa(agent.Id)
	}
	where += " AND state = 1 AND soft_sn != ''"
	var devSofts []model.DevSoft
	dao.DbSlave.Where(where).Find(&devSofts)
	var softSns []string
	for _, d := range devSofts {
		if d.SoftSn != "" {
			softSns = append(softSns, d.SoftSn)
		}
	}
	return softSns
}

func GetAgentTotalSoftSns(agent model.AgentBase) []string {
	where := ""
	if agent.Level == 1 {
		where = "agent_id = " + strconv.Itoa(agent.Id)
	} else {
		where = "agent_sub_id = " + strconv.Itoa(agent.Id)
	}
	where += " AND state in (1, 2)"
	var devSofts []model.DevSoft
	dao.DbSlave.Where(where).Or("state = 0 AND ec_agent_id = ?", agent.Id).Find(&devSofts)
	var softSns []string
	for _, d := range devSofts {
		softSns = append(softSns, d.SoftSn)
	}
	return softSns
}

// 获取服务商设备活跃率
func GetAgentActiveRate(agent model.AgentBase, softSns []string, startTime int, endTime int) float64 {
	if len(softSns) == 0 {
		return 0
	}

	//缩小取数范围
	where := "id > 0"
	tempIdx := model.DevSoftOpt{}
	dao.DbSlave.Model(&model.DevSoftOpt{}).Where("start_time > " + strconv.Itoa(startTime)).Take(&tempIdx)
	if tempIdx.Id > 0 {
		where = "id >" + strconv.Itoa(tempIdx.Id)
	}

	if agent.Level == 1 {
		where += " AND agent_id = " + strconv.Itoa(agent.Id)
	} else {
		where += " AND agent_sub_id = " + strconv.Itoa(agent.Id)
	}

	if startTime > 0 {
		where += " AND start_time > " + strconv.Itoa(startTime)
	}
	if endTime > 0 {
		where += " AND start_time < " + strconv.Itoa(endTime)
	}

	var activeNum int64
	where += " AND soft_sn in ?"
	dao.DbSlave.Model(&model.DevSoftOpt{}).Where(where, softSns).Distinct("soft_sn").Count(&activeNum)
	devNum := len(softSns)
	var activeRate float64
	if activeNum == 0 {
		activeRate = 0
	} else {
		activeRate, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(activeNum)/float64(devNum)), 64)
	}
	return activeRate
}

// 获取服务商设备付费率
func GetAgentPaidRate(agent model.AgentBase, softSns []string, startTime int, endTime int) float64 {
	if len(softSns) == 0 {
		return 0
	}

	//缩小取数范围
	where := "id > 0"
	tempIdx := model.DevSoftOpt{}
	dao.DbSlave.Model(&model.DevSoftOpt{}).Where("start_time > " + strconv.Itoa(startTime)).Take(&tempIdx)
	if tempIdx.Id > 0 {
		where = "id >" + strconv.Itoa(tempIdx.Id)
	}

	if agent.Level == 1 {
		where += " AND agent_id = " + strconv.Itoa(agent.Id)
	} else {
		where += " AND agent_sub_id = " + strconv.Itoa(agent.Id)
	}

	if startTime > 0 {
		where += " AND start_time > " + strconv.Itoa(startTime)
	}
	if endTime > 0 {
		where += " AND start_time < " + strconv.Itoa(endTime)
	}

	var paidNum int64
	where += " AND soft_sn in ? AND is_pay = 1"
	dao.DbSlave.Model(&model.DevSoftOpt{}).Where(where, softSns).Distinct("soft_sn").Count(&paidNum)
	devNum := len(softSns)
	var paidRate float64
	if paidNum == 0 {
		paidRate = 0
	} else {
		paidRate, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(paidNum)/float64(devNum)), 64)
	}
	return paidRate
}
