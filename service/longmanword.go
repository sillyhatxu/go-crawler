package service

import (
	"encoding/json"
	"fmt"
	"github.com/sillyhatxu/go-crawler/config"
	longmanword "github.com/sillyhatxu/go-crawler/crawler/longman"
	"github.com/sillyhatxu/retry-utils"
	"github.com/sillyhatxu/word-backend/grpc/grpcclient"
	"github.com/sillyhatxu/word-backend/grpc/longman"
	"github.com/sillyhatxu/word-backend/grpc/word"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"time"
)

func AutoCrawlerLongmanWord() {
	rows := 370103
	pageSize := 5000
	pageSum := (rows + pageSize - 1) / pageSize
	for page := 1; page < pageSum; page++ {
		offset := page * pageSize
		limit := pageSize
		wordArray, err := FindWordList(int64(offset), int64(limit))
		if err != nil {
			continue
		}
		for _, w := range wordArray {
			logrus.Infof("search word : %s", w.Word)
			description := ""
			status := longman.Status_Error
			vocabulary, err := longmanword.QueryVocabulary(w.Word)
			if err != nil {
				logrus.Errorf("QueryVocabulary error. %v", err)
				description = err.Error()
			}
			if vocabulary != nil {
				vocabularyJSON, err := json.Marshal(vocabulary)
				if err != nil {
					logrus.Errorf("vocabulary to json error. %v", err)
					description = err.Error()
				}
				if description == "" {
					description = string(vocabularyJSON)
					status = longman.Status_Success
				}
			}
			err = AddVocabulary(&longman.AddRequest{
				VocabularyId: w.Id,
				Description:  description,
				Status:       status,
			})
			if err != nil {
				logrus.Errorf("add vocabulary error. %v", err)
			}
			time.Sleep(5 * time.Second)
		}
		//fmt.Println(array[start:end])
	}
}

func AddVocabulary(request *longman.AddRequest) error {
	client := grpcclient.NewGRPCWordClient(config.Conf.EnvConfig.InternalWordGrpcHost)
	res, err := client.AddLongmanWord(request)
	if err != nil {
		logrus.Errorf("add longman word error. %v", err)
	}
	if res.Code != uint32(codes.OK) {
		return fmt.Errorf("add longman word code error. response : %v", res)
	}
	return nil
}
func FindWordList(offset, limit int64) ([]*word.Word, error) {
	client := grpcclient.NewGRPCWordClient(config.Conf.EnvConfig.InternalWordGrpcHost)
	var response *word.WordListResponse
	err := retry.Do(func() error {
		res, err := client.WordList(&word.ListRequest{Offset: offset, Limit: limit})
		if err != nil {
			return err
		}
		response = res
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response.Data, nil
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
