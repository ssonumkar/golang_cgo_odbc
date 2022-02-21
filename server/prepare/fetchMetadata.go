package prepare

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ssonumkar/test-odbc-snowflake/api"
	"github.com/ssonumkar/test-odbc-snowflake/server/coredbutils"
)

func describeColumn(h api.SQLHSTMT, idx int, namebuf []byte) (namelen int, sqltype api.SQLSMALLINT, ret api.SQLRETURN) {
	var l, decimal, nullable api.SQLSMALLINT
	var size api.SQLULEN
	ret = api.SQLDescribeCol(h, api.SQLUSMALLINT(idx+1),
		(*api.SQLCHAR)(unsafe.Pointer(&namebuf[0])),
		api.SQLSMALLINT(len(namebuf)), &l,
		&sqltype, &size, &decimal, &nullable)
	return int(l), sqltype, ret
}

func describeParam(h api.SQLHSTMT, idx int) (sqltype api.SQLSMALLINT, nullable api.SQLSMALLINT, ret api.SQLRETURN) {
	var decimal api.SQLSMALLINT
	var size api.SQLULEN
	ret = api.SQLDescribeParam(h, api.SQLUSMALLINT(idx+1),
		&sqltype, &size, &decimal, &nullable)
	return sqltype, nullable, ret
}
func getInputSet(stmt api.SQLHSTMT) (*Input, error) {
	var in Input
	var nParams api.SQLSMALLINT
	ret := api.SQLNumParams(stmt, &nParams)
	in.numInput = int(nParams)
	if coredbutils.IsError(ret) {
		return nil, fmt.Errorf("errorwhile getInputSet")
	}
	var is InputSet
	for i := 0; i < int(nParams); i++ {
		var param Parameter
		sqltype, nullable, ret := describeParam(stmt, i)
		if coredbutils.IsError(ret) {
			return nil, fmt.Errorf("errorwhile getInputSet")
		}
		param.Name = fmt.Sprintf("Param[%d]", i+1)
		param.Nullable = !(int(nullable) == 0)
		param.SqlType = sqltype
		is.Parameters = param
		in.InputSets = append(in.InputSets, is)
	}
	return &in, nil
}

func getOutputSet(stmt api.SQLHSTMT) (*Output, error) {
	var out Output
	var nCols api.SQLSMALLINT
	ret := api.SQLNumResultCols(stmt, &nCols)
	if coredbutils.IsError(ret) {
		return nil, fmt.Errorf("errorwhile getOutputSet--1")
	}
	out.numOutput = int(nCols)
	var os OutputSet
	for i := 0; i < int(nCols); i++ {
		var col Column
		namebuf := make([]byte, 150)
		namelen, sqltype, ret := describeColumn(stmt, i, namebuf)
		if ret == api.SQL_SUCCESS_WITH_INFO && namelen > len(namebuf) {
			// try again with bigger buffer
			namebuf = make([]byte, namelen)
			namelen, sqltype, ret = describeColumn(stmt, i, namebuf)
		}
		if coredbutils.IsError(ret) {
			return nil, fmt.Errorf("error while getOutputSet--2")
		}
		if namelen > len(namebuf) {
			// still complaining about buffer size
			return nil, errors.New("failed to allocate column name buffer")
		}
		col.Name = string(namebuf[:namelen])
		col.SqlType = sqltype
		os.Columns = col
		out.OutputSets = append(out.OutputSets, os)
	}
	return &out, nil
}

func GetQueryMetadata(stmt api.SQLHSTMT) (*QueryMetadata, error) {
	oset, err := getOutputSet(stmt)
	if err != nil {
		return nil, err
	}
	iset, err := getInputSet(stmt)
	if err != nil {
		return nil, err
	}
	// logger.Log("INFO", "output: ", oset, "input: ", iset)
	return &QueryMetadata{
		OutputDT: oset,
		InputDT:  iset}, nil
}
