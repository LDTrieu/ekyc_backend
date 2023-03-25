package portal

import (
	"ekyc-app/internal/fsdb"
	"ekyc-app/internal/wUtil"
	"ekyc-app/library/net"
	"ekyc-app/library/qrcode"
	"ekyc-app/package/wlog"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Reg(router *gin.Engine) {

	// Websocket - ping test
	router.GET("/ws/connect", websocketConnection)
	router.GET("/ping", ping)

	// Signin - Signup
	router.GET("/login/qr/download/:qrId/:reqId", downloadQR)
	router.GET("/login/qr/rend/:reqId", rendQRLogin)
	router.POST("/login/auth/:reqId", loginBasic)
	router.POST("/signup/auth/:reqId", signupBasic)

	// User admin
	router.GET("/portal/user/list/:reqId", filterListUser)

	// Student
	router.GET("/portal/student/list/:reqId", filterListStudent)
	//router.GET("/portal/user/detail/:reqId", filterListUser)
	router.GET("/portal/student/detail/:reqId", studentDetails)
	router.POST("/portal/student/create/:reqId", createStudentProfile)
	router.POST("/portal/student/update-ekyc/:reqId", updateStudentEkyc)
	router.POST("/portal/student/update/:reqId", updateStudent)
	//router.GET("/portal/student/list-pdf/:reqId", downloadStudentList)
	//router.POST("/portal/student/update-submit/:reqId", createStudentProfile)

	// Face auth session
	router.GET("/portal/session/list/:reqId", filterListSession)

	// Device
	router.GET("/portal/device/list/:reqId", filterListDevice)
	// Upload file
	router.POST("/portal/file/upload/face-image/:studentId/:reqId", uploadFaceImage)
	router.POST("/portal/file/upload/national-id-card/:reqId", uploadNationalIdImage)
	router.POST("/portal/file/upload/face-reg/:reqId", uploadFaceRegImage)

	// router.POST("/portal/file/upload/student-card/:reqId", uploadStudentCardImage)

}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OKE",
	})
}

