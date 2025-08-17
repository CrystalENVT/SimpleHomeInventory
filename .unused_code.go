import (
	"encoding/json"
	"fmt"
	"github.com/vgpc/upc"                       // https://github.com/vgpc/upc -- good for UPC parsing, but probably won't use long term
)

function upc_decode() {
	//This section allows for breaking a UPC value into parts. Not necessarily needed, but saving just in case
	u, err := upc.Parse(UPC_String)
	if err == upc.ErrInvalidCheckDigit {
		fmt.Println("There's a typo in the UPC")
	} else if err != nil {
		fmt.Printf("Something's wrong with the UPC: %s", err)
	} else {
		fmt.Println("UPC is: " + u.String())
		fmt.Printf("Number system: %d\n", u.NumberSystem())
		fmt.Printf("Check digit: %d\n", u.CheckDigit())
		if u.IsGlobalProduct() {
			fmt.Println("Manufacturer code: " + u.Manufacturer())
			fmt.Printf("Product code: %d\n", u.Product())
			fmt.Printf("Annotated UPC: %d|%s|%d|%d\n", u.NumberSystem(), u.Manufacturer(), u.Product(), u.CheckDigit())
		} else if u.IsDrug() {
			fmt.Printf("Drug code: %d\n", u.Ndc())
		} else if u.IsLocal() {
			fmt.Println("UPC intended only for local use")
		} else if u.IsCoupon() {
			fmt.Println("Manufacturer code: " + u.Manufacturer())
			fmt.Printf("Family code: %d\n", u.Family())
			fmt.Printf("Coupon value: $0.%02d\n", u.Value())
		} else {
			panic("Preceeding categories are exhaustive")
		}
	}
}