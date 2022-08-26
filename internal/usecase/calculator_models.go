package usecase

type statementSummary struct {
	total        float64
	avgCredit    float64
	avgDebit     float64
	monthSummary map[string]int
}

type calculator struct {
	dirPath    string
	apikey     string
	templateID string
}

type templateData struct {
	Email          string           `json:"email"`
	Name           string           `json:"name"`
	TotalBalance   string           `json:"total_balance"`
	FirstMonthYear string           `json:"first_month_year"`
	LastMonthYear  string           `json:"last_month_year"`
	AvgDebit       string           `json:"avg_debit"`
	AvgCredit      string           `json:"avg_credit"`
	MonthSummary   []monthOperation `json:"month_summary"`
}
type monthOperation struct {
	Month        string `json:"month"`
	Transactions int    `json:"transactions"`
}
