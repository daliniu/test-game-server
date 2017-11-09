package stk

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TransactionFunc 创建事务函数
func TransactionFunc(d *sqlx.DB) func(fn func(*sqlx.Tx) error) error {
	return func(fn func(*sqlx.Tx) error) (err error) {
		tx, err := d.Beginx()
		if err != nil {
			return err
		}

		defer func() {
			if e := recover(); e != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
					err = fmt.Errorf("事务出现异常: %v, 回滚失败: %v", e, rollbackErr)
					return
				}
				err = fmt.Errorf("事务出现异常: %v", e)
			}
		}()

		if err = fn(tx); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
		return tx.Commit()
	}
}
