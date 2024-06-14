package lang

type Iso639 int

const (
	AF Iso639 = iota
	AR
	AZ
	BE
	BG
	BN
	BS
	CA
	CS
	CY
	DA
	DE
	EL
	EN
	EO
	ES
	ET
	EU
	FA
	FI
	FR
	GA
	GU
	HE
	HI
	HR
	HU
	HY
	ID
	IS
	IT
	JA
	KA
	KK
	KO
	LA
	LG
	LT
	LV
	MI
	MK
	MN
	MR
	MS
	NB
	NL
	NN
	PA
	PL
	PT
	RO
	RU
	SK
	SL
	SN
	SO
	SQ
	SR
	ST
	SV
	SW
	TA
	TE
	TH
	TL
	TN
	TR
	TS
	UK
	UR
	VI
	XH
	YO
	ZH
	ZU
)

var languageToEnglishMap = map[Iso639]string{
	AF: "Afrikaans",
	AR: "Arabic",
	AZ: "Azerbaijani",
	BE: "Belarusian",
	BG: "Bulgarian",
	BN: "Bengali",
	BS: "Bosnian",
	CA: "Catalan",
	CS: "Czech",
	CY: "Welsh",
	DA: "Danish",
	DE: "German",
	EL: "Greek",
	EN: "English",
	EO: "Esperanto",
	ES: "Spanish",
	ET: "Estonian",
	EU: "Basque",
	FA: "Persian",
	FI: "Finnish",
	FR: "French",
	GA: "Irish",
	GU: "Gujarati",
	HE: "Hebrew",
	HI: "Hindi",
	HR: "Croatian",
	HU: "Hungarian",
	HY: "Armenian",
	ID: "Indonesian",
	IS: "Icelandic",
	IT: "Italian",
	JA: "Japanese",
	KA: "Georgian",
	KK: "Kazakh",
	KO: "Korean",
	LA: "Latin",
	LG: "Ganda",
	LT: "Lithuanian",
	LV: "Latvian",
	MI: "Maori",
	MK: "Macedonian",
	MN: "Mongolian",
	MR: "Marathi",
	MS: "Malay",
	NB: "Bokmal",
	NL: "Dutch",
	NN: "Nynorsk",
	PA: "Punjabi",
	PL: "Polish",
	PT: "Portuguese",
	RO: "Romanian",
	RU: "Russian",
	SK: "Slovak",
	SL: "Slovene",
	SN: "Shona",
	SO: "Somali",
	SQ: "Albanian",
	SR: "Serbian",
	ST: "Sotho",
	SV: "Swedish",
	SW: "Swahili",
	TA: "Tamil",
	TE: "Telugu",
	TH: "Thai",
	TL: "Tagalog",
	TN: "Tswana",
	TR: "Turkish",
	TS: "Tsonga",
	UK: "Ukrainian",
	UR: "Urdu",
	VI: "Vietnamese",
	XH: "Xhosa",
	YO: "Yoruba",
	ZH: "Chinese",
	ZU: "Zulu",
}

type Detector interface {
	Detect(input string) Iso639
}

func (language Iso639) String() string {
	return languageToEnglishMap[language]
}
