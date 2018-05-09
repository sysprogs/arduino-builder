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
	RelatedLibraryName      string
	RelatedLibraryDirectory string
}

type KnownHeader struct {
	Name      string
	Libraries []*KnownLibrary
}

type CodeModelBuilder struct {
	Core              *CodeModelLibrary
	Sketch            *CodeModelLibrary
	Libraries         []*CodeModelLibrary
	KnownHeaders      []*KnownHeader
	Prototypes        []*Prototype
	LinkerCommandLine string
}
