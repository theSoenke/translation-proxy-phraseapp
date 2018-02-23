package api

import (
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/phrase/phraseapp-go/phraseapp"
)

// Run translation proxy API
func Run(client *phraseapp.Client) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	if err != nil {
		log.Fatal(err)
	}

	l := locales{
		Client: client,
		Cache:  cache,
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Translation Proxy is running")
	})

	api := router.Group("/api/v2")
	{
		api.GET("/projects/:project_id/locales/:id/download", l.downloadLocale)
	}

	router.Run(":8080")
}

func (l *locales) downloadLocale(c *gin.Context) {
	projectID := c.Param("project_id")
	localeID := c.Param("id")

	var params downloadParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locale, err := l.getLocale(projectID, localeID, &params)
	if err != nil {
		log.Printf("error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, string(locale))
}
