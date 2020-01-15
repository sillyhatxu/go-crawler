package springfieldspringfield

import (
	"fmt"
	"testing"
)

func TestCrawlerScript(t *testing.T) {
	result := CrawlerScriptURL("https://www.springfieldspringfield.co.uk/episode_scripts.php?tv-show=good-luck-charlie")
	fmt.Println(fmt.Printf("%#v", result))
}
