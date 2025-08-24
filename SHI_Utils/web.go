package shiutils

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2" // web server
)

// RenderForm renders the HTML form.
func RenderForm(c *fiber.Ctx) error {
	return c.Render("form", fiber.Map{})
}

func saveImageToFile(fiberContext *fiber.Ctx, file *multipart.FileHeader) (savedFileName string, err error) {
	fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

	// Save the file to disk, to allow opening later for UPC decode:
	split_filename := strings.Split(file.Filename, ".")
	filename_extension := split_filename[len(split_filename)-1]
	savedFileName = fmt.Sprintf("./barcode.%s", filename_extension)

	if err := fiberContext.SaveFile(file, savedFileName); err != nil {
		return "", err
	}

	return savedFileName, nil
}

// ProcessForm processes the form submission.
func ProcessForm(fiberContext *fiber.Ctx) error {

	// return value, to be displayed in web form
	displayValue := ""

	form, form_err := fiberContext.MultipartForm()

	if form_err != nil {
		fmt.Println("Process Form Error: ", form_err)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", form_err.Error())})
	}

	// Only 1 image uploaded at a time, so grab the first file
	dst_filename, save_image_err := saveImageToFile(fiberContext, form.File["UPC_Image"][0])

	if save_image_err != nil {
		fmt.Println("Process Form Error: ", save_image_err)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", save_image_err.Error())})
	}

	upc_string, upc_err := valueFromUPC(dst_filename)

	if upc_err != nil {
		fmt.Println("Process Form Error: ", upc_err)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", upc_err.Error())})
	}

	fmt.Println("UPC_String:", upc_string)
	var productFound bool
	var minimalProduct MinimalProduct
	var upc_lookup_err error

	productFound, minimalProduct, upc_lookup_err = upcLookupViaDB(upc_string)

	if upc_lookup_err != nil {
		fmt.Println("Process Form Error: ", upc_lookup_err)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", upc_lookup_err.Error())})
	}

	fmt.Println("Product Found:", productFound)

	if productFound {
		println("Pretty Printing from Mongo")
		displayValue = prettyPrint(minimalProduct)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": displayValue})
	}

	minimalProduct, upc_lookup_err = upcLookupViaOFF(upc_string)
	if upc_lookup_err != nil {
		fmt.Println("Process Form Error: ", upc_lookup_err)
		return fiberContext.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", upc_lookup_err.Error())})
	}

	db_write_err := writeUPCtoDB(minimalProduct)
	if db_write_err != nil {
		fmt.Println("UPC write to UPC Cache failed: ", db_write_err)
	}

	displayValue = prettyPrint(minimalProduct)
	return fiberContext.Render("display_results", fiber.Map{"DisplayResults": displayValue})
}
