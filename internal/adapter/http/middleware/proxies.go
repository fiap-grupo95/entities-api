package middleware

import "github.com/gin-gonic/gin"

var proxies = []string{
	"0.0.0.0/0",
}

func SetTrustedProxies(eng *gin.Engine) {
	// Define os proxies confiáveis
	// Use []string{"0.0.0.0/0"} para confiar em todos (não recomendado em produção)
	// Ou especifique o IP/cidr exato do proxy (ex: "192.168.1.1" ou "10.0.0.0/8")
	err := eng.SetTrustedProxies(proxies)
	if err != nil {
		panic(err)
	}
}
