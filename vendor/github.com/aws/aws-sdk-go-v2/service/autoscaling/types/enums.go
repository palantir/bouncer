// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type AcceleratorManufacturer string

// Enum values for AcceleratorManufacturer
const (
	AcceleratorManufacturerNvidia            AcceleratorManufacturer = "nvidia"
	AcceleratorManufacturerAmd               AcceleratorManufacturer = "amd"
	AcceleratorManufacturerAmazonWebServices AcceleratorManufacturer = "amazon-web-services"
	AcceleratorManufacturerXilinx            AcceleratorManufacturer = "xilinx"
)

// Values returns all known values for AcceleratorManufacturer. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AcceleratorManufacturer) Values() []AcceleratorManufacturer {
	return []AcceleratorManufacturer{
		"nvidia",
		"amd",
		"amazon-web-services",
		"xilinx",
	}
}

type AcceleratorName string

// Enum values for AcceleratorName
const (
	AcceleratorNameA100          AcceleratorName = "a100"
	AcceleratorNameV100          AcceleratorName = "v100"
	AcceleratorNameK80           AcceleratorName = "k80"
	AcceleratorNameT4            AcceleratorName = "t4"
	AcceleratorNameM60           AcceleratorName = "m60"
	AcceleratorNameRadeonProV520 AcceleratorName = "radeon-pro-v520"
	AcceleratorNameVu9p          AcceleratorName = "vu9p"
)

// Values returns all known values for AcceleratorName. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AcceleratorName) Values() []AcceleratorName {
	return []AcceleratorName{
		"a100",
		"v100",
		"k80",
		"t4",
		"m60",
		"radeon-pro-v520",
		"vu9p",
	}
}

type AcceleratorType string

// Enum values for AcceleratorType
const (
	AcceleratorTypeGpu       AcceleratorType = "gpu"
	AcceleratorTypeFpga      AcceleratorType = "fpga"
	AcceleratorTypeInference AcceleratorType = "inference"
)

// Values returns all known values for AcceleratorType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AcceleratorType) Values() []AcceleratorType {
	return []AcceleratorType{
		"gpu",
		"fpga",
		"inference",
	}
}

type BareMetal string

// Enum values for BareMetal
const (
	BareMetalIncluded BareMetal = "included"
	BareMetalExcluded BareMetal = "excluded"
	BareMetalRequired BareMetal = "required"
)

// Values returns all known values for BareMetal. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (BareMetal) Values() []BareMetal {
	return []BareMetal{
		"included",
		"excluded",
		"required",
	}
}

type BurstablePerformance string

// Enum values for BurstablePerformance
const (
	BurstablePerformanceIncluded BurstablePerformance = "included"
	BurstablePerformanceExcluded BurstablePerformance = "excluded"
	BurstablePerformanceRequired BurstablePerformance = "required"
)

// Values returns all known values for BurstablePerformance. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (BurstablePerformance) Values() []BurstablePerformance {
	return []BurstablePerformance{
		"included",
		"excluded",
		"required",
	}
}

type CapacityDistributionStrategy string

// Enum values for CapacityDistributionStrategy
const (
	CapacityDistributionStrategyBalancedOnly       CapacityDistributionStrategy = "balanced-only"
	CapacityDistributionStrategyBalancedBestEffort CapacityDistributionStrategy = "balanced-best-effort"
)

// Values returns all known values for CapacityDistributionStrategy. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CapacityDistributionStrategy) Values() []CapacityDistributionStrategy {
	return []CapacityDistributionStrategy{
		"balanced-only",
		"balanced-best-effort",
	}
}

type CapacityReservationPreference string

// Enum values for CapacityReservationPreference
const (
	CapacityReservationPreferenceCapacityReservationsOnly  CapacityReservationPreference = "capacity-reservations-only"
	CapacityReservationPreferenceCapacityReservationsFirst CapacityReservationPreference = "capacity-reservations-first"
	CapacityReservationPreferenceNone                      CapacityReservationPreference = "none"
	CapacityReservationPreferenceDefault                   CapacityReservationPreference = "default"
)

// Values returns all known values for CapacityReservationPreference. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CapacityReservationPreference) Values() []CapacityReservationPreference {
	return []CapacityReservationPreference{
		"capacity-reservations-only",
		"capacity-reservations-first",
		"none",
		"default",
	}
}

type CpuManufacturer string

