/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package impl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportAPIPolicyFromEnv function is used with export policy rate-limiting command
func ExportAPIPolicyFromEnv(accessToken string, exportEnvironment string, apiPolicyName string, apiPolicyVersion string) (*resty.Response, error) {
	apiPolicyEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	// var query string
	apiPolicyEndpoint = utils.AppendSlashToString(apiPolicyEndpoint)
	// apiPolicyResource := "api-policies/c86da87e-da70-4977-bed2-57cb089c115f" + "/content"
	apiPolicyResource := "operation-policies/export?"

	query := `name=` + apiPolicyName + `&version=` + apiPolicyVersion

	apiPolicyResource += query
	url := apiPolicyEndpoint + apiPolicyResource
	utils.Logln(utils.LogPrefixInfo+"ExportAPIPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WriteAPIPolicyToFile writes the policy to a specified location
func WriteAPIPolicyToFile(exportLocationPath string, resp *resty.Response, exportAPIPolicyVersion string, exportAPIPolicyName string,
	runningExportThrottlePolicyCommand bool) {
	err := utils.CreateDirIfNotExist(exportLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+exportLocationPath, err)
	}
	zipFileName := exportAPIPolicyName + "_" + exportAPIPolicyVersion + ".zip"
	zipFile := filepath.Join(exportLocationPath, zipFileName)

	err = ioutil.WriteFile(zipFile, resp.Body(), 0644)
	if err != nil {
		return
	}

	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported API", err)
	}

	if runningExportThrottlePolicyCommand {
		fmt.Println("Successfully exported API Policy!")
		fmt.Println("Find the exported API Policies at " +
			utils.AppendSlashToString(exportLocationPath) + zipFileName)
	}
}
