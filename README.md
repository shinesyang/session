
	说明:
	为什么要自己写一个session中间件
	网上使用最多的session中间件是: github.com/gin-contrib/sessions(这里不做源码分析了)
	默认情况下，github.com/gin-contrib/sessions 验证调的是两个方法:
	返回header: set-cookie  字段
	获取header: cookie      字段,

	存在的问题:
	1. 在vue3作为前后端分离的相关项目中,无法通过:axios.interceptors.response.use 获取到set-cookie这个字段,
		那么无法判断用户登录之后的界面状态
	2. vue3作为前后端分离的相关项目中使用google调试时,不管是在开发环境还是编译之后的正式环境都无法自动传cookie这个header,
		手动指定时又会有报错提示: Refused to set unsafe header "Cookie"

	那么解决:
	1. 用户可以指定自己的请求头和响应头进行登录验证

	使用:
	```go
        /*
           创建: 基于cookie的存储引擎，secret 参数是用于加密的密钥，可以随便填写
        */
        store := session.NewStore([]byte("seasas"))

        store.Options(session.Option{
            MaxAge:    30 * 86400,
            SetHeader: "Set-Auth", // 自定义头
            Header:    "Auth",
        })

        /*
            应用: 设置session中间件，参数SESSIONID，指的是session的名字，也是cookie的名字，store是前面创建的存储引擎
        */
        router.Use(session.Sessions("SESSIONID", store))
	```
	其他的调用方法和： github.com/gin-contrib/sessions基本一致

	注意点:
	1. 此程序没有实现github.com/gin-contrib/sessions 将session存储数据库相关功能

	跨域问题:
	1. 设置自定义的请求/响应头部需要在跨域里面放行

	```go
        func Cors() gin.HandlerFunc {
        	return func(c *gin.Context) {
        		method := c.Request.Method
        		origin := c.Request.Header.Get("Origin")
        		if origin != "" {
        			c.Header("Access-Control-Allow-Origin", origin)
        			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        			//c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
        			c.Header("Access-Control-Allow-Credentials", "true")
        			c.Set("content-type", "application/json")
        			// 可以携带的请求头
        			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept,Auth")
        			// 返回给前端的请求头
        			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control,Set-Auth")
        		}
        		//拦截所有OPTIONS方法
        		if method == "OPTIONS" {
        			c.AbortWithStatus(http.StatusNoContent)
        		}
        		c.Next()
        	}
        }

       /*
       	1. Access-Control-Allow-Headers 必须设置session.Option.Header值一致
       	2. Access-Control-Expose-Headers 必须设置session.Option.SetHeader值一致
       */
	```