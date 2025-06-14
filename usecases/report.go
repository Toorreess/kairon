package usecases

import (
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
	"log"
	"time"
)

type ReportUsecase interface {
	GenerateFinancialReport(startDate, endDate time.Time) (*model.FinancialReport, error)
}

type ReportUsecaseImp struct {
	OrderRepo repositories.OrderRepository
	// MemberRepo repositories.MemberRepository
}

func NewReportUsecase(or repositories.OrderRepository /*, mr repositories.MemberRepository*/) ReportUsecase {
	return &ReportUsecaseImp{
		OrderRepo: or,
		// MemberRepo: mr,
	}
}

func (uc *ReportUsecaseImp) GenerateFinancialReport(startDate, endDate time.Time) (*model.FinancialReport, error) {
	qo := infrastructure.QueryOpts{
		Limit:       5000,
		Offset:      0,
		QueryString: "status:paid",
	}

	allOrders, err := uc.OrderRepo.List(qo)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return nil, err
	}

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	var filteredOrders []model.Order
	for _, order := range allOrders {
		if order.Created >= startDate.Unix() && order.Created <= endDate.Unix() {
			filteredOrders = append(filteredOrders, order)
		}
	}

	totalSalesIncome := 0.0
	for _, order := range filteredOrders {
		totalSalesIncome += order.Amount
	}

	return &model.FinancialReport{
		StartDate:        startDate.Format("2006-01-02"),
		EndDate:          endDate.Format("2006-01-02"),
		TotalSales:       len(filteredOrders),
		TotalSalesIncome: totalSalesIncome,
	}, nil
}
