package product

// ToDto converts a Product model to its DTO representation
func (p *Product) ToDto() *DTO {
	return &DTO{
		ID:          p.ID.String(),
		ProductName: p.ProductName,
		Kcal:        p.Kcal,
		Proteins:    p.Proteins,
		Carbs:       p.Carbs,
		Fats:        p.Fats,
	}
}

// ToDto converts a slice of Product models to a slice of DTOs
func (ps Products) ToDto() []*DTO {
	dtos := make([]*DTO, len(ps))
	for i, p := range ps {
		dtos[i] = p.ToDto()
	}
	return dtos
}

// ToModel converts a Form to a Product model (ID needs to be set separately)
func (f *Form) ToModel() *Product {
	return &Product{
		ProductName: f.ProductName,
		Kcal:        f.Kcal,
		Proteins:    f.Proteins,
		Carbs:       f.Carbs,
		Fats:        f.Fats,
	}
}
