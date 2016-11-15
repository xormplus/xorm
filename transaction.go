package xorm

import (
	"sync"

	"github.com/xormplus/core"
)

const (
	PROPAGATION_REQUIRED      = 0 //如果当前没有事务，就新建一个事务，如果已经存在一个事务中，加入到这个事务中。这是最常见的选择。
	PROPAGATION_SUPPORTS      = 1 //支持当前事务，如果当前没有事务，就以非事务方式执行。
	PROPAGATION_MANDATORY     = 2 //使用当前的事务，如果当前没有事务，就抛出异常。
	PROPAGATION_REQUIRES_NEW  = 3 //新建事务，如果当前存在事务，把当前事务挂起。
	PROPAGATION_NOT_SUPPORTED = 4 //以非事务方式执行操作，如果当前存在事务，就把当前事务挂起。
	PROPAGATION_NEVER         = 5 //以非事务方式执行，如果当前存在事务，则抛出异常。
	PROPAGATION_NESTED        = 6 //如果当前存在事务，则在嵌套事务内执行。如果当前没有事务，则执行与 PROPAGATION_REQUIRED 类似的操作。它使用了一个单独的事务，这个事务拥有多个可以回滚的保存点。内部事务的回滚不会对外部事务造成影响。
)

type Transaction struct {
	TxSession             *Session
	transactionDefinition int
	isNested              bool
	savePointID           string
}

func (transaction *Transaction) TransactionDefinition() int {
	return transaction.transactionDefinition
}

func (transaction *Transaction) IsExistingTransaction() bool {
	if transaction.TxSession.Tx == nil {
		return false
	} else {
		return true
	}
}

func (transaction *Transaction) GetSavePointID() string {
	return transaction.savePointID
}

func (transaction *Transaction) Session() *Session {
	return transaction.TxSession
}

func (transaction *Transaction) Do(doFunc func(params ...interface{}), params ...interface{}) {
	if transaction.isNested {
		go doFunc(params...)
	} else {
		doFunc(params...)
	}
}

func (transaction *Transaction) WaitToDo(doFunc func(params ...interface{}), params ...interface{}) {
	if transaction.isNested {
		var w sync.WaitGroup
		w.Add(1)
		go func() {
			doFunc(params...)
			w.Done()
		}()
		w.Wait()
	} else {
		doFunc(params...)
	}
}

func (session *Session) Begin(transactionDefinition ...int) (*Transaction, error) {
	var tx *Transaction
	if len(transactionDefinition) == 0 {
		tx = session.transaction(PROPAGATION_REQUIRED)
	} else {
		tx = session.transaction(transactionDefinition[0])
	}

	err := tx.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (session *Session) transaction(transactionDefinition int) *Transaction {
	if transactionDefinition > 6 || transactionDefinition < 0 {
		return &Transaction{TxSession: session, transactionDefinition: PROPAGATION_REQUIRED}
	}
	return &Transaction{TxSession: session, transactionDefinition: transactionDefinition}
}

func (transaction *Transaction) Begin() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED: //如果当前没有事务，就新建一个事务，如果已经存在一个事务中，加入到这个事务中。这是最常见的选择。
		if !transaction.IsExistingTransaction() {
			if err := transaction.TxSession.begin(); err != nil {
				return err
			}
		} else {
			if transaction.TxSession.currentTransaction != nil {
				transaction.savePointID = transaction.TxSession.currentTransaction.savePointID
			}
			transaction.isNested = true
		}
		transaction.TxSession.currentTransaction = transaction
		return nil
	case PROPAGATION_SUPPORTS: //支持当前事务，如果当前没有事务，就以非事务方式执行。
		if transaction.IsExistingTransaction() {
			transaction.isNested = true
			if transaction.TxSession.currentTransaction != nil {
				transaction.savePointID = transaction.TxSession.currentTransaction.savePointID
			}
			transaction.TxSession.currentTransaction = transaction
		}
		return nil
	case PROPAGATION_MANDATORY: //使用当前的事务，如果当前没有事务，就抛出异常。
		if !transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		} else {
			if transaction.TxSession.currentTransaction != nil {
				transaction.savePointID = transaction.TxSession.currentTransaction.savePointID
			}
			transaction.isNested = true
			transaction.TxSession.currentTransaction = transaction
		}
		return nil
	case PROPAGATION_REQUIRES_NEW: //新建事务，如果当前存在事务，把当前事务挂起。
		transaction.TxSession = transaction.TxSession.Engine.NewSession()
		if err := transaction.TxSession.begin(); err != nil {
			return err
		}
		transaction.isNested = false
		transaction.TxSession.currentTransaction = transaction
		return nil
	case PROPAGATION_NOT_SUPPORTED: //以非事务方式执行操作，如果当前存在事务，就把当前事务挂起。
		transaction.TxSession = transaction.TxSession.Engine.NewSession()
		if transaction.IsExistingTransaction() {
			transaction.isNested = true
		}
		return nil
	case PROPAGATION_NEVER: //以非事务方式执行，如果当前存在事务，则抛出异常。
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED: //如果当前存在事务，则在嵌套事务内执行。如果当前没有事务，则执行与 PROPAGATION_REQUIRED 类似的操作。
		if !transaction.IsExistingTransaction() {
			if err := transaction.TxSession.begin(); err != nil {
				return err
			}
		} else {
			transaction.isNested = true
			dbtype := transaction.TxSession.Engine.Dialect().DBType()
			if dbtype == core.MSSQL {
				transaction.savePointID = "xorm" + NewShortUUID().String()
			} else {
				transaction.savePointID = "xorm" + NewV1().WithoutDashString()
			}

			if err := transaction.SavePoint(transaction.savePointID); err != nil {
				return err
			}
			transaction.TxSession.IsAutoCommit = false
			transaction.TxSession.IsCommitedOrRollbacked = false
			transaction.TxSession.currentTransaction = transaction

		}
		return nil
	default:
		return ErrTransactionDefinition
	}

}

