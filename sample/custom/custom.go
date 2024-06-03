package custom

import apimaker "github.com/yasinsaee/api_maker"

func MyCustomAfterSaveFunction(model apimaker.Model) error {
	// Custom logic here
	return nil
}

func MyCustomBeforeSaveFunction(model apimaker.Model) error {
	// Custom logic here
	return nil
}
