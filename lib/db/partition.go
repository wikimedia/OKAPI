package db

// Available partitioning schemas
const (
	ListPartition  PartitionBy = "LIST"
	RangePartition PartitionBy = "RANGE"
	HashPartition  PartitionBy = "HASH"
)

// PartitionBy value for partitioning
type PartitionBy string

// Partition is struct for partitioning configuration
type Partition struct {
	By    PartitionBy
	Field string
}
