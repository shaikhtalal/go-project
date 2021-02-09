package hasura

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/xerrors"
)

type User struct {
	ID               int    `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Slug             string `json:"slug"`
	Password         string `json:"password"`
	BadgeID          string `json:"badge_id"`
	ExpiryDate       string `json:"expiry_data"`
	LocationID       string `json:"location_id"`
	VerificationCode string `json:"verification_code"`
	UserVerified     bool   `json:"user_verified"`
	RememberToken    string `json:"remember_token"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	IsDisabled       bool   `json:"is_disabled"`
}

func (c *GqlClient) GetUser(ctx context.Context) error {
	var variables = make(map[string]interface{})

	variables["email"] = "johnsnow@example.com"

	body, err := c.doGraphql(ctx, `query getUsers($email: String!) {
		users(where: {email: {_eq: $email}}) {
		  id
		  first_name
		  last_name
		}
	  }`, variables, "")
	if err != nil {
		return xerrors.Errorf("something went wrong:%v", err)
	}
	defer body.Close()
	response, err := ioutil.ReadAll(body)
	if err != nil {
		return xerrors.Errorf("reader failed: %v", err)
	}
	var user struct {
		Data struct {
			User []User `json:"users"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response, &user); err != nil {
		return xerrors.Errorf("unmarshalling failed: %v", err)
	}
	c.log.Infof("%+v", user)
	return nil
}

func (c *GqlClient) GetUsers(ctx context.Context) error {
	body, err := c.doGraphql(ctx, `query getUsers {
		users{
		  id
		  first_name
		  last_name
		}
	  }`, nil, "")
	if err != nil {
		return xerrors.Errorf("something went wrong:%v", err)
	}
	defer body.Close()

	response, err := ioutil.ReadAll(body)
	if err != nil {
		return xerrors.Errorf("reader failed: %v", err)
	}
	var user struct {
		Data struct {
			User []User `json:"users"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response, &user); err != nil {
		return xerrors.Errorf("unmarshalling failed: %v", err)
	}
	c.log.Infof("%+v", user)
	return nil
}
