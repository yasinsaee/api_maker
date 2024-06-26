/*
	@package   product
	@version   1.0.0
	@summary   Form definitions for adding a product
	@details   Provides form struct and methods to bind form data to the product model.
	@date      2024-06-03
	@author    YasinSaee
*/

package product

import (
	"errors"
	"fmt"

	apimaker "github.com/yasinsaee/api_maker"
)

// AddProductForm defines the fields required to add a new product.
type AddProductForm struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

// Validate implements the CustomValidator interface for custom validation logic
func (a AddProductForm) Validate() error {
	if a.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

// Bind binds the form data to the Product model.
// you can do anything here
// like you can pass default values.
func (a AddProductForm) Bind(model apimaker.Model) error {
	// Type assertion to ensure the model is a Product
	product, ok := model.(*Product)
	if !ok {
		return fmt.Errorf("invalid model type")
	}

	// Bind the form data to the Product model
	product.Name = "yasin"
	return nil
}