// Enum values for CpuManufacturer
const (
	CpuManufacturerIntel             CpuManufacturer = "intel"
	CpuManufacturerAmd               CpuManufacturer = "amd"
	CpuManufacturerAmazonWebServices CpuManufacturer = "amazon-web-services"
	CpuManufacturerApple             CpuManufacturer = "apple"
)

// Values returns all known values for CpuManufacturer. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CpuManufacturer) Values() []CpuManufacturer {
	return []CpuManufacturer{
		"intel",
		"amd",
		"amazon-web-services",
		"apple",
	}
}

type ImpairedZoneHealthCheckBehavior string

// Enum values for ImpairedZoneHealthCheckBehavior
const (
	ImpairedZoneHealthCheckBehaviorReplaceUnhealthy ImpairedZoneHealthCheckBehavior = "ReplaceUnhealthy"
	ImpairedZoneHealthCheckBehaviorIgnoreUnhealthy  ImpairedZoneHealthCheckBehavior = "IgnoreUnhealthy"
)

// Values returns all known values for ImpairedZoneHealthCheckBehavior. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ImpairedZoneHealthCheckBehavior) Values() []ImpairedZoneHealthCheckBehavior {
	return []ImpairedZoneHealthCheckBehavior{
		"ReplaceUnhealthy",
		"IgnoreUnhealthy",
	}
}

type InstanceGeneration string

// Enum values for InstanceGeneration
const (
	InstanceGenerationCurrent  InstanceGeneration = "current"
	InstanceGenerationPrevious InstanceGeneration = "previous"
)

// Values returns all known values for InstanceGeneration. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (InstanceGeneration) Values() []InstanceGeneration {
	return []InstanceGeneration{
		"current",
		"previous",
	}
}

type InstanceMetadataEndpointState string

// Enum values for InstanceMetadataEndpointState
const (
	InstanceMetadataEndpointStateDisabled InstanceMetadataEndpointState = "disabled"
	InstanceMetadataEndpointStateEnabled  InstanceMetadataEndpointState = "enabled"
)

// Values returns all known values for InstanceMetadataEndpointState. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (InstanceMetadataEndpointState) Values() []InstanceMetadataEndpointState {
	return []InstanceMetadataEndpointState{
		"disabled",
		"enabled",
	}
}

type InstanceMetadataHttpTokensState string

// Enum values for InstanceMetadataHttpTokensState
const (
	InstanceMetadataHttpTokensStateOptional InstanceMetadataHttpTokensState = "optional"
	InstanceMetadataHttpTokensStateRequired InstanceMetadataHttpTokensState = "required"
)

// Values returns all known values for InstanceMetadataHttpTokensState. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (InstanceMetadataHttpTokensState) Values() []InstanceMetadataHttpTokensState {
	return []InstanceMetadataHttpTokensState{
		"optional",
		"required",
	}
}

type InstanceRefreshStatus string

// Enum values for InstanceRefreshStatus
const (
	InstanceRefreshStatusPending            InstanceRefreshStatus = "Pending"
	InstanceRefreshStatusInProgress         InstanceRefreshStatus = "InProgress"
	InstanceRefreshStatusSuccessful         InstanceRefreshStatus = "Successful"
	InstanceRefreshStatusFailed             InstanceRefreshStatus = "Failed"
	InstanceRefreshStatusCancelling         InstanceRefreshStatus = "Cancelling"
	InstanceRefreshStatusCancelled          InstanceRefreshStatus = "Cancelled"
	InstanceRefreshStatusRollbackInProgress InstanceRefreshStatus = "RollbackInProgress"
	InstanceRefreshStatusRollbackFailed     InstanceRefreshStatus = "RollbackFailed"
	InstanceRefreshStatusRollbackSuccessful InstanceRefreshStatus = "RollbackSuccessful"
	InstanceRefreshStatusBaking             InstanceRefreshStatus = "Baking"
)

// Values returns all known values for InstanceRefreshStatus. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (InstanceRefreshStatus) Values() []InstanceRefreshStatus {
	return []InstanceRefreshStatus{
		"Pending",
		"InProgress",
		"Successful",
		"Failed",
		"Cancelling",
		"Cancelled",
		"RollbackInProgress",
		"RollbackFailed",
		"RollbackSuccessful",
		"Baking",
	}
}

type LifecycleState string

