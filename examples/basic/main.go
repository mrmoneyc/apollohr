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
	location := os.Getenv("LOCATION")

	c := apollohr.NewClient(companyCode, employeeNo, password)

	// Get user info
	userInfo, _ := c.GetUserInfo()
	fmt.Printf("%+v\n\n", userInfo)

	// Punch for start work
	// punchResult, err := c.Punch(1, location)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n\n", punchResult)

	// Punch for off duty
	punchResult, err := c.Punch(2, location)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n\n", punchResult)
}
