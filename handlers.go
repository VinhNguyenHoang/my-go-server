package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*EmailServer) Main() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "HTTP Server by Gin",
		})
	}
}

// this handler only handle request from SendGrid
func (s *EmailServer) Webhook() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		jsonData, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Fatalf("error get body: %+v", err)
			return
		}
		err = HandleSGAuthentication(ctx.Request.Header, jsonData)
		if err != nil {
			log.Println(fmt.Sprintf("failed handle SG Authentication: %+v", err))
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		s.Repo.InsertMany(s.DbConn, string(jsonData))
		log.Println(jsonData)
		ctx.IndentedJSON(http.StatusOK, string(jsonData))
	}
}

func (s *EmailServer) WebhookTest() func(c *gin.Context) {
	return func(c *gin.Context) {
		s.SendGridClient.TriggerWebhookTest()
	}
}

func (s *EmailServer) SendEmail() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		toAddress := ctx.Query("to_address")
		if toAddress == "" {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "address is empty",
			})
			return
		}
		resp, err := s.SendGridClient.SendEmail(toAddress)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": resp,
		})
	}
}
