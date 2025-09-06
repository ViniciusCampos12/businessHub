package entities

import (
	"fmt"
	"testing"
)

func TestCompanyPwdPercentageCalcCorrectly(t *testing.T) {
	companyDummie := &Company{}

	testCases := []struct {
		Value    int
		Expected float64
	}{
		{Value: 100, Expected: 0.02},
		{Value: 201, Expected: 0.03},
		{Value: 501, Expected: 0.04},
		{Value: 1001, Expected: 0.05},
		{Value: 0, Expected: 0},
	}

	for k := range testCases {
		testName := fmt.Sprintf("value: %v, expected: %v", testCases[k].Value, testCases[k].Expected)
		t.Run(testName, func(t *testing.T) {
			res := companyDummie.pwdPercentage(testCases[k].Value)
			assertEqual := (testCases[k].Expected == res)
			if !assertEqual {
				t.Errorf("expected result %v, got %v", testCases[k].Expected, res)
			}
		})
	}
}

func TestCompanyCheckPWDQuotaCalcCorrectly(t *testing.T) {
	companyDummie := &Company{}

	t.Run("pwd is ok", func(t *testing.T) {
		err := companyDummie.CheckPWDQuota(100, 2)

		if err != nil {
			t.Error("pwd quota expected nil, but got error")
		}
	})

	t.Run("pwd is not ok", func(t *testing.T) {
		err := companyDummie.CheckPWDQuota(100, 1)

		if err == nil {
			t.Error("pwd quota expected error, but got nil")
		}
	})
}
