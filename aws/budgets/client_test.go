package budgets

import (
	"os"
	"testing"
)

// go test ./aws/budgets

var (
	roleArn         = os.Getenv("ROLE_ANR")
	roleSessionName = os.Getenv("ROLE_SESSION_NAME")
	budgetAccountID = os.Getenv("BUDGET_ACCOUNT_ID")
)

func TestNew(t *testing.T) {
	c, err := New(roleArn, roleSessionName)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}
