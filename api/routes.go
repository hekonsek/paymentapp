package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hekonsek/paymentapp/payments"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func generateResponse(c *gin.Context, err error, body interface{}) {
	if err != nil {
		generateErrorResponse(c, err)
	} else {
		c.JSON(200, body)
	}
}

func generateErrorResponse(c *gin.Context, err error) {
	if err == payments.NoSuchElementErr {
		c.JSON(404, nil)
	} else if err != nil {
		c.JSON(500, nil)
		log.Debugf("Internal server error: %s", err.Error())
	}
}

func routes(apiServer *ApiServer, router *gin.Engine) {
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/payments", func(c *gin.Context) {
		paymentsList, err := apiServer.Store.List(0, 10)
		generateResponse(c, err, gin.H{
			"data": paymentsList,
		})
	})
	router.GET("/payments/count", func(c *gin.Context) {
		count, err := apiServer.Store.Count()
		generateResponse(c, err, gin.H{
			"count": count,
		})
	})
	router.POST("/payments", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		generateErrorResponse(c, err)
		var payment payments.Payment
		err = json.Unmarshal(body, &payment)
		generateErrorResponse(c, err)
		id, err := apiServer.Store.Create(&payment)
		generateResponse(c, err, gin.H{
			"id": id,
		})
	})
	router.PUT("/payments", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		generateErrorResponse(c, err)
		var payment payments.Payment
		err = json.Unmarshal(body, &payment)
		generateErrorResponse(c, err)
		err = apiServer.Store.Update(&payment)
		generateResponse(c, err, nil)
	})
	router.DELETE("/payment/:id", func(c *gin.Context) {
		err := apiServer.Store.Delete(c.Param("id"))
		generateResponse(c, err, nil)
	})
	router.GET("/payment/:id", func(c *gin.Context) {
		payment, err := apiServer.Store.FindById(c.Param("id"))
		generateResponse(c, err, payment)
	})
}
