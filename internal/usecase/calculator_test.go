package usecase

import (
	"errors"
	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestNewCalculator(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "must instantiate a Calculator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCalculator(tt.args.options...)
			var i interface{} = c
			_, ok := i.(Calculator)
			assert.True(t, ok)
		})
	}
}

func TestWithAPIKey(t *testing.T) {
	type args struct {
		apikey string
	}
	tests := []struct {
		name string
		args args
		calc calculator
	}{
		{
			name: "must set the API key for the calculator",
			args: args{
				apikey: "abcdefg",
			},
			calc: calculator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithAPIKey(tt.args.apikey)
			opt(&tt.calc)
			assert.Equal(t, tt.calc.apikey, tt.args.apikey)
		})
	}
}

func TestWithDirPath(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name string
		args args
		calc calculator
	}{
		{
			name: "must set the dir path for calculator",
			args: args{
				dirPath: "some/path",
			},
			calc: calculator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				opt := WithDirPath(tt.args.dirPath)
				opt(&tt.calc)
				assert.Equal(t, tt.calc.dirPath, tt.args.dirPath)
			})
		})
	}
}

func TestWithTemplateID(t *testing.T) {
	type args struct {
		templateID string
	}
	tests := []struct {
		name string
		args args
		calc calculator
	}{
		{
			name: "must set the templateID for calculator",
			args: args{
				templateID: "template-id",
			},
			calc: calculator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTemplateID(tt.args.templateID)
			opt(&tt.calc)
			assert.Equal(t, tt.calc.templateID, tt.args.templateID)
		})
	}
}

