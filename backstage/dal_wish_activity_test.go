package backstage

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestPaginateForInt64(t *testing.T) {
	pageSize, curPage := SetMgoPage(10, 2)
	ids := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	newIds := PaginateForInt64(ids, curPage*pageSize, pageSize)
	logs.Info(newIds)
}