func (transaction *Transaction) Commit() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED: //如果当前没有事务，就新建一个事务，如果已经存在一个事务中，加入到这个事务中。这是最常见的选择。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.TxSession.commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_SUPPORTS: //支持当前事务，如果当前没有事务，就以非事务方式执行。
		if transaction.IsExistingTransaction() {
			if !transaction.isNested {
				err := transaction.TxSession.commit()
				if err != nil {
					return err
				}
			}
		}
		return nil
	case PROPAGATION_MANDATORY: //使用当前的事务，如果当前没有事务，就抛出异常。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.TxSession.commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_REQUIRES_NEW: //新建事务，如果当前存在事务，把当前事务挂起。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.TxSession.commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_NOT_SUPPORTED: //以非事务方式执行操作，如果当前存在事务，就把当前事务挂起
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NEVER: //以非事务方式执行，如果当前存在事务，则抛出异常。
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED: //如果当前存在事务，则在嵌套事务内执行。如果当前没有事务，则执行与 PROPAGATION_REQUIRED 类似的操作。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.TxSession.commit()
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return ErrTransactionDefinition
	}
}

func (transaction *Transaction) Rollback() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED: //如果当前没有事务，就新建一个事务，如果已经存在一个事务中，加入到这个事务中。这是最常见的选择。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		err := transaction.TxSession.rollback()
		if err != nil {
			return err
		}
		return nil
	case PROPAGATION_SUPPORTS: //支持当前事务，如果当前没有事务，就以非事务方式执行。
		if transaction.IsExistingTransaction() {
			err := transaction.TxSession.rollback()
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	case PROPAGATION_MANDATORY: //使用当前的事务，如果当前没有事务，就抛出异常。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if transaction.savePointID != "" {
			if err := transaction.RollbackToSavePoint(transaction.savePointID); err != nil {
				return err
			}
			return nil
		} else {
			err := transaction.TxSession.rollback()
			if err != nil {
				return err
			}
			return nil
		}

	case PROPAGATION_REQUIRES_NEW: //新建事务，如果当前存在事务，把当前事务挂起。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		err := transaction.TxSession.rollback()
		if err != nil {
			return err
		}
		return nil
	case PROPAGATION_NOT_SUPPORTED: //以非事务方式执行操作，如果当前存在事务，就把当前事务挂起
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NEVER: //以非事务方式执行，如果当前存在事务，则抛出异常。
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED: //如果当前存在事务，则在嵌套事务内执行。如果当前没有事务，则执行与 PROPAGATION_REQUIRED 类似的操作。
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if transaction.isNested {
			if err := transaction.RollbackToSavePoint(transaction.savePointID); err != nil {
				return err
			}
			return nil
		} else {
			err := transaction.TxSession.rollback()
			if err != nil {
				return err
			}
			return nil
		}
	default:
		return ErrTransactionDefinition
	}
}

func (transaction *Transaction) SavePoint(savePointID string) error {
	if transaction.TxSession.Tx == nil {
		return ErrNotInTransaction
	}

	var lastSQL string
	dbtype := transaction.TxSession.Engine.Dialect().DBType()
	if dbtype == core.MSSQL {
		lastSQL = "save tran " + savePointID
	} else {
		lastSQL = "SAVEPOINT " + savePointID + ";"
	}

	transaction.TxSession.saveLastSQL(lastSQL)
	if _, err := transaction.TxSession.Tx.Exec(lastSQL); err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) RollbackToSavePoint(savePointID string) error {
	if transaction.TxSession.Tx == nil {
		return ErrNotInTransaction
	}

	var lastSQL string
	dbtype := transaction.TxSession.Engine.Dialect().DBType()
	if dbtype == core.MSSQL {
		lastSQL = "rollback tran " + savePointID
	} else {
		lastSQL = "ROLLBACK TO SAVEPOINT " + transaction.savePointID + ";"
	}

	transaction.TxSession.saveLastSQL(lastSQL)
	if _, err := transaction.TxSession.Tx.Exec(lastSQL); err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) SetISOLATION() error {

	var lastSQL string
	dbtype := transaction.TxSession.Engine.Dialect().DBType()
	if dbtype == core.MSSQL {
		lastSQL = "SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED"
	} else {
		lastSQL = "SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED"
	}

	transaction.TxSession.saveLastSQL(lastSQL)

	if transaction.TxSession.Tx == nil {
		if _, err := transaction.TxSession.Exec(lastSQL); err != nil {
			return err
		}
	} else {
		if _, err := transaction.TxSession.Tx.Exec(lastSQL); err != nil {
			return err
		}
	}

	return nil
}
