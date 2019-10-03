package longman

type Vocabulary struct {
	Word     string              `json:"word"`
	Explains []VocabularyExplain `json:"explains"`
}

type VocabularyExplain struct {
	Index          int               `json:"index"`
	Phonetic       string            `json:"phonetic"`
	PronounceUKURL string            `json:"pronounce_uk_url"`
	PronounceUSURL string            `json:"pronounce_us_url"`
	POS            string            `json:"pos"` //词性 Part of Speech
	Senses         []VocabularySense `json:"senses"`
}

type VocabularySense struct {
	Index      int                      `json:"index"`
	Definition string                   `json:"definition"`
	Examples   []VocabularySenseExample `json:"example"`
}

type VocabularySenseExample struct {
	Example  string `json:"example"`
	SoundURL string `json:"sound_url"`
}
