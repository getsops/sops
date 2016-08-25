package kms

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/aws/request"
import "github.com/aws/aws-sdk-go/service/kms"

type MockKMSAPI struct {
	mock.Mock
}

// CancelKeyDeletionRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CancelKeyDeletionRequest(_a0 *kms.CancelKeyDeletionInput) (*request.Request, *kms.CancelKeyDeletionOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.CancelKeyDeletionInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.CancelKeyDeletionOutput
	if rf, ok := ret.Get(1).(func(*kms.CancelKeyDeletionInput) *kms.CancelKeyDeletionOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.CancelKeyDeletionOutput)
		}
	}

	return r0, r1
}

// CancelKeyDeletion provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CancelKeyDeletion(_a0 *kms.CancelKeyDeletionInput) (*kms.CancelKeyDeletionOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.CancelKeyDeletionOutput
	if rf, ok := ret.Get(0).(func(*kms.CancelKeyDeletionInput) *kms.CancelKeyDeletionOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.CancelKeyDeletionOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.CancelKeyDeletionInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAliasRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateAliasRequest(_a0 *kms.CreateAliasInput) (*request.Request, *kms.CreateAliasOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.CreateAliasInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.CreateAliasOutput
	if rf, ok := ret.Get(1).(func(*kms.CreateAliasInput) *kms.CreateAliasOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.CreateAliasOutput)
		}
	}

	return r0, r1
}

// CreateAlias provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateAlias(_a0 *kms.CreateAliasInput) (*kms.CreateAliasOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.CreateAliasOutput
	if rf, ok := ret.Get(0).(func(*kms.CreateAliasInput) *kms.CreateAliasOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.CreateAliasOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.CreateAliasInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateGrantRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateGrantRequest(_a0 *kms.CreateGrantInput) (*request.Request, *kms.CreateGrantOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.CreateGrantInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.CreateGrantOutput
	if rf, ok := ret.Get(1).(func(*kms.CreateGrantInput) *kms.CreateGrantOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.CreateGrantOutput)
		}
	}

	return r0, r1
}

// CreateGrant provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateGrant(_a0 *kms.CreateGrantInput) (*kms.CreateGrantOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.CreateGrantOutput
	if rf, ok := ret.Get(0).(func(*kms.CreateGrantInput) *kms.CreateGrantOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.CreateGrantOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.CreateGrantInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateKeyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateKeyRequest(_a0 *kms.CreateKeyInput) (*request.Request, *kms.CreateKeyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.CreateKeyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.CreateKeyOutput
	if rf, ok := ret.Get(1).(func(*kms.CreateKeyInput) *kms.CreateKeyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.CreateKeyOutput)
		}
	}

	return r0, r1
}

// CreateKey provides a mock function with given fields: _a0
func (_m *MockKMSAPI) CreateKey(_a0 *kms.CreateKeyInput) (*kms.CreateKeyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.CreateKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.CreateKeyInput) *kms.CreateKeyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.CreateKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.CreateKeyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DecryptRequest(_a0 *kms.DecryptInput) (*request.Request, *kms.DecryptOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.DecryptInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.DecryptOutput
	if rf, ok := ret.Get(1).(func(*kms.DecryptInput) *kms.DecryptOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.DecryptOutput)
		}
	}

	return r0, r1
}

// Decrypt provides a mock function with given fields: _a0
func (_m *MockKMSAPI) Decrypt(_a0 *kms.DecryptInput) (*kms.DecryptOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.DecryptOutput
	if rf, ok := ret.Get(0).(func(*kms.DecryptInput) *kms.DecryptOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DecryptOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DecryptInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAliasRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DeleteAliasRequest(_a0 *kms.DeleteAliasInput) (*request.Request, *kms.DeleteAliasOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.DeleteAliasInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.DeleteAliasOutput
	if rf, ok := ret.Get(1).(func(*kms.DeleteAliasInput) *kms.DeleteAliasOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.DeleteAliasOutput)
		}
	}

	return r0, r1
}

