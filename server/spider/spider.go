package spider

import (
	"blog-admin-api/core"
	"blog-admin-api/pkg/httplib"
	"bytes"
	"crypto/tls"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Img struct {
	ImgId  string
	Title  string
	Url    string
	Praise string
}

type CountRes struct {
	Success int
	Failed  int
	List    int
	Page    int
	Exist   int
}

type Conf struct {
	waitGroup *sync.WaitGroup
	maxCh     chan int
	count     CountRes
	cookie    string
}

func Get(c *core.Context, cook string) error {
	conf := new(Conf)
	conf.waitGroup = new(sync.WaitGroup)
	conf.maxCh = make(chan int, 100)
	conf.count = CountRes{}
	conf.cookie = cook

	return GetList(c, conf)
}

func GetList(c *core.Context, conf *Conf) error {
	bookmarkReq, _ := http.NewRequest("GET", bookmark, nil)
	bookmarkReq.Header.Set("cookie", conf.cookie)
	proxy, _ := url.Parse("socks5://127.0.0.1:1080")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	bookmarkResp, err := httplib.NewClient(httplib.WithTimeout(time.Second*60), httplib.WithTransport(tr)).Do(bookmarkReq)
	if err != nil {
		return err
	}
	if bookmarkResp == nil {
		return errors.New("bookmarkResp is nil")
	}

	var buf []byte
	buf, _ = io.ReadAll(bookmarkResp.Body)
	content := string(buf)
	c.Info("获取到bookmarkResp.Body", nil, content)

	allContent := content
	pageExpInfos, err := regexp.Compile(`w&amp;p=(\d+)[\s\S]*s="next"`)
	if err != nil {
		c.ErrorL("bookmarkResp正则匹配准备失败", nil, err.Error())
		return err
	}
	if len(pageExpInfos.FindStringSubmatch(content)) == 0 {
		return errors.New("pageExpInfos.FindStringSubmatch(content) == 0")
	}
	page, err := strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
	if err != nil {
		c.ErrorL("bookmarkResp正则匹配失败", content, err.Error())
		return err
	}
	if page == 0 {
		page = 1
	}
	p := 1
	for {
		if p > 1 {
			bookmarkReq, _ = http.NewRequest("GET", bookmark+"?rest=show&p="+strconv.Itoa(p), nil)
			bookmarkReq.Header.Set("cookie", conf.cookie)
			bookmarkResp, err = httplib.NewClient(httplib.WithTimeout(time.Second*60), httplib.WithTransport(tr)).Do(bookmarkReq)
			if bookmarkResp == nil {
				c.Info("获取bookmark返回nil", p, nil)
				return errors.New("bookmarkResp is nil")
			}
			buf, _ = io.ReadAll(bookmarkResp.Body)
			content = string(buf)
			c.Info("获取到bookmarkResp.Body", p, content)

			allContent += content
			pageExpInfos, _ = regexp.Compile(`w&amp;p=\d+[\s\S]*s="">(.+?)<[\s\S]*s="next"`)
			page, _ = strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
			if page == 0 {
				page = 1
			}
		}
		p = p + 1
		if p > page {
			break
		}
	}
	conf.count.Page = p
	conf.count.Success = 0
	conf.count.Failed = 0
	conf.count.Exist = 0
	defer func() {
		_ = bookmarkResp.Body.Close()
	}()

	size := (page + 1) * 20
	k := 0
	imgSlice := make([]Img, size)
	r, _ := regexp.Compile(`data-id="(.+?)".+?title="(.+?)".+?e"></i>(.+?)</a>`)
	imgExpInfos := r.FindAllStringSubmatch(allContent, size)

	conf.waitGroup.Add(len(imgExpInfos))
	for _, v := range imgExpInfos {
		imgSlice[k].ImgId = v[1]
		imgSlice[k].Title = v[2]
		imgSlice[k].Url = "https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + v[1]
		imgSlice[k].Praise = v[3]
		conf.maxCh <- 1
		go GetDetail(c, conf, imgSlice[k], false)
		k++
	}
	conf.waitGroup.Wait()
	conf.count.List = k

	c.Info("同步pixiv图片结束", nil, conf)

	return nil
}

func GetDetail(c *core.Context, conf *Conf, img Img, try bool) {
	defer conf.waitGroup.Done()
	defer func() {
		<-conf.maxCh
	}()
	logMap := make(map[string]interface{})
	logMap["conf"] = conf
	logMap["img"] = img
	logMap["try"] = try

	c.Info("开始获取图片资源", logMap, nil)

	req, _ := http.NewRequest("GET", img.Url, nil)
	req.Header.Set("cookie", conf.cookie)
	proxy, _ := url.Parse("socks5://127.0.0.1:1080")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	res, err := httplib.NewClient(httplib.WithTimeout(time.Second*60), httplib.WithTransport(tr)).Do(req)
	if err != nil {
		c.ErrorL("获取图片资源失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}
	if res == nil {
		c.ErrorL("获取图片资源返回nil", logMap, nil)
		conf.count.Failed += 1
		return
	}
	defer func() { _ = res.Body.Close() }()

	var buf []byte
	var content string
	buf, _ = io.ReadAll(res.Body)
	exp, err := regexp.Compile(`nal":"(.+?)"}`)
	contentArr := exp.FindStringSubmatch(string(buf))
	if len(contentArr) > 1 {
		content = contentArr[1]
		logMap["content"] = content
		c.Info("获取网页内容成功", logMap, nil)
	} else {
		c.ErrorL("未匹配到网页内容", logMap, string(buf))
		conf.count.Failed += 1
		return
	}

	exp, _ = regexp.Compile(`\\`)
	src := exp.ReplaceAllString(content, "")
	var suffix string
	exp, err = regexp.Compile(`p0(.+)`)
	suffixArr := exp.FindStringSubmatch(src)
	if len(suffixArr) > 1 {
		suffix = suffixArr[1]
		logMap["suffix"] = suffix
		c.Info("获取网页内容成功", logMap, nil)
	} else {
		c.ErrorL("未匹配到类型后缀", logMap, src)
		conf.count.Failed += 1
		return
	}

	exp, _ = regexp.Compile(`/`)
	effecTitle := exp.ReplaceAllString(img.Title+img.ImgId, "-")
	bucket, err := Bucket()
	if err != nil {
		c.ErrorL("打开bucket失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}
	isExist, err := bucket.IsObjectExist(effecTitle + suffix)
	if err != nil {
		c.ErrorL(effecTitle+"判断图片是否存在失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}
	if isExist == true {
		c.Info(effecTitle+suffix+"已存在", logMap, nil)
		conf.count.Exist += 1
		return
	}

	imgreq, _ := http.NewRequest("GET", src, nil)
	imgreq.Header.Set("cookie", conf.cookie)
	imgreq.Header.Set("Accept", accept)
	imgreq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	imgreq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7")
	imgreq.Header.Set("Referer", img.Url)
	imgreq.Header.Set("pragma", "no-cache")
	imgreq.Header.Set("Cache-Control", "no-cache")
	imgreq.Header.Set("User-Agent", userAgent)
	if try {
		var createDate string
		exp, _ = regexp.Compile(`createDate":"(.+?)",`)
		createDateArr := exp.FindStringSubmatch(string(buf))
		if len(createDateArr) > 1 {
			createDate = createDateArr[1]
			logMap["createDate"] = createDate
			c.Info("获取创建时间成功", logMap, nil)
		} else {
			c.ErrorL(effecTitle+"未匹配到创建时间", logMap, nil)
			conf.count.Failed += 1
			return
		}
		exp, _ = regexp.Compile(`T`)
		effCreateDate := exp.ReplaceAllString(createDate, " ")
		exp, _ = regexp.Compile(`\+.+`)
		date := exp.ReplaceAllString(effCreateDate, "")
		timestamp, _ := time.Parse("2006-01-02 15:04:05", date)
		GMTtime := timestamp.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		imgreq.Header.Set("Upgrade-Insecure-Requests", "1")
		imgreq.Header.Set("If-Modified-Since", GMTtime)
	}

	imgRes, err := httplib.NewClient(httplib.WithTimeout(time.Second*60), httplib.WithTransport(tr)).Do(imgreq)
	if err != nil {
		c.ErrorL(effecTitle+"imgres获取图片资源失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}
	if imgRes == nil {
		c.ErrorL(effecTitle+"imgres获取图片资源返回nil", logMap, nil)
		conf.count.Failed += 1
		return
	}
	defer func() { _ = imgRes.Body.Close() }()

	if imgRes.ContentLength == 0 {
		c.ErrorL("imgRes.ContentLength == 0", logMap, nil)
		return
	}

	imgBytes, err := io.ReadAll(imgRes.Body)
	if err != nil {
		c.ErrorL(effecTitle+"读取imgBytes失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}

	err = bucket.PutObject(effecTitle+suffix, bytes.NewReader(imgBytes))
	if err != nil {
		c.ErrorL(effecTitle+suffix+"上传oss失败", logMap, err.Error())
		conf.count.Failed += 1
		return
	}
	c.Info(effecTitle+suffix+"上传成功", logMap, nil)
	conf.count.Success += 1
}
