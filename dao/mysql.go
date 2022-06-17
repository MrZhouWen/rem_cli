package dao

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

var mysqlConfigMap map[string]interface{}

var (
	DbSlave  *gorm.DB
	Db       *gorm.DB
	DbS      *gorm.DB
	DbSSlave *gorm.DB
)

func InitMysql() {
	initMysql("master", &Db)
	initMysql("slave", &DbSlave)
	initMysql("s_master", &DbS)
	initMysql("s_slave", &DbSSlave)
}

func initMysql(node string, db **gorm.DB) {
	var err error
	mysqlConfigMap = make(map[string]interface{}, 10)
	mysqlConfigMap = viper.GetStringMap("mysql." + node)
	*db, err = gorm.Open(mysql.Open(getMysqlUrl()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func getMysqlUrl() string {
	return mysqlConfigMap["user"].(string) + ":" + mysqlConfigMap["password"].(string) + "@tcp(" + mysqlConfigMap["host"].(string) + ":" + strconv.FormatInt(mysqlConfigMap["port"].(int64), 10) + ")/" + mysqlConfigMap["name"].(string)
}