// DeleteAlias provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DeleteAlias(_a0 *kms.DeleteAliasInput) (*kms.DeleteAliasOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.DeleteAliasOutput
	if rf, ok := ret.Get(0).(func(*kms.DeleteAliasInput) *kms.DeleteAliasOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DeleteAliasOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DeleteAliasInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeKeyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DescribeKeyRequest(_a0 *kms.DescribeKeyInput) (*request.Request, *kms.DescribeKeyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.DescribeKeyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.DescribeKeyOutput
	if rf, ok := ret.Get(1).(func(*kms.DescribeKeyInput) *kms.DescribeKeyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.DescribeKeyOutput)
		}
	}

	return r0, r1
}

// DescribeKey provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DescribeKey(_a0 *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.DescribeKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.DescribeKeyInput) *kms.DescribeKeyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DescribeKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DescribeKeyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DisableKeyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DisableKeyRequest(_a0 *kms.DisableKeyInput) (*request.Request, *kms.DisableKeyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.DisableKeyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.DisableKeyOutput
	if rf, ok := ret.Get(1).(func(*kms.DisableKeyInput) *kms.DisableKeyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.DisableKeyOutput)
		}
	}

	return r0, r1
}

// DisableKey provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DisableKey(_a0 *kms.DisableKeyInput) (*kms.DisableKeyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.DisableKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.DisableKeyInput) *kms.DisableKeyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DisableKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DisableKeyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DisableKeyRotationRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DisableKeyRotationRequest(_a0 *kms.DisableKeyRotationInput) (*request.Request, *kms.DisableKeyRotationOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.DisableKeyRotationInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.DisableKeyRotationOutput
	if rf, ok := ret.Get(1).(func(*kms.DisableKeyRotationInput) *kms.DisableKeyRotationOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.DisableKeyRotationOutput)
		}
	}

	return r0, r1
}

// DisableKeyRotation provides a mock function with given fields: _a0
func (_m *MockKMSAPI) DisableKeyRotation(_a0 *kms.DisableKeyRotationInput) (*kms.DisableKeyRotationOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.DisableKeyRotationOutput
	if rf, ok := ret.Get(0).(func(*kms.DisableKeyRotationInput) *kms.DisableKeyRotationOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DisableKeyRotationOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.DisableKeyRotationInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnableKeyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) EnableKeyRequest(_a0 *kms.EnableKeyInput) (*request.Request, *kms.EnableKeyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.EnableKeyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.EnableKeyOutput
	if rf, ok := ret.Get(1).(func(*kms.EnableKeyInput) *kms.EnableKeyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.EnableKeyOutput)
		}
	}

	return r0, r1
}

// EnableKey provides a mock function with given fields: _a0
func (_m *MockKMSAPI) EnableKey(_a0 *kms.EnableKeyInput) (*kms.EnableKeyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.EnableKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.EnableKeyInput) *kms.EnableKeyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.EnableKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.EnableKeyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnableKeyRotationRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) EnableKeyRotationRequest(_a0 *kms.EnableKeyRotationInput) (*request.Request, *kms.EnableKeyRotationOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.EnableKeyRotationInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.EnableKeyRotationOutput
	if rf, ok := ret.Get(1).(func(*kms.EnableKeyRotationInput) *kms.EnableKeyRotationOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.EnableKeyRotationOutput)
		}
	}

	return r0, r1
}

