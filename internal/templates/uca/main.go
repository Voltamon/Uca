package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"{{APP_NAME}}/uca"
)

func main() {
	app := pocketbase.New()

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		err := e.Next()
		if err != nil {
			return err
		}
		uca.RunMigrations(app)
		return nil
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
    	for serviceName, methods := range uca.Registry {
        	for method, handler := range methods {
            	path := "/api/" + serviceName
             	switch method {
             		case "GET":
                		se.Router.GET(path, handler)
             		case "POST":
                		se.Router.POST(path, handler)
             		case "PUT":
                		se.Router.PUT(path, handler)
             		case "DELETE":
                		se.Router.DELETE(path, handler)
             	}
        	}
    	}

    	for agentName, agentPort := range uca.Agents {
        	name := agentName
        	port := agentPort
        	se.Router.GET("/api/chat/"+name, func(e *core.RequestEvent) error {
            	message := e.Request.URL.Query().Get("message")
            	agentURL := "http://127.0.0.1:" + port + "/chat?message=" + url.QueryEscape(message)

            	resp, err := http.Get(agentURL)
            	if err != nil {
                	return e.JSON(http.StatusInternalServerError, map[string]string{"error": "agent unavailable"})
             	}
              	defer resp.Body.Close()

              	e.Response.Header().Set("Content-Type", "text/event-stream")
              	e.Response.Header().Set("Cache-Control", "no-cache")
              	e.Response.Header().Set("Connection", "keep-alive")

               	flusher, ok := e.Response.(http.Flusher)
               	if !ok {
               		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "streaming not supported"})
               	}

               	buf := make([]byte, 64)
               	for {
               		n, err := resp.Body.Read(buf)
               		if n > 0 {
               			e.Response.Write(buf[:n])
               			flusher.Flush()
               		}
               		if err == io.EOF {
               			break
               		}
               		if err != nil {
               			break
               		}
               	}

               	return nil
         	})
     	}

      	if _, err := os.Stat("dist"); err == nil {
        	se.Router.GET("/{path...}", func(e *core.RequestEvent) error {
            	path := e.Request.PathValue("path")
            	filePath := filepath.Join("dist", path)

            	if _, err := os.Stat(filePath); os.IsNotExist(err) {
                	filePath = filepath.Join("dist", "index.html")
            	}

            	http.ServeFile(e.Response, e.Request, filePath)
            	return nil
        	})
    	}

    	fmt.Println("Uca dev server running on http://localhost:8090")
    	return se.Next()
	})

	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	if err := apis.Serve(app, apis.ServeConfig{
		HttpAddr:        "127.0.0.1:8090",
		ShowStartBanner: false,
	}); err != nil {
		log.Fatal(err)
	}
}
