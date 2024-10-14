package commands

// SPACES
type (
	GroupCommand struct {
		Name        string
		Description string
	}

	GroupByNameCommand struct {
		Name string
	}

	GroupByCommand struct {
		ID int
	}

	JoinLeaveGroupCommand struct {
		GroupName string
		UserID    int
	}

	UpdateGroupCommand struct {
		ID          int
		Name        string
		Description string
	}
)
