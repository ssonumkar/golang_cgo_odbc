package prepare

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ssonumkar/test-odbc-snowflake/api"
	"github.com/ssonumkar/test-odbc-snowflake/server/connect"
	"github.com/ssonumkar/test-odbc-snowflake/server/coredbutils"
)

func Prepare(conn_string string, query string) {

	var out api.SQLHANDLE
	//get conncetion handle h_dbc which has the connection
	h_dbc, err := connect.Connect(conn_string)

	if err != nil {
		fmt.Println(err)
		return
	}
	//allocate the STMT handle by passing h_dbc as input handle
	ret := api.SQLAllocHandle(api.SQL_HANDLE_STMT, api.SQLHANDLE(h_dbc), &out)
	if coredbutils.IsError(ret) {

		fmt.Println("error creating handle")
	}
	//Get the SQLHSTMT handle
	h := api.SQLHSTMT(out)
	b := syscall.StringByteSlice(query)
	//call SQLPrepare api using the STMT handle and the query
	ret = api.SQLPrepare(h, (*api.SQLCHAR)(unsafe.Pointer(&b[0])), api.SQL_NTS)
	if coredbutils.IsError(ret) {
		defer coredbutils.ReleaseHandle(h)
		fmt.Println("errorfor prepare: ", coredbutils.NewError("SQLPrepare", h))
	}
	qmd, err := GetQueryMetadata(h)
	if err != nil {
		defer coredbutils.ReleaseHandle(h)
		fmt.Println("error for getQueryMetadata: ", err)
	}
	defer coredbutils.ReleaseHandle(h)
	fmt.Printf("query metadata: %v", qmd.OutputDT)
}