// Enum values for LifecycleState
const (
	LifecycleStatePending                  LifecycleState = "Pending"
	LifecycleStatePendingWait              LifecycleState = "Pending:Wait"
	LifecycleStatePendingProceed           LifecycleState = "Pending:Proceed"
	LifecycleStateQuarantined              LifecycleState = "Quarantined"
	LifecycleStateInService                LifecycleState = "InService"
	LifecycleStateTerminating              LifecycleState = "Terminating"
	LifecycleStateTerminatingWait          LifecycleState = "Terminating:Wait"
	LifecycleStateTerminatingProceed       LifecycleState = "Terminating:Proceed"
	LifecycleStateTerminated               LifecycleState = "Terminated"
	LifecycleStateDetaching                LifecycleState = "Detaching"
	LifecycleStateDetached                 LifecycleState = "Detached"
	LifecycleStateEnteringStandby          LifecycleState = "EnteringStandby"
	LifecycleStateStandby                  LifecycleState = "Standby"
	LifecycleStateWarmedPending            LifecycleState = "Warmed:Pending"
	LifecycleStateWarmedPendingWait        LifecycleState = "Warmed:Pending:Wait"
	LifecycleStateWarmedPendingProceed     LifecycleState = "Warmed:Pending:Proceed"
	LifecycleStateWarmedTerminating        LifecycleState = "Warmed:Terminating"
	LifecycleStateWarmedTerminatingWait    LifecycleState = "Warmed:Terminating:Wait"
	LifecycleStateWarmedTerminatingProceed LifecycleState = "Warmed:Terminating:Proceed"
	LifecycleStateWarmedTerminated         LifecycleState = "Warmed:Terminated"
	LifecycleStateWarmedStopped            LifecycleState = "Warmed:Stopped"
	LifecycleStateWarmedRunning            LifecycleState = "Warmed:Running"
	LifecycleStateWarmedHibernated         LifecycleState = "Warmed:Hibernated"
)

// Values returns all known values for LifecycleState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LifecycleState) Values() []LifecycleState {
	return []LifecycleState{
		"Pending",
		"Pending:Wait",
		"Pending:Proceed",
		"Quarantined",
		"InService",
		"Terminating",
		"Terminating:Wait",
		"Terminating:Proceed",
		"Terminated",
		"Detaching",
		"Detached",
		"EnteringStandby",
		"Standby",
		"Warmed:Pending",
		"Warmed:Pending:Wait",
		"Warmed:Pending:Proceed",
		"Warmed:Terminating",
		"Warmed:Terminating:Wait",
		"Warmed:Terminating:Proceed",
		"Warmed:Terminated",
		"Warmed:Stopped",
		"Warmed:Running",
		"Warmed:Hibernated",
	}
}

type LocalStorage string

// Enum values for LocalStorage
const (
	LocalStorageIncluded LocalStorage = "included"
	LocalStorageExcluded LocalStorage = "excluded"
	LocalStorageRequired LocalStorage = "required"
)

// Values returns all known values for LocalStorage. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LocalStorage) Values() []LocalStorage {
	return []LocalStorage{
		"included",
		"excluded",
		"required",
	}
}

type LocalStorageType string

// Enum values for LocalStorageType
const (
	LocalStorageTypeHdd LocalStorageType = "hdd"
	LocalStorageTypeSsd LocalStorageType = "ssd"
)

// Values returns all known values for LocalStorageType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LocalStorageType) Values() []LocalStorageType {
	return []LocalStorageType{
		"hdd",
		"ssd",
	}
}

type MetricStatistic string

// Enum values for MetricStatistic
const (
	MetricStatisticAverage     MetricStatistic = "Average"
	MetricStatisticMinimum     MetricStatistic = "Minimum"
	MetricStatisticMaximum     MetricStatistic = "Maximum"
	MetricStatisticSampleCount MetricStatistic = "SampleCount"
	MetricStatisticSum         MetricStatistic = "Sum"
)

// Values returns all known values for MetricStatistic. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (MetricStatistic) Values() []MetricStatistic {
	return []MetricStatistic{
		"Average",
		"Minimum",
		"Maximum",
		"SampleCount",
		"Sum",
	}
}

type MetricType string

// Enum values for MetricType
const (
	MetricTypeASGAverageCPUUtilization MetricType = "ASGAverageCPUUtilization"
	MetricTypeASGAverageNetworkIn      MetricType = "ASGAverageNetworkIn"
	MetricTypeASGAverageNetworkOut     MetricType = "ASGAverageNetworkOut"
	MetricTypeALBRequestCountPerTarget MetricType = "ALBRequestCountPerTarget"
)

// Values returns all known values for MetricType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (MetricType) Values() []MetricType {
	return []MetricType{
		"ASGAverageCPUUtilization",
		"ASGAverageNetworkIn",
		"ASGAverageNetworkOut",
		"ALBRequestCountPerTarget",
	}
}

