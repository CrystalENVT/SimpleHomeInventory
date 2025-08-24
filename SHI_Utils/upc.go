package shiutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/jinzhu/copier"                  // copy struct data
	"github.com/makiuchi-d/gozxing"             // barcode processing
	"github.com/makiuchi-d/gozxing/oned"        // barcode processing
	"github.com/openfoodfacts/openfoodfacts-go" // Open Food Facts official go API
)

// Subset of the Spec from here: https://pkg.go.dev/github.com/openfoodfacts/openfoodfacts-go#Product
type MinimalProduct struct {
	Id                  string      `json:"id" bson:"_id"`
	Code                string      `json:"code" copier:"-"` // We'll always set this UPC_String that we are searching
	Brands              string      `json:"brands"`
	BrandsTags          []string    `json:"brands_tags"`
	GenericName         string      `json:"generic_name"`
	ImageFrontURL       string      `json:"image_front_url" bson:"image_front_url,omitempty" copier:"-"`
	ImageIngredientsURL string      `json:"image_ingredients_url" bson:"image_ingredients_url,omitempty" copier:"-"`
	ImageNutritionURL   string      `json:"image_nutrition_url" bson:"image_nutrition_url,omitempty" copier:"-"`
	ImageURL            string      `json:"image_url" bson:"image_url,omitempty" copier:"-"`
	Keywords            []string    `json:"_keywords"`
	ProductName         string      `json:"product_name"`
	ProductNameEn       string      `json:"product_name_en"`
	Quantity            string      `json:"quantity"`
	ScansNumber         int         `json:"scans_n"`
	ServingQuantity     json.Number `json:"serving_quantity"`
	ServingSize         string      `json:"serving_size"`
}

func productToMinimalProduct(product openfoodfacts.Product, upc_string string) (results MinimalProduct) {
	ret := MinimalProduct{}
	// Copy the majority of values over
	copier.Copy(&ret, &product)
	// Set our minimal product as the same as our UPC
	ret.Code = upc_string
	// save all URLs as a String instead of as a broken apart json blob
	ret.ImageFrontURL = product.ImageFrontURL.String()
	ret.ImageIngredientsURL = product.ImageIngredientsURL.String()
	ret.ImageNutritionURL = product.ImageIngredientsURL.String()
	ret.ImageURL = product.ImageURL.String()
	return ret
}

func upcLookupViaOFF(upc_string string) (results MinimalProduct, err error) {
	api := openfoodfacts.NewClient("world", "", "")
	product, err := api.Product(upc_string)
	minimalProduct := MinimalProduct{}
	if err == nil {
		// Copy a subset of the full Product data from OFF, for more efficient local DB storage / cache
		minimalProduct := productToMinimalProduct(*product, upc_string)
		fmt.Printf("Minimized Print:\n%+v\n\n", minimalProduct)
		fmt.Println("Pretty Print:\n" + prettyPrint(minimalProduct))
		return minimalProduct, nil
	} else {
		err_msg := fmt.Sprintln("Error from OFF API, Product Likely does not exist, or there was a manual entry error\n\tError message from OFF:", err)
		fmt.Println("Error from OFF API, Product Likely does not exist, or there was a manual entry error")
		fmt.Println("\tError message from OFF:", err)
		return minimalProduct, errors.New(err_msg)
	}
}

func upcFromValue(upc_string string) {
	// Generate a barcode image (*BitMatrix)
	enc := oned.NewUPCAWriter()
	img, _ := enc.Encode(upc_string, gozxing.BarcodeFormat_UPC_A, 250, 50, nil)

	file, _ := os.Create("./barcode.png")
	defer file.Close()

	// *BitMatrix implements the image.Image interface,
	// so it is able to be passed to png.Encode directly.
	_ = png.Encode(file, img)
}

func valueFromUPC(inputFileName string) (results string, err error) {
	fmt.Println("Form FileName:", inputFileName)

	inputFile, inputFileOpenErr := os.Open(inputFileName)
	if inputFileOpenErr != nil {
		fmt.Println("FileErr:", inputFileOpenErr)
	}

	img, _, imgOpenErr := image.Decode(inputFile)
	if imgOpenErr != nil {
		fmt.Println("ImgErr: ", imgOpenErr)
	}

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	upcReader := oned.NewUPCAReader()
	decodingHints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ALSO_INVERTED: true,
		gozxing.DecodeHintType_TRY_HARDER:    true,
	}
	result, bcDecodingErr := upcReader.Decode(bmp, decodingHints)

	inputFileCloseErr := inputFile.Close()
	if inputFileCloseErr != nil {
		fmt.Println(inputFileCloseErr)
	}

	inputFileRemoveErr := os.Remove(inputFileName)
	if inputFileRemoveErr != nil {
		fmt.Println(inputFileRemoveErr)
	}

	if bcDecodingErr != nil {
		fmt.Println(bcDecodingErr)
		return "", errors.New("No Barcode found in image, Please Retry")
	}

	fmt.Println(result)
	return result.String(), nil
}
