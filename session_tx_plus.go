package xorm

import (
	"sync"

	"github.com/xormplus/core"
)

const (
	PROPAGATION_REQUIRED      = 0 //Support a current transaction; create a new one if none exists.
	PROPAGATION_SUPPORTS      = 1 //Support a current transaction; execute non-transactionally if none exists.
	PROPAGATION_MANDATORY     = 2 //Support a current transaction; return an error if no current transaction exists.
	PROPAGATION_REQUIRES_NEW  = 3 //Create a new transaction, suspending the current transaction if one exists.
	PROPAGATION_NOT_SUPPORTED = 4 //Do not support a current transaction; rather always execute non-transactionally.
	PROPAGATION_NEVER         = 5 //Do not support a current transaction; return an error if a current transaction exists.
	PROPAGATION_NESTED        = 6 //Execute within a nested transaction if a current transaction exists, behave like PROPAGATION_REQUIRED else.
	PROPAGATION_NOT_REQUIRED  = 7
)

type Transaction struct {
	txSession             *Session
	transactionDefinition int
	isNested              bool
	savePointID           string
}

func (transaction *Transaction) TransactionDefinition() int {
	return transaction.transactionDefinition
}

func (transaction *Transaction) IsExistingTransaction() bool {
	if transaction.txSession.tx == nil {
		return false
	} else {
		return true
	}
}

func (transaction *Transaction) GetSavePointID() string {
	return transaction.savePointID
}

func (transaction *Transaction) Session() *Session {
	return transaction.txSession
}

func (transaction *Transaction) Do(doFunc func(params ...interface{}), params ...interface{}) {
	if transaction.isNested {
		go doFunc(params...)
	} else {
		doFunc(params...)
	}
}

func (transaction *Transaction) WaitForDo(doFunc func(params ...interface{}), params ...interface{}) {
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

func (session *Session) BeginTrans(transactionDefinition ...int) (*Transaction, error) {
	var tx *Transaction
	if len(transactionDefinition) == 0 {
		tx = session.transaction(PROPAGATION_REQUIRED)
	} else {
		tx = session.transaction(transactionDefinition[0])
	}

	err := tx.BeginTrans()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (session *Session) transaction(transactionDefinition int) *Transaction {
	if transactionDefinition > 6 || transactionDefinition < 0 {
		return &Transaction{txSession: session, transactionDefinition: PROPAGATION_REQUIRED}
	}
	return &Transaction{txSession: session, transactionDefinition: transactionDefinition}
}

// Begin a transaction
func (transaction *Transaction) BeginTrans() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED:
		if !transaction.IsExistingTransaction() {
			if err := transaction.txSession.Begin(); err != nil {
				return err
			}
		} else {
			if transaction.txSession.currentTransaction != nil {
				transaction.savePointID = transaction.txSession.currentTransaction.savePointID
			}
			transaction.isNested = true
		}
		transaction.txSession.currentTransaction = transaction
		return nil
	case PROPAGATION_SUPPORTS:
		if transaction.IsExistingTransaction() {
			transaction.isNested = true
			if transaction.txSession.currentTransaction != nil {
				transaction.savePointID = transaction.txSession.currentTransaction.savePointID
			}
			transaction.txSession.currentTransaction = transaction
		}
		return nil
	case PROPAGATION_MANDATORY:
		if !transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		} else {
			if transaction.txSession.currentTransaction != nil {
				transaction.savePointID = transaction.txSession.currentTransaction.savePointID
			}
			transaction.isNested = true
			transaction.txSession.currentTransaction = transaction
		}
		return nil
	case PROPAGATION_REQUIRES_NEW:
		transaction.txSession = transaction.txSession.engine.NewSession()
		if err := transaction.txSession.Begin(); err != nil {
			return err
		}
		transaction.isNested = false
		transaction.txSession.currentTransaction = transaction
		return nil
	case PROPAGATION_NOT_SUPPORTED:
		if transaction.IsExistingTransaction() {
			transaction.isNested = true
			transaction.txSession = transaction.txSession.engine.NewSession()
		}
		return nil
	case PROPAGATION_NEVER:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED:
		if !transaction.IsExistingTransaction() {
			if err := transaction.txSession.Begin(); err != nil {
				return err
			}
		} else {
			transaction.isNested = true
			dbtype := transaction.txSession.engine.Dialect().DBType()
			if dbtype == core.MSSQL {
				transaction.savePointID = "xorm" + NewShortUUID().String()
			} else {
				transaction.savePointID = "xorm" + NewV1().WithoutDashString()
			}

			if err := transaction.SavePoint(transaction.savePointID); err != nil {
				return err
			}
			transaction.txSession.isAutoCommit = false
			transaction.txSession.isCommitedOrRollbacked = false
			transaction.txSession.currentTransaction = transaction

		}
		return nil
	case PROPAGATION_NOT_REQUIRED:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}

		if err := transaction.txSession.Begin(); err != nil {
			return err
		}
		return nil
	default:
		return ErrTransactionDefinition
	}

}

