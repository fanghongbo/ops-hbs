package http

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/cache"
	"github.com/fanghongbo/ops-hbs/common/g"
	"github.com/fanghongbo/ops-hbs/common/model"
	"net/http"
	"os"
	"time"
)

func init() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    "ok",
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    g.VersionInfo(),
		})
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var err error

			if err = g.ReloadConfig(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     err.Error(),
					"data":    nil,
				})
			} else {
				RenderJson(w, map[string]interface{}{
					"success": true,
					"msg":     "reload success",
					"data":    nil,
				})
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "exited success",
				"data":    nil,
			})
			go func() {
				time.Sleep(time.Second)
				dlog.Warning("exited..")
				os.Exit(0)
			}()
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var result []model.Expression

			result = cache.ExpressionCache.Get()
			RenderJson(w, result)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var result []string

			result = cache.AgentsCache.Keys()
			RenderJson(w, result)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var data map[string]model.Host

			data = make(map[string]model.Host, len(cache.MonitorHostCache.Get()))
			for k, v := range cache.MonitorHostCache.Get() {
				data[fmt.Sprint(k)] = v
			}
			RenderJson(w, data)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var data map[string]model.Template

			data = make(map[string]model.Template, len(cache.TemplateCache.GetMap()))
			for k, v := range cache.TemplateCache.GetMap() {
				data[fmt.Sprint(k)] = v
			}
			RenderJson(w, data)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var (
				result   []string
				hostname string
				err      error
			)

			result = []string{}
			if err = r.ParseForm(); err != nil {
				dlog.Errorf("parse url err: %s", err)
				RenderJson(w, result)
				return
			}

			values := r.Form["hostname"]
			if len(values) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     "miss url params",
					"data":    nil,
				})
				return
			}

			hostname = values[0]
			result = cache.GetPlugins(hostname)
			RenderJson(w, result)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var data map[string]model.Strategy

			data = make(map[string]model.Strategy, len(cache.StrategiesCache.GetMap()))
			for k, v := range cache.StrategiesCache.GetMap() {
				data[fmt.Sprint(k)] = v
			}
			RenderJson(w, data)
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})
}

