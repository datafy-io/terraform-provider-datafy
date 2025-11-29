package datafy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	expected := Account{AccountId: "acc-123", AccountName: "my-account", ParentAccountId: "parent-001"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.CreateAccount(context.Background(), &CreateAccountRequest{AccountName: expected.AccountName})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.Account)
}

func TestGetAccount(t *testing.T) {
	expected := Account{AccountId: "acc-123", AccountName: "my-account", ParentAccountId: "parent-001"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.GetAccount(context.Background(), &GetAccountRequest{AccountId: expected.AccountId})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.Account)
}

func TestUpdateAccount(t *testing.T) {
	expected := Account{AccountId: "acc-123", AccountName: "updated-account", ParentAccountId: "parent-001"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.UpdateAccount(context.Background(), &UpdateAccountRequest{AccountId: expected.AccountId, AccountName: expected.AccountName})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.Account)
}

func TestDeleteAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.DeleteAccount(context.Background(), &DeleteAccountRequest{AccountId: "acc-123"})

	assert.NoError(t, err)
	assert.NotNil(t, out)
}