type PredefinedLoadMetricType string

// Enum values for PredefinedLoadMetricType
const (
	PredefinedLoadMetricTypeASGTotalCPUUtilization     PredefinedLoadMetricType = "ASGTotalCPUUtilization"
	PredefinedLoadMetricTypeASGTotalNetworkIn          PredefinedLoadMetricType = "ASGTotalNetworkIn"
	PredefinedLoadMetricTypeASGTotalNetworkOut         PredefinedLoadMetricType = "ASGTotalNetworkOut"
	PredefinedLoadMetricTypeALBTargetGroupRequestCount PredefinedLoadMetricType = "ALBTargetGroupRequestCount"
)

// Values returns all known values for PredefinedLoadMetricType. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PredefinedLoadMetricType) Values() []PredefinedLoadMetricType {
	return []PredefinedLoadMetricType{
		"ASGTotalCPUUtilization",
		"ASGTotalNetworkIn",
		"ASGTotalNetworkOut",
		"ALBTargetGroupRequestCount",
	}
}

type PredefinedMetricPairType string

// Enum values for PredefinedMetricPairType
const (
	PredefinedMetricPairTypeASGCPUUtilization PredefinedMetricPairType = "ASGCPUUtilization"
	PredefinedMetricPairTypeASGNetworkIn      PredefinedMetricPairType = "ASGNetworkIn"
	PredefinedMetricPairTypeASGNetworkOut     PredefinedMetricPairType = "ASGNetworkOut"
	PredefinedMetricPairTypeALBRequestCount   PredefinedMetricPairType = "ALBRequestCount"
)

// Values returns all known values for PredefinedMetricPairType. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PredefinedMetricPairType) Values() []PredefinedMetricPairType {
	return []PredefinedMetricPairType{
		"ASGCPUUtilization",
		"ASGNetworkIn",
		"ASGNetworkOut",
		"ALBRequestCount",
	}
}

type PredefinedScalingMetricType string

// Enum values for PredefinedScalingMetricType
const (
	PredefinedScalingMetricTypeASGAverageCPUUtilization PredefinedScalingMetricType = "ASGAverageCPUUtilization"
	PredefinedScalingMetricTypeASGAverageNetworkIn      PredefinedScalingMetricType = "ASGAverageNetworkIn"
	PredefinedScalingMetricTypeASGAverageNetworkOut     PredefinedScalingMetricType = "ASGAverageNetworkOut"
	PredefinedScalingMetricTypeALBRequestCountPerTarget PredefinedScalingMetricType = "ALBRequestCountPerTarget"
)

// Values returns all known values for PredefinedScalingMetricType. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PredefinedScalingMetricType) Values() []PredefinedScalingMetricType {
	return []PredefinedScalingMetricType{
		"ASGAverageCPUUtilization",
		"ASGAverageNetworkIn",
		"ASGAverageNetworkOut",
		"ALBRequestCountPerTarget",
	}
}

type PredictiveScalingMaxCapacityBreachBehavior string

// Enum values for PredictiveScalingMaxCapacityBreachBehavior
const (
	PredictiveScalingMaxCapacityBreachBehaviorHonorMaxCapacity    PredictiveScalingMaxCapacityBreachBehavior = "HonorMaxCapacity"
	PredictiveScalingMaxCapacityBreachBehaviorIncreaseMaxCapacity PredictiveScalingMaxCapacityBreachBehavior = "IncreaseMaxCapacity"
)

// Values returns all known values for PredictiveScalingMaxCapacityBreachBehavior.
// Note that this can be expanded in the future, and so it is only as up to date as
// the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PredictiveScalingMaxCapacityBreachBehavior) Values() []PredictiveScalingMaxCapacityBreachBehavior {
	return []PredictiveScalingMaxCapacityBreachBehavior{
		"HonorMaxCapacity",
		"IncreaseMaxCapacity",
	}
}

type PredictiveScalingMode string

// Enum values for PredictiveScalingMode
const (
	PredictiveScalingModeForecastAndScale PredictiveScalingMode = "ForecastAndScale"
	PredictiveScalingModeForecastOnly     PredictiveScalingMode = "ForecastOnly"
)

// Values returns all known values for PredictiveScalingMode. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PredictiveScalingMode) Values() []PredictiveScalingMode {
	return []PredictiveScalingMode{
		"ForecastAndScale",
		"ForecastOnly",
	}
}

type RefreshStrategy string

