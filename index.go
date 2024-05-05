package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/baidubce/app-builder/go/appbuilder"
)

func main() {
    os.Setenv("APPBUILDER_TOKEN", "bce-v3/ALTAK-oH2mqX2fCfis7BoQ54Dc3/439538993fe13b7e098f43a63e3d2f8894e2c661")
    os.Setenv("GATEWAY_URL_V2", "")
    config, err := appbuilder.NewSDKConfig("", "")
    if err != nil {
        fmt.Println("new config failed: ", err)
        return
    }
    appID := "71c36ed9-7749-40e2-9a00-1a9690726b31"
    builder, err := appbuilder.NewAppBuilderClient(appID, config)
    if err != nil {
        fmt.Println("new agent builder failed: ", err)
        return
    }

    http.HandleFunc("/send-message", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        var message struct {
            Message string `json:"message"`
        }
        if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        conversationID, err := builder.CreateConversation()
        if err != nil {
            fmt.Println("create conversation failed: ", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        i, err := builder.Run(conversationID, message.Message, nil, true)
        if err != nil {
            fmt.Println("run failed: ", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        var completedAnswer string
        for answer, err := i.Next(); err == nil; answer, err = i.Next() {
            completedAnswer += answer.Answer
        }

        w.Header().Set("Content-Type", "application/json")
        response := struct {
            Answer string `json:"answer"`
        }{completedAnswer}
        json.NewEncoder(w).Encode(response)
    })

    fmt.Println("Server listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}
