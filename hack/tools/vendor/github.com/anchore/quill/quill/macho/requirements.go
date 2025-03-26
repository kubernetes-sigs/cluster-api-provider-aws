package macho

type RequirementType uint32

const (
	HostRequirementType       RequirementType = 1 /* what hosts may run us */
	GuestRequirementType      RequirementType = 2 /* what guests we may run */
	DesignatedRequirementType RequirementType = 3 /* designated requirement */ // this is the only supported type
	LibraryRequirementType    RequirementType = 4 /* what libraries we may link against */
	PluginRequirementType     RequirementType = 5 /* what plug-ins we may load */
)

type RequirementsHeader struct {
	Count  uint32 // TODO: what is this field?? ("count" is inferred)
	Type   RequirementType
	Offset uint32
}

type Requirements struct {
	RequirementsHeader
	// followed by dynamic content as located by offset fields above
	Payload []byte
}
