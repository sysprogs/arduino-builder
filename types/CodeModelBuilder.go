package types

type CodeModelGCCInvocation struct {
	GCC        string
	InputFile  string
	ObjectFile string
	Arguments  []string
}

type CodeModelLibrary struct {
	Name            string
	SourceDirectory string
	ArchiveFile     string
	Invocations     []*CodeModelGCCInvocation
}

type KnownLibrary struct {
	Folder        string
	SrcFolder     string
	UtilityFolder string
	Layout        LibraryLayout
	Name          string
	RealName      string
	IsLegacy      bool
	Version       string
	Author        string
	Maintainer    string
	Sentence      string
	Paragraph     string
	URL           string
	Category      string
	License       string
}

type KnownHeader struct {
	Name               string
	LibraryDirectories []string
}

type CommonCodeModelProperties struct {
	OutputFile   string
	MaxFLASHSize string
	MaxRAMSize   string
}

type CodeModelBuilder struct {
	Core              *CodeModelLibrary
	Sketch            *CodeModelLibrary
	Libraries         []*CodeModelLibrary
	KnownHeaders      []*KnownHeader
	Prototypes        []*Prototype
	KnownLibraries    []*KnownLibrary
	LinkerCommandLine string
	Common            CommonCodeModelProperties
}
