package core_domain

type PatchTaskCommand struct {
	Title       *string
	Description *string
}

func NewPatchTaskCommand(title *string, description *string) PatchTaskCommand {
	return PatchTaskCommand{
		Title:       title,
		Description: description,
	}
}
