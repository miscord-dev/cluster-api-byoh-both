module.exports = async () => {
    const byohVersion = (await $`cat go.mod | grep github.com/vmware-tanzu/cluster-api-provider-bringyourownhost`)
        .stdout
        .match(/(v\d+\.\d+\.\d+)/)[0]
    const byohDir = `./byoh`;

    await $`git clone https://github.com/vmware-tanzu/cluster-api-provider-bringyourownhost.git -b ${byohVersion} --depth=1 ${byohDir}`

    await $`cp metadata.yaml ./${byohDir}/metadata.yaml`
    await $`cd ${byohDir} && IMG="ghcr.io/miscord-win/cluster-api-byoh-controller:${byohVersion}" make build-release-artifacts`

    await $`cd ${byohDir} && cat << EOF > Makefile
host-agent-binaries-arm64: ## Builds the binaries for the host-agent
    RELEASE_BINARY=./byoh-hostagent GOOS=linux GOARCH=arm64 GOLDFLAGS="\$(LDFLAGS) \$(STATIC)" \
    HOST_AGENT_DIR=./\$(HOST_AGENT_DIR) \$(MAKE) host-agent-binary`
    await $`cd ${byohDir} && make host-agent-binaries-arm64`

    await $`ls byoh/_dist`
}
