package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/gitlab"
)

var oauthConf = &oauth2.Config{
    ClientID:     os.Getenv("GITLAB_CLIENT_ID"),
    ClientSecret: os.Getenv("GITLAB_CLIENT_SECRET"),
    RedirectURL:  os.Getenv("GITLAB_REDIRECT_URI"),
    Scopes:       []string{"read_user"},
    Endpoint:     gitlab.Endpoint,
}

func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/login", login)
    http.HandleFunc("/callback", callback)

    fmt.Println("Server started at http://0.0.0.0:8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `<a href="/login">Login with GitLab</a>`)
}

func login(w http.ResponseWriter, r *http.Request) {
    url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusFound)
}

func callback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "No code in request", http.StatusBadRequest)
        return
    }

    token, err := oauthConf.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to get token", http.StatusInternalServerError)
        return
    }

    client := oauthConf.Client(context.Background(), token)
    resp, err := client.Get("https://gitlab.com/api/v4/user")
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userInfo)
}
