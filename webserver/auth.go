package webserver

import (
	"net"
	"net/http"
	"strings"

	"github.com/FlowKeeper/FlowServer/v2/config"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/httpResponse"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
)

//authorizationMiddleware should check the "ScraperUUID" header and determine if the client is allowed to send http requests to this agent
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matched := false

		for _, k := range config.Config.AllowedIPs {
			clientIP := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
			if k.Contains(clientIP) {
				matched = true
			}
		}

		if !matched {
			httpResponse.UserError(w, 401, "You are not allowed to access this server")
			logger.Warning(loggingArea, "Someone tried to access this server from outside the allowed IPs:"+r.RemoteAddr)
			return
		}

		next.ServeHTTP(w, r)

	})
}
