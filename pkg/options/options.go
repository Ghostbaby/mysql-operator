/*
Copyright 2018 Pressinfra SRL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"

	"github.com/presslabs/mysql-operator/pkg/util"
)

// nolint: unparam
func getFromEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	return value
}

// Options is the data structure that contains information about mysql operator configuration
type Options struct {
	HelperImage string

	MetricsExporterImage string

	ImagePullSecretName string
	ImagePullPolicy     corev1.PullPolicy

	OrchestratorURI              string
	OrchestratorTopologyPassword string
	OrchestratorTopologyUser     string

	JobCompleteSuccessGraceTime time.Duration

	HTTPServeAddr string

	LeaderElectionNamespace string
	LeaderElectionID        string
}

type pullpolicy corev1.PullPolicy

func (pp *pullpolicy) String() string {
	return string(*pp)
}

func (pp *pullpolicy) Set(value string) error {
	*pp = pullpolicy(value)
	return nil
}

func (pp *pullpolicy) Type() string {
	return "v1.PullPolicy"
}

// nolint: unparam
func newPullPolicyValue(defaultValue corev1.PullPolicy, v *corev1.PullPolicy) *pullpolicy {
	*v = defaultValue
	return (*pullpolicy)(v)
}

const (
	defaultExporterImage = "prom/mysqld-exporter:latest"

	defaultImagePullPolicy = corev1.PullIfNotPresent

	defaultOrchestratorTopologyUser     = ""
	defaultOrchestratorTopologyPassword = ""

	defaultHTTPServerAddr = ":80"

	defaultLeaderElectionNamespace = "default"
	defaultLeaderElectionID        = ""
)

var (
	defaultHelperImage  = "quay.io/presslabs/mysql-helper:" + util.AppVersion
	defaultJobGraceTime = 24 * time.Hour
)

// AddFlags registers all mysql-operator needed flags
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.HelperImage, "helper-image", defaultHelperImage,
		"The image that instrumentate mysql.")

	fs.StringVar(&o.MetricsExporterImage, "metrics-exporter-image", defaultExporterImage,
		"The image for mysql metrics exporter.")
	fs.StringVar(&o.ImagePullSecretName, "image-pull-secret", "",
		"The secret name for used as pull secret.")

	fs.VarP(newPullPolicyValue(defaultImagePullPolicy, &o.ImagePullPolicy),
		"image-pull-policy", "", "Set image pull policy.")

	fs.StringVar(&o.OrchestratorURI, "orchestrator-uri", "",
		"The orchestrator uri")
	fs.StringVar(&o.OrchestratorTopologyPassword, "orchestrator-topology-password", defaultOrchestratorTopologyUser,
		"The orchestrator topology password. Can also be set as ORC_TOPOLOGY_PASSWORD environment variable.")
	fs.StringVar(&o.OrchestratorTopologyUser, "orchestrator-topology-user", defaultOrchestratorTopologyPassword,
		"The orchestrator topology user. Can also be set as ORC_TOPOLOGY_USER environment variable.")
	fs.DurationVar(&o.JobCompleteSuccessGraceTime, "job-grace-time", defaultJobGraceTime,
		"The time in hours how jobs after completion are keept.")

	fs.StringVar(&o.HTTPServeAddr, "http-serve-addr", defaultHTTPServerAddr,
		"The address for http server.")

	fs.StringVar(&o.LeaderElectionNamespace, "leader-election-namespace", defaultLeaderElectionNamespace,
		"The leader election namespace.")
	fs.StringVar(&o.LeaderElectionID, "leader-election-id", defaultLeaderElectionID,
		"The leader election id.")
}

var instance *Options
var once sync.Once

// GetOptions returns a singlethon that contains all options
func GetOptions() *Options {
	once.Do(func() {
		instance = &Options{
			HelperImage:          defaultHelperImage,
			MetricsExporterImage: defaultExporterImage,

			ImagePullPolicy:             defaultImagePullPolicy,
			JobCompleteSuccessGraceTime: defaultJobGraceTime,

			OrchestratorTopologyUser:     "",
			OrchestratorTopologyPassword: "",

			HTTPServeAddr: defaultHTTPServerAddr,
		}
	})

	return instance
}

// Validate validate the command line values
func (o *Options) Validate() error {

	if len(o.OrchestratorTopologyUser) == 0 {
		o.OrchestratorTopologyUser = getFromEnvOrDefault("ORC_TOPOLOGY_USER", "")
	}
	if len(o.OrchestratorTopologyPassword) == 0 {
		o.OrchestratorTopologyPassword = getFromEnvOrDefault("ORC_TOPOLOGY_PASSWORD", "")
	}
	return nil
}