// Enum values for RefreshStrategy
const (
	RefreshStrategyRolling RefreshStrategy = "Rolling"
)

// Values returns all known values for RefreshStrategy. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RefreshStrategy) Values() []RefreshStrategy {
	return []RefreshStrategy{
		"Rolling",
	}
}

type ScaleInProtectedInstances string

// Enum values for ScaleInProtectedInstances
const (
	ScaleInProtectedInstancesRefresh ScaleInProtectedInstances = "Refresh"
	ScaleInProtectedInstancesIgnore  ScaleInProtectedInstances = "Ignore"
	ScaleInProtectedInstancesWait    ScaleInProtectedInstances = "Wait"
)

// Values returns all known values for ScaleInProtectedInstances. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ScaleInProtectedInstances) Values() []ScaleInProtectedInstances {
	return []ScaleInProtectedInstances{
		"Refresh",
		"Ignore",
		"Wait",
	}
}

type ScalingActivityStatusCode string

// Enum values for ScalingActivityStatusCode
const (
	ScalingActivityStatusCodePendingSpotBidPlacement         ScalingActivityStatusCode = "PendingSpotBidPlacement"
	ScalingActivityStatusCodeWaitingForSpotInstanceRequestId ScalingActivityStatusCode = "WaitingForSpotInstanceRequestId"
	ScalingActivityStatusCodeWaitingForSpotInstanceId        ScalingActivityStatusCode = "WaitingForSpotInstanceId"
	ScalingActivityStatusCodeWaitingForInstanceId            ScalingActivityStatusCode = "WaitingForInstanceId"
	ScalingActivityStatusCodePreInService                    ScalingActivityStatusCode = "PreInService"
	ScalingActivityStatusCodeInProgress                      ScalingActivityStatusCode = "InProgress"
	ScalingActivityStatusCodeWaitingForELBConnectionDraining ScalingActivityStatusCode = "WaitingForELBConnectionDraining"
	ScalingActivityStatusCodeMidLifecycleAction              ScalingActivityStatusCode = "MidLifecycleAction"
	ScalingActivityStatusCodeWaitingForInstanceWarmup        ScalingActivityStatusCode = "WaitingForInstanceWarmup"
	ScalingActivityStatusCodeSuccessful                      ScalingActivityStatusCode = "Successful"
	ScalingActivityStatusCodeFailed                          ScalingActivityStatusCode = "Failed"
	ScalingActivityStatusCodeCancelled                       ScalingActivityStatusCode = "Cancelled"
	ScalingActivityStatusCodeWaitingForConnectionDraining    ScalingActivityStatusCode = "WaitingForConnectionDraining"
)

// Values returns all known values for ScalingActivityStatusCode. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ScalingActivityStatusCode) Values() []ScalingActivityStatusCode {
	return []ScalingActivityStatusCode{
		"PendingSpotBidPlacement",
		"WaitingForSpotInstanceRequestId",
		"WaitingForSpotInstanceId",
		"WaitingForInstanceId",
		"PreInService",
		"InProgress",
		"WaitingForELBConnectionDraining",
		"MidLifecycleAction",
		"WaitingForInstanceWarmup",
		"Successful",
		"Failed",
		"Cancelled",
		"WaitingForConnectionDraining",
	}
}

type StandbyInstances string

// Enum values for StandbyInstances
const (
	StandbyInstancesTerminate StandbyInstances = "Terminate"
	StandbyInstancesIgnore    StandbyInstances = "Ignore"
	StandbyInstancesWait      StandbyInstances = "Wait"
)

// Values returns all known values for StandbyInstances. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (StandbyInstances) Values() []StandbyInstances {
	return []StandbyInstances{
		"Terminate",
		"Ignore",
		"Wait",
	}
}

type WarmPoolState string

// Enum values for WarmPoolState
const (
	WarmPoolStateStopped    WarmPoolState = "Stopped"
	WarmPoolStateRunning    WarmPoolState = "Running"
	WarmPoolStateHibernated WarmPoolState = "Hibernated"
)

// Values returns all known values for WarmPoolState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (WarmPoolState) Values() []WarmPoolState {
	return []WarmPoolState{
		"Stopped",
		"Running",
		"Hibernated",
	}
}

type WarmPoolStatus string

// Enum values for WarmPoolStatus
const (
	WarmPoolStatusPendingDelete WarmPoolStatus = "PendingDelete"
)

// Values returns all known values for WarmPoolStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (WarmPoolStatus) Values() []WarmPoolStatus {
	return []WarmPoolStatus{
		"PendingDelete",
	}
}
