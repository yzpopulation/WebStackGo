package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type JsLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

type JsConfig struct {
	Title          string `json:"title"`
	Url            string `json:"url"`
	Port           int    `json:"port"`
	Keywords       string `json:"keywords"`
	Description    string `json:"description"`
	IsFbOpenGraph  bool   `json:"isFbOpenGraph"`
	IsTwitterCards bool   `json:"isTwitterCards"`
	Recordcode     string `json:"recordcode"`
	Footer         string `json:"footer"`
}

type JsMenu struct {
	Menu string   `json:"menu"`
	Name string   `json:"name"`
	Icon string   `json:"icon"`
	Sub  []JsMenu `json:"sub"`
	Url  string   `json:"url"`
}

type JsClassItem struct {
	Url  string `json:"url"`
	Img  string `json:"img"`
	Name string `json:"name"`
	Mark string `json:"mark"`
}

type JsClass struct {
	Name string        `json:"name"`
	Rows []JsClassItem `json:"rows"`
}

type JsWebStack struct {
	Menu  []JsMenu  `json:"Menu"`
	Class []JsClass `json:"Class"`
}

var (
	Login    JsLogin
	Config   JsConfig
	WebStack JsWebStack
)

func main() {
	var err error
	//如果没有json文件目录就创建一个
	if _, err := os.Stat("./json/"); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll("./json/", os.ModePerm)
		}
	}
	if !checkFileIsExist("./json/login.json") {
		var d1 = []byte("{\"username\":\"admin\",\"password\":\"21232f297a57a5a743894a0e4a801fc3\",\"path\":\"/login\"}")
		err = ioutil.WriteFile("./json/login.json", d1, 0666)
		if err != nil {
			fmt.Print("初始化登陆文件login.json错误：")
			fmt.Println(err)
			return
		}
	}
	if !checkFileIsExist("./json/config.json") {
		var d1 = []byte("{\"title\":\"WebStackGo - 轻量级网址导航\",\"url\":\"0.0.0.0\",\"port\":80,\"keywords\":\"UI设计,UI设计素材,设计导航,网址导航,设计资源,创意导航,创意网站导航,设计师网址大全,设计素材大全,设计师导航,UI设计资源,优秀UI设计欣赏,设计师导航,设计师网址大全,设计师网址导航,产品经理网址导航,交互设计师网址导航,WebStackGo\",\"description\":\"WebStackGo - 收集国内外优秀设计网站、UI设计资源网站、灵感创意网站、素材资源网站，定时更新分享优质产品设计书签。\",\"isFbOpenGraph\":true,\"isTwitterCards\":true,\"recordcode\":\"备案号\",\"footer\":\"\\u0026copy; 2017-2021 \\u003ca href=\\\"/about.html\\\"\\u003e \\u003cstrong\\u003e轻量级网址导航\\u003c/strong\\u003e \\u003c/a\\u003e design by \\u003ca href=\\\"https://github.com/WebStackPage/WebStackPage.github.io\\\" target=\\\"_blank\\\"\\u003e \\u003cstrong\\u003eWebStackPage\\u003c/strong\\u003e \\u003c/a\\u003e \\u0026 \\u003ca href=\\\"https://github.com/xiaoxinpro/WebStackGo\\\" target=\\\"_blank\\\"\\u003e \\u003cstrong\\u003eWebStackGo\\u003c/strong\\u003e \\u003c/a\\u003e\"}")
		err = ioutil.WriteFile("./json/config.json", d1, 0666)
		if err != nil {
			fmt.Print("初始化登陆文件config.json错误：")
			fmt.Println(err)
			return
		}
	}
	if !checkFileIsExist("./json/webstack.json") {
		var d1 = []byte("{\"Menu\": [{\"menu\": \"smooth\", \"name\": \"\\u5e38\\u7528\\u63a8\\u8350\", \"icon\": \"linecons-star\", \"sub\": [], \"url\": \"#\\u5e38\\u7528\\u63a8\\u8350\"}, {\"menu\": \"smooth\", \"name\": \"\\u793e\\u533a\\u8d44\\u8baf\", \"icon\": \"linecons-doc\", \"sub\": [], \"url\": \"#\\u793e\\u533a\\u8d44\\u8baf\"}, {\"menu\": \"smooth\", \"name\": \"\\u7075\\u611f\\u91c7\\u96c6\", \"icon\": \"linecons-lightbulb\", \"sub\": [{\"menu\": \"smooth\", \"name\": \"\\u53d1\\u73b0\\u4ea7\\u54c1\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u53d1\\u73b0\\u4ea7\\u54c1\"}, {\"menu\": \"smooth\", \"name\": \"\\u754c\\u9762\\u7075\\u611f\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u754c\\u9762\\u7075\\u611f\"}, {\"menu\": \"smooth\", \"name\": \"\\u7f51\\u9875\\u7075\\u611f\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u7f51\\u9875\\u7075\\u611f\"}], \"url\": \"#\\u7075\\u611f\\u91c7\\u96c6\"}, {\"menu\": \"smooth\", \"name\": \"\\u7d20\\u6750\\u8d44\\u6e90\", \"icon\": \"linecons-thumbs-up\", \"sub\": [{\"menu\": \"smooth\", \"name\": \"\\u56fe\\u6807\\u7d20\\u6750\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u56fe\\u6807\\u7d20\\u6750\"}, {\"menu\": \"smooth\", \"name\": \"LOGO\\u8bbe\\u8ba1\", \"icon\": \"\", \"sub\": [], \"url\": \"#LOGO\\u8bbe\\u8ba1\"}, {\"menu\": \"smooth\", \"name\": \"\\u5e73\\u9762\\u7d20\\u6750\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u5e73\\u9762\\u7d20\\u6750\"}, {\"menu\": \"smooth\", \"name\": \"UI\\u8d44\\u6e90\", \"icon\": \"\", \"sub\": [], \"url\": \"#UI\\u8d44\\u6e90\"}, {\"menu\": \"smooth\", \"name\": \"Sketch\\u8d44\\u6e90\", \"icon\": \"\", \"sub\": [], \"url\": \"#Sketch\\u8d44\\u6e90\"}, {\"menu\": \"smooth\", \"name\": \"\\u5b57\\u4f53\\u8d44\\u6e90\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u5b57\\u4f53\\u8d44\\u6e90\"}, {\"menu\": \"smooth\", \"name\": \"Mockup\", \"icon\": \"\", \"sub\": [], \"url\": \"#Mockup\"}, {\"menu\": \"smooth\", \"name\": \"\\u6444\\u5f71\\u56fe\\u5e93\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u6444\\u5f71\\u56fe\\u5e93\"}, {\"menu\": \"smooth\", \"name\": \"PPT\\u8d44\\u6e90\", \"icon\": \"\", \"sub\": [], \"url\": \"#PPT\\u8d44\\u6e90\"}], \"url\": \"#\\u7d20\\u6750\\u8d44\\u6e90\"}, {\"menu\": \"smooth\", \"name\": \"\\u5e38\\u7528\\u5de5\\u5177\", \"icon\": \"linecons-diamond\", \"sub\": [{\"menu\": \"smooth\", \"name\": \"\\u56fe\\u5f62\\u521b\\u610f\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u56fe\\u5f62\\u521b\\u610f\"}, {\"menu\": \"smooth\", \"name\": \"\\u754c\\u9762\\u8bbe\\u8ba1\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u754c\\u9762\\u8bbe\\u8ba1\"}, {\"menu\": \"smooth\", \"name\": \"\\u4ea4\\u4e92\\u52a8\\u6548\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u4ea4\\u4e92\\u52a8\\u6548\"}, {\"menu\": \"smooth\", \"name\": \"\\u5728\\u7ebf\\u914d\\u8272\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u5728\\u7ebf\\u914d\\u8272\"}, {\"menu\": \"smooth\", \"name\": \"\\u5728\\u7ebf\\u5de5\\u5177\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u5728\\u7ebf\\u5de5\\u5177\"}, {\"menu\": \"smooth\", \"name\": \"Chrome\\u63d2\\u4ef6\", \"icon\": \"\", \"sub\": [], \"url\": \"#Chrome\\u63d2\\u4ef6\"}], \"url\": \"#\\u5e38\\u7528\\u5de5\\u5177\"}, {\"menu\": \"smooth\", \"name\": \"\\u5b66\\u4e60\\u6559\\u7a0b\", \"icon\": \"linecons-pencil\", \"sub\": [{\"menu\": \"smooth\", \"name\": \"\\u8bbe\\u8ba1\\u89c4\\u8303\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u8bbe\\u8ba1\\u89c4\\u8303\"}, {\"menu\": \"smooth\", \"name\": \"\\u89c6\\u9891\\u6559\\u7a0b\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u89c6\\u9891\\u6559\\u7a0b\"}, {\"menu\": \"smooth\", \"name\": \"\\u8bbe\\u8ba1\\u6587\\u7ae0\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u8bbe\\u8ba1\\u6587\\u7ae0\"}, {\"menu\": \"smooth\", \"name\": \"\\u8bbe\\u8ba1\\u7535\\u53f0\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u8bbe\\u8ba1\\u7535\\u53f0\"}, {\"menu\": \"smooth\", \"name\": \"\\u4ea4\\u4e92\\u8bbe\\u8ba1\", \"icon\": \"\", \"sub\": [], \"url\": \"#\\u4ea4\\u4e92\\u8bbe\\u8ba1\"}], \"url\": \"#\\u5b66\\u4e60\\u6559\\u7a0b\"}, {\"menu\": \"smooth\", \"name\": \"UED\\u56e2\\u961f\", \"icon\": \"linecons-user\", \"sub\": [], \"url\": \"#UED\\u56e2\\u961f\"}, {\"menu\": \"url\", \"name\": \"\\u5173\\u4e8e\\u672c\\u7ad9\", \"icon\": \"linecons-heart\", \"sub\": [], \"url\": \"about.html\"}], \"Class\": [{\"name\": \"\\u5e38\\u7528\\u63a8\\u8350\", \"rows\": [{\"url\": \"https://dribbble.com/\", \"img\": \"../assets/images/logos/dribbble.png\", \"name\": \"Dribbble\", \"mark\": \"\\u5168\\u7403UI\\u8bbe\\u8ba1\\u5e08\\u4f5c\\u54c1\\u5206\\u4eab\\u5e73\\u53f0\\u3002\"}, {\"url\": \"https://behance.net/\", \"img\": \"../assets/images/logos/behance.png\", \"name\": \"Behance\", \"mark\": \"Adobe\\u65d7\\u4e0b\\u7684\\u8bbe\\u8ba1\\u5e08\\u4ea4\\u6d41\\u5e73\\u53f0\\uff0c\\u6765\\u81ea\\u4e16\\u754c\\u5404\\u5730\\u7684\\u8bbe\\u8ba1\\u5e08\\u5728\\u8fd9\\u91cc\\u5206\\u4eab\\u81ea\\u5df1\\u7684\\u4f5c\\u54c1\\u3002\"}, {\"url\": \"http://www.ui.cn/\", \"img\": \"../assets/images/logos/uicn.png\", \"name\": \"UI\\u4e2d\\u56fd\", \"mark\": \"\\u56fe\\u5f62\\u4ea4\\u4e92\\u4e0e\\u754c\\u9762\\u8bbe\\u8ba1\\u4ea4\\u6d41\\u3001\\u4f5c\\u54c1\\u5c55\\u793a\\u3001\\u5b66\\u4e60\\u5e73\\u53f0\\u3002\"}, {\"url\": \"http://www.zcool.com.cn/\", \"img\": \"../assets/images/logos/zcool.png\", \"name\": \"\\u7ad9\\u9177\", \"mark\": \"\\u4e2d\\u56fd\\u4eba\\u6c14\\u8bbe\\u8ba1\\u5e08\\u4e92\\u52a8\\u5e73\\u53f0\"}, {\"url\": \"https://www.pinterest.com/\", \"img\": \"../assets/images/logos/pinterest.png\", \"name\": \"Pinterest\", \"mark\": \"\\u5168\\u7403\\u7f8e\\u56fe\\u6536\\u85cf\\u91c7\\u96c6\\u7ad9\"}, {\"url\": \"http://huaban.com/\", \"img\": \"../assets/images/logos/huaban.png\", \"name\": \"\\u82b1\\u74e3\", \"mark\": \"\\u6536\\u96c6\\u7075\\u611f,\\u4fdd\\u5b58\\u6709\\u7528\\u7684\\u7d20\\u6750\"}, {\"url\": \"https://medium.com/\", \"img\": \"../assets/images/logos/medium.png\", \"name\": \"Medium\", \"mark\": \"\\u9ad8\\u8d28\\u91cf\\u8bbe\\u8ba1\\u6587\\u7ae0\"}, {\"url\": \"http://www.uisdc.com/\", \"img\": \"../assets/images/logos/uisdc.png\", \"name\": \"\\u4f18\\u8bbe\", \"mark\": \"\\u8bbe\\u8ba1\\u5e08\\u4ea4\\u6d41\\u5b66\\u4e60\\u5e73\\u53f0\"}, {\"url\": \"https://www.producthunt.com/\", \"img\": \"../assets/images/logos/producthunt.png\", \"name\": \"Producthunt\", \"mark\": \"\\u53d1\\u73b0\\u65b0\\u9c9c\\u6709\\u8da3\\u7684\\u4ea7\\u54c1\"}, {\"url\": \"https://www.youtube.com/\", \"img\": \"../assets/images/logos/youtube.png\", \"name\": \"Youtube\", \"mark\": \"\\u5168\\u7403\\u6700\\u5927\\u7684\\u5b66\\u4e60\\u5206\\u4eab\\u5e73\\u53f0\"}, {\"url\": \"https://www.google.com/\", \"img\": \"../assets/images/logos/google.png\", \"name\": \"Google\", \"mark\": \"\\u5168\\u7403\\u6700\\u5927\\u7684UI\\u5b66\\u4e60\\u5206\\u4eab\\u5e73\\u53f0\"}]}, {\"name\": \"\\u754c\\u9762\\u7075\\u611f\", \"rows\": [{\"url\": \"https://www.aliyun.com/\", \"img\": \"../assets/images/logos/aliyun.png\", \"name\": \"\\u963f\\u91cc\\u4e91\", \"mark\": \"\\u70b9\\u51fb\\u9886\\u53d62000\\u5143\\u9650\\u91cf\\u4e91\\u4ea7\\u54c1\\u4f18\\u60e0\\u5238\"}, {\"url\": \"https://www.gulusucai.com/\", \"img\": \"../assets/images/logos/gulusucai.png\", \"name\": \"\\u5495\\u565c\\u7d20\\u6750\", \"mark\": \"\\u8d28\\u91cf\\u5f88\\u9ad8\\u7684\\u8bbe\\u8ba1\\u7d20\\u6750\\u7f51\\u7ad9\\uff08\\u63a8\\u8350\\uff09\"}, {\"url\": \"https://xiyou4you.us\", \"img\": \"../assets/images/logos/xiyou.png\", \"name\": \"\\u897f\\u6e38-\\u79d1\\u5b66\\u4e0a\\u7f51\", \"mark\": \"\\u4f18\\u79c0\\u7684\\u79d1\\u5b66\\u4e0a\\u7f51\\uff08\\u7565\\u8d35\\uff0c\\u4f46\\u662f\\u8d3c\\u7a33\\uff0c\\u70b9\\u51fb\\u6ce8\\u518c\\u9886\\u53d6\\u4f18\\u60e0\\u5238\\uff09\"}, {\"url\": \"https://www.shejizhoukan.com/\", \"img\": \"../assets/images/logos/shejizhoukan.png\", \"name\": \"\\u8bbe\\u8ba1\\u5468\\u520a\", \"mark\": \"\\u6700\\u65b0\\u8bbe\\u8ba1\\u8d44\\u8baf\\uff08\\u6bcf\\u5929\\u66f4\\u65b0\\uff09\"}, {\"url\": \"https://www.ziticangku.com/\", \"img\": \"../assets/images/logos/ziticangku.png\", \"name\": \"\\u5b57\\u4f53\\u4ed3\\u5e93\", \"mark\": \"\\u6700\\u5168\\u7684\\u514d\\u8d39\\u5546\\u7528\\u5b57\\u4f53\\u5e93\"}]}, {\"name\": \"\\u793e\\u533a\\u8d44\\u8baf\", \"rows\": [{\"url\": \"https://www.leiphone.com/\", \"img\": \"../assets/images/logos/leiphone.png\", \"name\": \"\\u96f7\\u950b\\u7f51\", \"mark\": \"\\u4eba\\u5de5\\u667a\\u80fd\\u548c\\u667a\\u80fd\\u786c\\u4ef6\\u9886\\u57df\\u7684\\u4e92\\u8054\\u7f51\\u79d1\\u6280\\u5a92\\u4f53\"}, {\"url\": \"http://36kr.com/\", \"img\": \"../assets/images/logos/36kr.png\", \"name\": \"36kr\", \"mark\": \"\\u521b\\u4e1a\\u8d44\\u8baf\\u3001\\u79d1\\u6280\\u65b0\\u95fb\"}, {\"url\": \"https://www.digitaling.com/\", \"img\": \"../assets/images/logos/digitaling.png\", \"name\": \"\\u6570\\u82f1\\u7f51\", \"mark\": \"\\u6570\\u5b57\\u5a92\\u4f53\\u53ca\\u804c\\u4e1a\\u62db\\u8058\\u7f51\\u7ad9\"}, {\"url\": \"http://www.lieyunwang.com/\", \"img\": \"../assets/images/logos/lieyunwang.png\", \"name\": \"\\u730e\\u4e91\\u7f51\", \"mark\": \"\\u4e92\\u8054\\u7f51\\u521b\\u4e1a\\u9879\\u76ee\\u63a8\\u8350\\u548c\\u521b\\u4e1a\\u521b\\u65b0\\u8d44\\u8baf\"}, {\"url\": \"http://www.woshipm.com/\", \"img\": \"../assets/images/logos/woshipm.png\", \"name\": \"\\u4eba\\u4eba\\u90fd\\u662f\\u4ea7\\u54c1\\u7ecf\\u7406\", \"mark\": \"\\u4ea7\\u54c1\\u7ecf\\u7406\\u3001\\u4ea7\\u54c1\\u7231\\u597d\\u8005\\u5b66\\u4e60\\u4ea4\\u6d41\\u5e73\\u53f0\"}, {\"url\": \"https://www.zaodula.com/\", \"img\": \"../assets/images/logos/zaodula.png\", \"name\": \"\\u4e92\\u8054\\u7f51\\u65e9\\u8bfb\\u8bfe\", \"mark\": \"\\u4e92\\u8054\\u7f51\\u884c\\u4e1a\\u6df1\\u5ea6\\u9605\\u8bfb\\u4e0e\\u5b66\\u4e60\\u5e73\\u53f0\"}, {\"url\": \"http://www.chanpin100.com/\", \"img\": \"../assets/images/logos/chanpin100.png\", \"name\": \"\\u4ea7\\u54c1\\u58f9\\u4f70\", \"mark\": \"\\u4e3a\\u4ea7\\u54c1\\u7ecf\\u7406\\u7231\\u597d\\u8005\\u63d0\\u4f9b\\u6700\\u4f18\\u8d28\\u7684\\u4ea7\\u54c1\\u8d44\\u8baf\\u3001\\u539f\\u521b\\u5185\\u5bb9\\u548c\\u76f8\\u5173\\u89c6\\u9891\\u8bfe\\u7a0b\"}, {\"url\": \"http://www.pmcaff.com/\", \"img\": \"../assets/images/logos/pmcaff.png\", \"name\": \"PMCAFF\", \"mark\": \"\\u4e2d\\u56fd\\u7b2c\\u4e00\\u4ea7\\u54c1\\u7ecf\\u7406\\u4eba\\u6c14\\u7ec4\\u7ec7\\uff0c\\u4e13\\u6ce8\\u4e8e\\u7814\\u7a76\\u4e92\\u8054\\u7f51\\u4ea7\\u54c1\"}, {\"url\": \"http://www.iyunying.org/\", \"img\": \"../assets/images/logos/iyunying.png\", \"name\": \"\\u7231\\u8fd0\\u8425\", \"mark\": \"\\u7f51\\u7ad9\\u8fd0\\u8425\\u4eba\\u5458\\u5b66\\u4e60\\u4ea4\\u6d41\\uff0c\\u4e13\\u6ce8\\u4e8e\\u7f51\\u7ad9\\u4ea7\\u54c1\\u8fd0\\u8425\\u7ba1\\u7406\\u3001\\u6dd8\\u5b9d\\u8fd0\\u8425\\u3002\"}, {\"url\": \"http://www.niaogebiji.com/\", \"img\": \"../assets/images/logos/niaogebiji.png\", \"name\": \"\\u9e1f\\u54e5\\u7b14\\u8bb0\", \"mark\": \"\\u79fb\\u52a8\\u4e92\\u8054\\u7f51\\u7b2c\\u4e00\\u5e72\\u8d27\\u5e73\\u53f0\"}, {\"url\": \"http://www.gtn9.com/\", \"img\": \"../assets/images/logos/gtn9.png\", \"name\": \"\\u53e4\\u7530\\u8def9\\u53f7\", \"mark\": \"\\u56fd\\u5185\\u4e13\\u4e1a\\u54c1\\u724c\\u521b\\u610f\\u5e73\\u53f0\"}, {\"url\": \"http://www.uigreat.com/\", \"img\": \"../assets/images/logos/uigreat.png\", \"name\": \"\\u4f18\\u9601\\u7f51\", \"mark\": \"UI\\u8bbe\\u8ba1\\u5e08\\u5b66\\u4e60\\u4ea4\\u6d41\\u793e\\u533a\"}]}, {\"name\": \"\\u53d1\\u73b0\\u4ea7\\u54c1\", \"rows\": [{\"url\": \"https://www.producthunt.com/\", \"img\": \"../assets/images/logos/producthunt.png\", \"name\": \"Producthunt\", \"mark\": \"\\u53d1\\u73b0\\u65b0\\u9c9c\\u6709\\u8da3\\u7684\\u4ea7\\u54c1\"}, {\"url\": \"http://next.36kr.com/posts\", \"img\": \"../assets/images/logos/NEXT.png\", \"name\": \"NEXT\", \"mark\": \"\\u4e0d\\u9519\\u8fc7\\u4efb\\u4f55\\u4e00\\u4e2a\\u65b0\\u4ea7\\u54c1\"}, {\"url\": \"https://sspai.com/\", \"img\": \"../assets/images/logos/sspai.png\", \"name\": \"\\u5c11\\u6570\\u6d3e\", \"mark\": \"\\u9ad8\\u54c1\\u8d28\\u6570\\u5b57\\u6d88\\u8d39\\u6307\\u5357\"}, {\"url\": \"http://liqi.io/\", \"img\": \"../assets/images/logos/liqi.png\", \"name\": \"\\u5229\\u5668\", \"mark\": \"\\u521b\\u9020\\u8005\\u548c\\u4ed6\\u4eec\\u7684\\u5de5\\u5177\"}, {\"url\": \"http://today.itjuzi.com/\", \"img\": \"../assets/images/logos/today.png\", \"name\": \"Today\", \"mark\": \"\\u4e3a\\u8eab\\u8fb9\\u7684\\u65b0\\u4ea7\\u54c1\\u559d\\u5f69\"}, {\"url\": \"https://faxian.appinn.com/\", \"img\": \"../assets/images/logos/appinn.png\", \"name\": \"\\u5c0f\\u4f17\\u8f6f\\u4ef6\", \"mark\": \"\\u5728\\u8fd9\\u91cc\\u53d1\\u73b0\\u66f4\\u591a\\u6709\\u8da3\\u7684\\u5e94\\u7528\"}]}, {\"name\": \"\\u754c\\u9762\\u7075\\u611f\", \"rows\": [{\"url\": \"https://www.pttrns.com/\", \"img\": \"../assets/images/logos/Pttrns.png\", \"name\": \"Pttrns\", \"mark\": \"Check out the finest collection of design patterns, resources, mobile apps and inspiration\"}, {\"url\": \"http://collectui.com/\", \"img\": \"../assets/images/logos/CollectUI.png\", \"name\": \"Collect UI\", \"mark\": \"Daily inspiration collected from daily ui archive and beyond.\"}, {\"url\": \"http://ui.uigreat.com/\", \"img\": \"../assets/images/logos/uiuigreat.png\", \"name\": \"UI uigreat\", \"mark\": \"APP\\u754c\\u9762\\u622a\\u56fe\\u53c2\\u8003\"}, {\"url\": \"http://androidniceties.tumblr.com/\", \"img\": \"../assets/images/logos/AndroidNiceties.png\", \"name\": \"Android Niceties\", \"mark\": \"A collection of screenshots encompassing some of the most beautiful looking Android apps.\"}]}, {\"name\": \"\\u7f51\\u9875\\u7075\\u611f\", \"rows\": [{\"url\": \"https://www.awwwards.com/\", \"img\": \"../assets/images/logos/awwwards.png\", \"name\": \"Awwwards\", \"mark\": \"Awwwards are the Website Awards that recognize and promote the talent and effort of the best developers, designers and web agencies in the world.\"}, {\"url\": \"https://www.cssdesignawards.com/\", \"img\": \"../assets/images/logos/CSSDesignAwards.png\", \"name\": \"CSS Design Awards\", \"mark\": \"Website Awards & Inspiration - CSS Gallery\"}, {\"url\": \"https://thefwa.com/\", \"img\": \"../assets/images/logos/fwa.png\", \"name\": \"The FWA\", \"mark\": \"FWA - showcasing innovation every day since 2000\"}, {\"url\": \"http://www.ecommercefolio.com/\", \"img\": \"../assets/images/logos/Ecommercefolio.png\", \"name\": \"Ecommercefolio\", \"mark\": \"Only the Best Ecommerce Design Inspiration\"}, {\"url\": \"http://www.lapa.ninja/\", \"img\": \"../assets/images/logos/Lapa.png\", \"name\": \"Lapa\", \"mark\": \"The best landing page design inspiration from around the web.\"}, {\"url\": \"http://reeoo.com/\", \"img\": \"../assets/images/logos/reeoo.png\", \"name\": \"Reeoo\", \"mark\": \"web design inspiration and website gallery\"}, {\"url\": \"https://designmunk.com/\", \"img\": \"../assets/images/logos/designmunk.png\", \"name\": \"Designmunk\", \"mark\": \"Best Homepage Design Inspiration\"}, {\"url\": \"https://bestwebsite.gallery/\", \"img\": \"../assets/images/logos/BWG.png\", \"name\": \"Best Websites Gallery\", \"mark\": \"Website Showcase Inspiration | Best Websites Gallery\"}, {\"url\": \"http://www.pages.xyz/\", \"img\": \"../assets/images/logos/pages.png\", \"name\": \"Pages\", \"mark\": \"Curated directory of the best Pages\"}, {\"url\": \"https://sitesee.co/\", \"img\": \"../assets/images/logos/SiteSee.png\", \"name\": \"SiteSee\", \"mark\": \"SiteSee is a curated gallery of beautiful, modern websites collections.\"}, {\"url\": \"https://www.siteinspire.com/\", \"img\": \"../assets/images/logos/siteInspire.png\", \"name\": \"Site Inspire\", \"mark\": \"A CSS gallery and showcase of the best web design inspiration.\"}, {\"url\": \"http://web.uedna.com/\", \"img\": \"../assets/images/logos/WebInspiration.png\", \"name\": \"WebInspiration\", \"mark\": \"\\u7f51\\u9875\\u8bbe\\u8ba1\\u6b23\\u8d4f,\\u5168\\u7403\\u9876\\u7ea7\\u7f51\\u9875\\u8bbe\\u8ba1\"}, {\"url\": \"http://navnav.co/\", \"img\": \"../assets/images/logos/navnav.png\", \"name\": \"navnav\", \"mark\": \"A ton of CSS, jQuery, and JavaScript responsive navigation examples, demos, and tutorials from all over the web.\"}, {\"url\": \"https://www.reallygoodux.io/\", \"img\": \"../assets/images/logos/ReallyGoodUX.png\", \"name\": \"Really Good UX\", \"mark\": null}]}, {\"name\": \"\\u56fe\\u6807\\u7d20\\u6750\", \"rows\": [{\"url\": \"https://www.iconfinder.com\", \"img\": \"../assets/images/logos/Iconfinder.png\", \"name\": \"Iconfinder\", \"mark\": \"2,100,000+ free and premium vector icons.\"}, {\"url\": \"http://www.iconfont.cn/\", \"img\": \"../assets/images/logos/iconfont.png\", \"name\": \"iconfont\", \"mark\": \"\\u963f\\u91cc\\u5df4\\u5df4\\u77e2\\u91cf\\u56fe\\u6807\\u5e93\"}, {\"url\": \"https://iconmonstr.com/\", \"img\": \"../assets/images/logos/iconmonstr.png\", \"name\": \"iconmonstr\", \"mark\": \"Free simple icons for your next project\"}, {\"url\": \"http://www.iconarchive.com/\", \"img\": \"../assets/images/logos/iconarchive.png\", \"name\": \"Icon Archive\", \"mark\": \"Search 590,912 free icons\"}, {\"url\": \"https://findicons.com/\", \"img\": \"../assets/images/logos/FindIcons.png\", \"name\": \"FindIcons\", \"mark\": \"Search through 300,000 free icons\"}, {\"url\": \"https://icomoon.io/app/\", \"img\": \"../assets/images/logos/IcoMoonApp.png\", \"name\": \"IcoMoonApp\", \"mark\": \"Icon Font, SVG, PDF & PNG Generator\"}, {\"url\": \"http://www.easyicon.net/\", \"img\": \"../assets/images/logos/easyicon.png\", \"name\": \"easyicon\", \"mark\": \"PNG\\u3001ICO\\u3001ICNS\\u683c\\u5f0f\\u56fe\\u6807\\u641c\\u7d22\\u3001\\u56fe\\u6807\\u4e0b\\u8f7d\\u670d\\u52a1\"}, {\"url\": \"https://www.flaticon.com/\", \"img\": \"../assets/images/logos/flaticon.png\", \"name\": \"flaticon\", \"mark\": \"634,000+ Free vector icons in SVG, PSD, PNG, EPS format or as ICON FONT.\"}, {\"url\": \"http://ui-cloud.com/\", \"img\": \"../assets/images/logos/UICloud.png\", \"name\": \"UICloud\", \"mark\": \"The largest user interface design database in the world.\"}, {\"url\": \"https://material.io/icons/\", \"img\": \"../assets/images/logos/Materialicons.png\", \"name\": \"Material icons\", \"mark\": \"Access over 900 material system icons, available in a variety of sizes and densities, and as a web font.\"}, {\"url\": \"fontawesomeicon\", \"img\": \"../assets/images/logos/fontawesomeicon.png\", \"name\": \"Font Awesome Icon\", \"mark\": \"The complete set of 675 icons in Font Awesome\"}, {\"url\": \"http://ionicons.com/\", \"img\": \"../assets/images/logos/ionicons.png\", \"name\": \"ion icons\", \"mark\": \"The premium icon font for Ionic Framework.\"}, {\"url\": \"http://simplelineicons.com/\", \"img\": \"../assets/images/logos/simplelineicons.png\", \"name\": \"Simpleline Icons\", \"mark\": \"Simple line Icons pack\"}]}, {\"name\": \"LOGO\\u8bbe\\u8ba1\", \"rows\": [{\"url\": \"http://www.iconsfeed.com/\", \"img\": \"../assets/images/logos/iconsfeed.png\", \"name\": \"Iconsfeed\", \"mark\": \"iOS icons gallery\"}, {\"url\": \"http://iosicongallery.com/\", \"img\": \"../assets/images/logos/iosicongallery.png\", \"name\": \"iOS Icon Gallery\", \"mark\": \"Showcasing beautiful icon designs from the iOS App Store\"}, {\"url\": \"https://worldvectorlogo.com/\", \"img\": \"../assets/images/logos/worldvectorlogo.png\", \"name\": \"World Vector Logo\", \"mark\": \"Brand logos free to download\"}, {\"url\": \"http://instantlogosearch.com/\", \"img\": \"../assets/images/logos/InstantLogoSearch.png\", \"name\": \"Instant Logo Search\", \"mark\": \"Search & download thousands of logos instantly\"}]}, {\"name\": \"\\u5e73\\u9762\\u7d20\\u6750\", \"rows\": [{\"url\": \"https://www.gulusucai.com/\", \"img\": \"../assets/images/logos/gulusucai.png\", \"name\": \"\\u5495\\u565c\\u7d20\\u6750\", \"mark\": \"\\u8d28\\u91cf\\u5f88\\u9ad8\\u7684\\u8bbe\\u8ba1\\u7d20\\u6750\\u7f51\\u7ad9\\uff08\\u826f\\u5fc3\\u63a8\\u8350\\uff09\"}, {\"url\": \"https://wallhalla.com/\", \"img\": \"../assets/images/logos/wallhalla.png\", \"name\": \"wallhalla\", \"mark\": \"Find awesome high quality wallpapers for desktop and mobile in one place.\"}, {\"url\": \"https://365psd.com/\", \"img\": \"../assets/images/logos/365PSD.png\", \"name\": \"365PSD\", \"mark\": \"Free PSD & Graphics, Illustrations\"}, {\"url\": \"https://medialoot.com/\", \"img\": \"../assets/images/logos/Medialoot.png\", \"name\": \"Medialoot\", \"mark\": \"Free & Premium Design Resources \\u2014 Medialoot\"}, {\"url\": \"http://www.58pic.com/\", \"img\": \"../assets/images/logos/qiantu.png\", \"name\": \"\\u5343\\u56fe\\u7f51\", \"mark\": \"\\u4e13\\u6ce8\\u514d\\u8d39\\u8bbe\\u8ba1\\u7d20\\u6750\\u4e0b\\u8f7d\\u7684\\u7f51\\u7ad9\"}, {\"url\": \"http://588ku.com/\", \"img\": \"../assets/images/logos/qianku.png\", \"name\": \"\\u5343\\u5e93\\u7f51\", \"mark\": \"\\u514d\\u8d39png\\u56fe\\u7247\\u80cc\\u666f\\u7d20\\u6750\\u4e0b\\u8f7d\"}, {\"url\": \"http://www.ooopic.com/\", \"img\": \"../assets/images/logos/wotu.png\", \"name\": \"\\u6211\\u56fe\\u7f51\", \"mark\": \"\\u6211\\u56fe\\u7f51,\\u63d0\\u4f9b\\u56fe\\u7247\\u7d20\\u6750\\u53ca\\u6a21\\u677f\\u4e0b\\u8f7d,\\u4e13\\u6ce8\\u6b63\\u7248\\u8bbe\\u8ba1\\u4f5c\\u54c1\\u4ea4\\u6613\"}, {\"url\": \"http://90sheji.com/\", \"img\": \"../assets/images/logos/90sheji.png\", \"name\": \"90\\u8bbe\\u8ba1\", \"mark\": \"\\u7535\\u5546\\u8bbe\\u8ba1\\uff08\\u6dd8\\u5b9d\\u7f8e\\u5de5\\uff09\\u5343\\u56fe\\u514d\\u8d39\\u6dd8\\u5b9d\\u7d20\\u6750\\u5e93\"}, {\"url\": \"http://www.nipic.com/\", \"img\": \"../assets/images/logos/nipic.png\", \"name\": \"\\u6635\\u56fe\\u7f51\", \"mark\": \"\\u539f\\u521b\\u7d20\\u6750\\u5171\\u4eab\\u5e73\\u53f0\"}, {\"url\": \"http://www.lanrentuku.com/\", \"img\": \"../assets/images/logos/lanrentuku.png\", \"name\": \"\\u61d2\\u4eba\\u56fe\\u5e93\", \"mark\": \"\\u61d2\\u4eba\\u56fe\\u5e93\\u4e13\\u6ce8\\u4e8e\\u63d0\\u4f9b\\u7f51\\u9875\\u7d20\\u6750\\u4e0b\\u8f7d\"}, {\"url\": \"http://so.ui001.com/\", \"img\": \"../assets/images/logos/sousucai.png\", \"name\": \"\\u7d20\\u6750\\u641c\\u7d22\", \"mark\": \"\\u8bbe\\u8ba1\\u7d20\\u6750\\u641c\\u7d22\\u805a\\u5408\"}, {\"url\": \"http://psefan.com/\", \"img\": \"../assets/images/logos/psefan.png\", \"name\": \"PS\\u996d\\u56e2\\u7f51\", \"mark\": \"\\u4e0d\\u4e00\\u6837\\u7684\\u8bbe\\u8ba1\\u7d20\\u6750\\u5e93\\uff01\\u8ba9\\u81ea\\u5df1\\u7684\\u8bbe\\u8ba1\\u4e0e\\u4f17\\u4e0d\\u540c\\uff01\"}, {\"url\": \"http://www.sccnn.com/\", \"img\": \"../assets/images/logos/sccnn.png\", \"name\": \"\\u7d20\\u6750\\u4e2d\\u56fd\", \"mark\": \"\\u514d\\u8d39\\u7d20\\u6750\\u5171\\u4eab\\u5e73\\u53f0\"}, {\"url\": \"https://www.freepik.com/\", \"img\": \"../assets/images/logos/freepik.png\", \"name\": \"freepik\", \"mark\": \"More than a million free vectors, PSD, photos and free icons.\"}]}, {\"name\": \"UI\\u8d44\\u6e90\", \"rows\": [{\"url\": \"https://www.gulusucai.com/\", \"img\": \"../assets/images/logos/gulusucai.png\", \"name\": \"\\u5495\\u565c\\u7d20\\u6750\", \"mark\": \"\\u8d28\\u91cf\\u5f88\\u9ad8\\u7684\\u8bbe\\u8ba1\\u7d20\\u6750\\u7f51\\u7ad9\\uff08\\u826f\\u5fc3\\u63a8\\u8350\\uff09\"}, {\"url\": \"https://freebiesbug.com/\", \"img\": \"../assets/images/logos/freebiesbug.png\", \"name\": \"Freebiesbug\", \"mark\": \"Hand-picked resources for web designer and developers, constantly updated.\"}, {\"url\": \"https://freebiesupply.com/\", \"img\": \"../assets/images/logos/freebiesupply.png\", \"name\": \"Freebie Supply\", \"mark\": \"Free Resources For Designers\"}, {\"url\": \"https://www.yrucd.com/\", \"img\": \"../assets/images/logos/yrucd.png\", \"name\": \"\\u4e91\\u745e\", \"mark\": \"\\u4f18\\u79c0\\u8bbe\\u8ba1\\u8d44\\u6e90\\u7684\\u5206\\u4eab\\u7f51\\u7ad9\"}, {\"url\": \"https://xituqu.com/\", \"img\": \"../assets/images/logos/xituqu.png\", \"name\": \"\\u7a00\\u571f\\u533a\", \"mark\": \"\\u4f18\\u8d28\\u8bbe\\u8ba1\\u5f00\\u53d1\\u8d44\\u6e90\\u5206\\u4eab\"}, {\"url\": \"https://ui8.net/\", \"img\": \"../assets/images/logos/ui8.png\", \"name\": \"ui8\", \"mark\": \"UI Kits, Wireframe Kits, Templates, Icons and More\"}, {\"url\": \"https://www.uplabs.com/\", \"img\": \"../assets/images/logos/uplabs.png\", \"name\": \"uplabs\", \"mark\": \"Daily resources for product designers & developers\"}, {\"url\": \"http://www.uikit.me/\", \"img\": \"../assets/images/logos/uikitme.png\", \"name\": \"UIkit.me\", \"mark\": \"\\u6700\\u4fbf\\u6377\\u65b0\\u9c9c\\u7684uikit\\u8d44\\u6e90\\u4e0b\\u8f7d\\u7f51\\u7ad9\"}, {\"url\": \"http://www.fribbble.com/\", \"img\": \"../assets/images/logos/Fribbble.png\", \"name\": \"Fribbble\", \"mark\": \"Free PSD files and other free design resources by Dribbblers.\"}, {\"url\": \"http://principlerepo.com/\", \"img\": \"../assets/images/logos/PrincipleRepo.png\", \"name\": \"PrincipleRepo\", \"mark\": \"Free, High Quality Principle Resources\"}, {\"url\": \"https://designmodo.com/\", \"img\": \"../assets/images/logos/Designmodo.png\", \"name\": \"Designmodo\", \"mark\": \"Web Design Blog and Shop\"}]}, {\"name\": \"Sketch\\u8d44\\u6e90\", \"rows\": [{\"url\": \"https://sketchapp.com/\", \"img\": \"../assets/images/logos/Sketch.png\", \"name\": \"Sketch\", \"mark\": \"The digital design toolkit\"}, {\"url\": \"http://utom.design/measure/\", \"img\": \"../assets/images/logos/SketchMeasure.png\", \"name\": \"Sketch Measure\", \"mark\": \"Friendly user interface offers you a more intuitive way of making marks.\"}, {\"url\": \"https://www.sketchappsources.com/\", \"img\": \"../assets/images/logos/sketchappsources.png\", \"name\": \"Sketch App Sources\", \"mark\": \"Free design resources and plugins - Icons, UI Kits, Wireframes, iOS, Android Templates for Sketch\"}, {\"url\": \"http://www.sketch.im/\", \"img\": \"../assets/images/logos/sketchIm.png\", \"name\": \"Sketch.im\", \"mark\": \"Sketch \\u76f8\\u5173\\u8d44\\u6e90\\u6c47\\u805a\"}, {\"url\": \"http://sketchhunt.com/\", \"img\": \"../assets/images/logos/sketchhunt.png\", \"name\": \"Sketch Hunt\", \"mark\": \"Sketch Hunt is an independent blog sharing gems in learning, plugins & design tools for fans of Sketch app.\"}, {\"url\": \"http://www.sketchcn.com/\", \"img\": \"../assets/images/logos/sketchcn.png\", \"name\": \"Sketch\\u4e2d\\u6587\\u7f51\", \"mark\": \"\\u5206\\u4eab\\u6700\\u65b0\\u7684Sketch\\u4e2d\\u6587\\u624b\\u518c\"}, {\"url\": \"https://awesome-sket.ch/\", \"img\": \"../assets/images/logos/AwesomeSketchPlugins.png\", \"name\": \"Awesome Sketch Plugins\", \"mark\": \"A collection of really useful Sketch plugins.\"}, {\"url\": \"https://www.sketchcasts.net/\", \"img\": \"../assets/images/logos/sketchcasts.png\", \"name\": \"Sketchcasts\", \"mark\": \"Learn Sketch Train your design skills with a weekly video tutorial\"}]}, {\"name\": \"\\u5b57\\u4f53\\u8d44\\u6e90\", \"rows\": [{\"url\": \"https://www.ziticangku.com/\", \"img\": \"../assets/images/logos/ziticangku.png\", \"name\": \"\\u5b57\\u4f53\\u4ed3\\u5e93\", \"mark\": \"\\u6700\\u5168\\u7684\\u514d\\u8d39\\u5546\\u7528\\u5b57\\u4f53\\u5e93\"}, {\"url\": \"https://fonts.google.com/\", \"img\": \"../assets/images/logos/googlefont.png\", \"name\": \"Google Font\", \"mark\": \"Making the web more beautiful, fast, and open through great typography\"}, {\"url\": \"https://typekit.com/\", \"img\": \"../assets/images/logos/typekit.png\", \"name\": \"Typekit\", \"mark\": \"Quality fonts from the world\\u2019s best foundries.\"}, {\"url\": \"http://www.foundertype.com/\", \"img\": \"../assets/images/logos/Fondertype.png\", \"name\": \"\\u65b9\\u6b63\\u5b57\\u5e93\", \"mark\": \"\\u65b9\\u6b63\\u5b57\\u5e93\\u5b98\\u65b9\\u7f51\\u7ad9\"}, {\"url\": \"http://ziticq.com/\", \"img\": \"../assets/images/logos/ziticq.png\", \"name\": \"\\u5b57\\u4f53\\u4f20\\u5947\\u7f51\", \"mark\": \"\\u4e2d\\u56fd\\u9996\\u4e2a\\u5b57\\u4f53\\u54c1\\u724c\\u8bbe\\u8ba1\\u5e08\\u4ea4\\u6d41\\u7f51\"}, {\"url\": \"https://www.fontsquirrel.com/\", \"img\": \"../assets/images/logos/fontsquirrel.png\", \"name\": \"Fontsquirrel\", \"mark\": \"FREE fonts for graphic designers\"}, {\"url\": \"https://www.urbanfonts.com/\", \"img\": \"../assets/images/logos/UrbanFonts.png\", \"name\": \"Urban Fonts\", \"mark\": \"Download Free Fonts and Free Dingbats.\"}, {\"url\": \"http://www.losttype.com/\", \"img\": \"../assets/images/logos/losttype.png\", \"name\": \"Lost Type\", \"mark\": \"Lost Type is a Collaborative Digital Type Foundry\"}, {\"url\": \"https://fonts2u.com/\", \"img\": \"../assets/images/logos/fonts2u.png\", \"name\": \"FONTS2U\", \"mark\": \"Download free fonts for Windows and Macintosh.\"}, {\"url\": \"http://www.fontex.org/\", \"img\": \"../assets/images/logos/fontex.png\", \"name\": \"Fontex\", \"mark\": \"Free Fonts to Download + Premium Typefaces\"}, {\"url\": \"http://fontm.com/\", \"img\": \"../assets/images/logos/FontM.png\", \"name\": \"FontM\", \"mark\": \"Free Fonts\"}, {\"url\": \"http://www.myfonts.com/\", \"img\": \"../assets/images/logos/MyFonts.png\", \"name\": \"My Fonts\", \"mark\": \"Fonts for Print, Products & Screens\"}, {\"url\": \"https://www.dafont.com/\", \"img\": \"../assets/images/logos/dafont.png\", \"name\": \"Da Font\", \"mark\": \"Archive of freely downloadable fonts.\"}, {\"url\": \"https://www.onlinewebfonts.com/\", \"img\": \"../assets/images/logos/OnlineWebFonts.png\", \"name\": \"OnlineWebFonts\", \"mark\": \"WEB Free Fonts for Windows and Mac / Font free Download\"}, {\"url\": \"http://www.abstractfonts.com/\", \"img\": \"../assets/images/logos/abstractfonts.png\", \"name\": \"Abstract Fonts\", \"mark\": \"Abstract Fonts (13,866 free fonts)\"}]}, {\"name\": \"Mockup\", \"rows\": [{\"url\": \"https://mockup.zone/\", \"img\": \"../assets/images/logos/MockupZone.png\", \"name\": \"MockupZone\", \"mark\": \"Mockup Zone is an online store where you can find free and premium PSD mockup files to show your designs in a professional way.\"}, {\"url\": \"http://dunnnk.com/\", \"img\": \"../assets/images/logos/Dunnnk.png\", \"name\": \"Dunnnk\", \"mark\": \" Generate Product Mockups For Free\"}, {\"url\": \"http://www.graphberry.com/\", \"img\": \"../assets/images/logos/graphberry.png\", \"name\": \"Graphberry\", \"mark\": \"Free design resources, Mockups, PSD web templates, Icons\"}, {\"url\": \"http://threed.io/\", \"img\": \"../assets/images/logos/threed.png\", \"name\": \"Threed\", \"mark\": \"Generate 3D Mockups right in your Browser\"}, {\"url\": \"https://free.lstore.graphics/\", \"img\": \"../assets/images/logos/mockupworld.png\", \"name\": \"Mockup World\", \"mark\": \"The best free Mockups from the Web\"}, {\"url\": \"https://free.lstore.graphics/\", \"img\": \"../assets/images/logos/lstore.png\", \"name\": \"Lstore\", \"mark\": \"Exclusive mindblowing freebies for designers and developers\"}, {\"url\": \"https://www.pixeden.com/\", \"img\": \"../assets/images/logos/pixeden.png\", \"name\": \"pixeden\", \"mark\": \"free web resources and graphic design templates.\"}, {\"url\": \"http://forgraphictm.com/\", \"img\": \"../assets/images/logos/forgraphictm.png\", \"name\": \"For Graphic TM\", \"mark\": \"High Quality PSD Mockups for Graphic Designers.\"}]}, {\"name\": \"\\u6444\\u5f71\\u56fe\\u5e93\", \"rows\": [{\"url\": \"https://unsplash.com/\", \"img\": \"../assets/images/logos/unsplash.png\", \"name\": \"Unsplash\", \"mark\": \"Beautiful, free photos.\"}, {\"url\": \"https://visualhunt.com/\", \"img\": \"../assets/images/logos/visualhunt.png\", \"name\": \"visualhunt\", \"mark\": \"100% Free High Quality Photos\"}, {\"url\": \"https://librestock.com/\", \"img\": \"../assets/images/logos/librestock.png\", \"name\": \"librestock\", \"mark\": \"65,084 high quality do-what-ever-you-want stock photos\"}, {\"url\": \"https://pixabay.com/\", \"img\": \"../assets/images/logos/pixabay.png\", \"name\": \"pixabay\", \"mark\": \"\\u53ef\\u5728\\u4efb\\u4f55\\u5730\\u65b9\\u4f7f\\u7528\\u7684\\u514d\\u8d39\\u56fe\\u7247\\u548c\\u89c6\\u9891\"}, {\"url\": \"https://www.splitshire.com/\", \"img\": \"../assets/images/logos/SplitShire.png\", \"name\": \"SplitShire\", \"mark\": \"Free Stock Photos and Videos for commercial use.\"}, {\"url\": \"https://stocksnap.io/\", \"img\": \"../assets/images/logos/StockSnap.png\", \"name\": \"StockSnap\", \"mark\": \"Beautiful free stock photos\"}, {\"url\": \"http://albumarium.com/\", \"img\": \"../assets/images/logos/albumarium.png\", \"name\": \"albumarium\", \"mark\": \"The best place to find & share beautiful images\"}, {\"url\": \"https://myphotopack.com/\", \"img\": \"../assets/images/logos/myphotopack.png\", \"name\": \"myphotopack\", \"mark\": \"A free photo pack just for you. Every month.\"}, {\"url\": \"http://notaselfie.com/\", \"img\": \"../assets/images/logos/notaselfie.png\", \"name\": \"Notaselfie\", \"mark\": \"Photos that happen along the way. You can use the images anyway you like. Have fun!\"}, {\"url\": \"http://papers.co/\", \"img\": \"../assets/images/logos/papers.png\", \"name\": \"papers\", \"mark\": \"Wallpapers Every Hour!Hand collected :)\"}, {\"url\": \"http://stokpic.com/\", \"img\": \"../assets/images/logos/stokpic.png\", \"name\": \"stokpic\", \"mark\": \"Free Stock Photos For Commercial Use\"}, {\"url\": \"https://55mm.co/visuals\", \"img\": \"../assets/images/logos/55mm.png\", \"name\": \"55mm\", \"mark\": \"Use our FREE photos to tell your story! \"}, {\"url\": \"http://thestocks.im/\", \"img\": \"../assets/images/logos/thestocks.png\", \"name\": \"thestocks\", \"mark\": \"Use our FREE photos to tell your story! \"}, {\"url\": \"http://freenaturestock.com/\", \"img\": \"../assets/images/logos/freenaturestock.png\", \"name\": \"freenaturestock\", \"mark\": \"Exclusive mindblowing freebies for designers and developers\"}, {\"url\": \"https://negativespace.co/\", \"img\": \"../assets/images/logos/negativespace.png\", \"name\": \"negativespace\", \"mark\": \"Beautiful, High-Resolution Free Stock Photos\"}, {\"url\": \"https://gratisography.com/\", \"img\": \"../assets/images/logos/gratisography.png\", \"name\": \"gratisography\", \"mark\": \"Free high-resolution pictures you can use on your personal and commercial projects, free of copyright restrictions. \"}, {\"url\": \"http://imcreator.com/free\", \"img\": \"../assets/images/logos/imcreator.png\", \"name\": \"imcreator\", \"mark\": \"A curated collection of free web design resources, all for commercial use.\"}, {\"url\": \"http://www.lifeofpix.com/\", \"img\": \"../assets/images/logos/lifeofpix.png\", \"name\": \"lifeofpix\", \"mark\": \"Free high resolution photography\"}, {\"url\": \"https://skitterphoto.com/\", \"img\": \"../assets/images/logos/skitterphoto.png\", \"name\": \"skitterphoto\", \"mark\": \"Free Stock Photos for Creative Professionals\"}, {\"url\": \"https://mmtstock.com/\", \"img\": \"../assets/images/logos/mmtstock.png\", \"name\": \"mmtstock\", \"mark\": \"Free photos for commercial use\"}, {\"url\": \"https://skitterphoto.com/\", \"img\": \"../assets/images/logos/skitterphoto.png\", \"name\": \"skitterphoto\", \"mark\": \"a place to find, show and share public domain photos\"}, {\"url\": \"https://magdeleine.co/browse/\", \"img\": \"../assets/images/logos/magdeleine.png\", \"name\": \"magdeleine\", \"mark\": \"HAND-PICKED FREE PHOTOS FOR YOUR INSPIRATION\"}, {\"url\": \"http://jeshoots.com/\", \"img\": \"../assets/images/logos/jeshoots.png\", \"name\": \"jeshoots\", \"mark\": \"New Free Photos & Mockups in to your Inbox!\"}, {\"url\": \"https://www.hdwallpapers.net\", \"img\": \"../assets/images/logos/hdwallpapers.png\", \"name\": \"hdwallpapers\", \"mark\": \"High Definition Wallpapers & Desktop Backgrounds\"}, {\"url\": \"http://publicdomainarchive.com/\", \"img\": \"../assets/images/logos/publicdomainarchive.png\", \"name\": \"publicdomainarchive\", \"mark\": \"New 100% Free Stock Photos. Every. Single. Week.\"}]}, {\"name\": \"PPT\\u8d44\\u6e90\", \"rows\": [{\"url\": \"http://www.officeplus.cn/Template/Home.shtml\", \"img\": \"../assets/images/logos/officeplus.png\", \"name\": \"OfficePLUS\", \"mark\": \"OfficePLUS\\uff0c\\u5fae\\u8f6fOffice\\u5b98\\u65b9\\u5728\\u7ebf\\u6a21\\u677f\\u7f51\\u7ad9\\uff01\"}, {\"url\": \"http://www.ypppt.com/\", \"img\": \"../assets/images/logos/ypppt.png\", \"name\": \"\\u4f18\\u54c1PPT\", \"mark\": \"\\u9ad8\\u8d28\\u91cf\\u7684\\u6a21\\u7248\\uff0c\\u800c\\u4e14\\u8fd8\\u6709PPT\\u56fe\\u8868\\uff0cPPT\\u80cc\\u666f\\u56fe\\u7b49\\u8d44\\u6e90\"}, {\"url\": \"http://www.pptplus.cn/\", \"img\": \"../assets/images/logos/pptplus.png\", \"name\": \"PPT+\", \"mark\": \"PPT\\u52a0\\u76f4\\u64ad\\u3001\\u5f55\\u5236\\u548c\\u5206\\u4eab\\u2014PPT+\\u8bed\\u97f3\\u5185\\u5bb9\\u5206\\u4eab\\u5e73\\u53f0\"}, {\"url\": \"http://www.pptmind.com/\", \"img\": \"../assets/images/logos/pptmind.png\", \"name\": \"PPTMind\", \"mark\": \"\\u5206\\u4eab\\u9ad8\\u7aefppt\\u6a21\\u677f\\u4e0ekeynote\\u6a21\\u677f\\u7684\\u6570\\u5b57\\u4f5c\\u54c1\\u4ea4\\u6613\\u5e73\\u53f0\"}, {\"url\": \"http://www.tretars.com/ppt-templates\", \"img\": \"../assets/images/logos/tretars.png\", \"name\": \"tretars\", \"mark\": \"The best free Mockups from the Web\"}, {\"url\": \"http://ppt.500d.me/\", \"img\": \"../assets/images/logos/500d.png\", \"name\": \"5\\u767e\\u4e01\", \"mark\": \"\\u4e2d\\u56fd\\u9886\\u5148\\u7684PPT\\u6a21\\u677f\\u5171\\u4eab\\u5e73\\u53f0\"}]}, {\"name\": \"\\u56fe\\u5f62\\u521b\\u610f\", \"rows\": [{\"url\": \"https://www.adobe.com/cn/products/photoshop.html\", \"img\": \"../assets/images/logos/photoshop.png\", \"name\": \"photoshop\", \"mark\": \"Photoshop\\u4e0d\\u9700\\u8981\\u89e3\\u91ca\"}, {\"url\": \"https://affinity.serif.com/\", \"img\": \"../assets/images/logos/AffinityDesigner.png\", \"name\": \"Affinity Designer\", \"mark\": \"\\u4e13\\u4e1a\\u521b\\u610f\\u8f6f\\u4ef6\"}, {\"url\": \"https://www.adobe.com/cn/products/illustrator/\", \"img\": \"../assets/images/logos/Illustrator.png\", \"name\": \"Illustrator\", \"mark\": \"\\u77e2\\u91cf\\u56fe\\u5f62\\u548c\\u63d2\\u56fe\\u3002\"}, {\"url\": \"http://www.adobe.com/cn/products/indesign.html\", \"img\": \"../assets/images/logos/INDESIGN .png\", \"name\": \"indesign\", \"mark\": \"\\u9875\\u9762\\u8bbe\\u8ba1\\u3001\\u5e03\\u5c40\\u548c\\u51fa\\u7248\\u3002\"}, {\"url\": \"https://www.maxon.net/en/products/cinema-4d/overview/\", \"img\": \"../assets/images/logos/cinema4d.png\", \"name\": \"cinema-4d\", \"mark\": \"Cinema 4D is the perfect package for all 3D artists who want to achieve breathtaking results fast and hassle-free.\"}, {\"url\": \"https://www.autodesk.com/products/3ds-max/overview\", \"img\": \"../assets/images/logos/3dsmax.png\", \"name\": \"3ds-max\", \"mark\": \"3D modeling, animation, and rendering software\"}, {\"url\": \"https://www.blender.org/\", \"img\": \"../assets/images/logos/blender.png\", \"name\": \"Blender\", \"mark\": \"Blender is the free and open source 3D creation suite.\"}]}, {\"name\": \"\\u754c\\u9762\\u8bbe\\u8ba1\", \"rows\": [{\"url\": \"https://sketchapp.com/\", \"img\": \"../assets/images/logos/sketchapp.png\", \"name\": \"Sketch\", \"mark\": \"The digital design toolkit\"}, {\"url\": \"http://www.adobe.com/products/xd.html\", \"img\": \"../assets/images/logos/ADOBEXDCC.png\", \"name\": \"Adobe XD\", \"mark\": \"Introducing Adobe XD. Design. Prototype. Experience.\"}, {\"url\": \"https://www.invisionapp.com/\", \"img\": \"../assets/images/logos/invisionapp.png\", \"name\": \"invisionapp\", \"mark\": \"Powerful design prototyping tools\"}, {\"url\": \"https://marvelapp.com/\", \"img\": \"../assets/images/logos/marvelapp.png\", \"name\": \"marvelapp\", \"mark\": \"Simple design, prototyping and collaboration\"}, {\"url\": \"https://creative.adobe.com/zh-cn/products/download/muse\", \"img\": \"../assets/images/logos/MuseCC.png\", \"name\": \"Muse CC\", \"mark\": \"\\u65e0\\u9700\\u5229\\u7528\\u7f16\\u7801\\u5373\\u53ef\\u8fdb\\u884c\\u7f51\\u7ad9\\u8bbe\\u8ba1\\u3002\"}, {\"url\": \"https://www.figma.com/\", \"img\": \"../assets/images/logos/figma.png\", \"name\": \"figma\", \"mark\": \"Design, prototype, and gather feedback all in one place with Figma.\"}]}, {\"name\": \"\\u4ea4\\u4e92\\u52a8\\u6548\", \"rows\": [{\"url\": \"https://www.adobe.com/cn/products/aftereffects/\", \"img\": \"../assets/images/logos/AdobeAfterEffectsCC.png\", \"name\": \"Adobe After Effects CC\", \"mark\": \"\\u7535\\u5f71\\u822c\\u7684\\u89c6\\u89c9\\u6548\\u679c\\u548c\\u52a8\\u6001\\u56fe\\u5f62\\u3002\"}, {\"url\": \"http://principleformac.com/\", \"img\": \"../assets/images/logos/principle.png\", \"name\": \"principle\", \"mark\": \"Animate Your Ideas, Design Better Apps\"}, {\"url\": \"https://www.flinto.com/\", \"img\": \"../assets/images/logos/flinto.png\", \"name\": \"flinto\", \"mark\": \"Flinto is a Mac app used by top designers around the world to create interactive and animated prototypes of their app designs.\"}, {\"url\": \"https://framer.com/\", \"img\": \"../assets/images/logos/framer.png\", \"name\": \"framer\", \"mark\": \"Design everything from detailed icons to high-fidelity interactions\\u2014all in one place.\"}, {\"url\": \"http://www.protopie.cn/\", \"img\": \"../assets/images/logos/protopie.png\", \"name\": \"ProtoPie\", \"mark\": \"\\u9ad8\\u4fdd\\u771f\\u4ea4\\u4e92\\u539f\\u578b\\u8bbe\\u8ba1\"}]}, {\"name\": \"\\u5728\\u7ebf\\u914d\\u8272\", \"rows\": [{\"url\": \"http://khroma.co/generator/\", \"img\": \"../assets/images/logos/khroma.png\", \"name\": \"khroma\", \"mark\": \"Khroma is the fastest way to discover, search, and save color combos you'll want to use.\"}, {\"url\": \"https://uigradients.com\", \"img\": \"../assets/images/logos/uigradients.png\", \"name\": \"uigradients\", \"mark\": \"Beautiful colored gradients\"}, {\"url\": \"http://gradients.io/\", \"img\": \"../assets/images/logos/gradients.png\", \"name\": \"gradients\", \"mark\": \"Curated gradients for designers and developers\"}, {\"url\": \"https://webkul.github.io/coolhue/\", \"img\": \"../assets/images/logos/Coolest.png\", \"name\": \"Coolest\", \"mark\": \"Coolest handpicked Gradient Hues for your next super \\u26a1 amazing stuff\"}, {\"url\": \"https://webgradients.com/\", \"img\": \"../assets/images/logos/webgradients.png\", \"name\": \"webgradients\", \"mark\": \"WebGradients is a free collection of 180 linear gradients that you can use as content backdrops in any part of your website. \"}, {\"url\": \"https://www.grabient.com/\", \"img\": \"../assets/images/logos/grabient.png\", \"name\": \"grabient\", \"mark\": \"2017 Grabient by unfold\"}, {\"url\": \"http://www.thedayscolor.com/\", \"img\": \"../assets/images/logos/thedayscolor.png\", \"name\": \"thedayscolor\", \"mark\": \"The daily color digest\"}, {\"url\": \"http://flatuicolors.com/\", \"img\": \"../assets/images/logos/flatuicolors.png\", \"name\": \"flatuicolors\", \"mark\": \"Copy Paste Color Pallette from Flat UI Theme\"}, {\"url\": \"https://coolors.co/\", \"img\": \"../assets/images/logos/coolors.png\", \"name\": \"coolors\", \"mark\": \"The super fast color schemes generator!\"}, {\"url\": \"http://www.colorhunt.co/\", \"img\": \"../assets/images/logos/colorhunt.png\", \"name\": \"colorhunt\", \"mark\": \"Beautiful Color Palettes\"}, {\"url\": \"https://color.adobe.com/zh/create/color-wheel\", \"img\": \"../assets/images/logos/AdobeColorCC.png\", \"name\": \"Adobe Color CC\", \"mark\": \"Create color schemes with the color wheel or browse thousands of color combinations from the Color community.\"}, {\"url\": \"http://www.flatuicolorpicker.com/\", \"img\": \"../assets/images/logos/flatuicolorpicker.png\", \"name\": \"flatuicolorpicker\", \"mark\": \"Best Flat Colors For UI Design\"}, {\"url\": \"http://qrohlf.com/trianglify-generator/\", \"img\": \"../assets/images/logos/trianglify.png\", \"name\": \"trianglify\", \"mark\": \"Trianglify Generator\"}, {\"url\": \"https://klart.co/colors/\", \"img\": \"../assets/images/logos/klart.png\", \"name\": \"klart\", \"mark\": \"Beautiful colors and designs to your inbox every week\"}, {\"url\": \"http://www.vanschneider.com/colors\", \"img\": \"../assets/images/logos/vanschneider.png\", \"name\": \"vanschneider\", \"mark\": \"Color Claim was created in 2012 by Tobias van Schneider with the goal to collect & combine unique colors for my future projects.\"}]}, {\"name\": \"\\u5728\\u7ebf\\u5de5\\u5177\", \"rows\": [{\"url\": \"https://tinypng.com/\", \"img\": \"../assets/images/logos/tinypng.png\", \"name\": \"tinypng\", \"mark\": \"Optimize your images with a perfect balance in quality and file size.\"}, {\"url\": \"http://goqr.me/\", \"img\": \"../assets/images/logos/goqr.png\", \"name\": \"goqr\", \"mark\": \"create QR codes for free (Logo, T-Shirt, vCard, EPS)\"}, {\"url\": \"https://ezgif.com\", \"img\": \"../assets/images/logos/ezgif.png\", \"name\": \"ezgif\", \"mark\": \"simple online GIF maker and toolset for basic animated GIF editing.\"}, {\"url\": \"http://inloop.github.io/shadow4android/\", \"img\": \"../assets/images/logos/Android9patch.png\", \"name\": \"Android 9 patch\", \"mark\": \"Android 9-patch shadow generator fully customizable shadows\"}, {\"url\": \"http://screensiz.es/\", \"img\": \"../assets/images/logos/screensizes.png\", \"name\": \"screen sizes\", \"mark\": \"Viewport Sizes and Pixel Densities for Popular Devices\"}, {\"url\": \"https://jakearchibald.github.io/svgomg/\", \"img\": \"../assets/images/logos/svgomg.png\", \"name\": \"svgomg\", \"mark\": \"SVG\\u5728\\u7ebf\\u538b\\u7f29\\u5e73\\u53f0\"}, {\"url\": \"https://www.gaoding.com\", \"img\": \"../assets/images/logos/gaoding.png\", \"name\": \"\\u7a3f\\u5b9a\\u62a0\\u56fe\", \"mark\": \"\\u514d\\u8d39\\u5728\\u7ebf\\u62a0\\u56fe\\u8f6f\\u4ef6,\\u56fe\\u7247\\u5feb\\u901f\\u6362\\u80cc\\u666f-\\u62a0\\u767d\\u5e95\\u56fe\"}]}, {\"name\": \"Chrome\\u63d2\\u4ef6\", \"rows\": [{\"url\": \"https://www.wappalyzer.com/\", \"img\": \"../assets/images/logos/wappalyzer.png\", \"name\": \"wappalyzer\", \"mark\": \"Identify technology on websites\"}, {\"url\": \"http://usepanda.com/\", \"img\": \"../assets/images/logos/usepanda.png\", \"name\": \"Panda\", \"mark\": \"A smart news reader built for productivity.\"}, {\"url\": \"https://sizzy.co/\", \"img\": \"../assets/images/logos/sizzy.png\", \"name\": \"sizzy\", \"mark\": \"A tool for developing responsive websites crazy-fast\"}, {\"url\": \"https://csspeeper.com/\", \"img\": \"../assets/images/logos/csspeeper.png\", \"name\": \"csspeeper\", \"mark\": \"Smart CSS viewer tailored for Designers.\"}, {\"url\": \"http://insight.io/\", \"img\": \"../assets/images/logos/insight.png\", \"name\": \"insight\", \"mark\": \"IDE-like code search and navigation, on the cloud\"}, {\"url\": \"http://mustsee.earth/\", \"img\": \"../assets/images/logos/mustsee.png\", \"name\": \"mustsee\", \"mark\": \"Discover the world's most beautiful places at every opened tab.\"}]}, {\"name\": \"\\u8bbe\\u8ba1\\u89c4\\u8303\", \"rows\": [{\"url\": \"http://designguidelines.co/\", \"img\": \"../assets/images/logos/designguidelines.png\", \"name\": \"Design Guidelines\", \"mark\": \"Design Guidelines \\u2014 The way products are built.\"}, {\"url\": \"https://github.com/alexpate/awesome-design-systems\", \"img\": \"../assets/images/logos/awesome_design_systems.png\", \"name\": \"Awesome design systems\", \"mark\": \" A collection of awesome design systems\"}, {\"url\": \"https://material.io/guidelines/\", \"img\": \"../assets/images/logos/Material_Design.png\", \"name\": \"Material Design\", \"mark\": \"Introduction - Material Design\"}, {\"url\": \"https://developer.apple.com/ios/human-interface-guidelines\", \"img\": \"../assets/images/logos/human_interface_guidelines.png\", \"name\": \"Human Interface Guidelines\", \"mark\": \"Human Interface Guidelines iOS\"}, {\"url\": \"http://viggoz.com/photoshopetiquette/\", \"img\": \"../assets/images/logos/photoshopetiquette.png\", \"name\": \"Photoshop Etiquette\", \"mark\": \"PS\\u793c\\u4eea-WEB\\u8bbe\\u8ba1\\u6307\\u5357\"}]}, {\"name\": \"\\u89c6\\u9891\\u6559\\u7a0b\", \"rows\": [{\"url\": \"http://www.photoshoplady.com/\", \"img\": \"../assets/images/logos/PhotoshopLady.png\", \"name\": \"Photoshop Lady\", \"mark\": \"Your Favourite Photoshop Tutorials in One Place\"}, {\"url\": \"http://doyoudo.com/\", \"img\": \"../assets/images/logos/doyoudo.png\", \"name\": \"doyoudo\", \"mark\": \"\\u521b\\u610f\\u8bbe\\u8ba1\\u8f6f\\u4ef6\\u5b66\\u4e60\\u5e73\\u53f0\"}, {\"url\": \"http://www.c945.com/web-ui-tutorial/\", \"img\": \"../assets/images/logos/web_ui_tutorial.png\", \"name\": \"\\u6ca1\\u4f4d\\u9053\", \"mark\": \"WEB UI\\u514d\\u8d39\\u89c6\\u9891\\u516c\\u5f00\\u8bfe\"}, {\"url\": \"https://www.imooc.com/\", \"img\": \"../assets/images/logos/imooc.png\", \"name\": \"\\u6155\\u8bfe\\u7f51\", \"mark\": \"\\u7a0b\\u5e8f\\u5458\\u7684\\u68a6\\u5de5\\u5382\\uff08\\u6709UI\\u8bfe\\u7a0b\\uff09\"}]}, {\"name\": \"\\u8bbe\\u8ba1\\u6587\\u7ae0\", \"rows\": [{\"url\": \"http://www.uisdc.com/\", \"img\": \"../assets/images/logos/uisdc.png\", \"name\": \"\\u4f18\\u8bbe\\u7f51\", \"mark\": \"\\u8bbe\\u8ba1\\u5e08\\u4ea4\\u6d41\\u5b66\\u4e60\\u5e73\\u53f0\"}, {\"url\": \"https://webdesignledger.com\", \"img\": \"../assets/images/logos/webdesignledger.png\", \"name\": \"Web Design Ledger\", \"mark\": \"Web Design Blog\"}, {\"url\": \"https://medium.com/\", \"img\": \"../assets/images/logos/medium.png\", \"name\": \"Medium\", \"mark\": \"Read, write and share stories that matter\"}]}, {\"name\": \"\\u8bbe\\u8ba1\\u7535\\u53f0\", \"rows\": [{\"url\": \"http://uxcoffee.co/\", \"img\": \"../assets/images/logos/uxcoffee.png\", \"name\": \"UX Coffee\", \"mark\": \"\\u300aUX Coffee \\u8bbe\\u8ba1\\u5496\\u300b\\u662f\\u4e00\\u6863\\u5173\\u4e8e\\u7528\\u6237\\u4f53\\u9a8c\\u7684\\u64ad\\u5ba2\\u8282\\u76ee\\u3002\\u6211\\u4eec\\u9080\\u8bf7\\u6765\\u81ea\\u7845\\u8c37\\u548c\\u56fd\\u5185\\u7684\\u5b66\\u8005\\u548c\\u804c\\u4eba\\u6765\\u804a\\u804a\\u300c\\u4ea7\\u54c1\\u8bbe\\u8ba1\\u300d\\u3001\\u300c\\u7528\\u6237\\u4f53\\u9a8c\\u300d\\u548c\\u300c\\u4e2a\\u4eba\\u6210\\u957f\\u300d\\u3002\"}, {\"url\": \"https://anyway.fm/\", \"img\": \"../assets/images/logos/anyway.png\", \"name\": \"Anyway.FM\", \"mark\": \"\\u8bbe\\u8ba1\\u6742\\u8c08 \\u2022 UI \\u8bbe\\u8ba1\\u5e08 JJ \\u548c Leon \\u4e3b\\u64ad\\u7684\\u8bbe\\u8ba1\\u64ad\\u5ba2\"}, {\"url\": \"https://www.yineng.fm\", \"img\": \"../assets/images/logos/yineng.png\", \"name\": \"\\u5f02\\u80fd\\u7535\\u53f0\", \"mark\": \"\\u5c06\\u5168\\u5b87\\u5b99\\u8bbe\\u8ba1\\u5e08\\u7684\\u6545\\u4e8b\\u8bb2\\u7ed9\\u4f60\\u542c\\u3002\"}]}, {\"name\": \"\\u4ea4\\u4e92\\u8bbe\\u8ba1\", \"rows\": [{\"url\": \"http://littlebigdetails.com/\", \"img\": \"../assets/images/logos/littlebigdetails.png\", \"name\": \"Little Big Details\", \"mark\": \"Little Big Details is a curated collection of the finer details of design, updated every day. \"}, {\"url\": \"https://www.smashingmagazine.com/category/user-experience\", \"img\": \"../assets/images/logos/smashingmagazine.png\", \"name\": \"Smashing Magazine\", \"mark\": \"Below you\\u2019ll find the best tips to take not only your UX design process but also the experiences you craft to the next level.\"}, {\"url\": \"https://www.nngroup.com/articles/\", \"img\": \"../assets/images/logos/nngroup.png\", \"name\": \"nngroup\", \"mark\": \"Evidence-Based User Experience Research, Training, and Consulting\"}, {\"url\": \"http://boxesandarrows.com/\", \"img\": \"../assets/images/logos/boxesandarrows.png\", \"name\": \"Boxes and Arrows\", \"mark\": \"Boxes and Arrows is devoted to the practice, innovation, and discussion of design; including graphic design, interaction design, information architecture and the design of business. \"}, {\"url\": \"http://uxdesignweekly.com/\", \"img\": \"../assets/images/logos/uxdesignweekly.png\", \"name\": \"UX Design Weekly\", \"mark\": \" get a hand picked list of the best user experience design links every week. \"}, {\"url\": \"http://uxren.cn/\", \"img\": \"../assets/images/logos/uxren.png\", \"name\": \"UX Ren\", \"mark\": \"\\u7528\\u6237\\u4f53\\u9a8c\\u4eba\\u7684\\u4e13\\u4e1a\\u793e\\u533a\"}, {\"url\": \"https://www.gulusucai.com/\", \"img\": \"../assets/images/logos/gulusucai.png\", \"name\": \"\\u5495\\u565c\\u7d20\\u6750\", \"mark\": \"\\u8d28\\u91cf\\u5f88\\u9ad8\\u7684\\u8bbe\\u8ba1\\u7d20\\u6750\\u7f51\\u7ad9\\uff08\\u826f\\u5fc3\\u63a8\\u8350\\uff09\"}]}, {\"name\": \"UED\\u56e2\\u961f\", \"rows\": [{\"url\": \"https://airbnb.design\", \"img\": \"../assets/images/logos/AirbnbDesign.png\", \"name\": \"Airbnb Design\", \"mark\": \"Airbnb Design\"}, {\"url\": \"http://facebook.design/\", \"img\": \"../assets/images/logos/FacebookDesign.png\", \"name\": \"Facebook Design\", \"mark\": \"Facebook Design\"}, {\"url\": \"https://design.google/\", \"img\": \"../assets/images/logos/GoogleDesign.png\", \"name\": \"Google Design\", \"mark\": \"Google Design\"}, {\"url\": \"http://eicodesign.com/\", \"img\": \"../assets/images/logos/eico.png\", \"name\": \"eico design\", \"mark\": \"\\u6570\\u5b57\\u5316\\u54a8\\u8be2\\u4e0e\\u4ea7\\u54c1\\u4e13\\u5bb6\"}, {\"url\": \"http://www.niceui.cn/\", \"img\": \"../assets/images/logos/niceui.png\", \"name\": \"nice design\", \"mark\": \"nicedesign\\u5948\\u601d\\u8bbe\\u8ba1\\u662f\\u9886\\u5148\\u7684\\u7528\\u6237\\u4f53\\u9a8c\\u8bbe\\u8ba1\\u4e0e\\u4e92\\u8054\\u7f51\\u54c1\\u724c\\u5efa\\u8bbe\\u516c\\u53f8\"}, {\"url\": \"http://cdc.tencent.com/\", \"img\": \"../assets/images/logos/cdc.png\", \"name\": \"\\u817e\\u8bafCDC\", \"mark\": \"\\u817e\\u8bafCDC\\u5173\\u6ce8\\u4e8e\\u4e92\\u8054\\u7f51\\u89c6\\u89c9\\u8bbe\\u8ba1\\u3001\\u4ea4\\u4e92\\u8bbe\\u8ba1\\u3001\\u7528\\u6237\\u7814\\u7a76\\u3001\\u524d\\u7aef\\u5f00\\u53d1\\u3002\"}, {\"url\": \"http://tgideas.qq.com/\", \"img\": \"../assets/images/logos/tgideas.png\", \"name\": \"TGideas\", \"mark\": \"TGideas\\u96b6\\u5c5e\\u4e8e\\u817e\\u8baf\\u516c\\u53f8\\u4e92\\u52a8\\u5a31\\u4e50\\u4e1a\\u52a1\\u7cfb\\u7edf\\u7684\\u4e13\\u4e1a\\u63a8\\u5e7f\\u7c7b\\u8bbe\\u8ba1\\u56e2\\u961f\"}, {\"url\": \"https://isux.tencent.com/\", \"img\": \"../assets/images/logos/isux.png\", \"name\": \"ISUX\", \"mark\": \"\\u817e\\u8baf\\u793e\\u4ea4\\u7528\\u6237\\u4f53\\u9a8c\\u8bbe\\u8ba1\\u90e8\"}, {\"url\": \"http://mxd.tencent.com/\", \"img\": \"../assets/images/logos/mxd.png\", \"name\": \"MXD\", \"mark\": \"\\u817e\\u8bafMIG\\u65e0\\u7ebf\\u4e92\\u8054\\u7f51\\u4e8b\\u4e1a\\u7fa4\\u8bbe\\u8ba1\\u56e2\\u961f\"}, {\"url\": \"http://www.aliued.com/\", \"img\": \"../assets/images/logos/aliued.png\", \"name\": \"Aliued\", \"mark\": \"\\u963f\\u91cc\\u5df4\\u5df4\\u56fd\\u9645UED\\u56e2\\u961f\"}, {\"url\": \"http://www.aliued.cn/\", \"img\": \"../assets/images/logos/aliuedcn.png\", \"name\": \"U\\u4e00\\u70b9\", \"mark\": \"\\u963f\\u91cc\\u5df4\\u5df4\\uff08\\u4e2d\\u56fd\\u7ad9\\uff09\\u7528\\u6237\\u4f53\\u9a8c\\u8bbe\\u8ba1\\u90e8\\u535a\\u5ba2U\\u4e00\\u70b9\\u8bbe\\u8ba1 UED\\u56e2\\u961f\"}, {\"url\": \"http://uedc.163.com/\", \"img\": \"../assets/images/logos/uedc.png\", \"name\": \"\\u7f51\\u6613uedc\", \"mark\": \"\\u7f51\\u6613\\u7528\\u6237\\u4f53\\u9a8c\\u8bbe\\u8ba1\\u4e2d\\u5fc3\\uff08User Experience Design Center\\uff09\"}, {\"url\": \"http://ued.baidu.com/\", \"img\": \"../assets/images/logos/uedbaidu.png\", \"name\": \"\\u767e\\u5ea6\\u7528\\u6237\\u4f53\\u9a8c\\u4e2d\\u5fc3\", \"mark\": \"\\u767e\\u5ea6\\u7528\\u6237\\u4f53\\u9a8c\\u4e2d\\u5fc3\"}, {\"url\": \"http://jdc.jd.com/\", \"img\": \"../assets/images/logos/jdc.png\", \"name\": \"\\u4eac\\u4e1c\\u8bbe\\u8ba1\\u4e2d\\u5fc3\", \"mark\": \"\\u4eac\\u4e1c\\u8bbe\\u8ba1\\u4e2d\\u5fc3\"}, {\"url\": \"http://eux.baidu.com/\", \"img\": \"../assets/images/logos/euxbaidu.png\", \"name\": \"\\u767e\\u5ea6\\u4f01\\u4e1a\\u4ea7\\u54c1\\u7528\\u6237\\u4f53\\u9a8c\\u4e2d\\u5fc3\", \"mark\": \"\\u767e\\u5ea6\\u4f01\\u4e1a\\u4ea7\\u54c1\\u7528\\u6237\\u4f53\\u9a8c\\u4e2d\\u5fc3\"}, {\"url\": \"http://ued.ctrip.com/\", \"img\": \"../assets/images/logos/ctrip.png\", \"name\": \"\\u643a\\u7a0b\\u8bbe\\u8ba1\\u59d4\\u5458\\u4f1a\", \"mark\": \"\\u643a\\u7a0b\\u8bbe\\u8ba1\\u59d4\\u5458\\u4f1a-Ctrip Design Committee\"}]}]}")
		err = ioutil.WriteFile("./json/webstack.json", d1, 0666)
		if err != nil {
			fmt.Print("初始化登陆文件webstack.json错误：")
			fmt.Println(err)
			return
		}
	}
	err = LoadJsonFile("./json/login.json", &Login)
	if err != nil {
		fmt.Print("加载登陆文件login.json错误：")
		fmt.Println(err)
		return
	}

	err = LoadJsonFile("./json/config.json", &Config)
	if err != nil {
		fmt.Print("加载配置文件config.json错误：")
		fmt.Println(err)
		return
	} else {
		if Config.Port <= 0 || Config.Port > 65535 {
			Config.Port = 2802
		}
		err = SaveJsonFile("./json/config.json", &Config)
	}

	err = LoadJsonFile("./json/webstack.json", &WebStack)
	if err != nil {
		fmt.Print("加载页面文件webstack.json错误： ")
		fmt.Println(err)
		return
	}

	r := gin.Default()
	r.Static("/assets", "./public")
	r.LoadHTMLGlob("views/**/*")
	r.GET("/", GetIndex)
	r.GET("/index.html", GetIndex)
	r.GET("/about.html", GetAbout)
	r.GET(Login.Path, GetLogin)
	r.POST(Login.Path, PostLogin)
	r.Use(AuthMiddleWare())
	{
		r.GET("/admin", AuthMiddleWare(), GetAdmin)
		r.POST("/admin", AuthMiddleWare(), PostAdmin)
		r.POST("/admin/upload", AuthMiddleWare(), PostAdminUpload)
		r.POST("/admin/uploadfile", AuthMiddleWare(), PostAdminUploadFile)
	}

	r.Run(fmt.Sprintf("%s:%d", Config.Url, Config.Port))
}

