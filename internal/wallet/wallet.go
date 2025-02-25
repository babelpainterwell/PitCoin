package wallet

// wallet.go contains features such as verifying transactions, signing transactions

type Wallet struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
}

// NewWallet creates a new wallet with a private and public key
// also generates new private and public keys
func NewWallet() *Wallet {
	return &Wallet{}
}

func (w *Wallet) SignTransaction() {
}

func (w *Wallet) VerifyTransaction() {
}