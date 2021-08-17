package appCoinTax

import "github.com/gin-gonic/gin"

func Router(r gin.IRouter) {
	r.POST("/app/coin/tax/query-list", AppCoinTaxQueryList)   
}
