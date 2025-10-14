package macho

// Definitions From: https://github.com/Apple-FOSS-Mirror/Security/blob/5bcad85836c8bbb383f660aaf25b555a805a48e4/OSX/sec/Security/Tool/codesign.c#L53-L89

const (
	CsSlotCodedirectory               SlotType = 0
	CsSlotInfoslot                    SlotType = 1 // Info.plist
	CsSlotRequirements                SlotType = 2 // internal requirements
	CsSlotResourcedir                 SlotType = 3 // resource directory
	CsSlotApplication                 SlotType = 4 // Application specific slot/Top-level directory list
	CsSlotEntitlements                SlotType = 5 // embedded entitlement configuration
	CsSlotRepSpecific                 SlotType = 6 // for use by disk rep
	CsSlotEntitlementsDer             SlotType = 7 // DER representation of entitlements
	CsSlotAlternateCodedirectories    SlotType = 0x1000
	CsSlotAlternateCodedirectoryMax            = 5
	CsSlotAlternateCodedirectoryLimit          = CsSlotAlternateCodedirectories + CsSlotAlternateCodedirectoryMax
	CsSlotCmsSignature                SlotType = 0x10000
	CsSlotIdentificationslot          SlotType = 0x10001
	CsSlotTicketslot                  SlotType = 0x10002
)

type BlobIndex struct {
	Type   SlotType // type of entry
	Offset uint32   // offset of entry (relative to superblob file offset)
}
