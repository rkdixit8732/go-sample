package utilities

type HealthCheckResponse struct {
	MongoDBStatus string `json:"mongodb_status"`
	RedisStatus   string `json:"redis_status"`
	AppStatus     string `json:"app_status"`
}

type ReadinessCheckResponse struct {
	MongoDBStatus string `json:"mongodb_status"`
	RedisStatus   string `json:"redis_status"`
	AppStatus     string `json:"app_status"`
}

type CreatePlayerRequest struct {
	Name    string `json:"name"`
	Credits int    `json:"credits"`
	Status  string `bson:"status"`
}

type PlayerResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Credits int    `json:"credits"`
	Status  string `json:"status"`
}