func Test_calculator_SendMail(t *testing.T) {
	type globals struct {
		sendgridGetRequest func(key, endpoint, host string) rest.Request
		sendgridAPI        func(request rest.Request) (*rest.Response, error)
	}
	type args struct {
		data templateData
	}
	tests := []struct {
		name    string
		args    args
		globals globals
		wantErr bool
	}{
		{
			name: "success: must send an email",
			args: args{
				data: templateData{
					Email:          "test@email.com",
					Name:           "test",
					TotalBalance:   "34.74",
					FirstMonthYear: "July 2021",
					LastMonthYear:  "August 2021",
					AvgDebit:       "35.25",
					AvgCredit:      "-15.38",
					MonthSummary: []monthOperation{
						{
							Month:        "July of 2021",
							Transactions: 2,
						},
						{
							Month:        "August of 2021",
							Transactions: 2,
						},
					},
				},
			},
			globals: globals{
				sendgridGetRequest: func(key, endpoint, host string) rest.Request {
					return rest.Request{}
				},
				sendgridAPI: func(request rest.Request) (*rest.Response, error) {
					body := "{\"personalizations\": [{\"to\": [{\"email\": \"test@email.com\"}],\"dynamic_template_data\":{\"email\":\"test@email.com\",\"name\":\"test\",\"total_balance\":\"34.74\",\"first_month_year\":\"July 2021\",\"last_month_year\":\"August 2021\",\"avg_debit\":\"35.25\",\"avg_credit\":\"-15.38\",\"month_summary\":[{\"month\":\"July of 2021\",\"transactions\":2},{\"month\":\"August of 2021\",\"transactions\":2}]}}],\"from\": {\"email\": \"vellonce@gmail.com\"},\"template_id\": \"\"}"
					assert.Equal(t, body, string(request.Body))
					return &rest.Response{
						StatusCode: http.StatusAccepted,
					}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "failure: returns an error if couldn't send an email",
			args: args{
				data: templateData{
					Email:          "test@email.com",
					Name:           "test",
					TotalBalance:   "34.74",
					FirstMonthYear: "July 2021",
					LastMonthYear:  "August 2021",
					AvgDebit:       "35.25",
					AvgCredit:      "-15.38",
					MonthSummary: []monthOperation{
						{
							Month:        "July of 2021",
							Transactions: 2,
						},
						{
							Month:        "August of 2021",
							Transactions: 2,
						},
					},
				},
			},
			globals: globals{
				sendgridGetRequest: func(key, endpoint, host string) rest.Request {
					return rest.Request{}
				},
				sendgridAPI: func(request rest.Request) (*rest.Response, error) {
					body := "{\"personalizations\": [{\"to\": [{\"email\": \"test@email.com\"}],\"dynamic_template_data\":{\"email\":\"test@email.com\",\"name\":\"test\",\"total_balance\":\"34.74\",\"first_month_year\":\"July 2021\",\"last_month_year\":\"August 2021\",\"avg_debit\":\"35.25\",\"avg_credit\":\"-15.38\",\"month_summary\":[{\"month\":\"July of 2021\",\"transactions\":2},{\"month\":\"August of 2021\",\"transactions\":2}]}}],\"from\": {\"email\": \"vellonce@gmail.com\"},\"template_id\": \"\"}"
					assert.Equal(t, body, string(request.Body))
					return &rest.Response{
						StatusCode: http.StatusBadRequest,
					}, errors.New("something bad happened")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendgridAPI = tt.globals.sendgridAPI
			sendgridGetRequest = tt.globals.sendgridGetRequest
			c := calculator{}
			if err := c.SendMail(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getAmount(t *testing.T) {
	type args struct {
		amnt string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "should get a positive amount",
			args: args{
				amnt: "+100.55",
			},
			want: float64(100.55),
		},
		{
			name: "should get a negative amount",
			args: args{
				amnt: "-100.55",
			},
			want: float64(-100.55),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAmount(tt.args.amnt); got != tt.want {
				t.Errorf("getAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAvg(t *testing.T) {
	type args struct {
		vals []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "must get the average of the values passed",
			args: args{
				vals: []float64{
					float64(3),
					float64(3.5),
					float64(4),
					float64(3.5),
				},
			},
			want: float64(3.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAvg(tt.args.vals); got != tt.want {
				t.Errorf("getAvg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEmailTemplateData(t *testing.T) {
	type args struct {
		filename  string
		statement statementSummary
	}
	tests := []struct {
		name    string
		args    args
		want    templateData
		wantErr bool
	}{
		{
			name: "must get a template data",
			args: args{
				filename: "user@mail.com.csv",
				statement: statementSummary{
					total:     50,
					avgCredit: 35.5,
					avgDebit:  -40.66,
					monthSummary: map[string]int{
						"2022-08": 3,
						"2022-09": 5,
					},
				},
			},
			want: templateData{
				Email:          "user@mail.com",
				Name:           "User",
				TotalBalance:   "50.00",
				FirstMonthYear: "August of 2022",
				LastMonthYear:  "September of 2022",
				AvgDebit:       "-40.66",
				AvgCredit:      "35.50",
				MonthSummary: []monthOperation{
					{
						Month:        "August of 2022",
						Transactions: 3,
					},
					{
						Month:        "September of 2022",
						Transactions: 5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "must return an error if a month is incorrect",
			args: args{
				filename: "user@mail.com.csv",
				statement: statementSummary{
					total:     50,
					avgCredit: 35.5,
					avgDebit:  -40.66,
					monthSummary: map[string]int{
						"2022-eee": 3,
						"2022-09":  5,
					},
				},
			},
			want:    templateData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEmailTemplateData(tt.args.filename, tt.args.statement)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEmailTemplateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEmailTemplateData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_humanizeMonth(t *testing.T) {
	type args struct {
		monthNumber string
		year        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "must return the month and year in human readable format",
			args: args{
				monthNumber: "10",
				year:        "2033",
			},
			want:    "October of 2033",
			wantErr: false,
		},
		{
			name: "must return an error if the month is not valid",
			args: args{
				monthNumber: "sup",
				year:        "2033",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := humanizeMonth(tt.args.monthNumber, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("humanizeMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("humanizeMonth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    statementSummary
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("processFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
