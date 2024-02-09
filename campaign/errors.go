package campaign

import "fmt"

type ValidationError struct {
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Reason)
}

func (e *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	if !ok {
		return false
	}
	return ok && t.Reason == e.Reason
}

type DateError struct {
	Reason string
}

func (e *DateError) Error() string {
	return fmt.Sprintf("date error: %s", e.Reason)
}

func (e *DateError) Is(target error) bool {
	t, ok := target.(*DateError)
	if !ok {
		return false
	}
	return ok && t.Reason == e.Reason
}

type LimitError struct {
	Reason string
}

func (e *LimitError) Error() string {
	return fmt.Sprintf("limit error: %s", e.Reason)
}

func (e *LimitError) Is(target error) bool {
	t, ok := target.(*LimitError)
	if !ok {
		return false
	}
	return ok && t.Reason == e.Reason
}

type CampaignTypeError struct {
	Reason string
}

func (e *CampaignTypeError) Error() string {
	return fmt.Sprintf("campaign type error: %s", e.Reason)
}

func (e *CampaignTypeError) Is(target error) bool {
	t, ok := target.(*CampaignTypeError)
	if !ok {
		return false
	}
	return ok && t.Reason == e.Reason
}
