const getVersion = () => {
    const byohVersion = 
        (`${process.env.GITHUB_REF}`.match(/(v\d+\.\d+\.\d+)/) ?? [])[0];

    return byohVersion;
}

exports.release = async () => {
    await $`make release IMG=ghcr.io/miscord-dev/cluster-api-byoh-both:${getVersion()}`
}
