package keyservice

// StackitKmsKey is used to serialize a STACKIT KMS key over the keyservice protocol.
// This is defined separately from the protobuf-generated types to avoid
// requiring protobuf regeneration. It follows the same pattern as other key types.
type StackitKmsKey struct {
	ResourceId string `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *StackitKmsKey) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

// Key_StackitKmsKey is a wrapper for StackitKmsKey to be used in Key.KeyType oneof.
type Key_StackitKmsKey struct {
	StackitKmsKey *StackitKmsKey `protobuf:"bytes,8,opt,name=stackit_kms_key,json=stackitKmsKey,proto3,oneof"`
}

func (*Key_StackitKmsKey) isKey_KeyType() {}

// GetStackitKmsKey returns the StackitKmsKey from the Key, if set.
func (x *Key) GetStackitKmsKey() *StackitKmsKey {
	if x, ok := x.GetKeyType().(*Key_StackitKmsKey); ok {
		return x.StackitKmsKey
	}
	return nil
}
