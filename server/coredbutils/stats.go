/*Package coredbutils
 * @author Abhijeet Padwal
 * @email apadwal@tibco.com
 * @create date 2020-02-27
 * @modify date 2020-02-27
 * @desc [description]
 */
package coredbutils

import (
	"fmt"
	"sync"

	"github.com/ssonumkar/test-odbc-snowflake/api"
)

type Stats struct {
	EnvCount  int
	ConnCount int
	StmtCount int
	mu        sync.Mutex
}

func (s *Stats) UpdateHandleCount(handleType api.SQLSMALLINT, change int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch handleType {
	case api.SQL_HANDLE_ENV:
		s.EnvCount += change
	case api.SQL_HANDLE_DBC:
		s.ConnCount += change
	case api.SQL_HANDLE_STMT:
		s.StmtCount += change
	default:
		return fmt.Errorf("unexpected handle type %d", handleType)
	}
	return nil
}
