package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}