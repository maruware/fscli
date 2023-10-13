package fscli

type Metacommand interface {
	Type() string
	MetacommandType() string
}

type BaseMetacommand struct {
}

func (m *BaseMetacommand) Type() string {
	return "Metacommand"
}

type MetacommandListCollections struct {
	BaseMetacommand
	baseDoc string
}

func (m *MetacommandListCollections) MetacommandType() string {
	return "ListCollection"
}

type MetacommandPager struct {
	BaseMetacommand
	on bool
}

func (m *MetacommandPager) MetacommandType() string {
	return "Pager"
}
