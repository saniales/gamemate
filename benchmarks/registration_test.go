package benchmarks

import (
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"testing"
)

func BenchmarkDevRegistrationRequestScalability(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	const (
		API_TOKEN     string = "5D170C3D25E269CE8FC98BEACBAB944F49125AD068EE239D6117A41EBF58B8904EF5F0F4BF747DCC9D46C033FF58E15B2CD0D352CDAACDAEA9FD942E891A88ED"
		URL           string = "http://gamemate.di.unito.it:8080/dev/register"
		USER_TEMPLATE string = "benchmark_user_"
	)
	form := url.Values{
		"Type":      {"DevRegistration"},
		"API_Token": {API_TOKEN},
		"Password":  {"benchmark"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.Add("Email", USER_TEMPLATE+strconv.Itoa(i))
		http.PostForm(URL, form)
		form.Del("Email")
	}
}
