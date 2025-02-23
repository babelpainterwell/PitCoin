package main

import (
	"bytes"
	"fmt"

	"github.com/babelpainterwell/shitcoin/internal/hashutil"
)



func main() {

	var buf bytes.Buffer

	fmt.Println(buf.Bytes())
	hashutil.EncodeInt32LE(&buf, 1)

	fmt.Println(buf.Bytes())
	hashutil.EncodeInt32LE(&buf, 254)

	fmt.Println(buf.Bytes())


}