// EnableKeyRotation provides a mock function with given fields: _a0
func (_m *MockKMSAPI) EnableKeyRotation(_a0 *kms.EnableKeyRotationInput) (*kms.EnableKeyRotationOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.EnableKeyRotationOutput
	if rf, ok := ret.Get(0).(func(*kms.EnableKeyRotationInput) *kms.EnableKeyRotationOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.EnableKeyRotationOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.EnableKeyRotationInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) EncryptRequest(_a0 *kms.EncryptInput) (*request.Request, *kms.EncryptOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.EncryptInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.EncryptOutput
	if rf, ok := ret.Get(1).(func(*kms.EncryptInput) *kms.EncryptOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.EncryptOutput)
		}
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: _a0
func (_m *MockKMSAPI) Encrypt(_a0 *kms.EncryptInput) (*kms.EncryptOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.EncryptOutput
	if rf, ok := ret.Get(0).(func(*kms.EncryptInput) *kms.EncryptOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.EncryptOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.EncryptInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateDataKeyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateDataKeyRequest(_a0 *kms.GenerateDataKeyInput) (*request.Request, *kms.GenerateDataKeyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.GenerateDataKeyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.GenerateDataKeyOutput
	if rf, ok := ret.Get(1).(func(*kms.GenerateDataKeyInput) *kms.GenerateDataKeyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.GenerateDataKeyOutput)
		}
	}

	return r0, r1
}

// GenerateDataKey provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateDataKey(_a0 *kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.GenerateDataKeyOutput
	if rf, ok := ret.Get(0).(func(*kms.GenerateDataKeyInput) *kms.GenerateDataKeyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.GenerateDataKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.GenerateDataKeyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateDataKeyWithoutPlaintextRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateDataKeyWithoutPlaintextRequest(_a0 *kms.GenerateDataKeyWithoutPlaintextInput) (*request.Request, *kms.GenerateDataKeyWithoutPlaintextOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.GenerateDataKeyWithoutPlaintextInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.GenerateDataKeyWithoutPlaintextOutput
	if rf, ok := ret.Get(1).(func(*kms.GenerateDataKeyWithoutPlaintextInput) *kms.GenerateDataKeyWithoutPlaintextOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.GenerateDataKeyWithoutPlaintextOutput)
		}
	}

	return r0, r1
}

// GenerateDataKeyWithoutPlaintext provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateDataKeyWithoutPlaintext(_a0 *kms.GenerateDataKeyWithoutPlaintextInput) (*kms.GenerateDataKeyWithoutPlaintextOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.GenerateDataKeyWithoutPlaintextOutput
	if rf, ok := ret.Get(0).(func(*kms.GenerateDataKeyWithoutPlaintextInput) *kms.GenerateDataKeyWithoutPlaintextOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.GenerateDataKeyWithoutPlaintextOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.GenerateDataKeyWithoutPlaintextInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateRandomRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateRandomRequest(_a0 *kms.GenerateRandomInput) (*request.Request, *kms.GenerateRandomOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.GenerateRandomInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.GenerateRandomOutput
	if rf, ok := ret.Get(1).(func(*kms.GenerateRandomInput) *kms.GenerateRandomOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.GenerateRandomOutput)
		}
	}

	return r0, r1
}

// GenerateRandom provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GenerateRandom(_a0 *kms.GenerateRandomInput) (*kms.GenerateRandomOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.GenerateRandomOutput
	if rf, ok := ret.Get(0).(func(*kms.GenerateRandomInput) *kms.GenerateRandomOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.GenerateRandomOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.GenerateRandomInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetKeyPolicyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GetKeyPolicyRequest(_a0 *kms.GetKeyPolicyInput) (*request.Request, *kms.GetKeyPolicyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.GetKeyPolicyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.GetKeyPolicyOutput
	if rf, ok := ret.Get(1).(func(*kms.GetKeyPolicyInput) *kms.GetKeyPolicyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.GetKeyPolicyOutput)
		}
	}

	return r0, r1
}

