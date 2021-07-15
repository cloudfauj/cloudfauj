package infrastructure

type AppInfra struct {
	App               string
	EcsTaskDefinition string
	TargetGroup       string
	AlbListenerRule   string
	DNSRecord         string
	ECSService        string
}
