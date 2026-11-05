package keyservice

// STACKITKmsKey is used to serialize a STACKIT KMS key over the keyservice protocol.
// This is defined separately from the protobuf-generated types to avoid
// requiring protobuf regeneration. It follows the same pattern as other key types.
type STACKITKmsKey struct {
	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *STACKITKmsKey) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

// Key_STACKITKmsKey is a wrapper for STACKITKmsKey to be used in Key.KeyType oneof.
type Key_STACKITKmsKey struct {
	STACKITKmsKey *STACKITKmsKey `protobuf:"bytes,8,opt,name=stackit_kms_key,json=stackitKmsKey,proto3,oneof"`
}

func (*Key_STACKITKmsKey) isKey_KeyType() {}

// GetSTACKITKmsKey returns the STACKITKmsKey from the Key, if set.
func (x *Key) GetSTACKITKmsKey() *STACKITKmsKey {
	if x, ok := x.GetKeyType().(*Key_STACKITKmsKey); ok {
		return x.STACKITKmsKey
	}
	return nil
}
