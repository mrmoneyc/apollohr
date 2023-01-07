package main

import (
	"fmt"
	"os"

	"github.com/mrmoneyc/apollohr"
)

func main() {
	companyCode := os.Getenv("COMPANY_CODE")
	employeeNo := os.Getenv("EMPLOYEE_NO")
	password := os.Getenv("PASSWORD")

	c := apollohr.NewClient(companyCode, employeeNo, password)

	// Print token
	fmt.Printf("Current token: %s\n\n", c.GetToken())

	// Get user info
	userInfo, _ := c.GetUserInfo()
	fmt.Printf("%+v\n\n", userInfo)

	// Refresh token
	err := c.RefreshToken()
	if err != nil {
		panic(err)
	}

	// Print new token
	fmt.Printf("New token: %s\n\n", c.GetToken())

	// Get user info, again
	userInfo, _ = c.GetUserInfo()
	fmt.Printf("%+v\n\n", userInfo)
}
