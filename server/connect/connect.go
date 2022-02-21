package connect

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ssonumkar/test-odbc-snowflake/api"
	"github.com/ssonumkar/test-odbc-snowflake/server/coredbutils"
)

var d *DbDrv = &DbDrv{&coredbutils.Drv}

func Connect(dsn string) (api.SQLHDBC, error) {
	var out api.SQLHANDLE
	//get the ENV handle
	ret := api.SQLAllocHandle(api.SQL_HANDLE_DBC, api.SQLHANDLE(d.H), &out)
	if coredbutils.IsError(ret) {

		fmt.Println("error creating handle")
		return nil, fmt.Errorf("error creating ALLOC handle")
	}
	//init handle for db connection
	h := api.SQLHDBC(out)
	b := syscall.StringByteSlice(dsn)
	//call SQLDriverConnect api by passing necessary params
	ret = api.SQLDriverConnect(h, 0,
		(*api.SQLCHAR)(unsafe.Pointer(&b[0])), api.SQL_NTS,
		nil, 0, nil, api.SQL_DRIVER_NOPROMPT)
	if coredbutils.IsError(ret) {
		defer coredbutils.ReleaseHandle(h)
		return nil, coredbutils.NewError("SQLDriverConnect", h)
	}
	fmt.Println("connection successfull")
	return h, nil
}
