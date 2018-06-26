FROM hashicorp/terraform:light
ENV SOPS_PLUGIN_VERSION=v0.0.1
RUN curl -sLo /bin/terraform-provider-sops_${SOPS_PLUGIN_VERSION} https://github.com/carlpett/terraform-sops/releases/download/${SOPS_PLUGIN_VERSION}/terraform-provider-sops_${SOPS_PLUGIN_VERSION}-linux-amd64 && \
	chmod +x /bin/terraform-provider-sops_${SOPS_PLUGIN_VERSION}
