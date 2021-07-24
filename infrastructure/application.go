package infrastructure

type AppInfra struct {
	App               string `json:"app"`
	EcsTaskDefinition string `json:"ecs_task_definition"`
	TargetGroup       string `json:"target_group"`
	AlbListenerRule   string `json:"alb_listener_rule"`
	DNSRecord         string `json:"dns_record"`
	ECSService        string `json:"ecs_service"`
}
