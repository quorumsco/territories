// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://quorumsco.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestFunctionScoreQuery(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		Add(NewTermFilter("name.last", "banon"), NewFactorFunction().BoostFactor(3)).
		AddScoreFunc(NewFactorFunction().BoostFactor(3)).
		AddScoreFunc(NewFactorFunction().BoostFactor(3)).
		Boost(3).
		MaxBoost(10).
		ScoreMode("avg")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":3,"functions":[{"boost_factor":3,"filter":{"term":{"name.last":"banon"}}},{"boost_factor":3},{"boost_factor":3}],"max_boost":10,"query":{"term":{"name.last":"banon"}},"score_mode":"avg"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithNilFilter(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("tag", "wow")).
		AddScoreFunc(NewRandomFunction()).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","max_boost":12,"query":{"term":{"tag":"wow"}},"random_score":{},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactor(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income")).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"},"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactorWithWeight(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income").Weight(2.5)).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"},"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max","weight":2.5}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactorWithMultipleScoreFuncsAndWeights(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income").Weight(2.5)).
		AddScoreFunc(NewScriptFunction("_score * doc['my_numeric_field'].value").Weight(1.25)).
		AddScoreFunc(NewWeightFactorFunction(0.5)).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","functions":[{"field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"},"weight":2.5},{"script_score":{"script":"_score * doc['my_numeric_field'].value"},"weight":1.25},{"weight":0.5}],"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithGaussScoreFunc(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewGaussDecayFunction().FieldName("pin.location").Origin("11, 12").Scale("2km").Offset("0km").Decay(0.33))
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"gauss":{"pin.location":{"decay":0.33,"offset":"0km","origin":"11, 12","scale":"2km"}},"query":{"term":{"name.last":"banon"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithGaussScoreFuncAndMultiValueMode(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewGaussDecayFunction().FieldName("pin.location").Origin("11, 12").Scale("2km").Offset("0km").Decay(0.33).MultiValueMode("avg"))
	data, err := json.Marshal(q.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"gauss":{"multi_value_mode":"avg","pin.location":{"decay":0.33,"offset":"0km","origin":"11, 12","scale":"2km"}},"query":{"term":{"name.last":"banon"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
