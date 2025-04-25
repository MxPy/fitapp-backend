package user

func (f *Form) ToModel() *User {
	return &User{
		Username:  f.Username,
		Full_name: f.Full_name,
	}
}

func (b *User) ToDto() *DTO {
	return &DTO{
		ID:        b.ID.String(),
		Username:  b.Username,
		Full_name: b.Full_name,
	}
}

func (bs Users) ToDto() []*DTO {
	dtos := make([]*DTO, len(bs))
	for i, v := range bs {
		dtos[i] = v.ToDto()
	}

	return dtos
}
