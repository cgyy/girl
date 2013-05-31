package girl

import (
	"encoding/json"
	"errors"
    "fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
)

type Invoker func(c *Context) View

type Context struct {
	Request     *http.Request
	params      map[string]string
	vars        map[string]string
	data        map[string]interface{}
	application *Girl
	http.ResponseWriter
}

type Girl struct {
	routes      map[string][]Route
	templateSet *template.Template
	RootDir     string
}

type Route struct {
	pattern string
	invoker Invoker
}

type View interface {
	Apply(w http.ResponseWriter, r *http.Request) error
}

type TemplateView struct {
	Template *template.Template
	Data     interface{}
	err      error
}

type JSONView struct {
	obj interface{}
}

type TextView struct {
	text string
}

func (t TemplateView) Apply(w http.ResponseWriter, r *http.Request) (err error) {
	if t.err != nil {
		err = t.err
	} else {
		err = t.Template.Execute(w, t.Data)
	}
	return
}

func (t JSONView) Apply(w http.ResponseWriter, r *http.Request) (err error) {
	var b []byte
	b, err = json.Marshal(t.obj)
	if err != nil {
		return
	}

	w.Write(b)
	return
}

func (t TextView) Apply(w http.ResponseWriter, r *http.Request) (err error) {
	w.Write([]byte(t.text))
	return nil
}

func (c *Context) Render(s string, data interface{}) View {
	tmpl := c.application.templateSet.Lookup(s)
	var err error
	if tmpl == nil {
		err = errors.New("No such file exists: " + c.application.RootDir + "/" + s + ".tpl")
		log.Println(err)
	}
	return TemplateView{tmpl, data, err}
}

func (c *Context) RenderJSON(obj interface{}) View {
    c.ResponseWriter.Header().Set("Content-Type", "application/json")
	return JSONView{obj}
}

func (c *Context) RenderText(text string) View {
	return TextView{text}
}

func (c *Context) Abort(status int, body string) View {
	c.ResponseWriter.WriteHeader(status)
	c.ResponseWriter.Write([]byte(body))
	return nil
}

// Temporarily Moved
func (c *Context) Redirect(url string) View {
	http.Redirect(c.ResponseWriter, c.Request, url, 302)
	return nil
}

func (c *Context) GetParam(s string) string {
    if c.params == nil {
        c.params = make(map[string]string)
        c.Request.ParseForm()
        c.Request.ParseMultipartForm(2 * 1024 * 1024)

        for k, v := range c.Request.Form {
            fmt.Println(k)
            c.params[k] = v[0]
        }

        if c.vars != nil {
            for k, v := range c.vars {
                c.params[k] = v
            }
        }

    }
    return c.params[s]
}

func (c *Context) GetNumParam(s string) int {
    v := c.GetParam(s)
    num, _ := strconv.Atoi(v)
    return num
}

func (c *Context) invoke(routeHander Invoker) (err error) {
    defer func() {
        if e := recover(); e != nil {
            err = errors.New(e.(string))
        }
    }()

    view := routeHander(c)
    if view != nil {
        err = view.Apply(c.ResponseWriter, c.Request)
    }
    return
}

func (g *Girl) Get(pattern string, invoker Invoker) {
    g.addRoute(pattern, "GET", invoker)
    g.addRoute(pattern, "HEAD", invoker)
}

func (g *Girl) Post(pattern string, invoker Invoker) {
    g.addRoute(pattern, "POST", invoker)
}

func (g *Girl) Put(pattern string, invoker Invoker) {
    g.addRoute(pattern, "PUT", invoker)
}

func (g *Girl) Delete(pattern string, invoker Invoker) {
    g.addRoute(pattern, "DELETE", invoker)
}

func (g *Girl) addRoute(pattern, method string, invoker Invoker) {
    signature := g.routes[method]
    g.routes[method] = append(signature, Route{pattern, invoker})
}

func matchUrl(scheme, path string) (bool, map[string]string) {
    schemeList := strings.Split(scheme, "/")
    pathList := strings.Split(path, "/")

    if len(schemeList) != len(pathList) {
        return false, nil
    }

    var vars map[string]string

    for i, s := range schemeList {
        p := pathList[i]
        if s != p {
            if len(s) > 1 && s[0] == ':' {
                if vars == nil {
                    vars = make(map[string]string)
                }
                vars[s[1:]] = p
            } else {
                return false, nil
            }
        }
    }

    return true, vars
}

func (g *Girl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rPath := r.URL.Path

    if rPath == "/favicon.ico" {
        http.ServeFile(w, r, g.RootDir+"/public"+rPath)
        return
    }

    if strings.HasPrefix(rPath, "/public/") {
        http.ServeFile(w, r, g.RootDir+rPath)
        return
    }

    log.Println(rPath)
    context := Context{
        Request:        r,
        params:         nil,
        vars:           nil,
        data:           nil,
        application:    g,
        ResponseWriter: w,
    }
    for _, router := range g.routes[r.Method] {
        scheme := router.pattern

        isMatch, vars := matchUrl(scheme, rPath)
        if isMatch {
            context.vars = vars
            err := context.invoke(router.invoker)

            if err != nil {
                stack := string(debug.Stack())
                log.Println(stack)
                context.Abort(http.StatusInternalServerError,
                "internal server error\n"+err.Error()+"\n"+stack)
            }
            return
        }
    }

    context.Abort(http.StatusNotFound, "not found")
}

func (g *Girl) Run(addr string) {
    log.Println("=> Girl server start at port: 9999")
    http.ListenAndServe(addr, g)
}

func (g *Girl) initTemplate() error {
    baseDir := path.Join(g.RootDir, "views")
    g.templateSet = template.New(baseDir)

    if !dirExists(baseDir) {
        log.Println("folder\"", baseDir, "\" not exists!\n")
        return nil
    }

    filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() && !strings.HasPrefix(filepath.Base(path), ".") {

            b, err := ioutil.ReadFile(path)
            if err != nil {
                return err
            }
            s := string(b)

            tmplName := path[len(baseDir)+1:]
            tmplName = strings.Split(tmplName, ".")[0]
            tmpl := g.templateSet.New(tmplName)

            _, err = tmpl.Parse(s)

            if err != nil {
                return err
            }

        }
        return nil
    })

    return nil
}

func rootDir() (root string, err error) {
    pwd, err := os.Getwd()
    if err != nil {
        log.Println("os.Getwd err:", err)
        return pwd, err
    }
    exePath := os.Args[0]
    parent, _ := path.Split(exePath)

    var exeDir string
    if filepath.IsAbs(exeDir) {
        exeDir = parent
    } else {
        exeDir = path.Join(pwd, parent)
    }

    if dirExists(path.Join(exeDir, "views")) && dirExists(path.Join(exeDir, "public")) {
        root = exeDir
    } else {
        root = pwd
    }

    return

}

func dirExists(dir string) bool {
    fi, err := os.Stat(dir)
    if err != nil {
        return false
    }
    return fi.IsDir()
}

func New() *Girl {
    root, _ := rootDir()
    log.Println("set project root dir:", root)

    routes := make(map[string][]Route)

    girl := Girl{routes, nil, root}
    err := girl.initTemplate()

    if err != nil {
        panic(err)
    }

    return &girl
}
