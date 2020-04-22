package http

import "github.com/gin-gonic/gin"

func (c *Controller) getCumulative() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := ctx.Param("state")
		result, err := c.entriesService.GetCumulative(state)
		if err != nil {
			ctx.JSON(403, gin.H{
				"status": "error",
				"err":    err})
		} else {
			ctx.JSON(200, gin.H{
				"status": "ok",
				"result": result})
		}
	}
}

func (c *Controller) getDeltas() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := ctx.Param("state")
		result, err := c.entriesService.GetDeltas(state)
		if err != nil {
			ctx.JSON(403, gin.H{
				"status": "error",
				"err":    err})
		} else {
			ctx.JSON(200, gin.H{
				"status": "ok",
				"result": result})
		}
	}
}
