package service

import (
	"github.com/sillyhatxu/word-backend/grpc/grpcclient"
	"github.com/sillyhatxu/word-backend/grpc/word"
)

func AutoCrawlerLongmanWord() {
	rows := 370103
	pageSize := 5000
	pageSum := (rows + pageSize - 1) / pageSize
	for page := 0; page < pageSum; page++ {
		offset := page * pageSize
		limit := pageSize

		//fmt.Println(array[start:end])
	}
}

func Find(address string, offset, limit int64) (*word.WordListResponse, error) {
	client := grpcclient.NewGRPCWordClient(address)
	return client.WordList(&word.ListRequest{Offset: offset, Limit: limit})
}

//limit %d,%d ", params.Offset, params.Limit)
//select * from vocabulary limit 0,10;
//select * from vocabulary limit 10,10;
//select * from vocabulary limit 20,10;
//select * from vocabulary limit 30,10;
//select * from vocabulary limit 40,10;

//type CrawlerLongmanWord struct {
//	Name string
//}
//
//func (clw CrawlerLongmanWord) Execute() {
//	logrus.Infof("name : %s start %v", clw.Name, time.Now().Format("2006-01-02T15:04:05"))
//}
