package main 

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"strconv"
	"crypto/sha1"
	"encoding/hex"
)

var (
	USER = "xxxxxxxxx"    //必填，飞鹅云后台注册的账号名
	UKEY = "xxxxxxxxxxxx" //必填，飞鹅云后台注册账号后生成的UKEY
	SN = "xxxxxxxxx"      //必填，打印机编号，必须要在管理后台里手动添加打印机或者通过API添加之后，才能调用API

	URL  = "http://api.feieyun.cn/Api/Open/"//不需要修改
)


//**********测试时，打开下面注释掉方法的即可,更多接口文档信息,请访问官网开放平台查看**********
func main(){

		//==================添加打印机接口（支持批量）==================
		//----------接口返回值说明----------
		//正确例子：{"msg":"ok","ret":0,"data":{"ok":["sn#key#remark#carnum","316500011#abcdefgh#快餐前台"],"no":["316500012#abcdefgh#快餐前台#13688889999  （错误：识别码不正确）"]},"serverExecutedTime":3}
		//错误：{"msg":"参数错误 : 该帐号未注册.","ret":-2,"data":null,"serverExecutedTime":37}

		//提示：打印机编号(必填) # 打印机识别码(必填) # 备注名称(选填) # 流量卡号码(选填)，多台打印机请换行（\n）添加新打印机信息，每次最多100行(台)。
		//snlist := "sn1#key1#remark1#carnum1\nsn2#key2#remark2#carnum2"
		//addprinter(snlist)



		//==================方法1.打印订单==================
		//----------接口返回值说明----------
		//正确例子：{"msg":"ok","ret":0,"data":"xxxx_xxxx_xxxxxxxxx","serverExecutedTime":6}
		//错误：{"msg":"错误信息.","ret":非零错误码,"data":null,"serverExecutedTime":5}
		
		//print(SN)
		
		
		
		//===========方法2.查询某订单是否打印成功=============
		//----------接口返回值说明----------
		//已打印：{"msg":"ok","ret":0,"data":true,"serverExecutedTime":6}
		//未打印：{"msg":"ok","ret":0,"data":false,"serverExecutedTime":6}

		//strorderid := "xxxxxxxx_xxxxxxxxx_xxxxxxxxx"//订单id，由方法1返回
		//queryOrderState(strorderid)

		
	
		//===========方法3.查询指定打印机某天的订单详情============
		//----------接口返回值说明----------
		//正确例子：{"msg":"ok","ret":0,"data":{"print":6,"waiting":1},"serverExecutedTime":9}
		//错误例子：{"msg":"参数错误 : 时间格式不正确。","ret":1001,"data":null,"serverExecutedTime":37}

		//strdate := "2017-04-02"//注意日期格式为yyyy-MM-dd
		//queryOrderInfoByDate(SN,strdate)




		//===========方法4.查询打印机的状态==========================
		//提示：由于获取到打印机状态有延时，不建议使用本接口作为发单的依据
		//如果有订单数据要打印，直接调用方法1传过来即可，不必先调用本接口获取打印机状态

		//----------接口返回值说明-----------
		//提示：返回的JOSN中文是编码过的
		//{"msg":"ok","ret":0,"data":"离线","serverExecutedTime":9}
		//{"msg":"ok","ret":0,"data":"在线，工作状态正常","serverExecutedTime":9}
		//{"msg":"ok","ret":0,"data":"在线，工作状态不正常","serverExecutedTime":9}

		//queryPrinterStatus(SN)
		

}


