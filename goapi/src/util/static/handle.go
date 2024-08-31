package static

type Path string

const (
	HEALTH_PATH   Path = "/health"
	DOMAIN_PATH   Path = "/domain/v1"
	USER_ID_PATH  Path = "/user/{userId}"
	QUEUE_ID_PATH Path = "/queue/{queueId}"
	CREATE_PATH   Path = "/create"
	STATUS_PATH   Path = "/status"
	USER_ID       Path = "userId"
	QUEUE_ID      Path = "queueId"
)
