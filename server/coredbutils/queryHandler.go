/*Package coredbutils ...
 * @author Abhijeet Padwal
 * @email apadwal@tibco.com
 * @create date 2020-02-27
 * @modify date 2020-02-27
 * @desc [description]
 */
package coredbutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_."
const terminators = " ;),<>+-*%/"
const starters = " =,(<>+-*%/"
const quotes = "\"'`"

func alphaOnly(s string, keyset string) bool {
	for _, char := range s {
		if !strings.Contains(keyset, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func QueryBasicCheck(s string) error {
	if alphaOnly(string(s[0]), quotes) || alphaOnly(string(s[len(s)-1]), quotes) {
		var dignostic []DBDiagnostic
		QueryError := ActionError{
			APIName: "HandleQuery",
			Diag: append(dignostic, DBDiagnostic{
				State:       "42000",
				NativeError: "1064",
				Messge:      "Query Syntax Error: Query should not be quoted, syntax error at or near " + string(s[0]) + " index:0",
			}),
		}
		return New(QueryError)
	}
	return nil
}

// HandleInsertStatement ...
func HandleInsertStatement(query string) (bool, string) {
	var tokens []string
	for _, item := range strings.Split(query, " ") {
		token := strings.TrimSpace(item)
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	if strings.ToLower(tokens[0]) == "insert" && strings.ToLower(tokens[1]) == "into" {
		return true, strings.Split(tokens[2], "(")[0]
	}
	return false, ""
}

// HandleQuery ...
func HandleQuery(query string) (string, []string, string, error) {
	re := regexp.MustCompile("\\n")
	query = re.ReplaceAllString(query, " ")
	re = regexp.MustCompile("\\t")
	query = re.ReplaceAllString(query, " ")
	query = strings.TrimSpace(query)
	if query[len(query)-1:] != ";" {
		query += ";"
	}
	odbcquery, param := "", ""
	endIndex := 0
	dq, sq, bt := "\"", "'", "`"
	dqm, sqm, btm := false, false, false
	var paramsarray []string
	paramMarker := false
	paramCounter := 0
	isInsert, tableName := HandleInsertStatement(query)
	err := QueryBasicCheck(string(query))
	if err != nil {
		return "", nil, "", err
	}
	for i, val := range query {
		chr := string(val)
		if !dqm || !sqm || !btm {
			if chr == dq && string(query[i-1]) != "\\" && !sqm && !btm {
				dqm = !dqm
				continue
			}
			if chr == sq && string(query[i-1]) != "\\" && !dqm && !btm {
				sqm = !sqm
				continue
			}
			if chr == bt && string(query[i-1]) != "\\" && !dqm && !sqm {
				btm = !btm
				continue
			}
			if !dqm && !sqm && !btm {
				if chr == "?" && alphaOnly(string(query[i-1]), starters) {
					if paramMarker {
						paramMarker = false
						param = ""
						continue
					}
					paramMarker = true
					param = ""
					continue
				}
				if paramMarker {
					if !alphaOnly(chr, alpha) && chr != "\n" {
						paramMarker = false
						if param == "" && string(query[i-1]) == "?" && alphaOnly(chr, terminators) {
							return "", nil, "", fmt.Errorf("Parameters can not be unnamed, hint: ?paramname")
							paramsarray = append(paramsarray, "Param"+strconv.Itoa(paramCounter+1))
							odbcquery += query[endIndex : i-len(param)]
							endIndex = i
							paramCounter++
							continue
						}
						if param == "" && string(query[i-1]) == "?" && !alphaOnly(chr, terminators) {
							continue
						}
						if alphaOnly(chr, terminators) {
							paramsarray = append(paramsarray, param)
							odbcquery += query[endIndex : i-len(param)]
							endIndex = i
							paramCounter++
							param = ""
							continue
						}
					}
					param = param + chr
				}
			}
		}
	}
	odbcquery += query[endIndex:]
	if isInsert {
		return odbcquery, paramsarray, tableName, nil
	}
	return odbcquery, paramsarray, "", nil
}
