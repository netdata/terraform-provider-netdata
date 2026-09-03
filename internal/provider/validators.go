package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	// Embeds the IANA timezone database in the provider binary so
	// ianaTimezoneValidator's time.LoadLocation does not depend on the host
	// having a zoneinfo database (absent on Windows without Go installed, and
	// in scratch/distroless/minimal container images).
	_ "time/tzdata"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ianaTimezoneValidator struct{}

func (v ianaTimezoneValidator) Description(_ context.Context) string {
	return "Value must be a valid IANA timezone (e.g. \"America/New_York\", \"UTC\")."
}

func (v ianaTimezoneValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ianaTimezoneValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	tz := strings.TrimSpace(req.ConfigValue.ValueString())
	if _, err := time.LoadLocation(tz); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid IANA Timezone",
			fmt.Sprintf("%q is not a valid IANA timezone: %s", tz, err),
		)
	}
}

// IANATimezone returns a validator that checks whether the string is a valid IANA timezone.
func IANATimezone() validator.String {
	return ianaTimezoneValidator{}
}

type rruleOnlyValidator struct{}

func (v rruleOnlyValidator) Description(_ context.Context) string {
	return "Value must not begin with a DTSTART line; use starts_at instead."
}

func (v rruleOnlyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v rruleOnlyValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	rrule := strings.TrimSpace(req.ConfigValue.ValueString())
	if strings.HasPrefix(strings.ToUpper(rrule), "DTSTART") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid RRULE",
			fmt.Sprintf("%q must not begin with a DTSTART line. The recurrence start is taken from starts_at; provide only the RRULE part, e.g. \"RRULE:FREQ=MONTHLY;INTERVAL=1\".", rrule),
		)
	}
}

// RRuleOnly returns a validator that requires the value to contain only the recurrence
// rule itself, rejecting values that lead with a schedule start (DTSTART) line.
func RRuleOnly() validator.String {
	return rruleOnlyValidator{}
}

type utcTimestampValidator struct{}

func (v utcTimestampValidator) Description(_ context.Context) string {
	return "Value must be an RFC3339 timestamp expressed in UTC (e.g. \"2026-09-01T00:00:00Z\")."
}

func (v utcTimestampValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v utcTimestampValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Unreachable while the attribute uses timetypes.RFC3339Type, which rejects
		// malformed values before attribute validators run. Kept so the check still
		// reports rather than silently passes if the custom type is ever removed.
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid RFC3339 Timestamp",
			fmt.Sprintf("%q is not a valid RFC3339 time: %s", value, err),
		)
		return
	}

	if _, offset := t.Zone(); offset != 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Timestamp Must Be UTC",
			fmt.Sprintf("%q uses a non-UTC offset. The API stores and returns timestamps in UTC, so a non-UTC value cannot round-trip. Express it in UTC instead: %q.", value, t.UTC().Format(time.RFC3339)),
		)
	}
}

// UTCTimestamp returns a validator that requires an RFC3339 timestamp with a UTC
// offset, rejecting values the API would normalize to UTC and hand back changed.
func UTCTimestamp() validator.String {
	return utcTimestampValidator{}
}

type uuidValidator struct{}

func (v uuidValidator) Description(_ context.Context) string {
	return "Value must be a valid UUID string (e.g. \"123e4567-e89b-12d3-a456-426614174000\")."
}

func (v uuidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uuidValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	err := uuid.Validate(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid UUID",
			fmt.Sprintf("%q is not a valid UUID: %s", value, err),
		)
		return
	}
}

// UUID returns a validator that requires a valid UUID string.
func UUID() validator.String {
	return uuidValidator{}
}
