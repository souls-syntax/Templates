package models

type Query struct {
	Query			string	`json:"query"`
	UserID		string	`json:"userid"`
}

type Response struct {
	Response		string
	Status			bool
}

type Decision struct {
	QueryHash		string
	QueryText		string

	Verdict			string
	Confidence	float64
	Decider			string
}

type Observation struct {
	Source				 		string
	ProcessingTimeMs	int64
}

type VerifyResponse struct {
	QueryHash 					string		`json:"query_hash"`
	QueryText						string		`json:"query_text"`

	DecisionVerdict			string		`json:"decision_verdict"`
	DecisionConfidence	string		`json:"decision_confidence"`
	DecisionDecider			string		`json:"decison_decider"`

	ObsSource						string		`json:"obs_source"`
	ObsProcessingTimeMs	string		`json:"obs_processing_time"`
}