// GetKeyPolicy provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GetKeyPolicy(_a0 *kms.GetKeyPolicyInput) (*kms.GetKeyPolicyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.GetKeyPolicyOutput
	if rf, ok := ret.Get(0).(func(*kms.GetKeyPolicyInput) *kms.GetKeyPolicyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.GetKeyPolicyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.GetKeyPolicyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetKeyRotationStatusRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GetKeyRotationStatusRequest(_a0 *kms.GetKeyRotationStatusInput) (*request.Request, *kms.GetKeyRotationStatusOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.GetKeyRotationStatusInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.GetKeyRotationStatusOutput
	if rf, ok := ret.Get(1).(func(*kms.GetKeyRotationStatusInput) *kms.GetKeyRotationStatusOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.GetKeyRotationStatusOutput)
		}
	}

	return r0, r1
}

// GetKeyRotationStatus provides a mock function with given fields: _a0
func (_m *MockKMSAPI) GetKeyRotationStatus(_a0 *kms.GetKeyRotationStatusInput) (*kms.GetKeyRotationStatusOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.GetKeyRotationStatusOutput
	if rf, ok := ret.Get(0).(func(*kms.GetKeyRotationStatusInput) *kms.GetKeyRotationStatusOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.GetKeyRotationStatusOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.GetKeyRotationStatusInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAliasesRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListAliasesRequest(_a0 *kms.ListAliasesInput) (*request.Request, *kms.ListAliasesOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ListAliasesInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ListAliasesOutput
	if rf, ok := ret.Get(1).(func(*kms.ListAliasesInput) *kms.ListAliasesOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ListAliasesOutput)
		}
	}

	return r0, r1
}

// ListAliases provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListAliases(_a0 *kms.ListAliasesInput) (*kms.ListAliasesOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ListAliasesOutput
	if rf, ok := ret.Get(0).(func(*kms.ListAliasesInput) *kms.ListAliasesOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListAliasesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListAliasesInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAliasesPages provides a mock function with given fields: _a0, _a1
func (_m *MockKMSAPI) ListAliasesPages(_a0 *kms.ListAliasesInput, _a1 func(*kms.ListAliasesOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*kms.ListAliasesInput, func(*kms.ListAliasesOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListGrantsRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListGrantsRequest(_a0 *kms.ListGrantsInput) (*request.Request, *kms.ListGrantsResponse) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ListGrantsInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ListGrantsResponse
	if rf, ok := ret.Get(1).(func(*kms.ListGrantsInput) *kms.ListGrantsResponse); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ListGrantsResponse)
		}
	}

	return r0, r1
}

