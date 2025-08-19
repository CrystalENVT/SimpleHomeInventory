package shiutils

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2" // web server
)

// RenderForm renders the HTML form.
func RenderForm(c *fiber.Ctx) error {
	return c.Render("form", fiber.Map{})
}

// ProcessForm processes the form submission.
func ProcessForm(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["UPC_Image"]

		// Only 1 image uploaded at a time, so grab the first file
		file := files[0]

		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		// Save the file to disk, to allow opening later for UPC decode:
		split_filename := strings.Split(file.Filename, ".")
		filename_extension := split_filename[len(split_filename)-1]
		dst_filename := fmt.Sprintf("./barcode.%s", filename_extension)
		if err := c.SaveFile(file, dst_filename); err != nil {
			return err
		}

		// return value, to be displayed in web form
		displayValue := ""

		UPC_String, UPC_Err := valueFromUPC(dst_filename)
		if UPC_Err != nil {
			displayValue = UPC_Err.Error()
		} else {
			fmt.Println("UPC_String:", UPC_String)
			minimalProduct, err := upcLookupViaOFF(UPC_String)
			if err == nil {
				displayValue = prettyPrint(minimalProduct)
			} else {
				displayValue = err.Error()
			}
		}
		return c.Render("display_results", fiber.Map{"DisplayResults": displayValue})
	} else {
		fmt.Println("Process Form Error: ", err)
		return c.Render("display_results", fiber.Map{"DisplayResults": fmt.Sprintln("Process Form Error: ", err.Error())})
	}
}
