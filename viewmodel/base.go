package viewmodel

type Data struct {
	Date  string  `json:"d"`
	Value float32 `json:"v"`
}

//Json es una estructura para devolver datos de dolar
type Json struct {
	Date      string  `json:"date"`
	Official  float32 `json:"officialdolar"`
	Blue      float32 `json:"bluedollar"`
	Variacion float32 `json:"variation"`
}