// ListGrants provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListGrants(_a0 *kms.ListGrantsInput) (*kms.ListGrantsResponse, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ListGrantsResponse
	if rf, ok := ret.Get(0).(func(*kms.ListGrantsInput) *kms.ListGrantsResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListGrantsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListGrantsInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGrantsPages provides a mock function with given fields: _a0, _a1
func (_m *MockKMSAPI) ListGrantsPages(_a0 *kms.ListGrantsInput, _a1 func(*kms.ListGrantsResponse, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*kms.ListGrantsInput, func(*kms.ListGrantsResponse, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListKeyPoliciesRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListKeyPoliciesRequest(_a0 *kms.ListKeyPoliciesInput) (*request.Request, *kms.ListKeyPoliciesOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ListKeyPoliciesInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ListKeyPoliciesOutput
	if rf, ok := ret.Get(1).(func(*kms.ListKeyPoliciesInput) *kms.ListKeyPoliciesOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ListKeyPoliciesOutput)
		}
	}

	return r0, r1
}

// ListKeyPolicies provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListKeyPolicies(_a0 *kms.ListKeyPoliciesInput) (*kms.ListKeyPoliciesOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ListKeyPoliciesOutput
	if rf, ok := ret.Get(0).(func(*kms.ListKeyPoliciesInput) *kms.ListKeyPoliciesOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListKeyPoliciesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListKeyPoliciesInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListKeyPoliciesPages provides a mock function with given fields: _a0, _a1
func (_m *MockKMSAPI) ListKeyPoliciesPages(_a0 *kms.ListKeyPoliciesInput, _a1 func(*kms.ListKeyPoliciesOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*kms.ListKeyPoliciesInput, func(*kms.ListKeyPoliciesOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListKeysRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListKeysRequest(_a0 *kms.ListKeysInput) (*request.Request, *kms.ListKeysOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ListKeysInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ListKeysOutput
	if rf, ok := ret.Get(1).(func(*kms.ListKeysInput) *kms.ListKeysOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ListKeysOutput)
		}
	}

	return r0, r1
}

// ListKeys provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListKeys(_a0 *kms.ListKeysInput) (*kms.ListKeysOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ListKeysOutput
	if rf, ok := ret.Get(0).(func(*kms.ListKeysInput) *kms.ListKeysOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListKeysOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListKeysInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListKeysPages provides a mock function with given fields: _a0, _a1
func (_m *MockKMSAPI) ListKeysPages(_a0 *kms.ListKeysInput, _a1 func(*kms.ListKeysOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*kms.ListKeysInput, func(*kms.ListKeysOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListRetirableGrantsRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListRetirableGrantsRequest(_a0 *kms.ListRetirableGrantsInput) (*request.Request, *kms.ListGrantsResponse) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ListRetirableGrantsInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ListGrantsResponse
	if rf, ok := ret.Get(1).(func(*kms.ListRetirableGrantsInput) *kms.ListGrantsResponse); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ListGrantsResponse)
		}
	}

	return r0, r1
}

// ListRetirableGrants provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ListRetirableGrants(_a0 *kms.ListRetirableGrantsInput) (*kms.ListGrantsResponse, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ListGrantsResponse
	if rf, ok := ret.Get(0).(func(*kms.ListRetirableGrantsInput) *kms.ListGrantsResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ListGrantsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ListRetirableGrantsInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutKeyPolicyRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) PutKeyPolicyRequest(_a0 *kms.PutKeyPolicyInput) (*request.Request, *kms.PutKeyPolicyOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.PutKeyPolicyInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.PutKeyPolicyOutput
	if rf, ok := ret.Get(1).(func(*kms.PutKeyPolicyInput) *kms.PutKeyPolicyOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.PutKeyPolicyOutput)
		}
	}

	return r0, r1
}

// PutKeyPolicy provides a mock function with given fields: _a0
func (_m *MockKMSAPI) PutKeyPolicy(_a0 *kms.PutKeyPolicyInput) (*kms.PutKeyPolicyOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.PutKeyPolicyOutput
	if rf, ok := ret.Get(0).(func(*kms.PutKeyPolicyInput) *kms.PutKeyPolicyOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.PutKeyPolicyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.PutKeyPolicyInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReEncryptRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ReEncryptRequest(_a0 *kms.ReEncryptInput) (*request.Request, *kms.ReEncryptOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ReEncryptInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ReEncryptOutput
	if rf, ok := ret.Get(1).(func(*kms.ReEncryptInput) *kms.ReEncryptOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ReEncryptOutput)
		}
	}

	return r0, r1
}

// ReEncrypt provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ReEncrypt(_a0 *kms.ReEncryptInput) (*kms.ReEncryptOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ReEncryptOutput
	if rf, ok := ret.Get(0).(func(*kms.ReEncryptInput) *kms.ReEncryptOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ReEncryptOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ReEncryptInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetireGrantRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) RetireGrantRequest(_a0 *kms.RetireGrantInput) (*request.Request, *kms.RetireGrantOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.RetireGrantInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.RetireGrantOutput
	if rf, ok := ret.Get(1).(func(*kms.RetireGrantInput) *kms.RetireGrantOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.RetireGrantOutput)
		}
	}

	return r0, r1
}

