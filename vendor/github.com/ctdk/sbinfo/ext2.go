package sbinfo

// Ext2Sb represents an ext2/3/4 superblock.
type Ext2Sb struct {
	SInodesCount          uint32
	SBlocksCount          uint32
	SRBlocksCount         uint32
	SFreeBlocksCount      uint32
	SFreeInodesCount      uint32
	SFirstDataBlock       uint32
	SLogBlockSize         uint32
	SLogClusterSize       uint32
	SBlocksPerGroup       uint32
	SClustersPerGroup     uint32
	SInodesPerGroup       uint32
	SMtime                uint32
	SWtime                uint32
	SMntCount             uint16
	SMaxMntCount          uint16
	SMagic                uint16
	SState                uint16
	SErrors               uint16
	SMinorRevLevel        uint16
	SLastcheck            uint32
	SCheckinterval        uint32
	SCreatorOs            uint32
	SRevLevel             uint32
	SDefResUID            uint16
	SDefResGID            uint16
	SFirstIno             uint32
	SInodeSize            uint16
	SBlockGroupNr         uint16
	SFeatureCompat        uint32
	SFeatureIncompat      uint32
	SFeatureROCompat      uint32
	SUUID                 [16]byte
	SVolumeName           [16]byte
	SLastMounted          [64]byte
	SAlgorithmUsageBitmap uint32
	SPreallocBlocks       uint8
	SPreallocDirBlocks    uint8
	SReservedGdtBlocks    uint16
	SJournalUUID          [16]byte
	SJournalInum          uint32
	SJournalDev           uint32
	SLastOrphan           uint32
	SHashSeed             [4]uint32
	SDefHashVersion       byte
	SJnlBackupType        byte
	SDefaultMountOpts     uint32
	SFirstMetaBg          uint32
	SMkfsTime             uint32
	SJnlBlocks            [17]uint32
	SBlocksCountHi        uint32
	SRBlocksCountHi       uint32
	SFreeBlocksCountHi    uint32
	SMinExtraIsize        uint16
	SWantExtraIsize       uint16
	SFlags                uint32
	SRaidStride           uint16
	SMmpInterval          uint16
	SMmpBlock             uint64
	SRaidStripeWidth      uint32
	SLogGroupsPerFlex     byte
	SChecksumType         byte
	SReservedPad          uint16
	SKbytesWritten        uint64
	SSnapshotInum         uint32
	SSnapshotId           uint32
	SSnapshotRBlocksCount uint64
	SSnapshotList         uint32
	SErrorCount           uint32
	SFirstErrorTime       uint32
	SFirstErrorIno        uint32
	SFirstErrorBlock      uint64
	SFirstErrorFunc       [32]byte
	SFirstErrorLine       uint32
	SLastErrorTime        uint32
	SLastErrorIno         uint32
	SLastErrorLine        uint32
	SLastErrorBlock       uint64
	SLastErrorFunc        [32]byte
	SMountOpts            [64]byte
	SUsrQuotaInum         uint32
	SGrpQuotaInum         uint32
	SOverheadBlocks       uint32
	SBackupBgs            [2]uint32
	SReserved             [106]uint32
	SChecksum             uint32
}