// Commit When using transaction, Commit will commit all operations.
func (transaction *Transaction) CommitTrans() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_SUPPORTS:
		if transaction.IsExistingTransaction() {
			if !transaction.isNested {
				err := transaction.txSession.Commit()
				if err != nil {
					return err
				}
			}
		}
		return nil
	case PROPAGATION_MANDATORY:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_REQUIRES_NEW:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		}
		return nil
	case PROPAGATION_NOT_SUPPORTED:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NEVER:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}

		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		} else if transaction.txSession.rollbackSavePointID == transaction.savePointID {
			if err := transaction.RollbackToSavePoint(transaction.savePointID); err != nil {
				transaction.txSession.rollbackSavePointID = ""
				return err
			}
		}
		return nil
	case PROPAGATION_NOT_REQUIRED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if !transaction.isNested {
			err := transaction.txSession.Commit()
			if err != nil {
				return err
			}
		} else {
			return ErrNestedTransaction
		}
		return nil
	default:
		return ErrTransactionDefinition
	}
}

// Rollback When using transaction, you can rollback if any error
func (transaction *Transaction) RollbackTrans() error {
	switch transaction.transactionDefinition {
	case PROPAGATION_REQUIRED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if transaction.savePointID == "" {
			err := transaction.txSession.Rollback()
			if err != nil {
				return err
			}
		} else {
			transaction.txSession.rollbackSavePointID = transaction.savePointID
		}

		return nil
	case PROPAGATION_SUPPORTS:
		if transaction.IsExistingTransaction() {
			if transaction.savePointID == "" {
				err := transaction.txSession.Rollback()
				if err != nil {
					return err
				}
			} else {
				transaction.txSession.rollbackSavePointID = transaction.savePointID
			}
			return nil
		}
		return nil
	case PROPAGATION_MANDATORY:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		if transaction.savePointID == "" {
			err := transaction.txSession.Rollback()
			if err != nil {
				return err
			}
		} else {
			transaction.txSession.rollbackSavePointID = transaction.savePointID
		}
		return nil

	case PROPAGATION_REQUIRES_NEW:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}
		err := transaction.txSession.Rollback()
		if err != nil {
			return err
		}
		return nil
	case PROPAGATION_NOT_SUPPORTED:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NEVER:
		if transaction.IsExistingTransaction() {
			return ErrNestedTransaction
		}
		return nil
	case PROPAGATION_NESTED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}

		if transaction.txSession.rollbackSavePointID == transaction.savePointID {
			return nil
		}

		if transaction.isNested {
			if err := transaction.RollbackToSavePoint(transaction.savePointID); err != nil {
				return err
			}
			return nil
		} else {
			err := transaction.txSession.Rollback()
			if err != nil {
				return err
			}
			return nil
		}
	case PROPAGATION_NOT_REQUIRED:
		if !transaction.IsExistingTransaction() {
			return ErrNotInTransaction
		}

		err := transaction.txSession.Rollback()
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrTransactionDefinition
	}
}

func (transaction *Transaction) SavePoint(savePointID string) error {
	if transaction.txSession.tx == nil {
		return ErrNotInTransaction
	}

	var lastSQL string
	dbtype := transaction.txSession.engine.Dialect().DBType()
	if dbtype == core.MSSQL {
		lastSQL = "save tran " + savePointID
	} else {
		lastSQL = "SAVEPOINT " + savePointID + ";"
	}

	transaction.txSession.saveLastSQL(lastSQL)
	if _, err := transaction.txSession.tx.Exec(lastSQL); err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) RollbackToSavePoint(savePointID string) error {
	if transaction.txSession.tx == nil {
		return ErrNotInTransaction
	}

	var lastSQL string
	dbtype := transaction.txSession.engine.Dialect().DBType()
	if dbtype == core.MSSQL {
		lastSQL = "rollback tran " + savePointID
	} else {
		lastSQL = "ROLLBACK TO SAVEPOINT " + transaction.savePointID + ";"
	}

	transaction.txSession.saveLastSQL(lastSQL)
	if _, err := transaction.txSession.tx.Exec(lastSQL); err != nil {
		return err
	}

	return nil
}
