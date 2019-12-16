package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Quiero que hagan un programa que para una fecha determinada
// que yo ingrese por consola me diga la cotización del dólar oficial,
// la del blue y calcule cuál fue porcentaje de variación.

// Por otro lado, para el rango de fechas 28-10-2019 hasta hoy,
// quiero que busquen cuál fue el mejor día para hacer puré, es decir,
// comprar dólar oficial, venderlo al blue y quedarme con una diferencia.
// A su vez, quiero saber cuál fue el mejor día para comprar dólar blue.

const (
	authToken = "BEARER eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNDEyMTYsInR5cGUiOiJleHRlcm5hbCIsInVzZXIiOiJkNjI5MTE0MUB1cmhlbi5jb20ifQ.ZWGZVQ9hioNs909HdOMf4j3NruH1Uu3D8y0zOgp1cGX0IO-SEl2xaulLUDiR4HRyA7fpm-9JdCrzbyxkEaNUQA"
)

type Data struct {
	Date  string  `json:"d"`
	Value float32 `json:"v"`
}

var varDollar []Data
var sblueDollar []Data
var sofficialDollar []Data
var blueDollar map[string]float32
var officialDollar map[string]float32

func main() {

	sofficialDollar = getInfo("official")
	sblueDollar = getInfo("blue")

	//I create maps to perform searches more efficiently
	officialDollar = make(map[string]float32)
	for i := 0; i < len(sofficialDollar); i++ {
		officialDollar[sofficialDollar[i].Date] = sofficialDollar[i].Value
	}

	blueDollar = make(map[string]float32)
	for i := 0; i < len(sblueDollar); i++ {
		blueDollar[sblueDollar[i].Date] = sblueDollar[i].Value
	}

	serverInit()

	//I added a functionality to choose in console.
	// var option int
	// fmt.Println("Que accion desea realizar? Ingrese:\n 1: Ver cotizaciones y variacion\n 2: Hacer Pure o Comprar Blue")
	// fmt.Scanln(&option)
	// switch option {
	// case 1:
	// 	DollarXDay()
	// case 2:
	// 	Pure()
	// }
}

func serverInit() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/cartera", cartera)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var s string
	for k := range values {
		fmt.Println(k, values[k])
		s = DollarXDay(k)
	}
	w.Write([]byte(values.Get("param1")))
	w.Write([]byte(s))
}

func cartera(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("👜"))
}

//DollarXDay devuelve las cotizaciones y el porcentaje de variacion.
func DollarXDay(fecha string) string {
	var response string
	existsO, existsB := false, false

	if value, exists := officialDollar[fecha]; exists {
		fmt.Println("Dolar Oficial: ", value)
		response = response + fmt.Sprint("Dolar Oficial: ", value, "\n")
		existsO = true
	} else {
		fmt.Println("No existe valor para Dolar Oficial en la fecha", fecha)
		response = response + "No existe valor para Dolar Oficial en la fecha\n"
	}

	if value, exists := blueDollar[fecha]; exists {
		fmt.Println("Dolar Blue: ", value)
		response = response + fmt.Sprint("Dolar Blue: ", value, "\n")
		existsB = true
	} else {
		fmt.Println("No existe valor para Dolar Blue en la fecha", fecha)
		response = response + "No existe valor para Dolar Blue en la fecha\n"
	}

	if existsB && existsO {
		variacion := blueDollar[fecha] - officialDollar[fecha]
		variacion = variacion / officialDollar[fecha]
		fmt.Println("El porcentaje de variacion para la fecha ", fecha, " fue de ", variacion*100, "%")
		response = response + fmt.Sprint("El porcentaje de variacion para la fecha ", fecha, " fue de ", variacion*100, "%\n")
	} else {

		if existsB || existsO {
			fmt.Println("Faltan datos para hacer el calculo de variacion")
			response = response + "Faltan datos para hacer el calculo de variacion\n"
		} else {
			fmt.Println("No hay datos para la fecha ingresada, verifique y reingrese")
			response = response + "No hay datos para la fecha ingresada, verifique y reingrese\n"
		}
	}
	return response
}

//Pure
func Pure() {
	fechai := "2019-10-28"
	i := 0
	j := 0

	//Busco fecha para oficial
	for {
		if sofficialDollar[i].Date == fechai {
			break
		}
		i++
		if i == len(sofficialDollar) {
			fmt.Println("Fecha Incorrecta (Sin datos Dolar Oficial")
			os.Exit(1)
		}
	}

	//Busco fecha para blue
	for {
		if sblueDollar[j].Date == fechai {
			break
		}
		j++
		if j == len(sblueDollar) {
			fmt.Println("Fecha Incorrecta (Sin datos Dolar Blue")
			os.Exit(1)
		}
	}

	var difmax float32
	var difmin float32
	var fmax string
	var fmin string
	difmin = 100
	i++
	j++
	for {
		aux1 := sofficialDollar[i].Value
		aux2 := sblueDollar[j].Value
		aux := aux2 - aux1
		aux = aux / aux1
		aux = aux * 100
		if aux > difmax {
			difmax = aux
			fmax = sofficialDollar[i].Date
		}
		if aux < difmin {
			difmin = aux
			fmin = sofficialDollar[i].Date
		}
		i++
		if i == len(sofficialDollar) {
			break
		}
		j++
		if j == len(sblueDollar) {
			break
		}
	}

	fmt.Println("El mejor dia para hacer Pure fue el: ", fmax, " con un porcentaje de ", difmax, "% de diferencia")
	fmt.Println("El mejor dia para comprar Blue fue el: ", fmin, "con un porcentaje de ", difmin, "% de diferencia")
}

func getInfo(option string) []Data {

	var url string

	switch option {
	case "official":
		url = "https://api.estadisticasbcra.com/usd_of"
	case "blue":
		url = "https://api.estadisticasbcra.com/usd"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ErrorHandlerUrl("Falló la creación del request a la URL '%s', dando el error %v", url, err)
	}

	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		ErrorHandlerUrl("Falló el acceso a la URL '%s', dando el error %v", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorHandlerUrl("Falló el acceso al body de la respuesta de '%s', dando el error %v", url, err)
	}

	var data []Data
	_ = json.Unmarshal(body, &data)

	return data
}

//ErrorHandlerUrl gestiona los errores de forma mas corta %s->url %v->error
func ErrorHandlerUrl(message string, url string, err error) {
	fmt.Printf(message, url, err.Error())
	os.Exit(1)
}
