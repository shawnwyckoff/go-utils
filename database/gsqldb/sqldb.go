package gsqldb

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/shawnwyckoff/gpkg/apputil/gerror"
	"github.com/shawnwyckoff/gpkg/container/gstring"
	"github.com/shawnwyckoff/gpkg/database/gdriver"
	"xorm.io/xorm"
)

type (
	SqlDB struct {
		ng *xorm.Engine
		dvr gdriver.Driver
	}
)

func Dial(dvr gdriver.Driver, connectString string) (*SqlDB, error) {
	r := &SqlDB{dvr:dvr}
	err := gerror.ErrNil

	switch dvr {
	case gdriver.MySQL:
		r.ng, err = xorm.NewEngine("mysql", connectString)
	case gdriver.Mssql:
		r.ng, err = xorm.NewEngine("mssql", connectString)
	case gdriver.PgSQL:
		r.ng, err = xorm.NewEngine("postgres", connectString)
	case gdriver.SQLite:
		r.ng, err = xorm.NewEngine("sqlite3", connectString)
	case gdriver.Oracle:
		r.ng, err = xorm.NewEngine("oracle", connectString)
	case gdriver.TiDB:
		r.ng, err = xorm.NewEngine("mysql", connectString)
	case gdriver.CockroachDB:
		r.ng, err = xorm.NewEngine("postgres", connectString)
	default:
		return nil, gerror.Errorf("unsupported database driver %s", dvr.String())
	}

	return r, err
}

func (s *SqlDB) Tables() ([]string, error) {
	tables, err := s.ng.DBMetas()
	if err != nil {
		return nil, err
	}

	var res []string
	for _, v := range tables {
		res = append(res, v.Name)
	}
	return res, nil
}

// 根据结构体中存在的非空数据来查询单条数据
func (s *SqlDB) SelectOne(condAndOutput interface{}) (bool, error) {
	return s.ng.Get(condAndOutput)
}

// 根据cond...结构体中存在的非空数据来查询全部数据
func (s *SqlDB) SelectAll(output interface{}, cond... interface{}) error {
	return s.ng.Find(output, cond...)
}

// 根据cond...结构体中存在的非空数据来Upsert单条数据
func (s *SqlDB) UpsertOne(data, cond interface{}) (int64, error) {
	n, err := s.ng.InsertOne(data)
	if s.dvr == gdriver.MySQL && n == 0 && gstring.StartWith(err.Error(), "Error 1062") { // Error 1062: duplicate primary key
		return s.ng.Update(data, cond)
	}
	return n, err
}

// 根据cond结构体中存在的非空数据来查询是否存在，同时cond也是要目标table名
// table: use to known which table to query
func (s *SqlDB) Exist(cond interface{}) (bool, error) {
	return s.ng.Exist(cond)
}

// 根据cond结构体中存在的非空数据来删除记录，同时cond也是要目标table名
func (s *SqlDB) Remove(cond interface{}) (int64, error) {
	return s.ng.Unscoped().Delete(cond)
}

func (s *SqlDB) Close() error {
	return s.ng.Close()
}