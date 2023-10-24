package equipment

type InnerEquipCommand struct {
	*InnerEquip
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

func (e *InnerEquip) NewExtractKeyWordsCmd() *InnerEquipCommand {
	return &InnerEquipCommand{
		e,
	}
}
