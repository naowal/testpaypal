package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/logpacker/PayPal-Go-SDK"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello naowal!") // send data to client side
}

func success(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Payment success!") // send data to client side
}

func deny(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Payment Deny!") // send data to client side
}

func payment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	r.ParseForm()
	t, _ := template.ParseFiles("payment.gtpl")
	t.Execute(w, nil)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("redirect.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("r Amount:", r.Form["amount"])

		OpenPayment(strings.Join(r.Form["amount"], "")) //Call OpenPayment function

	}
}

func OpenPayment(balance string) {
	// In sandbox , Add my Own clientID and secretID
	c, err := paypalsdk.NewClient("AV2aAlW78rnvuU8EV92wLVsQTnusENXSJJCLSCxo6kUd0nU84ZWdjOoAkt1JNuPP7bk5t3jQ-Wky3on-",
		"EMqdZcKB21emXae8R97HuKzOOISIrnppX06ILsZetgTllfO9hakMr7MvLE548LkqPsKq-Khv7c1UFRMz", paypalsdk.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	accessToken, err := c.GetAccessToken()

	if err != nil {
		panic(err)
	}

	// Try to set DirectPaypalPayment
	amount := paypalsdk.Amount{
		Total:    balance, //parse form amount field
		Currency: "USD",
	}
	redirectURI := "/success"
	cancelURI := "/deny"
	description := "Leaptips following payment"
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

	// Just debug in console by printing
	fmt.Println()
	fmt.Println("Token: ", accessToken.Token)
	fmt.Println(paymentResult.Links[0].Rel)
	fmt.Println(paymentResult.Links[0].Href)
	fmt.Println(paymentResult.Links[0].Method)
	fmt.Println(paymentResult.Links[0].Enctype)
	fmt.Println()
	fmt.Println(paymentResult.Links[1].Rel)
	fmt.Println(paymentResult.Links[1].Href)
	fmt.Println(paymentResult.Links[1].Method)
	fmt.Println(paymentResult.Links[1].Enctype)
	fmt.Println()
	fmt.Println(paymentResult.Links[2].Rel)
	fmt.Println(paymentResult.Links[2].Href)
	fmt.Println(paymentResult.Links[2].Method)
	fmt.Println(paymentResult.Links[2].Enctype)

	// open approvel url -> paypal for payment
	exec.Command("xdg-open", paymentResult.Links[1].Href).Run()

	// So, Next ?? How am I will do. with accessToken and paymentResult
}

func main() {
	http.HandleFunc("/", sayhelloName) // set router
	http.HandleFunc("/login", login)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/success", success)
	http.HandleFunc("/deny", deny)
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
