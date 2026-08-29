package diagnosis_test

import (
	"testing"

	"github.com/mabuhasna8/IntelliOps/internal/contracts"
	"github.com/mabuhasna8/IntelliOps/internal/diagnosis"
)

func validAdapterResult() contracts.AdapterResult {
	var result contracts.AdapterResult

	// Avoid depending on the concrete classification type.
	result.Classification.Code = "test.classification"
	result.Facts = map[string]any{
		"severity": "warning",
		"source":   "unit-test",
	}

	return result
}

func TestBasicDiagnoser_RejectsMissingClassificationCode(t *testing.T) {
	d := diagnosis.BasicDiagnoser{}

	_, err := d.Diagnose(contracts.AdapterResult{})
	if err == nil {
		t.Fatal("expected an error for a missing classification code")
	}
}

func TestBasicDiagnoser_ProducesDiagnosis(t *testing.T) {
	d := diagnosis.BasicDiagnoser{}

	result, err := d.Diagnose(validAdapterResult())
	if err != nil {
		t.Fatalf("Diagnose returned an error: %v", err)
	}

	if result.Code != "test.classification" {
		t.Fatalf("unexpected diagnosis code: got %q", result.Code)
	}

	if result.Severity != "warning" {
		t.Fatalf("unexpected diagnosis severity: got %q", result.Severity)
	}

	if result.Details == nil {
		t.Fatal("expected diagnosis details to be populated")
	}

	if got := result.Details["source"]; got != "unit-test" {
		t.Fatalf("unexpected source detail: got %v", got)
	}
}

func TestBasicDiagnoser_DefaultsSeverityToInfo(t *testing.T) {
	d := diagnosis.BasicDiagnoser{}

	result := validAdapterResult()
	result.Facts = nil

	diagnosisResult, err := d.Diagnose(result)
	if err != nil {
		t.Fatalf("Diagnose returned an error: %v", err)
	}

	if diagnosisResult.Code != "test.classification" {
		t.Fatalf("unexpected diagnosis code: got %q", diagnosisResult.Code)
	}

	if diagnosisResult.Severity != "info" {
		t.Fatalf("expected default severity %q, got %q",
			"info",
			diagnosisResult.Severity,
		)
	}
}

func TestBasicDiagnoser_IsDeterministic(t *testing.T) {
	d := diagnosis.BasicDiagnoser{}
	input := validAdapterResult()

	first, err := d.Diagnose(input)
	if err != nil {
		t.Fatalf("first Diagnose returned an error: %v", err)
	}

	second, err := d.Diagnose(input)
	if err != nil {
		t.Fatalf("second Diagnose returned an error: %v", err)
	}

	if first.Code != second.Code {
		t.Fatalf(
			"diagnosis code is not deterministic: first=%q second=%q",
			first.Code,
			second.Code,
		)
	}

	if first.Severity != second.Severity {
		t.Fatalf(
			"diagnosis severity is not deterministic: first=%q second=%q",
			first.Severity,
			second.Severity,
		)
	}
}
