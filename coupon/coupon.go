package coupon

import (
	"github.com/hibrid/coupons/generator"
)

type CouponConfig struct {
	Prefix        string `json:"prefix"`
	Suffix        string `json:"suffix"`
	Pattern       string `json:"pattern"`
	Divider       string `json:"divider"`
	CharacterSet  string `json:"characterSet"`
	Count         int    `json:"count"`
	MinimumLength int    `json:"minimumLength"`
	//VanityName       string `json:"vanityName,omitempty"`
	PatternCharacter string `json:"patternCharacter,omitempty"`
	SHA256Index      int    `json:"sha256Index,omitempty"`
}

type Coupon interface {
	GenerateCoupon() (string, error)
	GenerateCoupons() ([]string, error)
	ValidateCoupon(coupon string) (bool, error)
	GetGenerator() generator.GeneratorInterface
}

// CouponContext is a struct that holds the coupon configuration and is used to generate and validate the coupon
type CouponContext struct {
	CouponConfig *CouponConfig
	generator    generator.GeneratorInterface
}

// NewCouponContext creates a new coupon context based on the provided configuration
func NewCouponContext(config *CouponConfig) (Coupon, error) {
	// Initialize a slice to hold the options
	var options []generator.Option

	// Conditionally append options if they are not empty or zero
	if config.MinimumLength > 0 {
		options = append(options, generator.SetMinimumLength(uint32(config.MinimumLength)))
	}
	if config.Count > 0 {
		options = append(options, generator.SetGenerateCount(uint32(config.Count)))
	}
	if config.Pattern != "" {
		options = append(options, generator.SetPattern(config.Pattern))
	}
	if config.Divider != "" {
		options = append(options, generator.SetPatternDivider(config.Divider))
	}
	if config.CharacterSet != "" {
		options = append(options, generator.SetCharset(config.CharacterSet))
	}
	if config.Prefix != "" {
		options = append(options, generator.SetPrefix(config.Prefix))
	}
	if config.Suffix != "" {
		options = append(options, generator.SetSuffix(config.Suffix))
	}
	if config.PatternCharacter != "" {
		options = append(options, generator.SetPatternCharacter(config.PatternCharacter))
	}
	if config.SHA256Index > 0 {
		options = append(options, generator.SetCheckCharacterSHA256Index(uint32(config.SHA256Index)))
	}

	// Generate the coupon based on the configuration with the conditional options
	g, err := generator.NewWithOptions(options...)
	if err != nil {
		return nil, err
	}

	return &CouponContext{
		CouponConfig: config,
		generator:    g,
	}, nil
}

func (c *CouponContext) GetGenerator() generator.GeneratorInterface {
	return c.generator
}

// GenerateCoupon generates a coupon based on the configuration
func (c *CouponContext) GenerateCoupon() (string, error) {

	// Generate the coupon based on the configuration
	coupons, err := c.GenerateCoupons()
	if err != nil {
		return "", err
	}
	return coupons[0], nil
}

func (c *CouponContext) GenerateCoupons() ([]string, error) {
	coupons, err := c.generator.Run()
	if err != nil {
		return []string{}, err
	}

	return coupons, nil
}

// ValidateCoupon validates the provided coupon
func (c *CouponContext) ValidateCoupon(coupon string) (bool, error) {
	// Validate the provided coupon
	validatedCoupon, err := c.generator.Validate(coupon)
	if err != nil {
		return false, err
	}
	return validatedCoupon == coupon, nil
}
