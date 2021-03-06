// Copyright 2017 The Kubernetes Dashboard Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package statefulset

import (
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/api"
	metricapi "github.com/kubernetes/dashboard/src/app/backend/integration/metric/api"
	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/pod"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
)

// StatefulSetDetail is a presentation layer view of Kubernetes Pet Set resource. This means
// it is Pet Set plus additional augmented data we can get from other sources
// (like services that target the same pods).
type StatefulSetDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Aggregate information about pods belonging to this Pet Set.
	PodInfo common.PodInfo `json:"podInfo"`

	// Detailed information about Pods belonging to this Pet Set.
	PodList pod.PodList `json:"podList"`

	// Container images of the Pet Set.
	ContainerImages []string `json:"containerImages"`

	// List of events related to this Pet Set.
	EventList common.EventList `json:"eventList"`
}

// GetStatefulSetDetail gets pet set details.
func GetStatefulSetDetail(client *k8sClient.Clientset, metricClient metricapi.MetricClient,
	namespace, name string) (*StatefulSetDetail, error) {

	log.Printf("Getting details of %s service in %s namespace", name, namespace)

	// TODO(floreks): Use channels.
	statefulSetData, err := client.AppsV1beta1().StatefulSets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podList, err := GetStatefulSetPods(client, metricClient, dataselect.DefaultDataSelectWithMetrics, name, namespace)
	if err != nil {
		return nil, err
	}

	podInfo, err := getStatefulSetPodInfo(client, statefulSetData)
	if err != nil {
		return nil, err
	}

	events, err := GetStatefulSetEvents(client, dataselect.DefaultDataSelect, statefulSetData.Namespace, statefulSetData.Name)
	if err != nil {
		return nil, err
	}

	statefulSet := getStatefulSetDetail(statefulSetData, metricClient, *events, *podList, *podInfo)
	return &statefulSet, nil
}

func getStatefulSetDetail(statefulSet *apps.StatefulSet, metricClient metricapi.MetricClient,
	eventList common.EventList, podList pod.PodList, podInfo common.PodInfo) StatefulSetDetail {

	return StatefulSetDetail{
		ObjectMeta:      api.NewObjectMeta(statefulSet.ObjectMeta),
		TypeMeta:        api.NewTypeMeta(api.ResourceKindStatefulSet),
		ContainerImages: common.GetContainerImages(&statefulSet.Spec.Template.Spec),
		PodInfo:         podInfo,
		PodList:         podList,
		EventList:       eventList,
	}
}
