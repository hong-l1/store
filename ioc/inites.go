package ioc

import "github.com/olivere/elastic/v7"

func InitEs() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}
	return client
}
