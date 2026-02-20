package request

// AdditionalRepairServiceItem represents the minimal information required
// to associate a service with an additional repair.
type AdditionalRepairServiceItem struct {
	ID uint `json:"id" binding:"required"`
}

// AdditionalRepairPartsSupplyItem represents the minimal information required
// to associate a parts supply with an additional repair.
type AdditionalRepairPartsSupplyItem struct {
	ID              uint `json:"id" binding:"required"`
	QuantityReserve int  `json:"quantity_reserve" binding:"required"`
}

// AdditionalRepairCreateRequest is the payload used to create a new
// additional repair.
type AdditionalRepairCreateRequest struct {
	ServiceOrderID uint                              `json:"service_order_id" binding:"required"`
	Description    string                            `json:"description" binding:"required"`
	Services       []AdditionalRepairServiceItem     `json:"services"`
	PartsSupplies  []AdditionalRepairPartsSupplyItem `json:"parts_supplies"`
}

// AdditionalRepairItemsRequest is the payload used to add or remove
// services and parts supplies from an existing additional repair.
type AdditionalRepairItemsRequest struct {
	ServiceOrderID uint                              `json:"service_order_id"`
	Description    string                            `json:"description"`
	Services       []AdditionalRepairServiceItem     `json:"services"`
	PartsSupplies  []AdditionalRepairPartsSupplyItem `json:"parts_supplies"`
}

// AdditionalRepairApprovalRequest is the payload used to capture
// customer approval decisions for an additional repair.
type AdditionalRepairApprovalRequest struct {
	ApprovalStatus string `json:"approval_status" binding:"required,oneof=APPROVED DENIED"`
}
