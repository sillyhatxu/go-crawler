package longman

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/sillyhatxu/retry-utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	url = "https://www.ldoceonline.com/dictionary"
)

func QueryVocabulary(word string) (*Vocabulary, error) {
	var vocabulary *Vocabulary
	err := retry.Do(func() error {
		v, err := findVocabulary(word)
		if err != nil {
			return nil
		}
		vocabulary = v
		return nil
	}, retry.Attempts(10), retry.Delay(3*time.Second), retry.ErrorCallback(func(n uint, err error) {
		logrus.Errorf("retry [%d] find vocabulary {%s} error. %v", n, word, err)
	}))
	return vocabulary, err
}

func findVocabulary(word string) (*Vocabulary, error) {
	vocabulary := &Vocabulary{}
	c := colly.NewCollector(
		colly.AllowedDomains("ldoceonline.com", "www.ldoceonline.com"),
	)
	//在请求之前调用
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	//如果在请求期间发生错误则调用
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	//收到回复后调用
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	//OnResponse如果接收到的内容是HTML ，则在之后调用
	c.OnHTML("div[class=entry_content]", func(e *colly.HTMLElement) {

		//e.Request.Visit(e.Attr("href"))
		//classes := e.ChildAttrs("div", "class")//取出attr的值 class="dictionary" 取出dictionary
		logrus.Infof("word : %s", e.ChildText("h1.pagetitle"))
		vocabulary.Word = e.ChildText("h1.pagetitle")
		var explains []VocabularyExplain
		e.ForEach(".dictentry", func(_ int, el *colly.HTMLElement) { //取出class="dictentry"的组件
			//logrus.Infof("txt : %s", el.Text)
			logrus.Infof("word num : %s", el.ChildText("span.HOMNUM"))
			logrus.Infof("word phonetic : %s", el.ChildText("span.PRON"))
			logrus.Infof("word phonetic POS: %s", el.ChildText("span.POS"))
			//el.DOM.Find("span[class=frequent]").ChildAttrs("span[class=speaker brefile fa fa-volume-up hideOnAmp]","data-src-mp3")
			logrus.Infof("word phonetic UK url: %s", el.ChildAttr("span.brefile", "data-src-mp3"))
			logrus.Infof("word phonetic US url: %s", el.ChildAttr("span.amefile", "data-src-mp3"))

			var senses []VocabularySense
			el.ForEach(".Sense", func(_ int, senseElement *colly.HTMLElement) {
				index, err := strconv.Atoi(senseElement.ChildText("span.sensenum"))
				if err != nil {
					panic(err)
				}
				logrus.Infof("sense num : %s", senseElement.ChildText("span.sensenum"))
				logrus.Infof("sense txt : %s", senseElement.ChildText("span.DEF"))

				var examples []VocabularySenseExample
				logrus.Infof("sense example : %s", senseElement.ChildText("span.EXAMPLE"))
				senseElement.ForEach(".GramExa", func(_ int, gramExaElement *colly.HTMLElement) {
					examples = append(examples, VocabularySenseExample{
						Example:  gramExaElement.ChildText("span.EXAMPLE"),
						SoundURL: gramExaElement.ChildAttr("span.speaker", "data-src-mp3"),
					})
				})

				senses = append(senses, VocabularySense{
					Index:      index,
					Definition: senseElement.ChildText("span.DEF"),
					Examples:   examples,
				})

			})

			el.ChildText("span.EXAMPLE")

			//logrus.Infof("word phonetic UK and US url: %s", el.ChildAttrs("span","data-src-mp3"))
			//logrus.Infof("word phonetic US url: %s", el.ChildAttr("span","amefile"))
			//mail := Mail{
			//	Title:   el.ChildText("td:nth-of-type(1)"),
			//	Link:    el.ChildAttr("td:nth-of-type(1)", "href"),
			//	Author:  el.ChildText("td:nth-of-type(2)"),
			//	Date:    el.ChildText("td:nth-of-type(3)"),
			//	Message: el.ChildText("td:nth-of-type(4)"),
			//}
			//threads[threadSubject] = append(threads[threadSubject], mail)
			index, err := strconv.Atoi(el.ChildText("span.HOMNUM"))
			if err != nil {
				panic(err)
			}
			explains = append(explains, VocabularyExplain{
				Index:          index,
				Phonetic:       el.ChildText("span.PRON"),
				POS:            el.ChildText("span.POS"),
				PronounceUKURL: el.ChildAttr("span.brefile", "data-src-mp3"),
				PronounceUSURL: el.ChildAttr("span.amefile", "data-src-mp3"),
				Senses:         senses,
			})
			logrus.Infof(".......... word end ..........")
		})
		vocabulary.Explains = explains
		vocabularyJSON, err := json.Marshal(vocabulary)
		if err != nil {
			panic(err)
		}
		fmt.Println("------------------------------")
		fmt.Println(string(vocabularyJSON))
		fmt.Println("------------------------------")
	})
	//c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
	//	fmt.Println("---- OnHTML -----")
	//	fmt.Println("---- OnHTML tr td:nth-of-type(1) -----")
	//	fmt.Println(e.Text)
	//})

	//OnHTML如果接收到的内容是HTML或XML ，则在之后调用
	//c.OnXML("//h1", func(e *colly.XMLElement) {
	//	fmt.Println(e.Text)
	//})

	//在OnXML回调之后调用
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	err := c.Visit(fmt.Sprintf("%s/%s", url, word))
	if err != nil {
		return nil, err
	}
	return vocabulary, nil
}
