package equipment

type InnerEquipCommand struct {
	*Inner
}

func (cmd *InnerEquipCommand) ExecEXT() (interface{}, error) {
	return cmd.ExtractKeyWords(), nil
}

func (cmd *InnerEquipCommand) ExecE() error {
	cmd.ExtractKeyWords()
	return nil
}

func (cmd *InnerEquipCommand) Exec() {
	cmd.ExtractKeyWords()
}

func (e *Inner) NewExtractKeyWordsCmd() *InnerEquipCommand {
	return &InnerEquipCommand{
		e,
	}
}