func rendQRLogin(c *gin.Context) {
	var (
		request = rendQRLoginRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	response, err := guestRendQRLogin(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}
	// Trace client and result
	response.traceField = request.traceField
	c.JSON(http.StatusOK, response)
}

func downloadQR(c *gin.Context) {
	var (
		qrdata []byte
	)
	switch c.Query("action") {
	case qrActionLoginWebPortal:
		token, ok, err := fsdb.LoginSessionDBC.
			GetToken(c, c.Param("qrId"))
		if err != nil || !ok {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		qrdata = qrcode.GenQRCode(token)

	default:
		c.AbortWithError(http.StatusNotFound,
			errors.New("action undefine"))
		return
	}
	//WRITE PHOTO
	c.Writer.Header().Set("Content-Type",
		http.DetectContentType(qrdata))
	if _, err := c.Writer.Write(qrdata); err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}
}

/* */
func loginBasic(c *gin.Context) {
	var (
		request = loginBasicRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __loginBasic(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func signupBasic(c *gin.Context) {
	var (
		request = signupBasicRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
		}
	)
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __signupBasic(c.Request.Context(),
		&request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func filterListUser(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListUserRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListUser(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func filterListStudent(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListStudentRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListStudent(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func createStudentProfile(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = createStudentProfileRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	if err := c.BindJSON(&request.Payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __submitStudentProfile(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func updateStudentEkyc(c *gin.Context) {
	// validate token
	status, _, auth_data, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = updateStudentEkycRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: auth_data,
		}
	)

	if err := c.BindJSON(&request.Payload); err != nil {
		log.Println("err: ", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __updateStudentEkyc(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func studentDetails(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(
		c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = studentDetailsRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit:    from,
			StudentId: c.Query("studentId"),
		}
	)
	resp, err := __studentDetails(
		c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}
	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func updateStudent(c *gin.Context) {
	// validate token
	status, _, auth_data, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	// is_blocked, err := strconv.ParseBool(c.Query("isBlocked"))
	// if err != nil {
	// 	c.AbortWithError(status, err)
	// 	return
	// }

	var (
		request = updateStudentRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: auth_data,
			// Payload: update_student_data{
			// 	IsBlocked: is_blocked,
			// },
		}
	)
	log.Println("request.Payload: ", request.Payload)
	if err := c.BindJSON(&request.Payload); err != nil {
		log.Println("err: ", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := __updateStudent(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func uploadFaceImage(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}

	var (
		request = uploadFaceImageRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
			Payload: face_image_req{
				StudentId: c.Param("studentId"),
				FileName:  fmt.Sprintf("%s_%s.bin", c.Param("studentId"), primitive.NewObjectID().Hex()),
			},
		}
	)
	// validate image
	// Read multipart form files
	multipart_form, err := net.NewMultipartForm(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filename, err := multipart_form.GetForm("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	file_name, file, err := multipart_form.GetFile("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	request.Payload.FileName = func() string {
		if len(filename) > 0 {
			return filename
		}
		return file_name
	}()
	request.Payload.File = file
	// call API to Django

	// upload to Google Bucket
	resp, err := __uploadFaceImage(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	resp.Payload.URL = fmt.Sprintf("%s%s", wUtil.GetHost(c), resp.Payload.Path)
	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func uploadNationalIdImage(c *gin.Context) {
	// validate token
	// status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	// if err != nil {
	// 	c.AbortWithError(status, err)
	// 	return
	// }

	var (
		request = uploadNationalIdImageRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			//Permit: from,
			Payload: national_id_image_req{
				//StudentId: c.Param("studentId"),
				FileName: fmt.Sprintf("%s_%s.bin", "studentId", primitive.NewObjectID().Hex()),
			},
		}
	)
	// validate image
	// Read multipart form files
	multipart_form, err := net.NewMultipartForm(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filename, err := multipart_form.GetForm("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	file_name, file, err := multipart_form.GetFile("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	request.Payload.FileName = func() string {
		if len(filename) > 0 {
			return filename
		}
		return file_name
	}()
	request.Payload.File = file
	resp, err := __uploadNationalIdImage(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	//resp.Payload.URL = fmt.Sprintf("%s%s", wUtil.GetHost(c), resp.Payload.Path)
	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
func uploadFaceRegImage(c *gin.Context) {
	// validate token
	// status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	// if err != nil {
	// 	c.AbortWithError(status, err)
	// 	return
	// }

	var (
		request = uploadFaceRegImageRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			//Permit: from,
			Payload: face_reg_image_req{
				//StudentId: c.Param("studentId"),
				FileName: fmt.Sprintf("%s_%s.bin", "studentId", primitive.NewObjectID().Hex()),
			},
		}
	)
	// validate image
	// Read multipart form files
	multipart_form, err := net.NewMultipartForm(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filename, err := multipart_form.GetForm("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	file_name, file, err := multipart_form.GetFile("filename")
	if err != nil {
		wlog.Error(c, err)
	}
	request.Payload.FileName = func() string {
		if len(filename) > 0 {
			return filename
		}
		return file_name
	}()
	request.Payload.File = file
	resp, err := __uploadFaceRegImage(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	//resp.Payload.URL = fmt.Sprintf("%s%s", wUtil.GetHost(c), resp.Payload.Path)
	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)
}

/* */
// downloadStudentList

/* */
func filterListSession(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListSessionRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListSession(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}

/* */
func filterListDevice(c *gin.Context) {
	// validate token
	status, _, from, err := validateBearer(c.Request.Context(), c.Request)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	var (
		request = filterListDeviceRequest{
			traceField: traceField{
				RequestId: c.Param("reqId"),
			},
			Permit: from,
		}
	)
	// if err := c.BindJSON(&request.Payload); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	resp, err := __filterListDevice(c.Request.Context(), &request)
	if err != nil {
		wlog.Error(c, err)
	}

	// Trace client and result
	resp.traceField = request.traceField
	c.JSON(http.StatusOK, resp)

}