func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"config":   Config,
		"webstack": WebStack,
		//"body": template.HTML("<body>I 'm body<body>"),
	})
}

func GetAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "index/about.html", gin.H{
		"config": Config,
	})
}

func GetAdmin(c *gin.Context) {
	cmd := c.DefaultQuery("cmd", "null")
	switch cmd {
	case "logout":
		c.SetCookie("webstackgo_token", "", -1, "/", strings.Split(c.Request.Host, ":")[0], false, true)
		c.JSON(http.StatusOK, gin.H{
			"cmd":     cmd,
			"message": "退出登陆成功",
			"error":   0,
		})
	case "webstack.json":
		c.JSON(http.StatusOK, WebStack)
	case "menu.json":
		c.JSON(http.StatusOK, WebStack.Menu)
	case "class.json":
		c.JSON(http.StatusOK, WebStack.Class)
	default:
		c.HTML(http.StatusOK, "admin/index.html", gin.H{
			"login":    Login,
			"config":   Config,
			"webstack": WebStack,
		})
	}
}

func PostAdmin(c *gin.Context) {
	cmd := c.DefaultQuery("cmd", "null")
	jsonMap := make(map[string]string)
	c.BindJSON(&jsonMap)
	ret := gin.H{
		"message": "OK",
		"error":   0,
	}
	var ok bool
	switch cmd {
	case "login_path":
		if _, ok = jsonMap["path"]; !ok {
			ret["message"] = "无效数据"
			ret["error"] = 100
		} else if len(jsonMap["path"]) < 2 {
			ret["message"] = "登陆入口不可为空"
			ret["error"] = 101
		} else if string([]byte(jsonMap["path"])[:1]) != "/" {
			ret["message"] = "登陆入口格式错误，必须以/开头。"
			ret["error"] = 102
		} else {
			Login.Path = jsonMap["path"]
			err := SaveJsonFile("./json/login.json", &Login)
			if err == nil {
				ret["message"] = "登陆入口修改成功，重启WebStaskGo服务后生效。"
				ret["error"] = 0
			} else {
				ret["message"] = err.Error()
				ret["error"] = 103
			}
		}
		c.JSON(http.StatusOK, ret)
	case "user":
		if IsJsonKey(jsonMap, "username") && IsJsonKey(jsonMap, "password") && IsJsonKey(jsonMap, "password2") {
			jsonMap["username"] = strings.TrimSpace(jsonMap["username"])
			jsonMap["password"] = strings.TrimSpace(jsonMap["password"])
			jsonMap["password2"] = strings.TrimSpace(jsonMap["password2"])
			if len(jsonMap["username"]) < 2 {
				ret["key"] = "username"
				ret["message"] = "登陆账号太短，请输入大于2个字符。"
				ret["error"] = 111
			} else if len(jsonMap["password"]) < 6 && jsonMap["password"] != "" {
				ret["key"] = "password"
				ret["message"] = "登陆密码太短，请输入大于6个字符。"
				ret["error"] = 111
			} else if jsonMap["password"] != jsonMap["password2"] {
				ret["key"] = "password2"
				ret["message"] = "确认密码与登陆密码不相同，请重新输入。"
				ret["error"] = 111
			} else {
				message := ""
				if Login.Username != jsonMap["username"] {
					Login.Username = jsonMap["username"]
					message += "登陆账号修改完成，"
				}
				if jsonMap["password"] != "" && Login.Password != GetMD5(jsonMap["password"]) {
					Login.Password = GetMD5(jsonMap["password"])
					message += "登陆密码修改完成，"
				}
				err := SaveJsonFile("./json/login.json", &Login)
				if err == nil {
					ret["message"] = message + "请重新前往登陆页面登陆。"
					ret["error"] = 0
				} else {
					ret["message"] = err.Error()
					ret["error"] = 112
				}
			}
		} else {
			ret["message"] = "缺少有效数据"
			ret["error"] = 110
		}
		c.JSON(http.StatusOK, ret)
	case "stack":
		if IsJsonKey(jsonMap, "title") {
			Config.Title = jsonMap["title"]
		}
		if IsJsonKey(jsonMap, "description") {
			Config.Description = jsonMap["description"]
		}
		if IsJsonKey(jsonMap, "keywords") {
			Config.Keywords = jsonMap["keywords"]
		}
		if IsJsonKey(jsonMap, "recordcode") {
			Config.Recordcode = jsonMap["recordcode"]
		}
		if IsJsonKey(jsonMap, "footer") {
			Config.Footer = jsonMap["footer"]
		}
		if IsJsonKey(jsonMap, "url") {
			Config.Url = strings.TrimSpace(jsonMap["url"])
		}
		if IsJsonKey(jsonMap, "port") {
			if port, err := strconv.Atoi(jsonMap["port"]); err != nil && port > 0 && port < 65535 {
				Config.Port = port
			}
		}
		err := SaveJsonFile("./json/config.json", &Config)
		if err == nil {
			ret["message"] = "网页设置保存完成"
			ret["error"] = 0
		} else {
			ret["message"] = err.Error()
			ret["error"] = 122
		}
		c.JSON(http.StatusOK, ret)
	case "web-add":
		if IsJsonKey(jsonMap, "class1_name") && IsJsonKey(jsonMap, "class2_name") {
			// 新增类别menu后，网址class中没有改类别网址，导致GetClassId返回-1
			// 先判断class1_name 是否存在
			class1name := jsonMap["class1_name"]
			class2name := jsonMap["class2_name"]
			classIndex := -1
			if IsExistInMenu(class1name) {
				if class2name != "" {
					if IsExistInSubMenu(class2name) {
						// 检查class中是否已存在class2_name
						classIndex = IsExistInClass(class2name)
						if classIndex > -1 {
							// 存在该class2name的class，直接在row中增加
							if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "url") && IsJsonKey(jsonMap, "mark") && IsJsonKey(jsonMap, "img") {
								if AddWebData(classIndex, jsonMap) {
									if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
										ret["message"] = "添加网址成功"
										ret["error"] = 0
									} else {
										ret["message"] = err.Error()
										ret["error"] = 133
									}
								} else {
									ret["message"] = "无效的分类名称"
									ret["error"] = 132
								}
							} else {
								ret["message"] = "上报数据不完整"
								ret["error"] = 131
							}
						} else {
							// 没有该class2name的class，需要新增class，同时增加改row
							if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "url") && IsJsonKey(jsonMap, "mark") && IsJsonKey(jsonMap, "img") {
								if AddNewClass(class2name, jsonMap) {
									if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
										ret["message"] = "添加网址成功"
										ret["error"] = 0
									} else {
										ret["message"] = err.Error()
										ret["error"] = 133
									}
								} else {
									ret["message"] = "无效的分类名称"
									ret["error"] = 132
								}
							} else {
								ret["message"] = "上报数据不完整"
								ret["error"] = 131
							}
						}
					} else {
						ret["message"] = "无效的子分类名称"
						ret["error"] = 132
					}
				} else {
					// 检查class中是否已存在class1name
					classIndex = IsExistInClass(class1name)
					if classIndex > -1 {
						// 存在该class1name的class，直接在row中增加
						if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "url") && IsJsonKey(jsonMap, "mark") && IsJsonKey(jsonMap, "img") {
							if AddWebData(classIndex, jsonMap) {
								if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
									ret["message"] = "添加网址成功"
									ret["error"] = 0
								} else {
									ret["message"] = err.Error()
									ret["error"] = 133
								}
							} else {
								ret["message"] = "无效的分类名称"
								ret["error"] = 132
							}
						} else {
							ret["message"] = "上报数据不完整"
							ret["error"] = 131
						}
					} else {
						// 没有该class1name的class，需要新增class，同时增加改row
						if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "url") && IsJsonKey(jsonMap, "mark") && IsJsonKey(jsonMap, "img") {
							if AddNewClass(class1name, jsonMap) {
								if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
									ret["message"] = "添加网址成功"
									ret["error"] = 0
								} else {
									ret["message"] = err.Error()
									ret["error"] = 133
								}
							} else {
								ret["message"] = "无效的分类名称"
								ret["error"] = 132
							}
						} else {
							ret["message"] = "上报数据不完整"
							ret["error"] = 131
						}
					}
				}
			} else {
				ret["message"] = "无效的主分类名称"
				ret["error"] = 132
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 130
		}
		c.JSON(http.StatusOK, ret)
	case "web-edit":
		if IsJsonKey(jsonMap, "index") && IsJsonKey(jsonMap, "class1_name") && IsJsonKey(jsonMap, "class2_name") {
			classid := GetClassId(jsonMap["class1_name"], jsonMap["class2_name"])
			//fmt.Println(classid, WebStack.Class[classid])
			if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "url") && IsJsonKey(jsonMap, "mark") && IsJsonKey(jsonMap, "img") {
				if EditWebData(classid, jsonMap) {
					if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
						ret["message"] = "编辑网址成功"
						ret["error"] = 0
					} else {
						ret["message"] = err.Error()
						ret["error"] = 143
					}
				} else {
					ret["message"] = "无效的网址源信息"
					ret["error"] = 142
				}
			} else {
				ret["message"] = "上报数据不完整"
				ret["error"] = 141
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 140
		}
		c.JSON(http.StatusOK, ret)
	case "web-delete":
		isWebDeleteOk := false
		if IsJsonKey(jsonMap, "index") {
			classid, webid := WebIndex2ID(jsonMap["index"])
			if DeleteWebData(classid, webid) {
				isWebDeleteOk = true
			} else {
				ret["message"] = "无效的网址源信息"
				ret["error"] = 151
			}
		} else if IsJsonKey(jsonMap, "indexArray") {
			var indexArray []string
			err := json.Unmarshal([]byte(jsonMap["indexArray"]), &indexArray)
			if err == nil {
				for i := 0; i < len(indexArray); i++ {
					classid, webid := WebIndex2ID(indexArray[i])
					DeleteWebData(classid, webid)
				}
				isWebDeleteOk = true
			} else {
				ret["message"] = "批量删除数据结构错误。"
				ret["error"] = 153
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 150
		}
		if isWebDeleteOk == true {
			if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
				ret["message"] = "删除网址成功"
				ret["error"] = 0
			} else {
				ret["message"] = err.Error()
				ret["error"] = 152
			}
		}
		c.JSON(http.StatusOK, ret)
	case "class-add":
		if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "icon") && IsJsonKey(jsonMap, "class_up") && IsJsonKey(jsonMap, "class_id") {
			classup, _ := strconv.Atoi(jsonMap["class_up"])
			if CheckClassName(jsonMap["name"], "") == false {
				ret["message"] = "分类名称冲突，请更改分类名称。"
				ret["error"] = 162
			} else if AddClassData(classup, jsonMap["name"], jsonMap["icon"]) {
				if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
					ret["message"] = "添加分类成功"
					ret["error"] = 0
				} else {
					ret["message"] = err.Error()
					ret["error"] = 163
				}
			} else {
				ret["message"] = "上报数据参数错误"
				ret["error"] = 161
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 160
		}
		c.JSON(http.StatusOK, ret)
	case "class-edit":
		if IsJsonKey(jsonMap, "name") && IsJsonKey(jsonMap, "icon") && IsJsonKey(jsonMap, "class_up") && IsJsonKey(jsonMap, "class_id") {
			classup, _ := strconv.Atoi(jsonMap["class_up"])
			if !CheckClassName(jsonMap["name"], jsonMap["class_id"]) {
				ret["message"] = "分类名称冲突，请更改分类名称。"
				ret["error"] = 172
			} else if EditClassData(jsonMap["class_id"], classup, jsonMap["name"], jsonMap["icon"]) {
				if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
					ret["message"] = "编辑分类成功"
					ret["error"] = 0
				} else {
					ret["message"] = err.Error()
					ret["error"] = 173
				}
			} else {
				ret["message"] = "上报数据参数错误"
				ret["error"] = 171
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 170
		}
		c.JSON(http.StatusOK, ret)
	case "class-delete":
		isClassDeleteOk := false
		if IsJsonKey(jsonMap, "index") {
			if DeleteClassData(jsonMap["index"]) {
				isClassDeleteOk = true
			} else {
				ret["message"] = "无效的网址源信息"
				ret["error"] = 181
			}
		}
		if isClassDeleteOk == true {
			if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
				ret["message"] = "删除网址成功"
				ret["error"] = 0
			} else {
				ret["message"] = err.Error()
				ret["error"] = 182
			}
		}
		c.JSON(http.StatusOK, ret)
	case "class-sort":
		if IsJsonKey(jsonMap, "webStack") {
			if LoadJsonString([]byte(jsonMap["webStack"]), &WebStack) == nil {
				if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
					ret["message"] = "保存分类排序成功"
					ret["error"] = 0
				} else {
					ret["message"] = err.Error()
					ret["error"] = 193
				}
			} else {
				ret["message"] = "上报数据格式错误"
				ret["error"] = 191
			}
		} else {
			ret["message"] = "上报数据不完整"
			ret["error"] = 190
		}
		c.JSON(http.StatusOK, ret)
	case "web-sort-by-class":
		// 检查[]class中的name是否有不在nemu中的

		// 根据当前WebStack中的Menu及子Menu顺序，排序Class clice 中的元素顺序
		if WebSortByClass(WebStack) {
			// 保存webstack
			if err := SaveJsonFile("./json/webstack.json", &WebStack); err == nil {
				ret["message"] = "保存网址排序成功"
				ret["error"] = 0
			} else {
				ret["message"] = err.Error()
				ret["error"] = 195
			}
		} else {
			ret["message"] = "网页依据Menu排序失败"
			ret["error"] = 194
		}
		c.JSON(http.StatusOK, ret)
	default:
		c.JSON(http.StatusFound, gin.H{
			"message": "Error 302",
			"error":   302,
		})
	}
}

