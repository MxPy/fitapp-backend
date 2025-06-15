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
		Kcal:        int(float64(f.Kcal) * 100 / float64(f.Grams)),
		Proteins:    int(float64(f.Proteins) * 100 / float64(f.Grams)),
		Carbs:       int(float64(f.Carbs) * 100 / float64(f.Grams)),
		Fats:        int(float64(f.Fats) * 100 / float64(f.Grams)),
	}
}