// RetireGrant provides a mock function with given fields: _a0
func (_m *MockKMSAPI) RetireGrant(_a0 *kms.RetireGrantInput) (*kms.RetireGrantOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.RetireGrantOutput
	if rf, ok := ret.Get(0).(func(*kms.RetireGrantInput) *kms.RetireGrantOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.RetireGrantOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.RetireGrantInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RevokeGrantRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) RevokeGrantRequest(_a0 *kms.RevokeGrantInput) (*request.Request, *kms.RevokeGrantOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.RevokeGrantInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.RevokeGrantOutput
	if rf, ok := ret.Get(1).(func(*kms.RevokeGrantInput) *kms.RevokeGrantOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.RevokeGrantOutput)
		}
	}

	return r0, r1
}

// RevokeGrant provides a mock function with given fields: _a0
func (_m *MockKMSAPI) RevokeGrant(_a0 *kms.RevokeGrantInput) (*kms.RevokeGrantOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.RevokeGrantOutput
	if rf, ok := ret.Get(0).(func(*kms.RevokeGrantInput) *kms.RevokeGrantOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.RevokeGrantOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.RevokeGrantInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScheduleKeyDeletionRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ScheduleKeyDeletionRequest(_a0 *kms.ScheduleKeyDeletionInput) (*request.Request, *kms.ScheduleKeyDeletionOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.ScheduleKeyDeletionInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.ScheduleKeyDeletionOutput
	if rf, ok := ret.Get(1).(func(*kms.ScheduleKeyDeletionInput) *kms.ScheduleKeyDeletionOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.ScheduleKeyDeletionOutput)
		}
	}

	return r0, r1
}

// ScheduleKeyDeletion provides a mock function with given fields: _a0
func (_m *MockKMSAPI) ScheduleKeyDeletion(_a0 *kms.ScheduleKeyDeletionInput) (*kms.ScheduleKeyDeletionOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.ScheduleKeyDeletionOutput
	if rf, ok := ret.Get(0).(func(*kms.ScheduleKeyDeletionInput) *kms.ScheduleKeyDeletionOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.ScheduleKeyDeletionOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.ScheduleKeyDeletionInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAliasRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) UpdateAliasRequest(_a0 *kms.UpdateAliasInput) (*request.Request, *kms.UpdateAliasOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.UpdateAliasInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.UpdateAliasOutput
	if rf, ok := ret.Get(1).(func(*kms.UpdateAliasInput) *kms.UpdateAliasOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.UpdateAliasOutput)
		}
	}

	return r0, r1
}

// UpdateAlias provides a mock function with given fields: _a0
func (_m *MockKMSAPI) UpdateAlias(_a0 *kms.UpdateAliasInput) (*kms.UpdateAliasOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.UpdateAliasOutput
	if rf, ok := ret.Get(0).(func(*kms.UpdateAliasInput) *kms.UpdateAliasOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.UpdateAliasOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.UpdateAliasInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateKeyDescriptionRequest provides a mock function with given fields: _a0
func (_m *MockKMSAPI) UpdateKeyDescriptionRequest(_a0 *kms.UpdateKeyDescriptionInput) (*request.Request, *kms.UpdateKeyDescriptionOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*kms.UpdateKeyDescriptionInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *kms.UpdateKeyDescriptionOutput
	if rf, ok := ret.Get(1).(func(*kms.UpdateKeyDescriptionInput) *kms.UpdateKeyDescriptionOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*kms.UpdateKeyDescriptionOutput)
		}
	}

	return r0, r1
}

// UpdateKeyDescription provides a mock function with given fields: _a0
func (_m *MockKMSAPI) UpdateKeyDescription(_a0 *kms.UpdateKeyDescriptionInput) (*kms.UpdateKeyDescriptionOutput, error) {
	ret := _m.Called(_a0)

	var r0 *kms.UpdateKeyDescriptionOutput
	if rf, ok := ret.Get(0).(func(*kms.UpdateKeyDescriptionInput) *kms.UpdateKeyDescriptionOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.UpdateKeyDescriptionOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*kms.UpdateKeyDescriptionInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
