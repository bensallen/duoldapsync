package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
)

func buildAdminClient(url string, proxy func(*http.Request) (*url.URL, error)) *admin.Client {
	ikey := "eyekey"
	skey := "esskey"
	host := strings.Split(url, "//")[1]
	userAgent := "GoTestClient"
	base := duoapi.NewDuoApi(ikey, skey, host, userAgent, duoapi.SetTimeout(1*time.Second), duoapi.SetInsecure(), duoapi.SetProxy(proxy))
	return admin.New(*base)
}

func TestCreateUser(t *testing.T) {
	const createUserResponse = `{
		"stat": "OK",
		"response": {
			"alias1": null,
			"alias2": null,
			"alias3": null,
			"alias4": null,
			"created": 1489612729,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"groups": [],
			"is_enrolled": false,
			"last_directory_sync": null,
			"last_login": null,
			"lastname": "Smith",
			"notes": "",
			"phones": [],
			"realname": "Joe Smith",
			"status": "active",
			"tokens": [],
			"u2ftokens": [],
			"user_id": "DU3RP9I2WOC59VZX672N",
			"username": "jsmith"
		}
	}`

	createUserResult := PostUsersResult{
		duoapi.StatResult{Stat: "OK"},
		admin.User{
			Created:   1489612729,
			Groups:    []admin.Group{},
			Email:     "jsmith@example.com",
			FirstName: "Joe",
			LastName:  "Smith",
			Phones:    []admin.Phone{},
			RealName:  "Joe Smith",
			Status:    "active",
			Tokens:    []admin.Token{},
			UserID:    "DU3RP9I2WOC59VZX672N",
			Username:  "jsmith",
		},
	}

	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, createUserResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	user := User{
		Username:  "jsmith",
		FullName:  "Joe Smith",
		Email:     "jsmith@example.com",
		FirstName: "Joe",
		LastName:  "Smith",
	}

	urlParams, err := user.urlValues()
	if err != nil {
		t.Errorf("CreateUser() got error when trying to take urlValues from example user struct: %s", err)
	}

	type args struct {
		client *admin.Client
		params url.Values
		dryRun bool
	}
	tests := []struct {
		name    string
		args    args
		want    *PostUsersResult
		wantErr bool
	}{
		{
			name:    "Create user",
			args:    args{client: duo, params: urlParams, dryRun: false},
			want:    &createUserResult,
			wantErr: false,
		},
		{
			name:    "Create user dryRun",
			args:    args{client: duo, params: urlParams, dryRun: true},
			want:    &PostUsersResult{duoapi.StatResult{Stat: "OK"}, admin.User{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUser(tt.args.client, tt.args.params, tt.args.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	deleteUserResponse := `{
		"stat": "OK",
		"response": ""
	}`

	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, deleteUserResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	type args struct {
		client *admin.Client
		userID string
		dryRun bool
	}
	tests := []struct {
		name    string
		args    args
		want    *duoapi.StatResult
		wantErr bool
	}{
		{
			name:    "Delete user",
			args:    args{client: duo, userID: "jsmith", dryRun: false},
			want:    &duoapi.StatResult{Stat: "OK"},
			wantErr: false,
		},
		{
			name:    "Delete user dryRun",
			args:    args{client: duo, userID: "jsmith", dryRun: true},
			want:    &duoapi.StatResult{Stat: "OK"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteUser(tt.args.client, tt.args.userID, tt.args.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnrollUser(t *testing.T) {
	enrollUserResponse := `{
		"stat": "OK",
		"response": "00d70e730b22cb66"
	}`

	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, enrollUserResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	user := User{
		Username:  "jsmith",
		FullName:  "Joe Smith",
		Email:     "jsmith@example.com",
		FirstName: "Joe",
		LastName:  "Smith",
	}

	urlParams, err := user.urlValues()
	if err != nil {
		t.Errorf("CreateUser() got error when trying to take urlValues from example user struct: %s", err)
	}

	type args struct {
		client *admin.Client
		params url.Values
		dryRun bool
	}
	tests := []struct {
		name    string
		args    args
		want    *duoapi.StatResult
		wantErr bool
	}{
		{
			name:    "Enroll user",
			args:    args{client: duo, params: urlParams, dryRun: false},
			want:    &duoapi.StatResult{Stat: "OK"},
			wantErr: false,
		},
		{
			name:    "Delete user dryRun",
			args:    args{client: duo, params: urlParams, dryRun: true},
			want:    &duoapi.StatResult{Stat: "OK"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnrollUser(tt.args.client, tt.args.params, tt.args.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnrollUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
