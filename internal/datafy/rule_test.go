package datafy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountAutoscalingRule(t *testing.T) {
	expectedRule := AutoscalingRule{
		AccountId: "acc-123",
		RuleId:    "rule-abc",
		Active:    true,
		Rule:      json.RawMessage(`{"max":10,"min":1}`),
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/autoscaling/rules" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(expectedRule)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.CreateAccountAutoscalingRule(context.Background(), &CreateAccountAutoscalingRuleRequest{
		AccountId: expectedRule.AccountId,
		Active:    expectedRule.Active,
		Rule:      expectedRule.Rule,
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedRule, out.AutoscalingRule)
}

func TestGetAccountAutoscalingRule(t *testing.T) {
	expectedRule := AutoscalingRule{
		AccountId: "acc-123",
		RuleId:    "rule-abc",
		Active:    true,
		Rule:      json.RawMessage(`{"max":10,"min":1}`),
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/autoscaling/rules/rule-abc" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedRule)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.GetAccountAutoscalingRule(context.Background(), &GetAccountAutoscalingRuleRequest{
		AccountId: expectedRule.AccountId,
		RuleId:    expectedRule.RuleId,
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedRule, out.AutoscalingRule)
}

func TestUpdateAccountAutoscalingRule(t *testing.T) {
	expectedRule := AutoscalingRule{
		AccountId: "acc-123",
		RuleId:    "rule-abc",
		Active:    true,
		Rule:      json.RawMessage(`{"max":20,"min":2}`),
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/autoscaling/rules/rule-abc" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedRule)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.UpdateAccountAutoscalingRule(context.Background(), &UpdateAccountAutoscalingRuleRequest{
		AccountId: expectedRule.AccountId,
		RuleId:    expectedRule.RuleId,
		Active:    expectedRule.Active,
		Rule:      expectedRule.Rule,
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedRule, out.AutoscalingRule)
}

func TestDeleteAccountAutoscalingRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/autoscaling/rules/rule-abc" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.DeleteAccountAutoscalingRule(context.Background(), &DeleteAccountAutoscalingRuleRequest{
		AccountId: "acc-123",
		RuleId:    "rule-abc",
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
}
