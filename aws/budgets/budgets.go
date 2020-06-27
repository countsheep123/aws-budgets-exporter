package budgets

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/budgets"
)

type Budget struct {
	Name            string
	AccountID       string
	BudgetLimit     float64
	ActualSpend     *float64
	ForecastedSpend *float64
	LastUpdatedTime *time.Time
}

func (c *Client) GetBudgets(accountID string) ([]*Budget, error) {
	bs := []*Budget{}

	svc := budgets.New(c.cfg)
	input := &budgets.DescribeBudgetsInput{
		AccountId:  &accountID,
		MaxResults: aws.Int64(100),
		NextToken:  nil,
	}
	for {
		ctx := context.TODO()
		result, err := svc.DescribeBudgetsRequest(input).Send(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range result.Budgets {
			b, err := budget(&v)
			if err != nil {
				return nil, err
			}
			b.AccountID = accountID
			bs = append(bs, b)
		}
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
	}

	return bs, nil
}

func budget(b *budgets.Budget) (*Budget, error) {
	budgetLimit, err := strconv.ParseFloat(aws.StringValue(b.BudgetLimit.Amount), 64)
	if err != nil {
		return nil, err
	}

	var actualSpend *float64
	if b.CalculatedSpend.ActualSpend != nil {
		v, err := strconv.ParseFloat(aws.StringValue(b.CalculatedSpend.ActualSpend.Amount), 64)
		if err != nil {
			return nil, err
		}
		actualSpend = aws.Float64(v)
	}

	var forecastedSpend *float64
	if b.CalculatedSpend.ForecastedSpend != nil {
		v, err := strconv.ParseFloat(aws.StringValue(b.CalculatedSpend.ForecastedSpend.Amount), 64)
		if err != nil {
			return nil, err
		}
		forecastedSpend = aws.Float64(v)
	}

	return &Budget{
		Name:            aws.StringValue(b.BudgetName),
		BudgetLimit:     budgetLimit,
		ActualSpend:     actualSpend,
		ForecastedSpend: forecastedSpend,
		LastUpdatedTime: b.LastUpdatedTime,
	}, nil
}
