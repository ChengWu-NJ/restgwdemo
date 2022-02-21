package pgutils

import (
	"fmt"

	"ire.com/pg"
	"ire.com/pg/orm"
	"ire.com/slog"
)

func Upsert(db *pg.DB, ukeyField string, modelInstance interface{}) error {
	var (
		sqlstr string
		err    error
	)
	defer func() {
		slog.Debugf(`Upsert got err:[%+v], sql:[%s]`, err, sqlstr)
	}()

	q := db.Model(modelInstance).OnConflict(fmt.Sprintf(`(%s) DO UPDATE`, ukeyField))

	sqlstr, err = InsertQueryString(q)
	if err != nil {
		return err
	}

	_, err = q.Insert()

	return err
}

func SelectOneByKey(db *pg.DB, modelInstanceWithKey interface{}) error {
	var (
		sqlstr string
		err    error
	)
	defer func() {
		slog.Debugf(`SelectOneByKey got err:[%+v], sql:[%s]`, err, sqlstr)
	}()

	q := db.Model(modelInstanceWithKey).WherePK()
	
	sqlstr, err = SelectQueryString(q)
	if err != nil {
		return err
	}

	err = q.Select()

	return err
}

func DeleteOneByKey(db *pg.DB, modelInstanceWithKey interface{}) error {
	var (
		sqlstr string
		err    error
	)
	defer func() {
		slog.Debugf(`DeleteOneByKey got err:[%+v], sql:[%s]`, err, sqlstr)
	}()

	q := db.Model(modelInstanceWithKey).WherePK()
	
	sqlstr, err = DeleteQueryString(q)
	if err != nil {
		return err
	}

	_, err = q.Delete()

	return err
}

func CreateTableQueryString(q *orm.Query, opt *orm.CreateTableOptions) (string, error) {
	qq := orm.NewCreateTableQuery(q, opt)
	return queryString(qq)
}

func SelectQueryString(q *orm.Query) (string, error) {
	qq := orm.NewSelectQuery(q)
	return queryString(qq)
}

func InsertQueryString(q *orm.Query) (string, error) {
	qq := orm.NewInsertQuery(q)
	return queryString(qq)
}

func DeleteQueryString(q *orm.Query) (string, error) {
	qq := orm.NewDeleteQuery(q)
	return queryString(qq)
}

func queryString(model orm.QueryAppender) (string, error) {
	fmter := orm.NewFormatter().WithModel(model)
	b, err := model.AppendQuery(fmter, nil)

	return string(b), err
}
