FROM hashicorp/terraform:0.12.0
ARG SOPS_PLUGIN_VERSION
ENV SOPS_PLUGIN_VERSION=${SOPS_PLUGIN_VERSION}
RUN wget https://github.com/carlpett/terraform-provider-sops/releases/download/${SOPS_PLUGIN_VERSION}/terraform-provider-sops_${SOPS_PLUGIN_VERSION}_linux_amd64.zip && \
	unzip terraform-provider-sops_${SOPS_PLUGIN_VERSION}_linux_amd64.zip && \
	mv terraform-provider-sops_${SOPS_PLUGIN_VERSION} /bin/terraform-provider-sops_${SOPS_PLUGIN_VERSION} && \
	chmod +x /bin/terraform-provider-sops_${SOPS_PLUGIN_VERSION}
