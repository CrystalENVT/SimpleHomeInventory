package main

/**
 * Notes of Abbreviations & Terms:
 *
 * OFF | OpenFoodFacts - Main source for UPC data & Food Data
 */

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/jinzhu/copier"                  // https://github.com/jinzhu/copier
	"github.com/openfoodfacts/openfoodfacts-go" // https://pkg.go.dev/github.com/openfoodfacts/openfoodfacts-go
)

// Subset of the Spec from here: https://pkg.go.dev/github.com/openfoodfacts/openfoodfacts-go#Product
type MinimalProduct struct {
	Id                  string      `json:"id"`
	Code                string      `json:"code"`
	Brands              string      `json:"brands"`
	BrandsTags          []string    `json:"brands_tags"`
	GenericName         string      `json:"generic_name"`
	ImageFrontURL       string      `json:"image_front_url" copier:"-"`
	ImageIngredientsURL string      `json:"image_ingredients_url" copier:"-"`
	ImageNutritionURL   string      `json:"image_nutrition_url" copier:"-"`
	ImageURL            string      `json:"image_url" copier:"-"`
	Keywords            []string    `json:"_keywords"`
	ProductName         string      `json:"product_name"`
	ProductNameEn       string      `json:"product_name_en"`
	Quantity            string      `json:"quantity"`
	ScansNumber         int         `json:"scans_n"`
	ServingQuantity     json.Number `json:"serving_quantity"`
	ServingSize         string      `json:"serving_size"`
}

func productToMinimalProduct(product openfoodfacts.Product) (results MinimalProduct) {
	ret := MinimalProduct{}
	// Copy the majority of values over
	copier.Copy(&ret, &product)
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
		minimalProduct := productToMinimalProduct(*product)
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

// pure GoLang Pretty Print - Source: https://stackoverflow.com/a/51270134
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// RenderForm renders the HTML form.
func RenderForm(c *fiber.Ctx) error {
	return c.Render("form", fiber.Map{})
}

// ProcessForm processes the form submission.
func ProcessForm(c *fiber.Ctx) error {
	UPC_String := c.FormValue("UPC_String")
	fmt.Println("UPC_String:", UPC_String)
	minimalProduct, err := upcLookupViaOFF(UPC_String)
	displayValue := ""
	if err == nil {
		displayValue = prettyPrint(minimalProduct)
	} else {
		displayValue = err.Error()
	}
	return c.Render("display_results", fiber.Map{"DisplayResults": displayValue})
}

func main() {
	// Hardcoding for testing, change this out with UPC scan or entry
	//UPC_String := "051000185600" // V8 28 can package
	//UPC_String := "815154025911" // NOS Zero Energy Drink

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Serve static files (HTML templates and stylesheets).
	app.Static("/", "./static")

	// Define routes.
	app.Get("/", RenderForm)
	app.Post("/submit", ProcessForm)

	// Start the Fiber app on port 8080.
	app.Listen(":8080")
}
