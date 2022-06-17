package dao

import (
	"github.com/spf13/viper"
	"gopkg.in/olivere/elastic.v5"
)

//var Es *elastic.Client

func Es() *elastic.Client {
	addrs := viper.GetStringSlice("es.addrs")
	//fmt.Println(addrs)
	client, err := elastic.NewClient(

		elastic.SetURL(addrs...),
	//elastic.SetHealthcheck(false),
	//elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	return client
}

/*func InitEs() {
	initEs(&Es)
}

func initEs(client **elastic.Client)  {
	var err error
	esConfigMap = make(map[string]interface{}, 1)
	esConfigMap = viper.GetStringMap("es")
	fmt.Println(esConfigMap)
	*client, err = elastic.NewClient(elastic.SetURL(esConfigMap["addr"].(string)))
	if err != nil {
		panic(err)
	}
}*/
