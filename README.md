# Simple Home Inventory

Currently this is a little POC that is being spun up over a weekend, but I am looking to expand this into a full simple inventory system for keeping track of groceries (& possibly more).

## What does this do currently

Currently this can only be ran from `go run main.go`. Edit / Swap out the `UPC_String` to try out calling the OpenFoodFacts (OFF) API

## To Do (in rough intended order)
~~1. Create Web Page to display results~~  
~~2. Add text prompt box for UPC number~~  
~~3. BONUS: See about UPC scanning from Phone camera~~  
4. Store cached data locally (Json -> MongoDB?) as a "cache" to prevent being spammy of OFF's servers  
5. Add / Remove items from inventory  
6. Add Locations concept (Eg. Pantry, Fridge, Freezer) to keep track of where items would be stored (optional feature)  
7. Add Expiration Date (as optional data)  
8. BONUS: Android app (~~Especially if the web page doesn't allow scanning~~ can take picture & submit image in form, but not "scan")
