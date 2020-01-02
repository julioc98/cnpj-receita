package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/osext"
)

type receita struct {
	AtividadePrincipal []struct {
		Text string `json:"text"`
	} `json:"atividade_principal"`
	DataSituacao          string `json:"data_situacao"`
	Complemento           string `json:"complemento"`
	Tipo                  string `json:"tipo"`
	Nome                  string `json:"nome"`
	Telefone              string `json:"telefone"`
	Email                 string `json:"email"`
	Situacao              string `json:"situacao"`
	Bairro                string `json:"bairro"`
	Logradouro            string `json:"logradouro"`
	Numero                string `json:"numero"`
	Cep                   string `json:"cep"`
	Municipio             string `json:"municipio"`
	Fantasia              string `json:"fantasia"`
	Porte                 string `json:"porte"`
	Abertura              string `json:"abertura"`
	NaturezaJuridica      string `json:"natureza_juridica"`
	Uf                    string `json:"uf"`
	Cnpj                  string `json:"cnpj"`
	UltimaAtualizacao     string `json:"ultima_atualizacao"`
	Status                string `json:"status"`
	Efr                   string `json:"efr"`
	MotivoSituacao        string `json:"motivo_situacao"`
	SituacaoEspecial      string `json:"situacao_especial"`
	DataSituacaoEspecial  string `json:"data_situacao_especial"`
	AtividadesSecundarias []struct {
		Text string `json:"text"`
	} `json:"atividades_secundarias"`
	CapitalSocial string `json:"capital_social"`
	URL           string
}

func receitaToCSV(data receita) []string {

	csv := []string{
		data.Cnpj,
		data.AtividadePrincipal[0].Text,
		data.Nome,
		data.Telefone,
		data.Email,
		data.AtividadesSecundarias[0].Text,
		data.Situacao,
		data.Logradouro,
		data.Numero,
		data.Complemento,
		data.Bairro,
		data.Municipio,
		data.Uf,
		data.Cep,
		data.Tipo,
		data.Fantasia,
		data.Porte,
		data.Abertura,
		data.NaturezaJuridica,
		data.Cnpj,
		data.UltimaAtualizacao,
		data.Status,
		data.DataSituacao,
		data.MotivoSituacao,
		data.SituacaoEspecial,
		data.DataSituacaoEspecial,
		data.CapitalSocial,
	}
	return csv
}

func getReceita(cnpj string) (receita, error) {
	var res receita

	resp, err := http.Get("https://www.receitaws.com.br/v1/cnpj/" + cnpj)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func connected() bool {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func waitConnect() {
	if connected() {
		return
	}
	fmt.Println("Wait for internet connection...")
	for {
		if connected() {
			return
		}
	}
}

func main() {
	fmt.Println("Start")
	waitConnect()

	exeAbsolutePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new output file
	writeFile, err := os.Create(exeAbsolutePath + "/cnpj_out.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer writeFile.Close()

	// Open the file
	readFile, err := os.Open(exeAbsolutePath + "/cnpj.csv")
	if err != nil {
		log.Fatalln("1 - Couldn't open the csv file", err)
	}
	defer readFile.Close()

	// Parse the file
	r := csv.NewReader(readFile)

	// Create writer
	writer := csv.NewWriter(writeFile)
	defer writer.Flush()

	// Add Header
	writer.Write([]string{"CNPJ", "atividade_principal", "nome", "telefone", "email", "atividades_secundarias", "situacao", "logradouro", "numero",
		"complemento", "bairro", "municipio", "uf", "cep", "tipo", "fantasia", "porte", "abertura", "natureza_juridica", "cnpj", "ultima_atualizacao",
		"status", "efr", "situacao", "data_situacao", "motivo_situacao", "situacao_especial", "data_situacao_especial", "capital_social"})

	if err != nil {
		log.Fatal(err)
	}

	count := 0

	// Iterate through the records
	for {

		count++

		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			fmt.Println("\nEnd\nPress CTRL + C to EXIT")
			for {
			}
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nTurn: %d \n", count)
		cnpj := record[0]
		if len(cnpj) < 14 {
			fmt.Printf("%s is not a CNPJ \n", cnpj)
			continue
		}
		waitConnect()
		fmt.Printf("CNPJ: %s processing...\n", cnpj)
		waitConnect()

		data, err := getReceita(cnpj)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(data)

		writer.Write(receitaToCSV(data))
		writer.Flush()
		if count%3 == 0 {
			time.Sleep(61 * time.Second)
		}
	}

}