func AddNewClass(classname string, classData map[string]string) bool {
	if len(WebStack.Class) < 0 {
		return false
	}
	WebStack.Class = append(WebStack.Class, JsClass{
		Name: classname,
		Rows: []JsClassItem{{
			Url:  classData["url"],
			Img:  classData["img"],
			Name: classData["name"],
			Mark: classData["mark"],
		}},
	})
	return true
}

func IsExistInClass(classname string) int {
	var index int = -1
	for i, class := range WebStack.Class {
		if class.Name == classname {
			index = i
			break
		}
	}
	return index
}

func IsExistInMenu(menuname string) bool {
	var b bool = false
	for _, menu := range WebStack.Menu {
		if menu.Name == menuname {
			b = true
			break
		}
	}
	return b
}

func IsExistInSubMenu(submenuname string) bool {
	b := false
	for _, menu := range WebStack.Menu {
		if len(menu.Sub) > 0 {
			for _, submenu := range menu.Sub {
				if submenu.Name == submenuname {
					b = true
					break
				}
			}
		}
	}
	return b
}

func PostAdminUpload(c *gin.Context) {
	// https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"error":   801,
		})
		return
	}

	//获取文件后缀
	existing := strings.ToLower(Ext(file.Filename))
	if existing == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "文件类型错误，无法上传。",
			"error":   802,
		})
		return
	}
	extStrSlice := []string{".jpg", ".png", "gif"}
	if !ContainArray(existing, extStrSlice) {
		c.JSON(http.StatusOK, gin.H{
			"message": "文件类型错误，请上传图片文件（jpg、png、gif）。",
			"error":   803,
		})
		return
	}
	u1 := uuid.NewV4()
	resetfilename := u1.String()

	filepath := "public/images/uploads/"
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}
	path := filepath + resetfilename + existing //路径+文件名上传

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"error":   804,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     "../assets/images/uploads/" + resetfilename + existing,
		"message": "upload file success.",
		"error":   0,
	})
}

