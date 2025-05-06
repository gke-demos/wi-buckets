package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func main() {

	// setup simple web server to list contents of a bucket
	r := gin.Default()

	r.GET("/bucket/:bucket", func(c *gin.Context) {
		bucket := c.Params.ByName("bucket")

		ctx := context.Background()
		client, err := storage.NewClient(ctx)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"bucket": bucket, "error": err.Error()})

		} else {

			// create a random object
			bkt := client.Bucket(bucket)
			objName := fmt.Sprintf("%d", time.Now().UnixNano())
			obj := bkt.Object(objName)
			w := obj.NewWriter(ctx)
			if _, err := fmt.Fprintf(w, objName+"\n"); err != nil {
				c.JSON(http.StatusOK, gin.H{"bucket": bucket, "error": err.Error()})
			}

			if err := w.Close(); err != nil {
				c.JSON(http.StatusOK, gin.H{"bucket": bucket, "error": err.Error()})
			}

			// list and return all objects in the bucket
			query := &storage.Query{Prefix: ""}

			var names []string
			var errMsg string
			it := bkt.Objects(ctx, query)
			for {
				attrs, err := it.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					errMsg = err.Error()
					break
				}
				names = append(names, attrs.Name)
			}

			c.JSON(http.StatusOK, gin.H{"bucket": bucket, "objects": names, "error": errMsg})
		}
	})

	r.Run(":8080")
}
