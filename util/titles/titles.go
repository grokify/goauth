package titles

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Department int

const (
	BusinessDevelopment Department = 1 + iota
	CorporateMarketing
	DeveloperRelations
	Finance
	HumanResources
	IT
	Legal
	Operations
	ProductManagement
	ProductMarketing
	Sales
	SalesEngineering
	SalesOps
	SDK
	SoftwareEngineering
	Support
	UX
)

var departments = [...]string{
	"Business Development",
	"Corporate Marketing",
	"Developer Relations",
	"Finance",
	"Human Resources",
	"IT",
	"Legal",
	"Operations",
	"Product Management",
	"Product Marketing",
	"Sales",
	"Sales Engineering",
	"Sales Ops",
	"SDK",
	"Software Engineering",
	"Support",
	"UX",
}

func (d Department) String() string { return departments[d] }

func ParseDepartment(s string) (Department, error) {
	for i, dept := range departments {
		if strings.EqualFold(strings.TrimSpace(s), dept) {
			return Department(i), nil
		}
	}
	return SoftwareEngineering, fmt.Errorf("cannot parse %s", s)
}

var parseOrder = []string{
	"Human Resources",
	"UX",
	"Product Marketing",
	"Corporate Marketing",
	"Developer Relations",
	"SDK",
	"Product Management",
	"Business Development",
	"IT",
	"Sales Engineering",
	"Sales Ops",
	"Support",
	"Software Engineering",
	"Sales",
	"Finance",
}

// patterns must be lower case
const RegularExpressionsRaw = `{
	"departmentPatterns": {
		"Corporate Marketing":[
			"\\bcmo\\b",
			"\\bcontent\\s+marketing",
			"\\bgrowth\\b",
			"\\bmarketing\\b"
		],
		"Developer Relations":[
			"\\badvocate\\b",
			"\\bdeveloper\\s+community",
			"\\bevangelist\\b",
			"\\bevangelism\\b"
		],
		"Product Management": [
			"\\bcpo\\b",
			"product\\s+manage",
			"\\bproduct\\b"
		],
		"Product Marketing": [
			"\\bproduct\\s+marketing\\b",
			"\\bpmm\\b"
		],
		"Business Development": [
			"\\bbusiness\\s+development",
			"\\bbusiness\\s+partners",
			"\\bcarrier\\s+relations",
			"\\bstrategic\\s+alliances",
			"\\bstrategic\\s+partnerships"
		],
		"Finance":[
			"\\baccounting\\b",
			"\\baccountant\\b",
			"\\baccounts\\s+payable",
			"\\baccounts\\s+receivable",
			"\\banalyst\\s+relations",
			"\\bbilling\\b",
			"\\bcfo\\b",
			"\\bcontroller\\b",
			"\\bfinance\\b",
			"\\bfinancial\\b",
			"\\bfp\u0026amp;a",
			"\\brevenue\\s+operations\\b",
			"\\binvestor\\s+relations\\b",
			"\\btreasurer\\b"
		],
		"Human Resources":[
			"\\bbenefits\\b",
			"\\bcompensation\\b",
			"\\bculture\\b",
			"\\brecruit",
			"\\bhr\\b",
			"\\btalent\\s+development",
			"\\btravel\\s+coordinator"
		],
		"IT": [
			"\\bit\\b"
		],
		"Sales":[
			"sales",
			"\\baccount\\s+executive",
			"\\baccount\\s+manager",
			"\\baccount\\s+management"
		],
		"Sales Engineering": [
			"sales\\s+engineer",
			"solutions\\s+architect"
		],
		"Sales Ops": [
			"\\bcrm\\b"
		],
		"SDK":[
			"\\bsdk\\b"
		],
		"Support": [
			"customer\\s+success",
			"support\\s+engineer",
			"support\\s+manager",
			"\\bsupport\\b"
		],
		"Software Engineering": [
			"\\barchitect\\b",
			"\\bdeveloper\\b",
			"software\\s+architect",
			"softare\\s+architect",
			"software\\s+developer",
			"software\\s+engineer",
			"software\\s+engeering",
			"engineering\\s+manager",
			"\\br\u0026amp;d",
			"\\btech\\s+lead",
			"\\btechnical\\s+lead",
			"\\bengineer"
		],
		"UX":[
			"\\bux\\b",
			"\\bvisual\\s+design",
			"\\bui\\s+developer"
		]
	}
}
`

type Parser struct {
	DepartmentPatterns    map[string][]string `json:"departmentPatterns,omitempty"`
	DepartmentExpressions map[string][]*regexp.Regexp
}

func NewParser() Parser {
	p := Parser{
		DepartmentPatterns:    map[string][]string{},
		DepartmentExpressions: map[string][]*regexp.Regexp{},
	}
	err := json.Unmarshal([]byte(RegularExpressionsRaw), &p)
	if err != nil {
		panic(fmt.Sprintf("cannot parse RegularExpressionsRaw: %v", err))
	}
	for deptString, mapPatterns := range p.DepartmentPatterns {
		for _, mapPattern := range mapPatterns {
			rx := regexp.MustCompile(mapPattern)
			if _, ok := p.DepartmentExpressions[deptString]; !ok {
				p.DepartmentExpressions[deptString] = []*regexp.Regexp{}
			}
			p.DepartmentExpressions[deptString] = append(p.DepartmentExpressions[deptString], rx)
		}
	}
	return p
}

var DepartmentRxs = map[Department]*regexp.Regexp{}

func (p *Parser) ParseTitle(s string) (Department, error) {
	sToLower := strings.ToLower(s)
	for _, deptString := range parseOrder {
		rxs, ok := p.DepartmentExpressions[deptString]
		if !ok {
			panic("Cannot find Parse Order Department")
		}
		for _, rx := range rxs {
			m := rx.FindString(sToLower)
			if len(m) > 0 {
				dept, err := ParseDepartment(deptString)
				if err != nil {
					panic(fmt.Sprintf("cannot parse parseOrder Department String: %v", deptString))
				}
				return dept, nil
			}
		}
	}
	return SoftwareEngineering, fmt.Errorf("cannot match title %s", s)
}
