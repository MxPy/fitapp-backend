package userday

const DateFormat = "2006-01-02" // Standardowy format daty Go dla YYYY-MM-DD

// ToDto converts a UserDay model to its DTO representation
func (ud *UserDay) ToDto() *DTO {
	return &DTO{
		ID:            ud.ID.String(),
		UserID:        ud.UserID.String(),
		UserDate:      ud.UserDate.Format(DateFormat), // Formatuj datę
		DailyKcal:     ud.DailyKcal,
		DailyProteins: ud.DailyProteins,
		DailyCarbs:    ud.DailyCarbs,
		DailyFats:     ud.DailyFats,
	}
}

// ToDto converts a slice of UserDay models to a slice of DTOs
func (uds UserDays) ToDto() []*DTO {
	dtos := make([]*DTO, len(uds))
	for i, ud := range uds {
		dtos[i] = ud.ToDto()
	}
	return dtos
}

// --- Uwaga: ToModel nie jest już używane do konwersji bezpośrednio z Form,
// --- ponieważ parsowanie UserID i UserDate odbywa się w handlerze dla lepszej obsługi błędów.
// --- Można by stworzyć funkcję pomocniczą, która przyjmuje sparsowane dane.

/*
// Example helper if needed:
func NewUserDayModel(userID uuid.UUID, userDate time.Time, kcal, proteins, carbs, fats int) *UserDay {
    return &UserDay{
        UserID:        userID,
        UserDate:      userDate,
        DailyKcal:     kcal,
        DailyProteins: proteins,
        DailyCarbs:    carbs,
        DailyFats:     fats,
    }
}
*/
