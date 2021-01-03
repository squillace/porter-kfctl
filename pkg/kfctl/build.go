package kfctl

import (
	"fmt"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the kfctl mixin in porter.yaml
// mixins:
// - kfctl:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
// const dockerfileLines = `RUN apt-get update && \
// apt-get install gnupg apt-transport-https lsb-release software-properties-common -y && \
// echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ stretch main" | \
//    tee /etc/apt/sources.list.d/azure-cli.list && \
// apt-key --keyring /etc/apt/trusted.gpg.d/Microsoft.gpg adv \
// 	--keyserver packages.microsoft.com \
// 	--recv-keys BC528686B50D79E339D3721CEB3E94ADBE1229CF && \
// apt-get update && apt-get install azure-cli
// `

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	const dockerfileLines = `RUN apt-get update && apt-get install gnupg apt-transport-https lsb-release software-properties-common -y
		curl 
	
	`
	/*

	   	fmt.Fprintln(m.Out, `RUN wget --no-check-certificate https://raw.githubusercontent.com/stedolan/jq/master/sig/jq-release.key -O /tmp/jq-release.key && \
	       wget --no-check-certificate https://raw.githubusercontent.com/stedolan/jq/master/sig/v${JQ_VERSION}/jq-linux64.asc -O /tmp/jq-linux64.asc && \
	       wget --no-check-certificate https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 -O /tmp/jq-linux64 && \
	       gpg --import /tmp/jq-release.key && \
	       gpg --verify /tmp/jq-linux64.asc /tmp/jq-linux64 && \
	       cp /tmp/jq-linux64 /usr/bin/jq && \
	       chmod +x /usr/bin/jq && \
	       rm -f /tmp/jq-release.key && \
	       rm -f /tmp/jq-linux64.asc && \
	   	rm -f /tmp/jq-linux64`)
	*/

	// curl https://github.com/kubeflow/kfctl/releases/v1.2.0 | grep -o -m 1 "kfctl_v1.2.0.*linux.tar.gz" // does result in the single download file name.

	fmt.Fprintf(m.Out, `RUN curl -L https://github.com/kubeflow/kfctl/releases/download/v1.2.0/kfctl_v1.2.0-0-gbc038f9_linux.tar.gz --output kfctl.tar.gz && \
		tar -xvf kfctl.tar.gz && \
		cp kfctl /usr/bin/ && \
		curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.6.14 TARGET_ARCH=x86_64 sh - && \
		cp istio-1.6.14/bin/istioctl /usr/bin/
		
	
	`)

	// Example of pulling and defining a client version for your mixin
	// fmt.Fprintf(m.Out, "\nRUN curl -L https://github.com/kubeflow/kfctl/releases/download/v1.2.0/kfctl_v1.2.0-0-gbc038f9_linux.tar.gz --output kfctl.tar.gz", m.ClientVersion)

	return nil
}