const (
	EXT3_FEATURE_COMPAT_HAS_JOURNAL    uint32 = 0x0004
	EXT2_FEATURE_RO_COMPAT_SPARSE_SUPER uint32 = 0x0001
	EXT2_FEATURE_RO_COMPAT_LARGE_FILE   uint32 = 0x0002
	EXT2_FEATURE_RO_COMPAT_BTREE_DIR    uint32 = 0x0004
	EXT4_FEATURE_RO_COMPAT_HUGE_FILE    uint32 = 0x0008
	EXT4_FEATURE_RO_COMPAT_GDT_CSUM     uint32 = 0x0010
	EXT4_FEATURE_RO_COMPAT_DIR_NLINK    uint32 = 0x0020
	EXT4_FEATURE_RO_COMPAT_EXTRA_ISIZE  uint32 = 0x0040
	EXT2_FEATURE_INCOMPAT_FILETYPE    uint32 = 0x0002
	EXT3_FEATURE_INCOMPAT_RECOVER     uint32 = 0x0004
	EXT3_FEATURE_INCOMPAT_JOURNAL_DEV  uint32 = 0x0008
	EXT2_FEATURE_INCOMPAT_META_BG      uint32 = 0x0010
	EXT4_FEATURE_INCOMPAT_EXTENTS     uint32 = 0x0040
	EXT4_FEATURE_INCOMPAT_64BIT       uint32 = 0x0080
	EXT4_FEATURE_INCOMPAT_MMP         uint32 = 0x0100
	EXT4_FEATURE_INCOMPAT_FLEX_BG      uint32 = 0x0200
)

const EXT2_FEATURE_RO_COMPAT_SUPP uint32 = EXT2_FEATURE_RO_COMPAT_SPARSE_SUPER | EXT2_FEATURE_RO_COMPAT_LARGE_FILE | EXT2_FEATURE_RO_COMPAT_BTREE_DIR
const EXT2_FEATURE_INCOMPAT_SUPP uint32 = EXT2_FEATURE_INCOMPAT_FILETYPE | EXT2_FEATURE_INCOMPAT_META_BG
const EXT2_FEATURE_RO_COMPAT_UNSUPPORTED uint32 = ^EXT2_FEATURE_RO_COMPAT_SUPP
const EXT2_FEATURE_INCOMPAT_UNSUPPORTED uint32 = ^EXT2_FEATURE_INCOMPAT_SUPP
const EXT3_FEATURE_RO_COMPAT_SUPP uint32 = EXT2_FEATURE_RO_COMPAT_SPARSE_SUPER |  EXT2_FEATURE_RO_COMPAT_LARGE_FILE | EXT2_FEATURE_RO_COMPAT_BTREE_DIR
const EXT3_FEATURE_INCOMPAT_SUPP uint32 = EXT2_FEATURE_INCOMPAT_SUPP | EXT3_FEATURE_INCOMPAT_RECOVER
const EXT3_FEATURE_RO_COMPAT_UNSUPPORTED uint32 = ^EXT3_FEATURE_RO_COMPAT_SUPP
const EXT3_FEATURE_INCOMPAT_UNSUPPORTED uint32 = ^EXT3_FEATURE_INCOMPAT_SUPP

func (sb *Ext2Sb) IsExt4() bool {
	ext3 := sb.SFeatureIncompat & EXT3_FEATURE_INCOMPAT_UNSUPPORTED
	ext2 := sb.SFeatureIncompat & EXT2_FEATURE_INCOMPAT_UNSUPPORTED
	if ext3 > 0 && ext2 > 0 {
		return true
	}
	return false
}

func (sb *Ext2Sb) IsExt3() bool {
	j := sb.SFeatureCompat & EXT3_FEATURE_COMPAT_HAS_JOURNAL
	if j == 0 {
		return false
	}
	y := sb.SFeatureROCompat & EXT3_FEATURE_RO_COMPAT_UNSUPPORTED
	z := sb.SFeatureIncompat & EXT3_FEATURE_INCOMPAT_UNSUPPORTED
	
	if z == 0 && y == 0{
		return true
	}
	return false
}

func (sb *Ext2Sb) IsExt2() bool {
	j := sb.SFeatureCompat & EXT3_FEATURE_COMPAT_HAS_JOURNAL
	if j > 0 {
		return false
	}
	y := sb.SFeatureROCompat & EXT2_FEATURE_RO_COMPAT_UNSUPPORTED
	z := sb.SFeatureIncompat & EXT2_FEATURE_INCOMPAT_UNSUPPORTED
	if z == 0 && y == 0 {
		return true
	}
	return false
}
