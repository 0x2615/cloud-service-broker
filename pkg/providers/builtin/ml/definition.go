// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ml

import (
	"code.cloudfoundry.org/lager"
	"github.com/pivotal/cloud-service-broker/pkg/broker"
	accountmanagers "github.com/pivotal/cloud-service-broker/pkg/providers/builtin/account_managers"
	"github.com/pivotal/cloud-service-broker/pkg/providers/builtin/base"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2/jwt"
)

// ServiceDefinition creates a new ServiceDefinition object for the ML service.
func ServiceDefinition() *broker.ServiceDefinition {
	roleWhitelist := []string{
		"ml.developer",
		"ml.viewer",
		"ml.modelOwner",
		"ml.modelUser",
		"ml.jobOwner",
		"ml.operationOwner",
	}

	return &broker.ServiceDefinition{
		Id:               "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
		Name:             "google-ml-apis",
		Description:      "Machine Learning APIs including Vision, Translate, Speech, and Natural Language.",
		DisplayName:      "Google Machine Learning APIs",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/machine-learning.svg",
		DocumentationUrl: "https://cloud.google.com/ml/",
		SupportUrl:       "https://cloud.google.com/support/",
		Tags:             []string{"gcp", "ml"},
		Bindable:         true,
		PlanUpdateable:   false,
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:          "be7954e1-ecfb-4936-a0b6-db35e6424c7a",
					Name:        "default",
					Description: "Machine Learning API default plan.",
					Free:        brokerapi.FreeValue(false),
				},
				ServiceProperties: map[string]string{},
			},
		},
		ProvisionInputVariables: []broker.BrokerVariable{},
		DefaultRoleWhitelist:    roleWhitelist,
		BindInputVariables:      accountmanagers.ServiceAccountWhitelistWithDefault(roleWhitelist, "ml.modelUser"),
		BindOutputVariables:     accountmanagers.ServiceAccountBindOutputVariables(),
		BindComputedVariables:   accountmanagers.ServiceAccountBindComputedVariables(),
		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create an account with developer access to your ML models.",
				PlanId:          "be7954e1-ecfb-4936-a0b6-db35e6424c7a",
				ProvisionParams: map[string]interface{}{},
				BindParams: map[string]interface{}{
					"role": "ml.developer",
				},
			},
		},
		ProviderBuilder: func(projectId string, auth *jwt.Config, logger lager.Logger) broker.ServiceProvider {
			bb := base.NewBrokerBase(projectId, auth, logger)
			return &ApiServiceBroker{BrokerBase: bb}
		},
		IsBuiltin: true,
	}
}
