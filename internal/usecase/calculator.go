package usecase

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sirupsen/logrus"
)

type Option func(c *calculator)

func WithDirPath(dirPath string) Option {
	return func(c *calculator) {
		c.dirPath = dirPath
	}
}

func WithAPIKey(apikey string) Option {
	return func(c *calculator) {
		c.apikey = apikey
	}
}

func WithTemplateID(templateID string) Option {
	return func(c *calculator) {
		c.templateID = templateID
	}
}

func NewCalculator(options ...Option) Calculator {
	c := &calculator{
		dirPath: "",
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c calculator) Run() {
	// read files in path
	files, err := ioutil.ReadDir(c.dirPath)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		logrus.Info("processing file: ", file.Name())
		// process each file
		sttmntSummary, err := processFile(path.Join(c.dirPath, file.Name()))
		if err != nil {
			logrus.Error(err)
		}
		tplData, err := getEmailTemplateData(file.Name(), sttmntSummary)
		err = c.SendMail(tplData)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Info("file processed OK")
		}
	}
}

func getEmailTemplateData(email string, statement statementSummary) (templateData, error) {
	var err error
	template := templateData{
		Email:        email,
		TotalBalance: fmt.Sprintf("%.2f", statement.total),
		AvgCredit:    fmt.Sprintf("%.2f", statement.avgCredit),
		AvgDebit:     fmt.Sprintf("%.2f", statement.avgDebit),
	}
	e := strings.Split(email, "@")
	template.Name = strings.Title(e[0])
	// sort monthSummary keys
	keys := make([]string, 0, len(statement.monthSummary))
	for k := range statement.monthSummary {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fMonth := strings.Split(keys[0], "-")
	template.FirstMonthYear, err = humanizeMonth(fMonth[1], fMonth[0])
	if err != nil {
		return templateData{}, err
	}

	lMonth := strings.Split(keys[len(keys)-1], "-")
	template.LastMonthYear, err = humanizeMonth(lMonth[1], lMonth[0])
	if err != nil {
		return templateData{}, err
	}

	for _, k := range keys {
		km := strings.Split(k, "-")
		month, err := humanizeMonth(km[1], km[0])
		if err != nil {
			return templateData{}, err
		}
		template.MonthSummary = append(template.MonthSummary, monthOperation{Month: month, Transactions: statement.monthSummary[k]})
	}
	return template, nil
}

func humanizeMonth(monthNumber, year string) (string, error) {
	month, err := strconv.Atoi(monthNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s of %s", time.Month(month).String(), year), nil
}

func (c calculator) SendMail(data templateData) error {
	request := sendgrid.GetRequest(c.apikey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	bodyTpl := `	
	{
		"personalizations": [
			{
				"to": [
					{
						"email": "%s"
					}
				],
				"dynamic_template_data":%s
			}
		],
		"from": {
			"email": "vellonce@gmail.com"
		},
		"template_id": "%s"
	}
	`
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	request.Body = []byte(fmt.Sprintf(bodyTpl, data.Email, string(jsonData), c.templateID))
	_, err = sendgrid.API(request)
	return err
}

func processFile(path string) (statementSummary, error) {
	var (
		firstLine = true
		debit     = []float64{}
		credit    = []float64{}
		result    = statementSummary{}
	)
	file, err := os.Open(path)
	if err != nil {
		return statementSummary{}, err
	}
	defer file.Close()
	parser := csv.NewReader(file)
	parser.FieldsPerRecord = 3
	result.monthSummary = make(map[string]int)
	// read the file line by line to save memory (versus reading all the file at once)
	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return statementSummary{}, err
		}
		// skip the first line (line with headers) of CSV
		if firstLine {
			firstLine = false
			continue
		}
		// add amounts to calculate total
		amount := getAmount(record[2])
		result.total += amount
		// distinguish between credit/debit to calculate averages
		if amount < 0 {
			credit = append(credit, amount)
		} else {
			debit = append(debit, amount)
		}
		// summarize operations per month
		date := string([]byte(record[1])[:7])
		_, found := result.monthSummary[date]
		if !found {
			result.monthSummary[date] = 1
		} else {
			result.monthSummary[date] += 1
		}
	}
	result.avgCredit = getAvg(credit)
	result.avgDebit = getAvg(debit)
	return result, nil
}

func getAmount(amnt string) float64 {
	operation := string([]byte(amnt)[0])
	amount, err := strconv.ParseFloat(string([]byte(amnt)[1:]), 64)
	if err != nil {
		logrus.Error(err)
	}
	if operation == "-" {
		amount *= -1
	}
	return amount
}

func getAvg(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