//==============================================================================

	func addprinter(snlist string){

		itime := time.Now().Unix()
		stime := strconv.FormatInt(itime,10)
		sig := SHA1(USER+UKEY+stime)//生成签名

		client := http.Client{}
		postValues := url.Values{}
		postValues.Add("user",USER)//账号名
		postValues.Add("stime",stime)//当前时间的秒数，请求时间
		postValues.Add("sig",sig)//签名
		postValues.Add("apiname","Open_printerAddlist")//固定
		postValues.Add("printerContent",snlist)//打印机

		res,_ := client.PostForm(URL, postValues)
		data,_ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))//服务器返回的JSON字符串，建议要当做日志记录起来
		res.Body.Close()
	}



	func print(sn string){
			//标签说明：
			//单标签:
			//"<BR>"为换行,"<CUT>"为切刀指令(主动切纸,仅限切刀打印机使用才有效果)
			//"<LOGO>"为打印LOGO指令(前提是预先在机器内置LOGO图片),"<PLUGIN>"为钱箱或者外置音响指令
			//成对标签：
			//"<CB></CB>"为居中放大一倍,"<B></B>"为放大一倍,"<C></C>"为居中,<L></L>字体变高一倍
			//<W></W>字体变宽一倍,"<QR></QR>"为二维码,"<BOLD></BOLD>"为字体加粗,"<RIGHT></RIGHT>"为右对齐
			//拼凑订单内容时可参考如下格式
			//根据打印纸张的宽度，自行调整内容的格式，可参考下面的样例格式

			content := "<CB>测试打印</CB><BR>"
			content += "名称　　　　　 单价  数量 金额<BR>"
			content += "--------------------------------<BR>"
			content += "饭　　　　　 　10.0   10  10.0<BR>"
			content += "炒饭　　　　　 10.0   10  10.0<BR>"
			content += "蛋炒饭　　　　 10.0   100 100.0<BR>"
			content += "鸡蛋炒饭　　　 100.0  100 100.0<BR>"
			content += "西红柿炒饭　　 1000.0 1   100.0<BR>"
			content += "西红柿蛋炒饭　 100.0  100 100.0<BR>"
			content += "西红柿鸡蛋炒饭 15.0   1   15.0<BR>"
			content += "备注：加辣<BR>"
			content += "--------------------------------<BR>"
			content += "合计：xx.0元<BR>"
			content += "送货地点：广州市南沙区xx路xx号<BR>"
			content += "联系电话：13888888888888<BR>"
			content += "订餐时间：2014-08-08 08:08:08<BR>"
			content += "<QR>http://www.dzist.com</QR>"

			itime := time.Now().Unix()
			stime := strconv.FormatInt(itime,10)
			sig := SHA1(USER+UKEY+stime)//生成签名

			client := http.Client{}
			postValues := url.Values{}
		 	postValues.Add("user",USER)//账号名
		 	postValues.Add("stime",stime)//当前时间的秒数，请求时间
			postValues.Add("sig",sig)//签名
			postValues.Add("apiname","Open_printMsg")//固定
			postValues.Add("sn",sn)//打印机编号
		 	postValues.Add("content",content)//打印内容
		 	postValues.Add("times","1")//打印次数
		
			res,_ := client.PostForm(URL, postValues)
			data,_ := ioutil.ReadAll(res.Body)
			fmt.Println(string(data))//服务器返回的JSON字符串，建议要当做日志记录起来
			res.Body.Close()
	}
	
	
	func queryOrderState(strorderid string){
		itime := time.Now().Unix()
		stime := strconv.FormatInt(itime,10)
		sig := SHA1(USER+UKEY+stime)//生成签名

		client := http.Client{}
		postValues := url.Values{}
		postValues.Add("user",USER)//账号名
		postValues.Add("stime",stime)//当前时间的秒数，请求时间
		postValues.Add("sig",sig)//签名
		postValues.Add("apiname","Open_queryOrderState")//固定
		postValues.Add("orderid",strorderid)//订单ID由方法1返回

		res,_ := client.PostForm(URL, postValues)
		data,_ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		res.Body.Close()
	}
	
	
	func queryOrderInfoByDate(sn string,strdate string){
		itime := time.Now().Unix()
		stime := strconv.FormatInt(itime,10)
		sig := SHA1(USER+UKEY+stime)//生成签名

		client := http.Client{}
		postValues := url.Values{}
		postValues.Add("user",USER)//账号名
		postValues.Add("stime",stime)//当前时间的秒数，请求时间
		postValues.Add("sig",sig)//签名
		postValues.Add("apiname","Open_queryOrderInfoByDate")//固定
		postValues.Add("sn",sn)//打印机编号
		postValues.Add("date",strdate)//日期字符串

		res,_ := client.PostForm(URL, postValues)
		data,_ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		res.Body.Close()
	}
	

	func queryPrinterStatus(sn string){
		itime := time.Now().Unix()
		stime := strconv.FormatInt(itime,10)
		sig := SHA1(USER+UKEY+stime)//生成签名

		client := http.Client{}
		postValues := url.Values{}
		postValues.Add("user",USER)//账号名
		postValues.Add("stime",stime)//当前时间的秒数，请求时间
		postValues.Add("sig",sig)//签名
		postValues.Add("apiname","Open_queryPrinterStatus")//固定
		postValues.Add("sn",sn)//打印机编号

		res,_ := client.PostForm(URL, postValues)
		data,_ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		res.Body.Close()
	}


	func SHA1(str string) string {
		s := sha1.Sum([]byte(str))
		strsha1:= hex.EncodeToString(s[:])
		return strsha1
	}