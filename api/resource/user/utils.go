package user

// Helper function to convert boolean Sex to string representation
func sexToString(sex bool) string {
	if sex {
		return "male"
	}
	return "female"
}

// ToDto converts a User model to its DTO representation
func (u *User) ToDto() *DTO {
	return &DTO{
		ID:       u.ID.String(),
		Username: u.Username,
		FullName: u.FullName,
		Sex:      sexToString(u.Sex), // Convert bool to string
		Height:   u.Height,
		Weight:   u.Weight,
		Age:      u.Age,
	}
}

// ToDto converts a slice of User models to a slice of DTOs
func (us Users) ToDto() []*DTO {
	dtos := make([]*DTO, len(us))
	for i, u := range us {
		dtos[i] = u.ToDto()
	}
	return dtos
}

// --- Note: ToModel logic is now handled within the handler for better validation and pointer handling ---
/*
// Example ToModel (if parsing handled differently):
func (f *Form) ToModel() (*User, error) {
	if f.Sex == nil {
		// Handle missing required field - this should ideally be caught by validation
		return nil, errors.New("sex field is required")
	}
	return &User{
		Username: f.Username,
		FullName: f.FullName,
		Sex:      *f.Sex, // Dereference pointer
		Height:   f.Height,
		Weight:   f.Weight,
		Age:      f.Age,
	}, nil
}
*/
