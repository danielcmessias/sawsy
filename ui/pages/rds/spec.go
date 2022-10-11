package rds

import (
	"github.com/danielcmessias/sawsy/ui/components/chart/histogram"
	"github.com/danielcmessias/sawsy/ui/components/gallery"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var rdsPageSpec = page.PageSpec{
	Name: "rds",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Databases",
				Icon: icons.DATABASE,
			},
			Columns: []table.Column{
				{
					Title: "Identifier",
				},
				{
					Title: "Engine",
				},
				{
					Title: "Region & AZ",
				},
				{
					Title: "Size",
				},
				{
					Title: "Status",
				},
				{
					Title: "VPC",
				},
				{
					Title: "Multi-AZ",
				},
			},
		},
	},
}

type metric struct {
	APIName   string
	HumanName string
	Formatter func(float64) float64
}

var metrics = []metric{
	{"CPUUtilization", "CPU Utilization", nil},
	{"DatabaseConnections", "Database Connections", nil},
	{"FreeStorageSpace", "Free Storage Space", utils.BytesToMB},
	{"FreeableMemory", "Freeable Memory", utils.BytesToMB},
	{"WriteIOPS", "Write IOPS", nil},
	{"ReadIOPS", "Read IOPS", nil},
	{"DiskQueueDepth", "Queue Depth", nil},
	{"ReplicaLag", "Replica Lag", nil},
	{"WriteThroughput", "Write Throughput", utils.BytesToMB},
	{"ReadThroughput", "Read Throughput", utils.BytesToMB},
	{"SwapUsage", "Swap Usage", utils.BytesToMB},
	{"WriteLatency", "Write Latency", nil},
	{"ReadLatency", "Read Latency", nil},
	{"NetworkReceiveThroughput", "Network Receive Throughput", utils.BytesToMB},
	{"NetworkTransmitThroughput", "Network Transmit Throughput", utils.BytesToMB},
	{"CPUCreditUsage", "CPU Credit Usage", nil},
	{"CPUCreditBalance", "CPU Credit Balance", nil},
	{"TransactionLogsDiskUsage", "Transaction Logs Disk Usage", utils.BytesToMB},
	{"TransactionLogsGeneration", "Transaction Logs Generation", nil},
	{"OldestReplicationSlotLag", "Oldest Replication Slot Lag", nil},
	{"BurstBalance", "Burst Balance", nil},
	{"ReplicationSlotDiskUsage", "Replication Slot Disk Usage", utils.BytesToMB},
	{"DatabaseConnectionIPV6", "Database Connection IPV6", nil},
	{"FreeLocalStorage", "Free Local Storage", utils.BytesToMB},
	{"ReadIOPSLocalStorage", "Read IOPS Local Storage", nil},
	{"ReadLatencyLocalStorage", "Read Latency Local Storage", nil},
	{"ReadThroughputLocalStorage", "Read Throughput Local Storage", utils.BytesToMB},
	{"WriteIOPSLocalStorage", "Write IOPS Local Storage", nil},
	{"WriteLatencyLocalStorage", "Write Latency Local Storage", nil},
	{"ThroughputLocalStorage", "Write Throughput Local Storage", utils.BytesToMB},
	{"MaximumUsedTransactionIDs", "Maximum Used Transaction IDs", nil},
}

var metricSpecs = func() []pane.PaneSpec {
	var arr []pane.PaneSpec = make([]pane.PaneSpec, len(metrics))
	for i, m := range metrics {
		arr[i] = histogram.HistogramSpec{
			BaseSpec: pane.BaseSpec{
				Name: m.HumanName,
			},
		}
	}
	return arr
}()

var instanceSpecPage = page.PageSpec{
	Name: "rds/instance",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Details",
				Icon: icons.INFO,
			},
			Columns: []table.Column{
				{
					Title: "Key",
				},
				{
					Title: "Value",
				},
			},
		},
		gallery.GallerySpec{
			BaseSpec: pane.BaseSpec{
				Name: "Monitoring",
				Icon: icons.CHART,
			},
			Rows:      3,
			Cols:      3,
			PaneSpecs: metricSpecs,
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Tags",
				Icon: icons.TAG,
			},
			Columns: []table.Column{
				{
					Title: "Key",
				},
				{
					Title: "Value",
				},
			},
		},
	},
}
