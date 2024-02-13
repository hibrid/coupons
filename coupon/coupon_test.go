package coupon

import (
	"reflect"
	"testing"

	"github.com/hibrid/coupons/generator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// MockGenerator is a mock of Generator interface
type MockGenerator struct {
	ctrl                        *gomock.Controller
	recorder                    *MockGeneratorMockRecorder
	RunFunc                     func() ([]string, error)
	CountUniqueCombinationsFunc func(pattern, patternChar, alphanumeric string) int
	ValidateFunc                func(code string) (string, error)
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock: mock}
	return mock
}

func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

func (m *MockGenerator) CountUniqueCombinations(pattern, patternChar, alphanumeric string) int {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CountUniqueCombinations", pattern, patternChar, alphanumeric)
	if m.CountUniqueCombinationsFunc != nil {
		return m.CountUniqueCombinationsFunc(pattern, patternChar, alphanumeric)
	}
	return 0
}

func (m *MockGenerator) Run() ([]string, error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
	if m.RunFunc != nil {
		return m.RunFunc()
	}
	return []string{"TEST-COUPON"}, nil
}

func (m *MockGenerator) Validate(code string) (string, error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Validate", code)
	if m.ValidateFunc != nil {
		return m.ValidateFunc(code)
	}
	return "", nil
}

func (mr *MockGeneratorMockRecorder) CountUniqueCombinations(pattern, patternChar, alphanumeric string) *gomock.Call {
	mr.mock.ctrl.T.Helper()

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUniqueCombinations", reflect.TypeOf((*MockGenerator)(nil).CountUniqueCombinations), pattern, patternChar, alphanumeric)
}

func (mr *MockGeneratorMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockGenerator)(nil).Run))
}

func (mr *MockGeneratorMockRecorder) Validate(code string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockGenerator)(nil).Validate), code)
}

// TestNewCouponContext tests the creation of a new CouponContext
func TestNewCouponContext(t *testing.T) {
	config := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		CharacterSet:     "ABC123",
		Count:            1,
		MinimumLength:    8,
		PatternCharacter: "#",
	}
	coupon, err := NewCouponContext(config)
	assert.NoError(t, err)
	assert.NotNil(t, coupon)
}

// TestGenerateCoupon tests generating a single coupon
func TestGenerateCoupon(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockGenerator := NewMockGenerator(ctrl)
	mockGenerator.EXPECT().Run().Return([]string{"PRE-XXXX-SUF"}, nil)

	couponConfig := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		CharacterSet:     "ABC123",
		Count:            1,
		MinimumLength:    5,
		PatternCharacter: "#",
	}

	realCoupon, err := NewCouponContext(couponConfig)
	assert.NoError(t, err)

	mockGenerator.RunFunc = realCoupon.GetGenerator().Run
	couponContext := &CouponContext{
		CouponConfig: couponConfig,
		generator:    mockGenerator,
	}

	coupon, err := couponContext.GenerateCoupon()
	assert.NoError(t, err)
	assert.Regexp(t, "PRE-\\w{4}-SUF", coupon)
}

// TestGenerateCoupon tests generating a single coupon
func TestBadWordPrefix(t *testing.T) {

	couponConfig := &CouponConfig{
		Prefix:           "BOOB",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		CharacterSet:     "ABC123",
		Count:            1,
		MinimumLength:    1,
		PatternCharacter: "#",
	}

	_, err := NewCouponContext(couponConfig)
	assert.ErrorIs(t, err, generator.ErrPrefixBadWord)

}

// TestGenerateCoupon tests generating a single coupon
func TestBadwordSuffixLowerCase(t *testing.T) {

	couponConfig := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "Boob",
		Pattern:          "####",
		Divider:          "-",
		CharacterSet:     "ABC123",
		Count:            1,
		MinimumLength:    1,
		PatternCharacter: "#",
	}

	_, err := NewCouponContext(couponConfig)
	assert.ErrorIs(t, err, generator.ErrSuffixBadWord)
}

