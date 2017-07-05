package integrationtest_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/widuu/gojson"
	"net/http"
	"strconv"
	"strings"

	"bytes"
	"fmt"
	. "github.com/servicecomb/service-center/tests/integrationtest"
	"io/ioutil"
)

var _ = Describe("MicroService Api Test", func() {
	var serviceName = "integrationtestInstances"
	var serviceId = ""
	var serviceAppId = "integrationtestAppIdInstance"
	var serviceVersion = "0.0.2"
	var serviceInstanceID = ""
	Context("Tesing MicroService Instances API's", func() {
		BeforeEach(func() {
			schema := []string{"testSchema"}
			properties := map[string]string{"attr1": "aa"}
			servicemap := map[string]interface{}{
				"serviceName": serviceName,
				"appId":       serviceAppId,
				"version":     serviceVersion,
				"description": "examples",
				"level":       "FRONT",
				"schemas":     schema,
				"status":      "UP",
				"properties":  properties,
			}
			bodyParams := map[string]interface{}{
				"service": servicemap,
			}
			body, _ := json.Marshal(bodyParams)
			bodyBuf := bytes.NewReader(body)
			req, _ := http.NewRequest(POST, SCURL+REGISTERMICROSERVICE, bodyBuf)
			req.Header.Set("X-tenant-name", "default")
			resp, err := scclient.Do(req)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			// Validate the service creation
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			respbody, _ := ioutil.ReadAll(resp.Body)
			serviceId = gojson.Json(string(respbody)).Get("serviceId").Tostring()
			serviceIdInt, _ := strconv.Atoi(serviceId)
			Expect(serviceIdInt).Should(BeNumerically(">", 1))

			//Register MicroService Instance
			endpoints := []string{"cse://127.0.0.1:9984"}
			propertiesInstance := map[string]interface{}{
				"_TAGS":  "A,B",
				"attr1":  "a",
				"nodeIP": "one",
			}
			healthcheck := map[string]interface{}{
				"mode":     "push",
				"interval": 30,
				"times":    2,
			}
			instance := map[string]interface{}{
				"endpoints":   endpoints,
				"hostName":    "cse",
				"status":      "UP",
				"stage":       "prod",
				"properties":  propertiesInstance,
				"healthCheck": healthcheck,
			}

			bodyParams = map[string]interface{}{
				"instance": instance,
			}
			url := strings.Replace(REGISTERINSTANCE, ":serviceId", serviceId, 1)
			body, _ = json.Marshal(bodyParams)
			bodyBuf = bytes.NewReader(body)
			req, _ = http.NewRequest(POST, SCURL+url, bodyBuf)
			req.Header.Set("X-tenant-name", "default")
			resp, err = scclient.Do(req)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			// Validate the instance registration
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			respbody, _ = ioutil.ReadAll(resp.Body)
			serviceInstanceID = gojson.Json(string(respbody)).Get("instanceId").Tostring()
			serviceIdInt, _ = strconv.Atoi(serviceInstanceID)
			Expect(serviceIdInt).Should(BeNumerically(">", 1))

		})

		AfterEach(func() {
			if serviceInstanceID != "" {
				url := strings.Replace(UNREGISTERINSTANCE, ":serviceId", serviceId, 1)
				url = strings.Replace(url, ":instanceId", serviceInstanceID, 1)
				req, _ := http.NewRequest(DELETE, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				resp, _ := scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}

			if serviceId != "" {
				url := strings.Replace(UNREGISTERMICROSERVICE, ":serviceId", serviceId, 1)
				req, _ := http.NewRequest(DELETE, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				resp, _ := scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}

		})

		By("Register MicroService Instance API", func() {
			It("Register MicroService Instance with invalid params", func() {
				instance := map[string]interface{}{
					"hostName": "cse",
					"status":   "UP",
					"stage":    "prod",
				}

				bodyParams := map[string]interface{}{
					"instancse": instance,
				}
				url := strings.Replace(REGISTERINSTANCE, ":serviceId", serviceId, 1)
				body, _ := json.Marshal(bodyParams)
				bodyBuf := bytes.NewReader(body)
				req, _ := http.NewRequest(POST, SCURL+url, bodyBuf)
				req.Header.Set("X-tenant-name", "default")
				resp, err := scclient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("Register MicroService Instance with duplicate Params", func() {
				endpoints := []string{"cse://127.0.0.1:9984"}
				propertiesInstance := map[string]interface{}{
					"_TAGS":  "A,B",
					"attr1":  "a",
					"nodeIP": "one",
				}
				healthcheck := map[string]interface{}{
					"mode":     "push",
					"interval": 30,
					"times":    2,
				}
				instance := map[string]interface{}{
					"endpoints":   endpoints,
					"hostName":    "cse",
					"status":      "UP",
					"stage":       "prod",
					"properties":  propertiesInstance,
					"healthCheck": healthcheck,
				}

				bodyParams := map[string]interface{}{
					"instance": instance,
				}
				url := strings.Replace(REGISTERINSTANCE, ":serviceId", serviceId, 1)
				body, _ := json.Marshal(bodyParams)
				bodyBuf := bytes.NewReader(body)
				req, _ := http.NewRequest(POST, SCURL+url, bodyBuf)
				req.Header.Set("X-tenant-name", "default")
				resp, err := scclient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				// Validate the instance registration
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				respbody, _ := ioutil.ReadAll(resp.Body)

				//Validate the instance id is same as the old one
				Expect(gojson.Json(string(respbody)).Get("instanceId").Tostring()).To(Equal(serviceInstanceID))
			})

			It("Register MicroService Instance with valid params", func() {
				endpoints := []string{"cse://127.0.0.1:9985"}
				propertiesInstance := map[string]interface{}{
					"_TAGS":  "A,B",
					"attr1":  "a",
					"nodeIP": "one",
				}
				healthcheck := map[string]interface{}{
					"mode":     "push",
					"interval": 30,
					"times":    2,
				}
				instance := map[string]interface{}{
					"endpoints":   endpoints,
					"hostName":    "cse",
					"status":      "UP",
					"stage":       "prod",
					"properties":  propertiesInstance,
					"healthCheck": healthcheck,
				}

				bodyParams := map[string]interface{}{
					"instance": instance,
				}
				url := strings.Replace(REGISTERINSTANCE, ":serviceId", serviceId, 1)
				body, _ := json.Marshal(bodyParams)
				bodyBuf := bytes.NewReader(body)
				req, _ := http.NewRequest(POST, SCURL+url, bodyBuf)
				req.Header.Set("X-tenant-name", "default")
				resp, err := scclient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				// Validate the instance registration
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				respbody, _ := ioutil.ReadAll(resp.Body)
				serviceInst := gojson.Json(string(respbody)).Get("instanceId").Tostring()

				//Verify the instanceID is different for two instance
				Expect(serviceInst).NotTo(Equal(serviceInstanceID))
				//Delete Instance
				url = strings.Replace(UNREGISTERINSTANCE, ":serviceId", serviceId, 1)
				url = strings.Replace(url, ":instanceId", serviceInst, 1)
				req, _ = http.NewRequest(DELETE, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				resp, _ = scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		By("Discover MicroService Instance API", func() {
			It("Find Micro-service Info by AppID", func() {
				req, _ := http.NewRequest(GET, SCURL+FINDINSTANCE+"?appId="+serviceAppId+"&serviceName="+serviceName+"&version="+serviceVersion, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				respbody, _ := ioutil.ReadAll(resp.Body)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				servicesStruct := map[string][]map[string]interface{}{}
				json.Unmarshal(respbody, &servicesStruct)
				foundMicroServiceInstance := false
				for _, services := range servicesStruct["instances"] {
					if services["instanceId"] == serviceInstanceID {
						foundMicroServiceInstance = true
						break
					}
				}
				Expect(foundMicroServiceInstance).To(Equal(true))
			})

			It("Find Micro-service Info by invalid AppID", func() {
				req, _ := http.NewRequest(GET, SCURL+FINDINSTANCE+"?appId=XXXX&serviceName="+serviceName+"&version="+serviceVersion, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("Find Micro-Service Instance by ServiceId", func() {
				url := strings.Replace(GETINSTANCE, ":serviceId", serviceId, 1)
				req, _ := http.NewRequest(GET, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				respbody, _ := ioutil.ReadAll(resp.Body)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				servicesStruct := map[string][]map[string]interface{}{}
				json.Unmarshal(respbody, &servicesStruct)
				foundMicroServiceInstance := false
				for _, services := range servicesStruct["instances"] {
					if services["instanceId"] == serviceInstanceID {
						foundMicroServiceInstance = true
						break
					}
				}
				Expect(foundMicroServiceInstance).To(Equal(true))
			})

			It("Find Micro-Service Instance by Invalid ServiceId", func() {
				url := strings.Replace(GETINSTANCE, ":serviceId", "XX", 1)
				req, _ := http.NewRequest(GET, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("Find MicroServiceInstance with Service and IstanceID", func() {
				url := strings.Replace(GETINSTANCEBYINSTANCEID, ":serviceId", serviceId, 1)
				url = strings.Replace(url, ":instanceId", serviceInstanceID, 1)
				req, _ := http.NewRequest(GET, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				respbody, _ := ioutil.ReadAll(resp.Body)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				servicesStruct := map[string]map[string]interface{}{}
				json.Unmarshal(respbody, &servicesStruct)
				instance := servicesStruct["instance"]
				Expect(instance["instanceId"]).To(Equal(serviceInstanceID))
			})

			It("Find Micro-Service Instance by Invalid InstanceID", func() {
				url := strings.Replace(GETINSTANCEBYINSTANCEID, ":serviceId", serviceId, 1)
				url = strings.Replace(url, ":instanceId", "XX", 1)
				req, _ := http.NewRequest(GET, SCURL+url, nil)
				req.Header.Set("X-tenant-name", "default")
				req.Header.Set("X-ConsumerId", serviceId)
				resp, _ := scclient.Do(req)
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

})
