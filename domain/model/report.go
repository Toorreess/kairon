package model

type FinancialReport struct {
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	TotalSales       int     `json:"total_sales"`
	TotalSalesIncome float64 `json:"total_sales_income"`
	// TotalMembershipIncome float64 `json:"total_membership_income"`
	// GrandTotalIncome      float64 `json:"grand_total_income"`
}