// TestGenerateCoupons tests generating multiple coupons
func TestGenerateCoupons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerator := NewMockGenerator(ctrl)
	mockGenerator.EXPECT().Run().Return([]string{"PRE-XXXX-SUF"}, nil)

	couponConfig := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		Count:            10,
		MinimumLength:    5,
		PatternCharacter: "#",
	}

	realCoupons, err := NewCouponContext(couponConfig)
	assert.NoError(t, err)

	mockGenerator.RunFunc = realCoupons.GetGenerator().Run
	couponContext := &CouponContext{
		CouponConfig: couponConfig,
		generator:    mockGenerator,
	}

	coupons, err := couponContext.GenerateCoupons()
	assert.NoError(t, err)
	assert.Len(t, coupons, 10)
	for _, coupon := range coupons {
		assert.Regexp(t, "PRE-\\w{4}-SUF", coupon)
	}
}

// TestValidateCoupon tests the coupon validation logic
func TestValidateCoupon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerator := NewMockGenerator(ctrl)

	couponConfig := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		Count:            10,
		MinimumLength:    5,
		PatternCharacter: "#",
	}

	realCoupons, err := NewCouponContext(couponConfig)
	assert.NoError(t, err)

	coupons := []string{
		"PRE-S9E5-SUF",
		"PRE-9UFD-SUF",
		"PRE-7B45-SUF",
		"PRE-4RRG-SUF",
		"PRE-A52Q-SUF",
		"PRE-67AK-SUF",
		"PRE-TQ2P-SUF",
		"PRE-QPR5-SUF",
		"PRE-CT13-SUF",
	}
	mockGenerator.ValidateFunc = realCoupons.GetGenerator().Validate
	for _, coupon := range coupons {
		mockGenerator.EXPECT().Validate(coupon).Return(coupon, nil)

		couponContext := &CouponContext{
			CouponConfig: couponConfig,
			generator:    mockGenerator,
		}

		valid, err := couponContext.ValidateCoupon(coupon)
		assert.NoError(t, err)
		assert.True(t, valid)
	}

}

// TestValidateCoupon tests the coupon validation logic
func TestValidateBadCoupons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerator := NewMockGenerator(ctrl)

	couponConfig := &CouponConfig{
		Prefix:           "PRE",
		Suffix:           "SUF",
		Pattern:          "####",
		Divider:          "-",
		Count:            10,
		MinimumLength:    5,
		PatternCharacter: "#",
	}

	realCoupons, err := NewCouponContext(couponConfig)
	assert.NoError(t, err)

	tests := []struct {
		coupon      string
		wantErr     bool
		wantErrType error // Add this field for checking error type
	}{
		{"PRE-S9E5-SUF", false, nil},
		{"boob-9UFD-SUF", true, generator.ErrPrefixNotFound},
		{"PRE-7B45-SUF", false, nil},
		{"PRE-4RRG-SUF", false, nil},
		{"PRE-A52Q-SUF", false, nil},
		{"PRE-67AK-SUF", false, nil},
		{"PRE-TQ2P-SUF", false, nil},
		{"PRE-QPR5-SUF", false, nil},
		{"PRE-CT13-SUF", false, nil},
	}

	mockGenerator.ValidateFunc = realCoupons.GetGenerator().Validate
	for _, test := range tests {
		mockGenerator.EXPECT().Validate(test.coupon).Return(test.coupon, nil)

		couponContext := &CouponContext{
			CouponConfig: couponConfig,
			generator:    mockGenerator,
		}

		valid, err := couponContext.ValidateCoupon(test.coupon)
		if test.wantErr {
			assert.Error(t, err)
			assert.ErrorIs(t, err, test.wantErrType)
		} else {
			assert.NoError(t, err)
			assert.True(t, valid)
		}
	}

}
