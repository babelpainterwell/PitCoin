package transaction

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



