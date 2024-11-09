package model

type Status string

const (
	Active    Status = "ACTIVE"
	Inactive  Status = "INACTIVE"
	Cancelled Status = "CANCELLED"
	Suspended Status = "SUSPENDED"
)

func ValidateStatus(status string) error {
	validStatuses := map[string]bool{
		string(Active):   true,
		string(Inactive): true,
	}

	if !validStatuses[status] {

		return ErrInvalid
	}

	return nil
}
