package main

import (
	"fmt"
	"github.com/3scale/3scale-istio-adapter/pkg/threescale"
	"github.com/Aestek/haproxy-connect/spoe"
	"net/url"
	"threescale_haproxy/pkg/threescale_authorizer"
	"time"
)

func main() {

	CacheTTL := 3 * time.Minute
	CacheRefreshInterval := 2 * time.Minute
	CacheUpdateRetries := 2
	CacheEntriesMax := 500

	proxyCache := threescale.NewProxyConfigCache(CacheTTL, CacheRefreshInterval, CacheUpdateRetries, CacheEntriesMax)
	err := proxyCache.StartRefreshWorker()
	if err != nil {
		panic(err)
	}

	authorizer := threescale_authorizer.NewAuthorizer(proxyCache)
	var authorized int
	agent := spoe.New(func(messages []spoe.Message) ([]spoe.Action, error) {
		for _, msg := range messages {
			if msg.Name != "authrep" {
				continue
			} else {
				authorized = 0
				values, _ := url.ParseQuery(fmt.Sprintf("%s", msg.Args["query"]))
				request := threescale_authorizer.AuthorizeRequest{
					Host:        fmt.Sprintf("%s", msg.Args["host"]),
					ServiceId:   fmt.Sprintf("%s", msg.Args["service_id"]),
					SystemUrl:   fmt.Sprintf("%s", msg.Args["system_url"]),
					AccessToken: fmt.Sprintf("%s", msg.Args["access_token"]),
					Path:        fmt.Sprintf("%s", msg.Args["path"]),
					Method:      fmt.Sprintf("%s", msg.Args["method"]),
					AppID:       values.Get("app_id"),
					AppKey:      values.Get("app_key"),
					UserKey:     values.Get("user_key"),
				}
				if authorizer.AuthRep(request) {
					authorized = 1
				}
			}
		}

		return []spoe.Action{
			spoe.ActionSetVar{
				Name:  "authorized",
				Scope: spoe.VarScopeTransaction,
				Value: authorized,
			},
		}, nil
	})

	panic(agent.ListenAndServe(":12345"))
}
