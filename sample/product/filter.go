/*
	@package   product
	@version   1.0.0
	@summary   Product filter definitions
	@details   Provides filter struct and methods to filter product queries.
	@date      2024-06-03
	@author    YasinSaee
*/

package product

// ProductFilter defines the fields by which products can be filtered.
type ProductFilter struct {
	Name string
}

// GetFilters returns a map of filters to be used in queries.
func (filter ProductFilter) GetFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if filter.Name != "" {
		filters["name"] = filter.Name
	}

	return filters
}
