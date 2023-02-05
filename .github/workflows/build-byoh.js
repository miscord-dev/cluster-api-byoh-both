module.exports = async () => {
    const byohVersion = (await $`cat go.mod | grep github.com/vmware-tanzu/cluster-api-provider-bringyourownhost`)
        .stdout
        .match(/(v\d+\.\d+\.\d+)/)[0]
    const byohDir = `./byoh`;

    await $`git clone https://github.com/vmware-tanzu/cluster-api-provider-bringyourownhost.git -b ${byohVersion} --depth=1 ${byohDir}`

    await $`cd ${byohDir} && IMG="ghcr.io/miscord-win/cluster-api-byoh-controller:${byohVersion}" make build-release-artifacts`

    await $`ls byoh/_dist`
}
