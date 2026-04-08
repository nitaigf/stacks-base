package services

type RequestMetadata struct {
	ActorUserID *string
	ActorName   string
	ActorEmail  string
	ActorRole   string
	Route       string
	Method      string
	IPAddress   string
	UserAgent   string
}
