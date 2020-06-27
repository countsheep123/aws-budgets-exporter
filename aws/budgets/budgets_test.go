package budgets

import (
	"testing"
)

func TestGetBudgets(t *testing.T) {
	c, err := New(roleArn, roleSessionName)
	if err != nil {
		t.Error(err)
	}

	budgets, err := c.GetBudgets(budgetAccountID)
	if err != nil {
		t.Error(err)
	}
	t.Log(budgets)
}
