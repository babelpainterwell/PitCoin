package transaction

import (
	"bytes"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
)

type TxIn struct {
}

type TxOut struct {
}

type Transaction struct {
	Version int32 
	TxIns []TxIn
	TxOuts []TxOut 
	LockTime uint32
}

func (tx *Transaction) SerializeTransaction() []byte {
	// serialize the transaction
	var buf bytes.Buffer

	// serialization of all fields
	//
	//
	//

	return buf.Bytes()
}


func (tx *Transaction) GetTxID() [32]byte {
	transactionBytes := tx.SerializeTransaction()
	return hashutil.DoubleSha256(transactionBytes)
}