func PostAdminUploadFile(c *gin.Context) {
	// https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"error":   801,
		})
		return
	}

	//获取文件后缀
	existing := strings.ToLower(Ext(file.Filename))
	if existing == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "文件类型错误，无法上传。",
			"error":   802,
		})
		return
	}
	u1 := uuid.NewV4()
	resetfilename := u1.String()

	filepath := "public/images/fileuploads/"
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}
	path := filepath + resetfilename + existing //路径+文件名上传

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"error":   804,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     "assets/images/fileuploads/" + resetfilename + existing,
		"message": "upload file success.",
		"error":   0,
	})
}

func GetLogin(c *gin.Context) {
	if c.FullPath() == Login.Path {
		c.HTML(http.StatusOK, "admin/login.html", gin.H{
			"config":  Config,
			"success": false,
			"message": "",
		})
	} else {
		c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{
			"error":   401,
			"message": "The login page has been modified.",
		})
	}
}

func PostLogin(c *gin.Context) {
	username := strings.TrimSpace(c.DefaultPostForm("username", ""))
	password := GetMD5(strings.TrimSpace(c.DefaultPostForm("password", "webstackgo")))
	if username == Login.Username && password == Login.Password {
		now := time.Now()
		token := GetToken(username, password, now.Unix())
		fmt.Println(token, now)
		c.SetCookie("webstackgo_token", token, 7200, "/", strings.Split(c.Request.Host, ":")[0], false, true)
		c.HTML(http.StatusOK, "admin/login.html", gin.H{
			"config":  Config,
			"success": true,
			"message": "登陆成功！",
		})
	} else {
		c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{
			"config":  Config,
			"success": false,
			"message": "登陆失败：用户名或密码错误。",
		})
	}

}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("webstackgo_token"); err == nil {
			token, _ := url.QueryUnescape(cookie.Value)
			arr := strings.Split(token, "|")
			//fmt.Println(token, arr)
			if len(arr) == 2 {
				if intNow, err2 := strconv.ParseInt(arr[1], 10, 64); err2 == nil && token == GetToken(Login.Username, Login.Password, intNow) {
					if time.Now().Unix()-intNow < 3600 {
						token = GetToken(Login.Username, Login.Password, time.Now().Unix())
						c.SetCookie("webstackgo_token", token, 7200, "/", strings.Split(c.Request.Host, ":")[0], false, true)
					}
					c.Next()
					return
				}
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func IsJsonKey(m map[string]string, k string) bool {
	_, ret := m[k]
	return ret
}

func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetToken(username, password string, now int64) string {
	return fmt.Sprintf("%s|%d", GetMD5(fmt.Sprintf("%s|%s|%d", username, password, now)), now)
}

func SaveFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, os.ModePerm)
	return err
}

func SaveJsonFile(path string, obj interface{}) error {
	content, err := json.Marshal(obj)
	if err == nil {
		err = SaveFile(path, content)
	}
	return err
}

func LoadFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, err
}

func LoadJsonString(content []byte, obj interface{}) error {
	err := json.Unmarshal(content, obj)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func LoadJsonFile(path string, obj interface{}) error {
	content, err := LoadFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return LoadJsonString(content, obj)
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WebIndex2ID(index string) (int, int) {
	arrIndex := strings.Split(index, "-")
	classId := -1
	webId := -1
	if len(arrIndex) == 2 {
		num1, _ := strconv.Atoi(arrIndex[0])
		webId, _ = strconv.Atoi(arrIndex[1])
		if num1 >= 0 && num1 < len(WebStack.Menu) {
			classId = GetClassId(WebStack.Menu[num1].Name, "")
		}
	} else if len(arrIndex) == 3 {
		num1, _ := strconv.Atoi(arrIndex[0])
		num2, _ := strconv.Atoi(arrIndex[1])
		webId, _ = strconv.Atoi(arrIndex[2])
		if num1 >= 0 && num1 < len(WebStack.Menu) {
			classId = GetClassId(WebStack.Menu[num1].Name, WebStack.Menu[num1].Sub[num2].Name)
		}
	}
	return classId, webId
}

func ClassIndex2ID(index string) (int, int) {
	arrIndex := strings.Split(index, "-")
	classUp := -1
	classId := -1
	if len(arrIndex) == 2 {
		classUp, _ = strconv.Atoi(arrIndex[0])
		classId, _ = strconv.Atoi(arrIndex[1])
	} else if len(arrIndex) == 1 {
		classId, _ = strconv.Atoi(arrIndex[0])
	}
	return classUp, classId
}

func CheckClassName(name string, index string) bool {
	if len(strings.TrimSpace(name)) == 0 {
		return false
	}
	classup, classid := -1, -1
	if index != "" {
		classup, classid = ClassIndex2ID(index)
	}
	for id, menu := range WebStack.Menu {
		if menu.Name == name {
			if -1 == classup && id == classid {
				continue
			}
			return false
		} else {
			for subId, subMenu := range menu.Sub {
				if subMenu.Name == name {
					if id == classup && subId == classid {
						continue
					}
					return false
				}
			}
		}
	}
	return true
}

func GetClassIndex(name string) int {
	for index, item := range WebStack.Class {
		if item.Name == name {
			return index
		}
	}
	return -1
}

func GetClassId(name1 string, name2 string) int {
	index := 0
	name := ""
	for _, menu := range WebStack.Menu {
		if menu.Name == name1 {
			if len(menu.Sub) > 0 {
				for _, subMenu := range menu.Sub {
					if subMenu.Name == name2 {
						name = name2
					}
				}
			} else {
				name = name1
			}
			break
		}
		index += len(menu.Sub)
	}
	if name != "" {
		for ; index < len(WebStack.Class); index++ {
			if WebStack.Class[index].Name == name {
				return index
			}
		}
	}
	return -1
}

func AddWebData(classid int, classData map[string]string) bool {
	if classid < 0 || classid > len(WebStack.Class) {
		return false
	}
	WebStack.Class[classid].Rows = append(WebStack.Class[classid].Rows, JsClassItem{
		Url:  classData["url"],
		Img:  classData["img"],
		Name: classData["name"],
		Mark: classData["mark"],
	})
	return true
}

func DeleteWebData(classid int, webid int) bool {
	if classid >= 0 && webid >= 0 && classid < len(WebStack.Class) && webid < len(WebStack.Class[classid].Rows) {
		WebStack.Class[classid].Rows = append(WebStack.Class[classid].Rows[:webid], WebStack.Class[classid].Rows[webid+1:]...)
		return true
	} else {
		return false
	}
}

func EditWebData(classid int, classData map[string]string) bool {
	oldClassId, oldWebId := WebIndex2ID(classData["index"])
	if oldClassId >= 0 && oldWebId >= 0 {
		if oldClassId == classid {
			WebStack.Class[classid].Rows[oldWebId].Name = classData["name"]
			WebStack.Class[classid].Rows[oldWebId].Url = classData["url"]
			WebStack.Class[classid].Rows[oldWebId].Img = classData["img"]
			WebStack.Class[classid].Rows[oldWebId].Mark = classData["mark"]
			return true
		} else if AddWebData(classid, classData) {
			return DeleteWebData(oldClassId, oldWebId)
		}
	}
	return false
}

func AddClassData(classup int, classname string, classicon string) bool {
	if classup == -1 {
		WebStack.Menu = append(WebStack.Menu, JsMenu{
			Menu: "smooth",
			Name: classname,
			Icon: classicon,
			Sub:  []JsMenu{},
			Url:  "#" + classname,
		})
		return true
	} else if classup >= 0 && classup < len(WebStack.Menu) {
		WebStack.Menu[classup].Sub = append(WebStack.Menu[classup].Sub, JsMenu{
			Menu: "smooth",
			Name: classname,
			Icon: classicon,
			Sub:  []JsMenu{},
			Url:  "#" + classname,
		})
		return true
	} else {
		return false
	}

}

func DeleteClassData(classIndex string) bool {
	oldClassUp, oldClassId := ClassIndex2ID(classIndex)
	if oldClassUp == -1 {
		WebStack.Menu = append(WebStack.Menu[:oldClassId], WebStack.Menu[oldClassId+1:]...)
	} else if oldClassUp >= 0 && oldClassUp < len(WebStack.Menu) {
		WebStack.Menu[oldClassUp].Sub = append(WebStack.Menu[oldClassUp].Sub[:oldClassId], WebStack.Menu[oldClassUp].Sub[oldClassId+1:]...)
	} else {
		return false
	}
	return true
}

func EditClassData(classIndex string, classup int, classname string, classicon string) bool {
	oldClassUp, oldClassId := ClassIndex2ID(classIndex)
	if oldClassUp == classup {
		if oldClassUp == -1 {
			WebStack.Menu[oldClassId].Name = classname
			WebStack.Menu[oldClassId].Icon = classicon
		} else if oldClassUp >= 0 && oldClassUp < len(WebStack.Menu) {
			WebStack.Menu[oldClassUp].Sub[oldClassId].Name = classname
			WebStack.Menu[oldClassUp].Sub[oldClassId].Icon = classicon
		} else {
			return false
		}
	} else {
		if AddClassData(classup, classname, classicon) {
			return DeleteClassData(classIndex)
		} else {
			return false
		}
	}
	return true
}

//Contain 判断obj是否在target中，target支持的类型array,slice,map   false:不在 true:在
func ContainArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

//获取文件的扩展名
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// 取得所有menu+子menu的个数
func GetMenuCount(ws JsWebStack) int {
	count := 0
	for _, v := range ws.Menu {
		if len(v.Sub) > 0 {
			count += len(v.Sub)
		} else {
			count++
		}
	}
	return count
}

func GetIndexFromMenuSliceByName(name string, ms []string) int {
	for i, v := range ms {
		if v == name {
			return i
		} else {
			continue
		}
	}
	return -1
}

func WebSortByClass(ws JsWebStack) bool {
	// 取得所有menu+子menu的个数
	count := GetMenuCount(ws)
	if count == 0 {
		return false
	}
	// 根据当前menu排序，取得所有分类menu的name排序id
	var menuNameSlice = make([]string, 0, count)
	for _, v := range ws.Menu {
		if len(v.Sub) > 0 {
			for _, v1 := range v.Sub {
				menuNameSlice = append(menuNameSlice, v1.Name)
			}
		} else {
			menuNameSlice = append(menuNameSlice, v.Name)
		}
	}
	// 将 []JsClass中的name替换成menuNameSlice中对应的index
	for i, v := range ws.Class {
		index := GetIndexFromMenuSliceByName(v.Name, menuNameSlice)
		if index < 0 {
			return false
		} else {
			ws.Class[i].Name = strconv.Itoa(index)
		}
	}
	// 对ws.Class排序
	sort.Slice(ws.Class, func(i, j int) bool {
		return ws.Class[i].Name < ws.Class[j].Name
	})
	for i, v := range ws.Class {
		index, _ := strconv.Atoi(v.Name)
		ws.Class[i].Name = menuNameSlice[index]
	}
	return true
}
